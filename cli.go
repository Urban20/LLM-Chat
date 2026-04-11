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

		i, _ := prompts.Enviar_prompt(prompt)

		fmt.Println(i.Response)

	}

}
