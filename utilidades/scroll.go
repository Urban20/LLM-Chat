package utilidades

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"
)

type Scroll struct {
	delta         int
	inicio        int
	fin           int
	renglones     []string
	fd            int
	tam_renglones int
}

func (s *Scroll) Renderizar(texto string) {

	s.renglones = strings.Split(texto, "\n")

	s.tam_renglones = len(s.renglones)

}

func imprimir_scroll(renglon []string) {

	fmt.Print(HOME)
	var t string

	for _, l := range renglon {

		t += l+"\n"
	}

	fmt.Println(t)

}

func teclas_scroll(s *Scroll, tecla []byte) {

	flechas := tecla[2]

	if flechas == KEY_ABAJO || tecla[0] == 's' {

		s.fin++
		s.inicio++

	} else if flechas == KEY_ARRIBA || tecla[0] == 'w' {

		s.inicio--
		s.fin--

	}
}

func detectar_tecla() []byte {

	buffer := make([]byte, 3)

	os.Stdin.Read(buffer)

	return buffer

}

func (s Scroll) Iniciar() {

	Limpieza_rapida()
	fmt.Print(OCULTAR_CURSOR)

	st, _ := term.MakeRaw(s.fd)
	defer term.Restore(s.fd, st)

	if s.fin >= s.tam_renglones {

		imprimir_scroll(s.renglones)

		return

	}

	imprimir_scroll(s.renglones[s.inicio:s.fin])

	for {

		tecla := detectar_tecla()

		teclas_scroll(&s, tecla)

		if s.fin >= s.tam_renglones-1 { //si se llega al final del scroll corta el bucle

			break
		}

		if s.inicio <= 0 {

			Limpieza_rapida() //borra cualquier residuo no deseado, evita que se deforme el texto
			s.inicio = 0
			s.fin = s.inicio + s.delta
		}

		imprimir_scroll(s.renglones[s.inicio:s.fin])

	}

}

func Crear_scroll(alto int) *Scroll {

	n_delta := alto - 5

	s := Scroll{delta: n_delta,
		inicio: 0,
		fin:    n_delta,
		fd:     Stdin_fd,
	}

	return &s

}
