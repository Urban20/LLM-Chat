package main

import (
	consola "LLM-Chat/ansi"
	"LLM-Chat/prompts"
	"LLM-Chat/utilidades"
	"bufio"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/rvfet/rich-go"
)

const LIMITE_MEMORIA = 100

var Host_default = "localhost"
var Puerto_default = 11434
var Content_type = "aplication/json"
var IA_default = "llama3"
var conserr = consola.Iniciar_ANSI()

var ia_selec = flag.String("modelo", IA_default, "modelo de ia a utilizar")
var host_selec = flag.String("host", Host_default, "url al enpoint de Ollama")
var puerto_selec = flag.Int("puerto", Puerto_default, "puerto donde se escucha el endpoint")

func input(input string) string {

	fmt.Print("\n\n[presionar TAB + ENTER para enviar]")

	fmt.Print(utilidades.CELESTE_CLARO)
	fmt.Printf("\n\n%s :\n", input)
	fmt.Print(utilidades.RESET)
	lector := bufio.NewReader(os.Stdin)
	texto, _ := lector.ReadString('\t')
	return strings.Trim(texto, "\t")

}

func iniciar_prompts(modelo, api_chat, content_type string) {

	for {
		// TODO : quiza modifique esto
		prompt := input("Prompt")

		switch prompt {
		case "salir":

			rich.Info("\nsaliendo ...")
			time.Sleep(time.Second * 2)
			return

		default:

			if len(utilidades.Memoria) >= LIMITE_MEMORIA {
				rich.Warning("Se llego al limite de la memoria: %d, la IA ya no puede recordar mas\n", LIMITE_MEMORIA)
				utilidades.Memoria = utilidades.Memoria[:LIMITE_MEMORIA]

			}

			if err := prompts.Comunicacion(prompt, modelo, api_chat, content_type); err != nil {

				rich.Error(err)
				break
			}
			//fmt.Println(prompts.Memoria)

		}

	}
}

func main() {

	if conserr != nil {
		rich.Error("Problema al habilitar ansi: %v\n", conserr)
		return
	}

	flag.Parse()

	Host := *host_selec
	Puerto := *puerto_selec
	IA_MODELO := *ia_selec

	var Api_chat = fmt.Sprintf("http://%s:%d/api/chat", Host, Puerto)

	instalado := utilidades.Ollama_instalado()

	if !instalado {
		time.Sleep(time.Second * 3)
		return
	}

	status := fmt.Sprintf("http://%s:%d/api/version", Host, Puerto)

	resp, err := http.Get(status)

	if err != nil || resp.StatusCode == 404 {

		rich.Error("servidor apagado o no disponible")

		return

	}

	utilidades.Limpieza_rapida()

	iniciar_prompts(IA_MODELO, Api_chat, Content_type)

}
