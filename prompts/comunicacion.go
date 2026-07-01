package prompts

import (
	"Cli-ia/utilidades"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

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

	utilidades.Guardar_en_memoria(json_respuesta.Message.Content, "LLM (IA)")

	resp.Body.Close()

	return nil
}

// envio el prompt desde el usuario al LLM
func enviar_prompt(prompt, Modelo, Api_chat, Content_type string) (*http.Response, error) {

	utilidades.Guardar_en_memoria(prompt, "user")

	json_prompt_usuario := fmt.Sprintf(`{
   "model": "%s",
   "messages": [%s],
   "stream":false
	}`, Modelo, strings.Join(utilidades.Memoria, ",")) // no puedo pasarlo a mapa y a json porque se pierde el formateo

	data := strings.NewReader(json_prompt_usuario)

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
