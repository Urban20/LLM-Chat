package prompts

import (
	"LLM-Chat/utilidades"
	"encoding/json"
	"fmt"
	"io"
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

	json_respuesta := Info{}

	b, berror := io.ReadAll(resp.Body)

	if berror != nil {
		return berror
	}

	if jsonerr := json.Unmarshal(b, &json_respuesta); jsonerr != nil {

		return jsonerr
	}

	if mderr := utilidades.Imprimir_markdown("# LLM:\n" + json_respuesta.Message.Content); mderr != nil {

		return mderr
	}

	Guardar_en_memoria(json_respuesta.Message.Content, "LLM (IA)")

	resp.Body.Close()

	return nil
}

// envio el prompt desde el usuario al LLM
func enviar_prompt(prompt, Modelo, Api_chat, Content_type string) (*http.Response, error) {

	Guardar_en_memoria(prompt, "user")

	json_prompt_usuario := Mensaje_usuario{
		Model:    Modelo,
		Messages: Memoria,
		Stream:   false,
	}

	msg_byte, _ := json.Marshal(&json_prompt_usuario)

	data := strings.NewReader(string(msg_byte))

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
func Comunicacion(prompt, modelo, api_chat, content_type string) error {

	resp, prompterr := enviar_prompt(prompt, modelo, api_chat, content_type)

	if prompterr != nil {

		return prompterr
	}

	if recerr := recibir_prompt(resp); recerr != nil {

		return recerr
	}

	return nil

}
