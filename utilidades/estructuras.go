package utilidades

import (
	"fmt"
	"runtime"
	"strconv"
)

type Box_info struct {
	Modelo      string
	Sistema_op  string
	Temperatura float64
	Ctx         int
	Host        string
	Puerto      int
}

func (b Box_info) Box_informacion() {

	Limpieza_rapida()

	contenido_box := map[string]string{

		"Modelo selecionado":  b.Modelo,
		"Host":                fmt.Sprintf("%s:%d", b.Host, b.Puerto),
		"Sistema operativo":   runtime.GOOS,
		"Temperatura del LLM": fmt.Sprintf("%.2f", b.Temperatura),
		"Contexto del LLM":    strconv.Itoa(b.Ctx),
	}
	contenidos := Formato_string_box(contenido_box)
	Box(contenidos...)
}

type Estructura_prompt struct {
	// da la informacion que se le envian automaticamente al LLM

	Fecha  string
	Prompt string
}

func (e Estructura_prompt) Formatear_prompt() string {

	data := map[string]string{
		// la info la pongo en ingles como lenguaje neutro para el LLM

		"[CURRENT DATETIME]": e.Fecha,
		"[PROMPT]":           e.Prompt,
	}

	var instruccion string

	for c, v := range data {

		instruccion += fmt.Sprintf("%s\n\n%s\n\n", c, v)

	}

	return instruccion

}

type Prompt_archivo struct { //struct para manejar el envio de informacion al llm
	Prompt   string
	Archivos []string
}

func (s *Prompt_archivo) Borrar_informacion() {

	s.Prompt = ""
	s.Archivos = []string{}

}

type Respuesta_LLM struct {
	Respuesta_raw string
	Tokens        int
	Prompt        string // prompt del usuario
	Modelo        string
}
