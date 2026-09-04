package menu

import (
	"LLM-Chat/utilidades"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"golang.org/x/term"
)

type Carga struct {
	estado_1 string
	estado_2 string
	cargando bool
	tiempo   float32
}

const (
	OCULTAR_CURSOR = "\033[?25l"
	MOSTRAR_CURSOR = "\033[?25h"
	VERSION        = "V1.1"
)

const (
	KEY_ARRIBA   = 65
	KEY_ABAJO    = 66
	ENTER        = 13
	TIEMPO_CARGA = 0.85
)

func Logo() {

	logo := ` ██╗     ██╗     ███╗   ███╗       ██████╗██╗  ██╗ █████╗ ████████╗
 ██║     ██║     ████╗ ████║      ██╔════╝██║  ██║██╔══██╗╚══██╔══╝
 ██║     ██║     ██╔████╔██║█████╗██║     ███████║███████║   ██║   
 ██║     ██║     ██║╚██╔╝██║╚════╝██║     ██╔══██║██╔══██║   ██║   
 ███████╗███████╗██║ ╚═╝ ██║      ╚██████╗██║  ██║██║  ██║   ██║   
 ╚══════╝╚══════╝╚═╝     ╚═╝       ╚═════╝╚═╝  ╚═╝╚═╝  ╚═╝   ╚═╝ %s`

	fmt.Print(utilidades.NEGRO_BLANCO)

	fmt.Printf("\n"+logo+"\n", VERSION)

	fmt.Print(utilidades.RESET)

}

func (p *Carga) Iniciar(wg *sync.WaitGroup) {

	fmt.Print("\n\n" + OCULTAR_CURSOR)
	defer fmt.Print(MOSTRAR_CURSOR)

	wg.Add(1)
	defer wg.Done()

	estados := []string{p.estado_1, p.estado_2}
	var i int

	for p.cargando {

		if i > len(estados)-1 {
			i = 0
		}

		fmt.Printf("\r%s", estados[i])
		i++
		time.Sleep(time.Second * time.Duration(p.tiempo))

	}

}

func (p *Carga) Detener(wg *sync.WaitGroup) {

	p.cargando = false
	wg.Wait()
	fmt.Print("\r" + strings.Repeat(" ", len(p.estado_1)))
}

func Crear_carga() Carga {

	c := Carga{estado_1: "◌◌◌",
		estado_2: "●●●",
		cargando: true,
		tiempo:   TIEMPO_CARGA}

	return c

}

func actualizar_seccion(n int) {

	for x := 0; x < n; x++ {

		fmt.Print("\033[F")

	}

}

func leer_tecla(i *int, tecla []byte) bool {
	os.Stdin.Read(tecla)
	flechas := tecla[2]

	if tecla[0] == ENTER {

		return true

	} else if flechas == KEY_ARRIBA || tecla[0] == 'w' {

		*i--
	} else if flechas == KEY_ABAJO || tecla[0] == 's' {

		*i++
	}

	return false

}

func desplegar_opcion(opciones []string) string {

	var i int
	var op_largo = len(opciones)

	for {
		tecla := make([]byte, 3)

		for _, op := range opciones {

			if i > op_largo-1 {
				i = 0

			} else if i < 0 {
				i = op_largo - 1
			}

			if op == opciones[i] { // opcion seleccionada
				fmt.Println(utilidades.NEGRO_BLANCO + "> " + op + utilidades.RESET + "\r")
			} else {
				fmt.Println("  " + op + "\r")
			}
		}

		if leer_tecla(&i, tecla) {

			return opciones[i]
		}

		actualizar_seccion(op_largo)

	}
}

func Menu(opciones ...string) (string, error) {

	fmt.Print(OCULTAR_CURSOR)
	defer fmt.Print(MOSTRAR_CURSOR)

	fmt.Print(utilidades.GRIS_AZUL + "\nOpciones disponibles:\n\n")
	fmt.Print(utilidades.AZUL_OSCURO + "navegar con ↑↓\n\n" + utilidades.RESET)

	fd := int(os.Stdin.Fd())

	st, rawerr := term.MakeRaw(fd)

	if rawerr != nil {
		return "", rawerr
	}

	defer term.Restore(fd, st)

	opciones = utilidades.Eliminar_repetidos(opciones)
	return desplegar_opcion(opciones), nil

}
