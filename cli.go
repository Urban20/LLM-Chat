package main

import (
	"Cli-ia/llama3"
	"Cli-ia/modelo"
	"Cli-ia/prompts"
	"bufio"
	"fmt"
	"os"
	"time"
)

const LIMITE_MEMORIA = 100

func input() string {
	fmt.Print("\n\nPrompt >> ")
	sc := bufio.NewScanner(os.Stdin)
	sc.Scan()
	return sc.Text()

}

func iniciar_prompts() {

	for {

		prompt := input()

		if prompt == "salir" {
			fmt.Println("\nsaliendo ...")
			time.Sleep(time.Second * 2)
			return

		}

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
