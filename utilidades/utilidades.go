package utilidades

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/glamour"
	"golang.org/x/term"
)

var Memoria = []string{}

func separador() {

	x, _, _ := term.GetSize(int(os.Stdout.Fd()))
	fmt.Println("\033[47m")
	fmt.Println(strings.Repeat("_", x))
	fmt.Println("\033[0m")

}

func Imprimir_markdown(txt string) error {

	separador()
	md, err := glamour.Render(txt, "dark")
	if err != nil {
		return err
	}
	fmt.Print(md)
	separador()

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
