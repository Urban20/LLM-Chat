package menu

import (
	"LLM-Chat/utilidades"
	"fmt"
	"os"

	"golang.org/x/term"
)

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

	fmt.Print("\033[?25l")
	defer fmt.Print("\033[?25h")

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
