package utilidades

import (
	"encoding/json"
	"fmt"

	"github.com/charmbracelet/glamour"
)

var Memoria = []string{}

func Imprimir_markdown(txt string) error {

	md, err := glamour.Render(txt, "dark")
	if err != nil {
		return err
	}
	fmt.Print(md)
	return nil
}

func Guardar_en_memoria(prompt, rol string) {

	mensaje_usuario := map[string]string{
		"role":    rol,
		"content": prompt,
	}

	msg, _ := json.Marshal(&mensaje_usuario)
	Memoria = append(Memoria, string(msg))

}
