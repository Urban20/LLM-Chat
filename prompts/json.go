package prompts

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

var url = "http://localhost:11434/api/generate"

/*
var prompt_base = `CONTEXTO: La IA debe funcionar obligatoriamente como una IA de alto rendimiento .
Responde siempre como IA (nunca como humano) y de forma breve y no tan amplia.
Sé preciso y orgánico; prioriza eficiencia.
Si falta información pide solo lo imprescindible con preguntas puntuales.
PROMPT: %s`
*/

func Enviar_prompt(prompt string) (Info, error) {

	json_respuesta := Info{}

	//prompt = strings.Replace(fmt.Sprintf(prompt_base, prompt), "\n", `\n`, -1)

	data := strings.NewReader(fmt.Sprintf(`{"model": "llama3",
								"prompt": "%s",
								"stream": false
								}`, prompt))

	resp, resperr := http.Post(url, "aplication/json", data)

	if resperr != nil {

		return json_respuesta, resperr
	}

	b, _ := io.ReadAll(resp.Body)

	if marsherr := json.Unmarshal(b, &json_respuesta); marsherr != nil {

		return json_respuesta, marsherr
	}

	return json_respuesta, nil

}
