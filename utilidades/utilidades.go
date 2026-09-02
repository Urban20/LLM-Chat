package utilidades

import (
	"bufio"
	_ "embed"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/charmbracelet/glamour"
	"github.com/pterm/pterm"
	"github.com/rvfet/rich-go"
	"golang.org/x/term"
)

const (
	TIEMPO_PAUSA = 4
	AMARILLO     = "\033[0;33m"
	RESET        = "\033[0m"
	VIOLETA      = "\033[38;2;146;68;219m"
	GRIS_AZUL    = "\033[38;2;90;112;176m"
	BLANCO       = "\033[47m"
	AZUL_OSCURO  = "\033[38;2;116;116;247m"
	FONDO_VERDE  = "\033[48;2;46;166;66m"
	NEGRO_BLANCO = "\033[107;30m"
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

func Eliminar_repetidos(elementos []string) []string {

	copia := []string{}

	for _, el := range elementos {

		if slices.Contains(copia, el) {
			continue
		}

		copia = append(copia, el)

	}

	return copia

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
	var archivos []string

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

	fmt.Print(AMARILLO + str + ": " + RESET)
	sc := bufio.NewScanner(os.Stdin)

	sc.Scan()

	return sc.Text()

}

func Input_multilinea(input string) string {

	fmt.Printf("\n\n%s[presionar TAB + ENTER para enviar]%s", NEGRO_BLANCO, RESET)

	fmt.Print(VIOLETA)
	fmt.Printf("\n\n%s :\n", input)
	fmt.Print(RESET)
	lector := bufio.NewReader(os.Stdin)
	texto, _ := lector.ReadString('\t')
	return strings.TrimSpace(strings.Trim(texto, "\t"))

}

//go:embed estilo.json
var Estilos string

func Formatear_input(msg string) ([]string, error) {

	// espera un mensaje para el input y se devuelve la salida en formato de lista de strings
	//para ser procesado por la ia (envio de archivos como texto plano , imagenes, etc)

	var arch_list []string

	fmt.Print("\n")
	archivos := Input(AMARILLO + msg + RESET)

	if archivos == "" {

		return arch_list, errors.New("input vacio") // esto realmente no se usa, es para que evite ejecutando

	}

	arch_list = strings.Split(archivos, " ")

	return arch_list, nil

}

func Imagen_a_base64(rutas ...string) ([]string, error) {

	var bases []string

	for _, ruta := range rutas {

		img, imgerr := os.Open(ruta)

		if imgerr != nil {

			return bases, imgerr
		}

		defer img.Close()

		data, dataerr := io.ReadAll(img)

		if dataerr != nil {
			return bases, dataerr
		}

		base := base64.StdEncoding.EncodeToString(data)

		bases = append(bases, base)
	}

	return bases, nil

}
