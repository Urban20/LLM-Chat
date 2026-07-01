package consola

import (
	"os"

	"golang.org/x/sys/windows"
)

// obtiene el estado actual del stdout y habilita la secuencia ansi para soportar el programa
func Iniciar_ANSI() error {

	var modo_original uint32

	stdout := windows.Handle(os.Stdout.Fd())
	// guardo el estado de consola actual
	if conserr := windows.GetConsoleMode(stdout, &modo_original); conserr != nil {
		return conserr

	}
	if winerr := windows.SetConsoleMode(stdout, modo_original|windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING); winerr != nil {
		//por ultimo hago que windows inteprete el codigo ansi, le sumo la flag para interpretar ansi
		return winerr

	}

	return nil
}
