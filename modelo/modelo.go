package modelo

import (
	"Cli-ia/utilidades"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

/*
en esta seccion voy a definir como quiero que se
comporte el modelo de ia, si el modelo ya fue creado no se vuelve a crear


{
  "from": "gemma3",
  "model": "alpaca",
  "system": "You are Alpaca, a helpful AI assistant. You only answer with Emojis."
}

*/

func Conectar() bool {

	resp, resperr := http.Post(utilidades.Info_modelo, utilidades.Content_type, utilidades.Json_modelo)

	if resperr != nil {
		//fmt.Println(resperr)
		return false
	}

	return resp.StatusCode == http.StatusOK

}

func Crear_modelo() {

	fmt.Println("creando modelo...")

	data := map[string]string{
		"from":   utilidades.IA,
		"model":  utilidades.Modelo,
		"system": utilidades.Instruccion,
	}

	instrucciones, _ := json.Marshal(&data)

	instruccion_post := strings.NewReader(string(instrucciones))

	http.Post(utilidades.Api_modelo, utilidades.Content_type, instruccion_post)

}
