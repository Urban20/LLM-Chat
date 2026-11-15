package menu

import (
	"LLM-Chat/utilidades"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/pterm/pterm"
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
