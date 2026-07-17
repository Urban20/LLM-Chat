package main

import (
	consola "LLM-Chat/ansi"
	"LLM-Chat/menu"
	"LLM-Chat/prompts"
	"LLM-Chat/utilidades"
	"bufio"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/rvfet/rich-go"
)

var Host_default = "localhost"
var Puerto_default = 11434
var ctx_default = 16_000
var temp_defalut = 0.5
var Content_type = "aplication/json"
var conserr = consola.Iniciar_ANSI()

var host_selec = flag.String("host", Host_default, "url al enpoint de Ollama")
var puerto_selec = flag.Int("puerto", Puerto_default, "puerto donde se escucha el endpoint")
var ctx = flag.Int("ctx", ctx_default, "cantidad contexto que usara el LLM")
var temp = flag.Float64("temp", temp_defalut, "temperatura del LLM")

func input(input string) string {

	fmt.Printf("\n\n%s[presionar TAB + ENTER para enviar]%s", utilidades.AMARILLO, utilidades.RESET)

	fmt.Print(utilidades.VIOLETA)
	fmt.Printf("\n\n%s :\n", input)
	fmt.Print(utilidades.RESET)
	lector := bufio.NewReader(os.Stdin)
	texto, _ := lector.ReadString('\t')
	return strings.TrimSpace(strings.Trim(texto, "\t"))

}

func iniciar_prompts(modelo, url, content_type string, ctx int, temp float64) {

	opciones := []string{"Volver", "Borrar contexto", "Ingresar prompt"}

	api_chat := fmt.Sprintf("%s/chat", url)

	for {
		// TODO : quiza modifique esto

		seleccion, _ := menu.Menu(opciones...)

		switch seleccion {

		case opciones[0]:

			prompts.Borrar_memoria()

			return

		case opciones[1]:
			utilidades.Limpieza_rapida()
			prompts.Borrar_memoria()
			fmt.Print("\n")
			rich.Info("la memoria del LLM fue borrada")

		case opciones[2]:

			prompt := input("Prompt")

			carga := menu.Crear_carga()

			wg := sync.WaitGroup{}

			go carga.Iniciar(&wg)

			if err := prompts.Comunicacion(prompt, modelo, api_chat, content_type, ctx, temp, &carga, &wg); err != nil {
				fmt.Print("\n")
				rich.Warning(err)
			}

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

			return

		}

		box_informacion(Opcion_modelo, Host, Puerto, Temp, Ctx)

		iniciar_prompts(Opcion_modelo, url, Content_type, Ctx, Temp)

	}
}
