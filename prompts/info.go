package prompts

// modulo que contiene la informacion cruda para procesar y extrae los datos

type Info struct {
	Model      string  `json:"model"`
	Created_at string  `json:"created_at"`
	Message    message `json:"message"`
	Done       bool    `json:"done"`
}

type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// estas structs se usan unicamente para parsear el json
// para comunicarse con la ia uso un mapa
