package utilidades

import (
	"bufio"
	_ "embed"
	"encoding/base64"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"time"

	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/pterm/pterm"
	"golang.org/x/term"
)

const (
	TIEMPO_PAUSA    = 4
	AMARILLO        = "\033[0;33m"
	RESET           = "\033[0m"
	VIOLETA         = "\033[38;2;146;68;219m"
	GRIS_AZUL       = "\033[38;2;90;112;176m"
	FONDO_BLANCO    = "\033[47m"
	AZUL_OSCURO     = "\033[38;2;116;116;247m"
	FONDO_VERDE     = "\033[48;2;46;166;66m"
	NEGRO_BLANCO    = "\033[107;30m"
	HOME            = "\033[H"
	BORRADO         = "\033[2J"
	ALTERNATE       = "\033[?1049h"
	ALTERNATE_RESET = "\033[?1049l"
	ROJO_BRILLANTE  = "\033[91m"
	VERDE_BRILLANTE = "\033[92m"
	ANIL            = "\033[48;2;43;43;69m"
)

const (
	OCULTAR_CURSOR = "\033[?25l"
	MOSTRAR_CURSOR = "\033[?25h"
	RETROCESO      = "\033[F"

	VERSION = "V1.1.0"
)

const (
	KEY_ARRIBA   = 65
	KEY_ABAJO    = 66
	KEY_DER      = 67
	KEY_IZQ      = 68
	ENTER        = 13
	TIEMPO_CARGA = 0.85
)

const CONTENT_TYPE = "application/json"

var (
	// se adapta al color de la terminal
	oscuros  = lipgloss.AdaptiveColor{Light: "#000000", Dark: "#383838"}
	claros   = lipgloss.AdaptiveColor{Light: "#b9b9b9", Dark: "#ffffff"}
	violetas = lipgloss.AdaptiveColor{Light: "#2B2B45", Dark: "#434368"}
	azules   = lipgloss.AdaptiveColor{Light: "#404f7f", Dark: "#6f80bc"}
	verdes   = lipgloss.AdaptiveColor{Light: "#1dcb5a", Dark: "#1dcb5a"}
)

var (
	Stdin_fd  = int(os.Stdin.Fd())
	Stdout_fd = int(os.Stdout.Fd())
)

func Centrar(texto string) {

	pterm.DefaultCenter.Println(texto)

}

func Ollama_instalado() bool {

	/*
		con esta funcion miro si Ollama esta en las variables de entorno del sistema

	*/

	ollama, err := exec.LookPath("Ollama")

	return err == nil || ollama != ""

}

func Salir() {

	fmt.Print(ALTERNATE_RESET)
	fmt.Print(MOSTRAR_CURSOR)

	os.Exit(0)

}

func Imprimir_markdown(r Respuesta_LLM) error {

	x, y, termerr := term.GetSize(Stdout_fd)

	if termerr != nil {

		return termerr
	}

	render, termerr := glamour.NewTermRenderer(
		glamour.WithWordWrap(x), glamour.WithStylesFromJSONBytes([]byte(Estilos)))

	if termerr != nil {

		return termerr
	}

	md, err := render.Render(fmt.Sprintf("%s **[USUARIO]**:\n%s\n\n# LLM (%s):\n %s", Tiempo_actual(), r.Prompt, r.Modelo, r.Respuesta_raw))

	if err != nil {
		return err
	}
	scroll := Crear_scroll(y)
	scroll.Renderizar(md)
	scroll.Iniciar()

	return nil
}

func Box(msgs ...string) {

	superficie := pterm.DefaultBox.WithHorizontalPadding(15).WithBottomPadding(1)

	Centrar(superficie.Sprintln(strings.Join(msgs, "\n")))

}

func Limpieza_rapida() {

	fmt.Print(BORRADO)
	fmt.Print(HOME)

}

func Logueo_simple(mensaje any) {
	Error(mensaje)
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

func Archivo_a_prompt(rutas []string) Prompt_archivo {

	/*
		funcion que lee archivos de texto plano como .txt, .md y los transforma en un string
		el cual se inyecta en el prompt para darle al modelo contexto de los archivos

	*/

	var prompt string
	var archivos []string

	for _, ruta := range rutas {

		ruta = filepath.Clean(ruta)

		archivos = append(archivos, ruta)

		vacio := fmt.Sprintf("file: [ %s ]\n\n**empty**\n\n", ruta) //lo escribo en ingles para que la ia lo tome como prompt independientemente del idioma
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

		prompt += fmt.Sprintf("file: [ %s ]\n\n%s\n\n", ruta, strings.TrimSpace(string(contenido)))
	}

	return Prompt_archivo{Prompt: prompt, Archivos: archivos}

}

func String_vacio(str string) bool {

	for _, c := range str {

		if !slices.Contains([]rune{' ', '\n', '\t', '\r'}, c) {

			return false
		}

	}

	return true

}

func Input(str string) string {

	fmt.Print(AMARILLO + str + ": " + RESET)
	sc := bufio.NewScanner(os.Stdin)

	sc.Scan()

	return sc.Text()

}

func Input_multilinea(input string) string {
	fmt.Print("\n")
	var valor string

	ejemplos := []string{
		"¿como era la vida en los 90s?",
		"comenzar a escribir ...",
		"¿hasta donde llega la inteligencia humana?",
		"¿por que programar es genial?",
		"¿como nacio internet?",
	}

	seleccion := ejemplos[rand.Intn(len(ejemplos))]

	base := huh.ThemeBase()

	base.Focused.Title = base.Focused.Title.Foreground(oscuros).Background(claros)
	base.Focused.TextInput.Cursor = base.Focused.Title.Foreground(violetas)
	base.Focused.Description = base.Focused.Description.Foreground(azules)
	base.Focused.Base = base.Focused.Base.BorderForeground(verdes)

	i := huh.NewText().Title(input)
	i.WithTheme(base)
	i.Description("[ctrl + j/ alt + ENTER] : nueva linea\n")
	i.Placeholder(seleccion)
	i.Value(&valor)
	i.Editor("") //dado que no pude deshabilitar el editor externo (da problemas con el alternate screen)
	// se me ocurrio esta idea hacky (no es perfecta, mete un bug visual)
	// TODO: ver si se puede corregir definitivamente
	i.Run()

	return strings.TrimSpace(valor)

}

//go:embed estilo.json
var Estilos string

func crear_selector(dir_actual string) *huh.FilePicker {

	dir := huh.NewFilePicker()
	dir.Picking(true)
	dir.ShowSize(false)
	dir.ShowSize(true)
	dir.ShowHidden(true)
	dir.ShowPermissions(false)
	dir.DirAllowed(true)
	dir.CurrentDirectory(dir_actual)
	dir.Height(10)

	tema := huh.ThemeBase()

	tema.Focused.File = tema.Focused.File.Foreground(claros)
	tema.Focused.File = tema.Focused.File.Background(violetas)
	tema.Focused.Base = tema.Focused.Base.BorderForeground(verdes)

	dir.WithTheme(tema)

	return dir
}

func Abrir_selector_archivos() []string {

	// espera un mensaje para el input y se devuelve la salida en formato de lista de strings
	//para ser procesado por la ia (envio de archivos como texto plano , imagenes, etc)

	var arch_list = []string{}
	fmt.Print("\n")

	var env string

	switch runtime.GOOS {

	case "windows":

		env = "USERPROFILE"

	default: // pensado para linux y mac

		env = "HOME"

	}

	dir_actual := os.Getenv(env)

	for {

		dir := crear_selector(dir_actual)

		if err := dir.Run(); err != nil {

			break

		}

		arch_actual := fmt.Sprintf("%v", dir.GetValue())

		if arch_actual == "" {

			break
		}

		dir_actual = filepath.Dir(arch_actual)

		arch_list = append(arch_list, arch_actual)

	}

	return arch_list

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

func Tiempo_actual() string {

	t := time.Now()
	hora, min, seg := t.Clock()

	return fmt.Sprintf("[%d:%d:%d]", hora, min, seg)

}

func Fecha_hora() string {

	t := time.Now()
	a, m, d := t.Date()
	tiempo := Tiempo_actual()
	dia := t.Weekday().String()

	return fmt.Sprintf("%s %d %s %d\t%s", dia, d, m, a, tiempo)

}

func Separar_lista(lista []string, limite int) [][]string {

	var elementos [][]string

	var inicio int
	incremento := limite

	if limite <= 0 { // sin esto, si el limite es cero o negativo puede saturar la memoria

		return elementos

	}

	largo := len(lista)

	for largo > limite {

		elementos = append(elementos, lista[inicio:limite])
		inicio = limite
		limite += incremento

	}

	elementos = append(elementos, lista[inicio:])

	return elementos

}

// para testing: compara dos listas de listas de string
func Slices_Iguales(sl1, sl2 [][]string) bool {

	var i int
	tam1 := len(sl1)
	tam2 := len(sl2)

	if tam1 != tam2 {
		return false
	}

	for i < tam1 {

		if !slices.Equal(sl1[i], sl2[i]) {

			return false
		}

		i++
	}

	return true

}

func Max(l []string) int {

	var maximo int

	for _, el := range l {

		largo := len(el)

		if largo > maximo {

			maximo = largo

		}
	}

	return maximo

}

// evitamos que la lista contenga elementos mal formateados, como strings vacios
func Limpiar_listas(l []string) []string {

	var copia []string

	for _, v := range l {

		if !String_vacio(v) {

			copia = append(copia, v)

		}

	}

	return copia

}

func Mostrar_archivos(l []string) {

	if len(l) == 0 {

		return
	}

	fmt.Print(AZUL_OSCURO + "\n\n(*) Archivos adjuntos:\n\n" + RESET)

	for _, ruta := range l {

		arch := filepath.Base(ruta)

		fmt.Println(arch)
	}

}
