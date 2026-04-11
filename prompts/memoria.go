package prompts

import "encoding/json"

func guardar_en_memoria(prompt string) {

	mensaje_usuario := map[string]string{
		"role":    "user",
		"content": prompt,
	}

	msg, _ := json.Marshal(&mensaje_usuario)
	Memoria = append(Memoria, string(msg))

}
