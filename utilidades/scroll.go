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
	habilitar     bool
	tam_renglones int
	eje_y         int
}

func (s *Scroll) Renderizar(texto string) {

	s.renglones = strings.Split(texto, "\n")

	s.tam_renglones = len(s.renglones)

}

func imprimir_scroll(renglon []string) {

	fmt.Print(HOME)

	for _, l := range renglon {

		fmt.Println(l)
	}

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

	for s.habilitar { //TODO: verificar que no haya errores y refaccionar

		tecla := detectar_tecla()

		teclas_scroll(&s, tecla)

		if s.fin >= s.tam_renglones-1 { //si se llega al final del scroll corta el bucle

			s.habilitar = false
			return
		}

		if s.inicio <= 0 {

			Limpieza_rapida() //borra cualquier residuo no deseado, evita que se deforme el texto
			s.inicio = 0
			s.fin = s.inicio + s.delta
		}

		imprimir_scroll(s.renglones[s.inicio:s.fin])

	}

}

func Crear_scroll() *Scroll {

	_, y, _ := term.GetSize(Stdout_fd)

	n_delta := y - 5 //TODO : ajustar para evitar bugs visuales

	s := Scroll{delta: n_delta,
		inicio:    0,
		fin:       n_delta,
		fd:        Stdin_fd,
		habilitar: true,
		eje_y:     y,
	}

	return &s

}
