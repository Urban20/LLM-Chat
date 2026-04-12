package main

import (
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

func main() {

	var prompt string
	modelo.Crear_modelo()

	for {

		prompt = input()

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
