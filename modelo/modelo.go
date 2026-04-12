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

func Modelo_existe() bool {

	info_modelo := strings.NewReader(fmt.Sprintf(`{"model":"%s"}`, strings.ReplaceAll(utilidades.Modelo, "\n", `\n`)))

	resp, resperr := http.Post(utilidades.Info_modelo, utilidades.Content_type, info_modelo)

	if resperr != nil {
		fmt.Println(resperr)
		return false
	}
	fmt.Println(resp.StatusCode)

	return resp.StatusCode == http.StatusOK

}

func Crear_modelo() {

	if Modelo_existe() {
		fmt.Println("el modelo ya fue creado previamente")
		return
	}

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
