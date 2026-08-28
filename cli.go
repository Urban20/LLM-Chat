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
	"strconv"
	"sync"
	"time"

	"github.com/rvfet/rich-go"
)

const (
	HOST_DEFAULT   = "localhost"
	PUERTO_DEFAULT = 11434
	CTX_DEFAULT    = 16_000
	TEMP_DEFAULT   = 0.5
	CONTENT_TYPE   = "aplication/json"
)

var (
	conserr      = consola.Iniciar_ANSI()
	host_selec   = flag.String("host", HOST_DEFAULT, "url al enpoint de Ollama")
	puerto_selec = flag.Int("puerto", PUERTO_DEFAULT, "puerto donde se escucha el endpoint")
	ctx          = flag.Int("ctx", CTX_DEFAULT, "cantidad contexto que usara el LLM")
	temp         = flag.Float64("temp", TEMP_DEFAULT, "temperatura del LLM")
)

func iniciar_conversacion(archivo_prompt utilidades.Prompt_archivo, modelo, endpoint, content_type string, ctx int, temp float64, chat bool, imagenes []string) error {

	prompt := utilidades.Input_multilinea("Prompt")

	if prompt == "" {

		return errors.New("prompt vacio")
	}

	carga := menu.Crear_carga()

	wg := sync.WaitGroup{}

	go carga.Iniciar(&wg)

	if err := prompts.Comunicacion(archivo_prompt.Prompt+"\n[prompt]\n\n"+prompt, modelo, endpoint, content_type, ctx, temp, &carga, &wg, chat, imagenes); err != nil {
		fmt.Print("\n")
		rich.Warning(err)
	}

	return nil

}

func iniciar_prompts(modelo, url, content_type string, ctx int, temp float64) {

	opciones := []string{"Volver", "Borrar contexto", "Adjuntar archivos de texto plano", "Eliminar archivos adjuntos", "Adjuntar imagen", "Ingresar prompt"}

	api_chat := fmt.Sprintf("%s/chat", url)
	api_generate := fmt.Sprintf("%s/generate", url)

	var archivo_prompt utilidades.Prompt_archivo

	for {
		// TODO : quiza modifique esto

		seleccion, _ := menu.Menu(opciones...)

		switch seleccion {

		case opciones[0]:

			prompts.Borrar_memoria()
			prompts.Descargar_modelo(modelo, content_type, api_chat)

			return

		case opciones[1]:

			utilidades.Limpieza_rapida()
			prompts.Borrar_memoria()
			fmt.Print("\n")
			rich.Info("la memoria del LLM fue borrada")

		case opciones[2]:

			arch_list, formaterr := utilidades.Formatear_input("ruta de los archivos")

			if formaterr != nil {

				continue
			}

			archivo_prompt = utilidades.Archivo_a_prompt(arch_list)

		case opciones[3]:

			archivo_prompt.Borrar_informacion()

		case opciones[4]:

			imgs, formaterr := utilidades.Formatear_input("ruta de las imagenes")

			if formaterr != nil {

				continue
			}

			imagenes, imgerr := utilidades.Imagen_a_base64(imgs...)

			if imgerr != nil {

				continue
			}

			if err := iniciar_conversacion(archivo_prompt, modelo, api_generate, content_type, ctx, temp, false, imagenes); err != nil {

				continue
			}

		case opciones[5]:

			archivo_prompt.Mostrar_archivos()

			if err := iniciar_conversacion(archivo_prompt, modelo, api_chat, content_type, ctx, temp, true, []string{}); err != nil {

				continue
			}

			archivo_prompt.Borrar_informacion()
		}
	}
}

func box_informacion(IA_MODELO, Host string, Puerto int, temp float64, ctx int) {

	utilidades.Limpieza_rapida()

	contenido_box := map[string]string{

		"Modelo selecionado":  IA_MODELO,
		"Host":                fmt.Sprintf("%s:%d", Host, Puerto),
		"Sistema operativo":   runtime.GOOS,
		"Temperatura del LLM": fmt.Sprintf("%.2f", temp),
		"Contexto del LLM":    strconv.Itoa(ctx),
	}
	contenidos := utilidades.Formato_string_box(contenido_box)
	utilidades.Box(contenidos...)

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

func checkear_status(url string) error {

	status := fmt.Sprintf("%s/status", url)

	resp, err := http.Get(status)

	if err != nil || resp.StatusCode == 404 {

		return errors.New("servidor apagado o no disponible")

	}

	return nil

}

func menu_modelos(modelos_disponibles []string) (string, error) {

	IA_MODELO, menuerr := menu.Menu(modelos_disponibles...)

	if menuerr != nil {
		return "", menuerr
	}
	return IA_MODELO, nil
}

func main() {

	if conserr != nil {
		rich.Error(fmt.Sprintf("Problema al habilitar ansi: %v\n", conserr))
		return
	}

	flag.Parse()

	Host := *host_selec
	Puerto := *puerto_selec
	Ctx := *ctx // el nivel de memoria de trabajo que puede maneja el LLM
	Temp := *temp

	var url = fmt.Sprintf("http://%s:%d/api", Host, Puerto)

	instalado := utilidades.Ollama_instalado()

	if !instalado {

		rich.Warning("ollama no fue encontrado en las variables de entorno")
		time.Sleep(time.Second * utilidades.TIEMPO_PAUSA)
	}

	if err := checkear_status(url); err != nil {

		utilidades.Logueo_simple(err)
		return

	}

	modelos_disponibles := listar_modelos_disponibles(url)

	if len(modelos_disponibles) == 0 {
		fmt.Print("\n\n")
		rich.Warning(`No hay modelos disponibles instalados actualmente, usa el comando "ollama pull (modelo)" para descargarlos`)
		fmt.Print("\n\n")
		time.Sleep(time.Second * utilidades.TIEMPO_PAUSA)
		return
	}

	// flujo del programa
	var opcion_salir string = "[Salir]"

	opciones_modelos := []string{opcion_salir}

	opciones_modelos = append(opciones_modelos, modelos_disponibles...)

	for {

		utilidades.Limpieza_rapida()
		//TODO: si el usuario tiene muchos modelos se puede buguear visualmente, quiza deba corregir eso
		Opcion_modelo, menuerr := menu_modelos(opciones_modelos)

		if menuerr != nil {

			rich.Error(menuerr)
			rich.Info("visitar https://ollama.com/search para mas info")
			time.Sleep(time.Second * utilidades.TIEMPO_PAUSA)
			return

		}

		if Opcion_modelo == opcion_salir {
			utilidades.Limpieza_rapida()
			return

		}

		box_informacion(Opcion_modelo, Host, Puerto, Temp, Ctx)

		iniciar_prompts(Opcion_modelo, url, CONTENT_TYPE, Ctx, Temp)

	}
}
