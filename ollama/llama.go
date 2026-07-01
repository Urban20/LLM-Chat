package ollama

import (
	"os/exec"
)

// inicializador de llama3

func Ollama_instalado() bool {

	/*
		con esta funcion miro si Ollama esta en las variables de entorno del sistema

	*/

	ollama, err := exec.LookPath("Ollama")

	return err == nil || ollama != ""

}
