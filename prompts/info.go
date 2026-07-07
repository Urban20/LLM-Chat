package prompts

// modulo que contiene la informacion cruda para procesar y extrae los datos

type Info struct {
	Model       string  `json:"model"`
	Created_at  string  `json:"created_at"`
	Message     message `json:"message"`
	Done        bool    `json:"done"`
	Done_reason string  `json:"done_reason"`
}

type message struct {
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

type Mensaje_usuario struct {
	Model    string
	Messages []message
	Stream   bool
	Options  Opciones
}

type Opciones struct {
	num_ctx     int //controla tokens totales (memoria de trabajo total)
	num_predict int // sin limite de generacion de tokens (limite de tokens)
	temperature float64
}
