package prompts

import "encoding/json"

func guardar_en_memoria(prompt, rol string) {

	mensaje_usuario := map[string]string{
		"role":    rol,
		"content": prompt,
	}

	msg, _ := json.Marshal(&mensaje_usuario)
	Memoria = append(Memoria, string(msg))

}
