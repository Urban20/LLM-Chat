package utilidades

import (
	"bufio"
	_ "embed"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
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

type Prompt_archivo struct { //struct para manejar el envio de informacion al llm
	Prompt   string
	Archivos []string
}

func (s Prompt_archivo) Mostrar_archivos() {

	if len(s.Archivos) == 0 {

		return
	}

	fmt.Print(AZUL_OSCURO + "\n\n(*) Archivos adjuntos:\n\n" + RESET)

	for _, arch := range s.Archivos {

		fmt.Println(arch)
	}

}

func (s *Prompt_archivo) Borrar_informacion() {

	s.Prompt = ""
	s.Archivos = []string{}

}

func Archivo_a_prompt(rutas []string) Prompt_archivo {

	/*
		funcion que lee archivos de texto plano como .txt, .md y los transforma en un string
		el cual se inyecta en el prompt para darle al modelo contexto de los archivos

	*/

	var prompt string
	var archivos = []string{}

	for _, ruta := range rutas {

		ruta = filepath.Clean(ruta)
		nombre_archivo := filepath.Base(ruta)

		archivos = append(archivos, nombre_archivo)

		vacio := fmt.Sprintf("file: [ %s ]\n\n**empty**\n\n", nombre_archivo) //lo escribo en ingles para que la ia lo tome como prompt independientemente del idioma
		// el ingles es el idioma base

		archivo, archerr := os.Open(ruta)

		if archerr != nil {

			prompt += vacio
			continue
		}

		defer archivo.Close()

		contenido, conterr := io.ReadAll(archivo)

		if conterr != nil {

			prompt += vacio
			continue
		}

		prompt += fmt.Sprintf("file: [ %s ]\n\n%s\n\n", nombre_archivo, strings.TrimSpace(string(contenido)))
	}

	return Prompt_archivo{Prompt: prompt, Archivos: archivos}

}

func Input(str string) string {

	fmt.Print(AMARILLO + str + RESET)
	sc := bufio.NewScanner(os.Stdin)

	sc.Scan()

	return sc.Text()

}

//go:embed estilo.json
var Estilos string
