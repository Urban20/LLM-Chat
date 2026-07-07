package prompts

import (
	"LLM-Chat/utilidades"
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

var Memoria = []message{}

func Borrar_memoria() {

	Memoria = []message{}

}

func Guardar_en_memoria(prompt, rol string) {

	mensaje_usuario := message{Role: rol, Content: prompt}

	Memoria = append(Memoria, mensaje_usuario)

}

// recibo el prompt desde el LLM al usuario
func recibir_prompt(resp *http.Response) error {

	var cuerpo string

	escaner := bufio.NewScanner(resp.Body)
	defer resp.Body.Close()

	fmt.Print(strings.Repeat("\n", 4))

	for escaner.Scan() {

		json_respuesta := Info{}

		if marsherr := json.Unmarshal(escaner.Bytes(), &json_respuesta); marsherr != nil {

			return marsherr
		}

		fmt.Print(json_respuesta.Message.Thinking) //depende del modelo que se use

		cuerpo += json_respuesta.Message.Content

	}

	if markerr := utilidades.Imprimir_markdown("# LLM:\n" + strings.TrimSpace(cuerpo)); markerr != nil {

		return markerr
	}

	return nil
}

// envio el prompt desde el usuario al LLM
func enviar_prompt(prompt, Modelo, Api_chat, Content_type string, ctx int, temp float64) (*http.Response, error) {

	Guardar_en_memoria(prompt, "user")

	opciones := Opciones{
		num_ctx:     ctx,
		num_predict: -1,
		temperature: temp,
	}

	json_prompt_usuario := Mensaje_usuario{

		Model:    Modelo,
		Messages: Memoria,
		Stream:   true,
		Options:  opciones,
	}

	msg_byte, jsonerr := json.Marshal(&json_prompt_usuario)

	if jsonerr != nil {
		return &http.Response{}, jsonerr
	}

	data := bytes.NewReader(msg_byte)

	resp, resperr := http.Post(Api_chat, Content_type, data)

	if resperr != nil {

		return resp, resperr
	}

	if resp.StatusCode != http.StatusOK {

		return resp, fmt.Errorf("hubo un problema con la solicitud post, codigo de estado: %d", resp.StatusCode)
	}

	return resp, nil

}

// esta funcion se ocupa del envio y recepcion de los mensajes
func Comunicacion(prompt, modelo, api_chat, content_type string, ctx int, temp float64) error {

	resp, prompterr := enviar_prompt(prompt, modelo, api_chat, content_type, ctx, temp)

	if prompterr != nil {

		return prompterr
	}

	if recerr := recibir_prompt(resp); recerr != nil {

		return recerr
	}

	return nil

}
