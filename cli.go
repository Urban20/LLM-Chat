package main

import (
	"Cli-ia/modelo"
	"Cli-ia/ollama"
	"Cli-ia/prompts"
	"Cli-ia/web"
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

const LIMITE_MEMORIA = 100

var Host_default = "localhost"
var Puerto_default = 11434
var Content_type = "aplication/json"
var IA_default = "llama3"
var Modelo = "IA-CLI"
var Json_modelo = strings.NewReader(fmt.Sprintf(`{"model":"%s"}`, Modelo))

var ia_selec = flag.String("modelo", IA_default, "modelo de ia a utilizar")
var host_selec = flag.String("host", Host_default, "url al enpoint de Ollama")
var puerto_selec = flag.Int("puerto", Puerto_default, "puerto donde se escucha el endpoint")

func input(input string) string {
	fmt.Printf("\n\n%s >> ", input)
	sc := bufio.NewScanner(os.Stdin)
	sc.Scan()
	return sc.Text()

}

func iniciar_prompts(modelo, api_chat, content_type string) {

	for {

		prompt := input("Prompt")

		switch prompt {
		case "salir":

			fmt.Println("\nsaliendo ...")
			time.Sleep(time.Second * 2)
			return

		case "url":

			url := input("Url")
			doc, weberr := web.Buscar(url)

			if weberr != nil {
				fmt.Println("hubo un problema al buscar la url: ", weberr)
				continue
			}
			prompt = input("Prompt")
			url_prompt := fmt.Sprintf("%s:\n %s", prompt, doc)
			prompts.Comunicacion(url_prompt, modelo, api_chat, content_type)

		default:

			if len(prompts.Memoria) >= LIMITE_MEMORIA {
				fmt.Printf("Se llego al limite de la memoria: %d, la IA ya no puede recordar mas\n", LIMITE_MEMORIA)
				prompts.Memoria = prompts.Memoria[:LIMITE_MEMORIA]

			}

			if err := prompts.Comunicacion(prompt, modelo, api_chat, content_type); err != nil {

				fmt.Println(err)
				break
			}
			//fmt.Println(prompts.Memoria)

		}

	}
}

func intentar_reconexiones(ia, modelo_, api_modelo string) {

	for {

		if err := modelo.Crear_modelo(ia, modelo_, api_modelo, Content_type); err != nil {
			time.Sleep(time.Second * 3)
			fmt.Println("reconectando...")
			continue
		}

		fmt.Println("exito!")
		return
	}

}

func main() {

	flag.Parse()

	Host := *host_selec
	Puerto := *puerto_selec
	IA_MODELO := *ia_selec

	var Api_chat = fmt.Sprintf("http://%s:%d/api/chat", Host, Puerto)
	var Api_modelo = fmt.Sprintf("http://%s:%d/api/create", Host, Puerto)

	ollama_, instalado := ollama.Ollama_instalado()

	if !instalado {
		time.Sleep(time.Second * 3)
		return
	}
	fmt.Println("Ollama detectado")

	ollama.Iniciar_Ollama(ollama_)
	intentar_reconexiones(IA_MODELO, Modelo, Api_modelo)

	iniciar_prompts(IA_MODELO, Api_chat, Content_type)

}
