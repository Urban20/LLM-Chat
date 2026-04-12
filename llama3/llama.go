package llama3

import (
	"fmt"
	"os"
	"os/exec"
)

// inicializador de llama3

var Output_ollama = "ollama-cmd.log"

func Ollama_instalado() (string, bool) {

	/*
		con esta funcion miro si Ollama esta en las variables de entorno del sistema

	*/

	ollama, err := exec.LookPath("Ollama")

	if err != nil || ollama == "" {
		fmt.Println("Ollama no detectado: ", err.Error())
		return "", false
	}

	return ollama, true

}

func Iniciar_Ollama(ruta string) {

	cmd := exec.Command(ruta, "serve")
	out, _ := os.Create(Output_ollama)
	cmd.Stdout = out
	cmd.Stderr = out
	cmd.Start()
}
