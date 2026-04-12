package main

import (
	"Cli-ia/llama3"
	"Cli-ia/modelo"
	"Cli-ia/prompts"
	"Cli-ia/web"
	"bufio"
	"fmt"
	"os"
	"time"
)

const LIMITE_MEMORIA = 100

func input(input string) string {
	fmt.Printf("\n\n%s >> ", input)
	sc := bufio.NewScanner(os.Stdin)
	sc.Scan()
	return sc.Text()

}

func iniciar_prompts() {

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
			prompts.Comunicacion(url_prompt)

		default:

			if len(prompts.Memoria) >= LIMITE_MEMORIA {
				fmt.Printf("Se llego al limite de la memoria: %d, la IA ya no puede recordar mas\n", LIMITE_MEMORIA)
				prompts.Memoria = prompts.Memoria[:LIMITE_MEMORIA]

			}

			if err := prompts.Comunicacion(prompt); err != nil {

				fmt.Println(err)
				break
			}
			//fmt.Println(prompts.Memoria)

		}

	}
}

func intentar_reconexiones() {

	for {

		if err := modelo.Crear_modelo(); err != nil {
			time.Sleep(time.Second * 3)
			fmt.Println("reconectando...")
			continue
		}

		fmt.Println("exito!")
		return
	}

}

func main() {

	ollama, instalado := llama3.Ollama_instalado()

	if !instalado {
		time.Sleep(time.Second * 3)
		return
	}
	fmt.Println("Ollama detectado")

	llama3.Iniciar_Ollama(ollama)
	intentar_reconexiones()

	iniciar_prompts()

}
