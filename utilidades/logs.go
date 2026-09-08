package utilidades

import (
	"fmt"
)

var formateo string = "%s\t%s %v\n" //hora, prefijo, mensaje

const (
	FLAG_INFO        = "INFO"
	FLAG_ERROR       = "ERROR"
	FLAG_ADVERTENCIA = "ADVERTENCIA"
)

func rendColor(color, valor string) string {

	return fmt.Sprintf("%s%v%s", color, valor, RESET)

}

func Info(mensaje any) {

	fmt.Printf(formateo, rendColor(RESET, Tiempo_actual()), rendColor(VERDE_BRILLANTE, FLAG_INFO), mensaje)

}

func Advertencia(mensaje any) {

	fmt.Printf(formateo, rendColor(RESET, Tiempo_actual()), rendColor(AMARILLO, FLAG_ADVERTENCIA), mensaje)

}

func Error(mensaje any) {

	fmt.Printf(formateo, rendColor(RESET, Tiempo_actual()), rendColor(ROJO_BRILLANTE, FLAG_ERROR), mensaje)
}
