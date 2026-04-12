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

func Conectar() bool { // TODO: quiza esta funcion me convenga eliminarla (ver que hago)

	resp, resperr := http.Post(utilidades.Info_modelo, utilidades.Content_type, utilidades.Json_modelo)

	if resperr != nil {
		//fmt.Println(resperr)
		return false
	}

	return resp.StatusCode == http.StatusOK

}

func enviar_instruccion(instrucciones []byte) error {

	instruccion_post := strings.NewReader(string(instrucciones))

	_, posterr := http.Post(utilidades.Api_modelo, utilidades.Content_type, instruccion_post)
	if posterr != nil {
		return posterr
	}
	return nil

}

func Crear_modelo() error {

	/*
		intenta crear comportamiento que va a tener la ia
		por el momento se crea cada vez que se inicia la CLI
	*/

	fmt.Println("iniciando modelo...")

	data := map[string]string{
		"from":   utilidades.IA,
		"model":  utilidades.Modelo,
		"system": utilidades.Instruccion,
	}

	instrucciones, marsherr := json.Marshal(&data)

	if marsherr != nil {
		return marsherr
	}

	if err := enviar_instruccion(instrucciones); err != nil {
		return err
	}

	return nil

}
