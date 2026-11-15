package menu

import (
	"LLM-Chat/utilidades"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/x/term"
)

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

func leer_tecla(i, pag *int, tecla []byte, paginar bool) bool {
	os.Stdin.Read(tecla)
	flechas := tecla[2]

	if paginar && flechas == utilidades.KEY_DER {

		*pag++
	}
	if paginar && flechas == utilidades.KEY_IZQ {

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

// evita que el paginado tome valores incorrectos
func paginado(pag, max_len int) int {

	if pag > max_len-1 {

		pag = max_len - 1

	} else if pag < 0 {

		pag = 0
	}

	return pag

}

func imprimir_opciones(op string, actual []string, i int) {

	selec := actual[i]

	if op == selec { // opcion seleccionada

		fmt.Printf("%s> %s%s%s\r\n", utilidades.NEGRO_BLANCO, op, strings.Repeat(" ", utilidades.Max(actual)-len(selec)), utilidades.RESET)
	} else {
		fmt.Printf("  %s\r\n", op)
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

func desplegar_opcion(opciones []string, limite int, paginar bool) string {

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

		if paginar {

			fmt.Printf("%spagina %d/%d%s\n\n\r", utilidades.NEGRO_BLANCO, pag+1, max_len, utilidades.RESET)
			margen = 2
		}

		for _, op := range actual {

			i = limitar_indice(i, largo_op)

			imprimir_opciones(op, actual, i)
		}

		if leer_tecla(&i, &pag, tecla, paginar) {

			return actual[i]
		}

		actualizar_seccion(largo_op+margen, // le sumo los margenes de pagina
			utilidades.Max(actual)+2) // le sumo el margen

	}
}

type Menu struct {
	instrucciones string
	paginar       bool
	limite        int //limite de opciones por pagina
	opciones      []string
}

func CrearMenu() *Menu {

	MenuDefault := Menu{
		paginar: true,
		limite:  5,
	}

	return MenuDefault.SetInstrucciones("navegar con ↑↓", "← → cambiar pagina")

}

// habilitar / deshabilitar navegacion de paginas
func (m Menu) SetPaginar(paginar bool) *Menu {

	m.paginar = paginar

	return &m

}

func (m Menu) SetInstrucciones(instrucciones ...string) *Menu {

	m.instrucciones = strings.Join(instrucciones, " | ")

	return &m

}

// setea el limite de opciones a mostrar, si el paginado esta deshabilitado y hay mas opciones se truncan
func (m Menu) SetLimite(limite int) *Menu {

	if limite > 0 {

		m.limite = limite
	}

	return &m

}

func (m Menu) SetOpciones(opciones []string) *Menu {

	m.opciones = opciones

	return &m

}

func (m Menu) Desplegar() (string, error) {

	fmt.Print(utilidades.OCULTAR_CURSOR)
	defer fmt.Print(utilidades.MOSTRAR_CURSOR)

	fmt.Print(utilidades.GRIS_AZUL, "\nOpciones disponibles:\n\n")
	fmt.Print(utilidades.AZUL_OSCURO, m.instrucciones, "\n\n", utilidades.RESET)

	st, rawerr := term.MakeRaw(uintptr(utilidades.Stdin_fd))

	if rawerr != nil {
		return "", rawerr
	}

	defer func() { rawerr = term.Restore(uintptr(utilidades.Stdin_fd), st) }()

	return desplegar_opcion(
		utilidades.Limpiar_listas(utilidades.Eliminar_repetidos(m.opciones)),
		m.limite,
		m.paginar), rawerr

}
