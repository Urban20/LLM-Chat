package utilidades

import (
	"fmt"

	"github.com/charmbracelet/glamour"
)

func Imprimir_markdown(txt string) error {

	md, err := glamour.Render(txt, "dark")
	if err != nil {
		return err
	}
	fmt.Print(md)
	return nil
}
