package prompts

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

var url = "http://localhost:11434/api/chat"

/*
var prompt_base = `CONTEXTO: La IA debe funcionar obligatoriamente como una IA de alto rendimiento .
Responde siempre como IA (nunca como humano) y de forma breve y no tan amplia.
Sé preciso y orgánico; prioriza eficiencia.
Si falta información pide solo lo imprescindible con preguntas puntuales.
PROMPT: %s`
*/

func recibir_prompt(resp *http.Response) error {

	json_respuesta := Info{}
	var respuesta_str string

	sc := bufio.NewScanner(resp.Body)

	for sc.Scan() {

		if marsherr := json.Unmarshal(sc.Bytes(), &json_respuesta); marsherr != nil {

			return fmt.Errorf("error en unmarshall: %s", marsherr.Error())
		}
		fmt.Print(json_respuesta.Message.Content)
		respuesta_str += json_respuesta.Message.Content
	}

	mensaje_IA := fmt.Sprintf(`{"role":"assistant",
					"content":"%s"}`, respuesta_str)

	Memoria = append(Memoria, mensaje_IA)

	fmt.Print("\n")
	return nil
}

func enviar_prompt(prompt string) (*http.Response, error) {

	mensaje_usuario := fmt.Sprintf(`{
       "role": "user",
       "content": "%s"
    }`, prompt)

	Memoria = append(Memoria, mensaje_usuario)

	json_prompt_usuario := fmt.Sprintf(`{
   "model": "llama3",
   "messages": [
    %s
  	]
	}`, strings.Join(Memoria, ","))

	data := strings.NewReader(json_prompt_usuario)

	resp, resperr := http.Post(url, "aplication/json", data)

	if resperr != nil {

		return resp, fmt.Errorf("error en post: %s", resperr.Error())
	}

	return resp, nil

}

func Comunicacion(prompt string) error {

	//prompt = strings.Replace(fmt.Sprintf(prompt_base, prompt), "\n", `\n`, -1)

	resp, prompterr := enviar_prompt(prompt)

	if prompterr != nil {

		return prompterr
	}

	if recerr := recibir_prompt(resp); recerr != nil {

		return recerr
	}

	return nil

}
