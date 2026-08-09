package prompts

// modulo que contiene la informacion cruda para procesar y extrae los datos

type Info struct {
	// para chat
	Model       string       `json:"model"`
	Created_at  string       `json:"created_at"`
	Message     message_chat `json:"message"`
	Done        bool         `json:"done"`
	Done_reason string       `json:"done_reason"`

	// para generate
	Response string `json:"response"`
	Thinking string `json:"thinking"`
}

type Mensaje_usuario_generate struct { // envio al server
	Model   string
	Prompt  string
	Images  []string //base64
	Options Opciones
}

type message_chat struct {
	Role     string `json:"role"`
	Content  string `json:"content"`
	Thinking string `json:"thinking"`
}

// estas structs se usan unicamente para parsear el json
// para comunicarse con la ia uso un mapa

type modelo struct { // esto lo uso con la api de tags para listar los modelos disponibles

	Name         string   `json:"name"`
	Model        string   `json:"model"`
	Capabilities []string `json:"capabilities"` // (capacidades de los LLMs) no lo voy a usar pero quiza en un futuro me sirve
}

type Modelos struct {
	Models []modelo `json:"models"`
}

type Mensaje_usuario_chat struct {
	Model    string
	Messages []message_chat
	Stream   bool
	Options  Opciones
}

type Opciones struct {
	Num_ctx     int     `json:"num_ctx"`     //controla tokens totales (memoria de trabajo total)
	Num_predict int     `json:"bum_predict"` // sin limite de generacion de tokens (limite de tokens)
	Temperature float64 `json:"temperature"`
}
