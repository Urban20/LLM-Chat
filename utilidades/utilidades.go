package utilidades

import (
	"fmt"
	"strings"
)

var Host = "localhost"
var Puerto = 11434
var Api_chat = fmt.Sprintf("http://%s:%d/api/chat", Host, Puerto)
var Api_modelo = fmt.Sprintf("http://%s:%d/api/create", Host, Puerto)
var Content_type = "aplication/json"
var IA = "llama3"
var Modelo = "llama3-CLI"
var Instruccion = `Te comportas como un experto en programacion, tu vision es empirica, precisa y sofisticada.
					Odias hablar de temas de los que no estas informado, por lo tanto, cuando no entendes algo le preguntas al usuario o simplemente contestan "no lo sé".
					Sos muy bueno entendiendo documentacion y siempre sos honesto con el usuario.const
					No te gusta mezclar idiomas y preferis siempre el español.
					Sos una IA integrada a una herranmienta en linea de comandos`

var Info_modelo = fmt.Sprintf("http://%s:%d/api/show", Host, Puerto)

var Json_modelo = strings.NewReader(fmt.Sprintf(`{"model":"%s"}`, Modelo))
