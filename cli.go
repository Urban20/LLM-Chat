package main

import (
	"Cli-ia/prompts"
	"bufio"
	"fmt"
	"os"
)

func input() string {
	fmt.Print("Prompt >> ")
	sc := bufio.NewScanner(os.Stdin)
	sc.Scan()
	return sc.Text()

}

func main() {

	var prompt string

	for prompt != "salir" {

		prompt = input()

		if err := prompts.Comunicacion(prompt); err != nil {

			fmt.Println(err)
			break
		}

	}

}
