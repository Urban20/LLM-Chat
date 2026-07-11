package utilidades

import (
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/charmbracelet/glamour"
	"github.com/pterm/pterm"
	"github.com/rvfet/rich-go"
	"golang.org/x/term"
)

const TIEMPO_PAUSA = 4

const (
	AMARILLO    = "\033[0;33m"
	RESET       = "\033[0m"
	VIOLETA     = "\033[38;2;146;68;219m"
	GRIS_AZUL   = "\033[38;2;90;112;176m"
	BLANCO      = "\033[47m"
	AZUL_OSCURO = "\033[38;2;116;116;247m"
)

func separador() {

	x, _, _ := term.GetSize(int(os.Stdout.Fd()))
	fmt.Println(BLANCO)
	fmt.Println(strings.Repeat(" ", x))
	fmt.Println(RESET)

}

func Ollama_instalado() bool {

	/*
		con esta funcion miro si Ollama esta en las variables de entorno del sistema

	*/

	ollama, err := exec.LookPath("Ollama")

	return err == nil || ollama != ""

}

func Imprimir_markdown(txt string) error {

	render, termerr := glamour.NewTermRenderer(glamour.WithStylesFromJSONBytes([]byte(Estilos)))

	if termerr != nil {

		return termerr
	}
	separador()
	md, err := render.Render("# LLM:\n" + txt)

	if err != nil {
		return err
	}
	fmt.Print(md)
	separador()

	return nil
}

func Box(msgs ...string) {

	superficie := pterm.DefaultBox.WithHorizontalPadding(5).WithBottomPadding(1)

	superficie.Println(strings.Join(msgs, "\n"))

}

func Limpieza_rapida() {

	fmt.Print("\033[2J")
	fmt.Print("\033[H")

}

func Logueo_simple(mensaje any) {
	rich.Error(mensaje)
	time.Sleep(time.Second * TIEMPO_PAUSA)

}

func Formato_string_box(cuerpo map[string]string) []string {

	var retorno []string

	for clave, valor := range cuerpo {

		elemento := VIOLETA + clave + RESET + " : " + valor

		retorno = append(retorno, elemento)

	}

	return retorno

}

//go:embed estilo.json
var Estilos string
