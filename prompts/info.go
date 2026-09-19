package prompts

// modulo que contiene la informacion cruda para procesar y extrae los datos

/*
"tool_calls": [

	{
	  "function": {
	    "name": "calculate_average",
	    "arguments": {"numbers": [10, 20, 30, 40]}
	  }
*/
/*
  "role": "tool",
      "tool_name": "calculate_average",
      "content": "25.0"	  */

type Datos_tool_call struct {
	Nombre_herramienta string `json:"tool_name"`
	Retorno            string `json:"content"`
}

type funcion_respuesta struct {

	//corresponde a la respuesta del LLM , no confundir
	Nombre     string         `json:"name"`
	Argumentos map[string]any `json:"arguments"`
}

type Tool_call struct {
	//corresponde a la respuesta del LLM , no confundir
	Funcion funcion_respuesta `json:"function"`
}

type Info struct {
	// ambos

	Num_tokens_prompt int `json:"prompt_eval_count"`
	Num_tokens_resp   int `json:"eval_count"`

	// para chat
	Model       string       `json:"model"`
	Created_at  string       `json:"created_at"`
	Message     message_chat `json:"message"`
	Done        bool         `json:"done"`
	Done_reason string       `json:"done_reason"`
	Tools_calls []Tool_call  `json:"tools_calls"`

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
	Role        string      `json:"role"`
	Content     string      `json:"content"`
	Thinking    string      `json:"thinking"`
	Tools_calls []Tool_call `json:"tools_calls"` //se usa para obtener la tool necesaria para el llm
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
	Model    string         `json:"model"`
	Messages []message_chat `json:"messages"`
	Stream   bool           `json:"stream"`
	Options  Opciones       `json:"options"`
	Tools    []Herramienta  `json:"tools"` //se envia para mostrarle al llm las tools disponibles

}

type Opciones struct {
	Num_ctx     int     `json:"num_ctx"`     //controla tokens totales (memoria de trabajo total)
	Num_predict int     `json:"bum_predict"` // sin limite de generacion de tokens (limite de tokens)
	Temperature float64 `json:"temperature"`
}

/*

"tools": [
    {
      "type": "function",
      "function": {
        "name": "get_temperature",
        "description": "Get the current temperature for a city",
        "parameters": {
          "type": "object",
          "required": ["city"],
          "properties": {
            "city": {"type": "string", "description": "The name of the city"}
          }
        }
      }
    }
  ]

*/

type Prop_parametro struct {
	Tipo        string `json:"type"`
	Descripcion string `json:"description"`
}

type parametros struct {
	Tipo        string                    `json:"type"`
	Requerido   []string                  `json:"required"`
	Propiedades map[string]Prop_parametro `json:"properties"`
}
type funcion struct {
	Nombre      string     `json:"name"`
	Descripcion string     `json:"description"`
	Parametros  parametros `json:"parameters"`
}

type Herramienta struct {
	Tipo    string  `json:"type"`
	Funcion funcion `json:"function"`
}
