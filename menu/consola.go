package menu

import (
	"LLM-Chat/utilidades"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/pterm/pterm"
	"golang.org/x/term"
)

type Carga struct {
	estados  []string
	cargando bool
	tiempo   float32
}

func Logo() {

	logo := ` ██╗     ██╗     ███╗   ███╗       ██████╗██╗  ██╗ █████╗ ████████╗
 ██║     ██║     ████╗ ████║      ██╔════╝██║  ██║██╔══██╗╚══██╔══╝
 ██║     ██║     ██╔████╔██║█████╗██║     ███████║███████║   ██║   
 ██║     ██║     ██║╚██╔╝██║╚════╝██║     ██╔══██║██╔══██║   ██║   
 ███████╗███████╗██║ ╚═╝ ██║      ╚██████╗██║  ██║██║  ██║   ██║   
 ╚══════╝╚══════╝╚═╝     ╚═╝       ╚═════╝╚═╝  ╚═╝╚═╝  ╚═╝   ╚═╝   `

	lineas := strings.Split(logo, "\n")

	fmt.Print("\n\n")

	for _, l := range lineas {

		pterm.DefaultCenter.Printf(" %s %s %s", utilidades.ANIL, l, utilidades.RESET)
	}
	utilidades.Centrar(utilidades.VERSION)

	fmt.Print(utilidades.RESET)

}

func (p *Carga) Iniciar(wg *sync.WaitGroup) {

	fmt.Print("\n\n" + utilidades.OCULTAR_CURSOR)
	defer fmt.Print(utilidades.MOSTRAR_CURSOR)

	wg.Add(1)
	defer wg.Done()

	var i int

	for p.cargando {

		if i > len(p.estados)-1 {
			i = 0
		}

		fmt.Printf("\r%s", p.estados[i])
		i++
		time.Sleep(time.Second * time.Duration(p.tiempo))

	}

}

func Esperar_tecla() {

	fmt.Print(utilidades.OCULTAR_CURSOR)
	defer fmt.Print(utilidades.MOSTRAR_CURSOR)
	info := "[Tab para volver al menu]"

	fmt.Print("\n")
	fmt.Print(info)

	ejecutando := true

	st, _ := term.MakeRaw(utilidades.Stdin_fd)

	defer term.Restore(utilidades.Stdin_fd, st)

	for ejecutando {

		b := make([]byte, 3)

		os.Stdin.Read(b)

		if b[0] == '\t' {

			ejecutando = false

		}

	}

}

func (p *Carga) Detener(wg *sync.WaitGroup) {

	p.cargando = false
	wg.Wait()
	fmt.Print("\r" + strings.Repeat(" ", len(p.estados[0])))
}

func Crear_carga() Carga {

	c := Carga{

		estados: []string{"|", "/", "─", "\\"},

		cargando: true,
		tiempo:   utilidades.TIEMPO_CARGA}

	return c

}

func retroceder(n int) {

	for x := 0; x < n; x++ {

		fmt.Print(utilidades.RETROCESO)

	}

}

func actualizar_seccion(n, rep int) {

	retroceder(n)

	for x := 0; x < n; x++ {

		fmt.Println(strings.Repeat(" ", rep), "\r")

	}

	retroceder(n)

}

func leer_tecla(i, pag *int, tecla []byte) bool {
	os.Stdin.Read(tecla)
	flechas := tecla[2]

	if flechas == utilidades.KEY_DER {

		*pag++
	}
	if flechas == utilidades.KEY_IZQ {

		*pag--
	}

	if tecla[0] == utilidades.ENTER {

		return true

	} else if flechas == utilidades.KEY_ARRIBA || tecla[0] == 'w' {

		*i--
	} else if flechas == utilidades.KEY_ABAJO || tecla[0] == 's' {

		*i++
	}

	return false

}

func paginado(pag, max_len int) int {

	if pag > max_len-1 {

		pag = max_len - 1

	} else if pag < 0 {

		pag = 0
	}

	return pag

}

func imprimir_opciones(op string, actual []string, i int) {

	if op == actual[i] { // opcion seleccionada
		fmt.Println(utilidades.NEGRO_BLANCO + "> " + op + utilidades.RESET + "\r")
	} else {
		fmt.Println("  " + op + "\r")
	}

}

func limitar_indice(i, largo_op int) int {

	if i > largo_op-1 {
		i = 0

	} else if i < 0 {
		i = largo_op - 1
	}

	return i
}

func desplegar_opcion(opciones []string, limite int) string {

	var i int
	var pag int

	fraccionado := utilidades.Separar_lista(opciones, limite)
	max_len := len(fraccionado)
	var margen int

	for {

		tecla := make([]byte, 3)

		pag = paginado(pag, max_len)

		actual := fraccionado[pag]
		largo_op := len(actual)

		if max_len > 1 {

			fmt.Printf("%spagina %d/%d%s\n\n\r", utilidades.NEGRO_BLANCO, pag+1, max_len, utilidades.RESET)
			margen = 2
		}

		for _, op := range actual {

			i = limitar_indice(i, largo_op)

			imprimir_opciones(op, actual, i)
		}

		if leer_tecla(&i, &pag, tecla) {

			return actual[i]
		}

		actualizar_seccion(largo_op+margen, // le sumo los margenes de pagina
			utilidades.Max(actual)+2) // le sumo el margen

	}
}

func Menu(limite int, opciones ...string) (string, error) {

	fmt.Print(utilidades.OCULTAR_CURSOR)
	defer fmt.Print(utilidades.MOSTRAR_CURSOR)

	fmt.Print(utilidades.GRIS_AZUL + "\nOpciones disponibles:\n\n")
	fmt.Print(utilidades.AZUL_OSCURO + "navegar con ↑↓ | ← → cambiar pagina\n\n" + utilidades.RESET)

	st, rawerr := term.MakeRaw(utilidades.Stdin_fd)

	if rawerr != nil {
		return "", rawerr
	}

	defer term.Restore(utilidades.Stdin_fd, st)

	opciones = utilidades.Eliminar_repetidos(opciones)
	opciones = utilidades.Limpiar_listas(opciones)

	return desplegar_opcion(opciones, limite), nil

}
