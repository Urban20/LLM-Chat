package main

import (
	"Cli-ia/prompts"
	"bufio"
	"fmt"
	"os"
)

const LIMITE_MEMORIA = 50

func input() string {
	fmt.Print("\nPrompt >> ")
	sc := bufio.NewScanner(os.Stdin)
	sc.Scan()
	return sc.Text()

}

func main() {

	var prompt string

	for {

		prompt = input()

		if prompt == "salir" {
			fmt.Println("\nsaliendo ...")
			return

		}

		if len(prompts.Memoria) >= LIMITE_MEMORIA {
			fmt.Printf("Se llego al limite de la memoria: %d\n", LIMITE_MEMORIA)
			prompts.Memoria = prompts.Memoria[:LIMITE_MEMORIA]

		}

		if err := prompts.Comunicacion(prompt); err != nil {

			fmt.Println(err)
			break
		}
		//fmt.Println(prompts.Memoria)

	}

}
