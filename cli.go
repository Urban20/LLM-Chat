package main

import (
	consola "LLM-Chat/ansi"
	"LLM-Chat/menu"
	"LLM-Chat/prompts"
	"LLM-Chat/utilidades"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"runtime"
	"sync"
	"time"
)

const (
	HOST_DEFAULT    = "localhost"
	PUERTO_DEFAULT  = 11434
	CTX_DEFAULT     = 16_000
	TEMP_DEFAULT    = 0.5
	CONTENT_TYPE    = "application/json"
	TIMEOUT_DEFAULT = 5.0
)

var (
	conserr      = consola.Iniciar_ANSI()
	host_selec   = flag.String("host", HOST_DEFAULT, "url al enpoint de Ollama")
	puerto_selec = flag.Int("puerto", PUERTO_DEFAULT, "puerto donde se escucha el endpoint")
	ctx          = flag.Int("ctx", CTX_DEFAULT, "cantidad de contexto que usa el LLM")
	temp         = flag.Float64("temp", TEMP_DEFAULT, "temperatura del LLM")
	timeout      = flag.Float64("timeout", TIMEOUT_DEFAULT, "tiempo de espera para la conexion inicial")
	Output       = flag.Bool("o", false, "guarda las respuestas del LLM en un archivo .md")
)

func iniciar_conversacion(archivo_prompt utilidades.Prompt_archivo, box utilidades.Box_info, endpoint, content_type string, chat bool, imagenes []string) error {

	prompt := utilidades.Input_multilinea("Prompt") //prompt del usuario

	if utilidades.String_vacio(prompt) {

		return errors.New("prompt vacio")
		// esta pensado para salir del input

	}

	carga := menu.Crear_carga()

	wg := sync.WaitGroup{}

	go carga.Iniciar(&wg)

	if err := prompts.Comunicacion(archivo_prompt.Prompt, prompt, box, endpoint, content_type, &carga, &wg, chat, imagenes); err != nil {
		fmt.Print("\n")
		utilidades.Advertencia(err)
		time.Sleep(time.Second * utilidades.TIEMPO_PAUSA)
	}

	return nil

}

func iniciar_prompts(url, content_type string, box utilidades.Box_info) {

	opciones := []string{fmt.Sprintf("%sVolver%s", utilidades.ANIL, utilidades.RESET), "Borrar contexto", "Adjuntar archivos de texto plano", "Eliminar archivos adjuntos", "Adjuntar imagen", "Ingresar prompt"}

	api_chat := fmt.Sprintf("%s/chat", url)
	api_generate := fmt.Sprintf("%s/generate", url)

	var archivo_prompt utilidades.Prompt_archivo

	for {
		// TODO : quiza modifique esto
		box.Box_informacion()

		seleccion, _ := menu.Menu(len(opciones), opciones...)

		switch seleccion {

		case opciones[0]:

			prompts.Borrar_memoria()
			prompts.Descargar_modelo(box.Modelo, content_type, api_chat)

			return

		case opciones[1]:

			utilidades.Limpieza_rapida()
			prompts.Borrar_memoria()

		case opciones[2]:

			arch_list := utilidades.Abrir_selector_archivos()

			if len(arch_list) == 0 {

				continue
			}

			arch_list = utilidades.Eliminar_repetidos(arch_list)

			archivo_prompt = utilidades.Archivo_a_prompt(arch_list)

		case opciones[3]:

			archivo_prompt.Borrar_informacion()

		case opciones[4]:

			imgs := utilidades.Abrir_selector_archivos()

			imagenes, imgerr := utilidades.Imagen_a_base64(imgs...)

			if len(imagenes) == 0 {

				continue
			}

			if imgerr != nil {

				utilidades.Logueo_simple(imgerr)
				continue
			}

			utilidades.Mostrar_archivos(imgs)

			if err := iniciar_conversacion(archivo_prompt, box, api_generate, content_type, false, imagenes); err != nil {

				utilidades.Logueo_simple(err)
				continue
			}

		case opciones[5]:

			utilidades.Mostrar_archivos(archivo_prompt.Archivos)

			defer archivo_prompt.Borrar_informacion()

			if err := iniciar_conversacion(archivo_prompt, box, api_chat, content_type, true, []string{}); err != nil {

				continue
			}

		}
	}
}

func listar_modelos_disponibles(url string) []string {

	tags := fmt.Sprintf("%s/tags", url)

	resp, resperr := http.Get(tags)

	modelos_disponibles := []string{}

	var modelos prompts.Modelos

	if resperr != nil {

		return modelos_disponibles
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {

		return modelos_disponibles
	}

	data, rderr := io.ReadAll(resp.Body)

	if rderr != nil {

		return modelos_disponibles
	}

	if jsonerr := json.Unmarshal(data, &modelos); jsonerr != nil {

		return modelos_disponibles
	}

	for _, modelo := range modelos.Models {

		modelos_disponibles = append(modelos_disponibles, modelo.Model)
	}

	return modelos_disponibles

}

func checkear_status(url string, tiempo time.Duration) error {

	status := fmt.Sprintf("%s/status", url)

	c := http.Client{

		Timeout: time.Second * tiempo,
	}

	resp, err := c.Get(status)

	if err != nil || resp.StatusCode == 404 {

		return errors.New("servidor apagado o no disponible")

	}

	return nil

}

func menu_modelos(modelos_disponibles []string) (string, error) {

	limite_muestreo := 5

	IA_MODELO, menuerr := menu.Menu(limite_muestreo, modelos_disponibles...)

	if menuerr != nil {
		return "", menuerr
	}
	return IA_MODELO, nil
}

func main() {

	if conserr != nil {
		utilidades.Error(fmt.Sprintf("Problema al habilitar ansi: %v\n", conserr))
		return
	}

	flag.Parse()

	Host := *host_selec
	Puerto := *puerto_selec
	Ctx := *ctx // el nivel de memoria de trabajo que puede maneja el LLM
	Temp := *temp
	Timeout := time.Duration(*timeout)
	Out := *Output

	var url = fmt.Sprintf("http://%s:%d/api", Host, Puerto)

	instalado := utilidades.Ollama_instalado()

	if !instalado {

		utilidades.Advertencia("Ollama no fue encontrado en las variables de entorno")
		time.Sleep(time.Second * utilidades.TIEMPO_PAUSA)
	}

	if err := checkear_status(url, Timeout); err != nil {

		utilidades.Logueo_simple(err)
		return

	}

	modelos_disponibles := listar_modelos_disponibles(url)

	if len(modelos_disponibles) == 0 {
		fmt.Print("\n\n")
		utilidades.Advertencia(`No hay modelos disponibles instalados actualmente, usa el comando "ollama pull (modelo)" para descargarlos`)
		fmt.Print("\n\n")
		time.Sleep(time.Second * utilidades.TIEMPO_PAUSA)
		return
	}

	// flujo del programa
	var opcion_salir string = fmt.Sprintf("%s[Salir]%s", utilidades.ANIL, utilidades.RESET)

	opciones_modelos := []string{opcion_salir}

	opciones_modelos = append(opciones_modelos, modelos_disponibles...)

	fmt.Print(utilidades.ALTERNATE)

	fmt.Print(utilidades.HOME)

	for {
		menu.Logo()
		//TODO: si el usuario tiene muchos modelos se puede buguear visualmente, quiza deba corregir eso
		Opcion_modelo, menuerr := menu_modelos(opciones_modelos)

		if menuerr != nil {

			utilidades.Error(menuerr)
			utilidades.Info("visitar https://ollama.com/search para mas info")
			time.Sleep(time.Second * utilidades.TIEMPO_PAUSA)
			return

		}

		if Opcion_modelo == opcion_salir {
			utilidades.Salir()
		}

		box := utilidades.Box_info{
			Modelo:      Opcion_modelo,
			Sistema_op:  runtime.GOOS,
			Temperatura: Temp,
			Ctx:         Ctx,
			Host:        Host,
			Puerto:      Puerto,
			Archivo:     Out,
		}

		iniciar_prompts(url, CONTENT_TYPE, box)
		utilidades.Limpieza_rapida()
	}
}
