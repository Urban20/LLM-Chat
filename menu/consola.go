package menu

import (
	"LLM-Chat/utilidades"
	"fmt"
	"os"
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
)

func (p *Carga) Iniciar(wg *sync.WaitGroup) {

	fmt.Print("\n\n" + OCULTAR_CURSOR)
	defer fmt.Print(MOSTRAR_CURSOR)

	wg.Add(1)
	defer wg.Done()

	for p.cargando {

		for _, estado := range []string{p.estado_1, p.estado_2} {

			fmt.Print("\r" + estado)
			time.Sleep(time.Second * time.Duration(p.tiempo))

		}
	}

}

func (p *Carga) Detener() {

	p.cargando = false
}

func Crear_carga() Carga {

	c := Carga{estado_1: "◌◌◌",
		estado_2: "●●●",
		cargando: true,
		tiempo:   0.85}

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

	if tecla[0] == 13 { // enter

		return true

	} else if flechas == 65 {

		*i--
	} else if flechas == 66 {

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
				fmt.Println(utilidades.VIOLETA + "> " + op + utilidades.RESET + "\r")
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

	return desplegar_opcion(opciones), nil

}
