package prompts

import (
	"Cli-ia/utilidades"
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func recibir_prompt(resp *http.Response) error {

	json_respuesta := Info{}
	var respuesta_str string

	defer resp.Body.Close()

	sc := bufio.NewScanner(resp.Body)

	for sc.Scan() {

		if marsherr := json.Unmarshal(sc.Bytes(), &json_respuesta); marsherr != nil {

			return fmt.Errorf("error en unmarshall: %s", marsherr.Error())
		}
		fmt.Print(json_respuesta.Message.Content)
		respuesta_str += json_respuesta.Message.Content
	}

	guardar_en_memoria(respuesta_str, "assistant")

	return nil
}

func enviar_prompt(prompt string) (*http.Response, error) {

	guardar_en_memoria(prompt, "user")

	json_prompt_usuario := fmt.Sprintf(`{
   "model": "%s",
   "messages": [%s]
	}`, utilidades.Modelo, strings.Join(Memoria, ","))

	data := strings.NewReader(json_prompt_usuario)

	resp, resperr := http.Post(utilidades.Api_chat, utilidades.Content_type, data)

	if resp.StatusCode != http.StatusOK {

		return resp, fmt.Errorf("hubo un problema con la solicitud post, codigo de estado: %d", resp.StatusCode)
	}

	if resperr != nil {

		return resp, fmt.Errorf("error en post: %s", resperr.Error())
	}

	return resp, nil

}

func Comunicacion(prompt string) error {

	resp, prompterr := enviar_prompt(prompt)

	if prompterr != nil {

		return prompterr
	}

	if recerr := recibir_prompt(resp); recerr != nil {

		return recerr
	}

	return nil

}
