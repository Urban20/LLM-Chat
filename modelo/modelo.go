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

*/

func enviar_instruccion(instrucciones []byte, api_modelo, content_type string) error {

	instruccion_post := strings.NewReader(string(instrucciones))

	_, posterr := http.Post(api_modelo, content_type, instruccion_post)
	if posterr != nil {
		return posterr
	}
	return nil

}

func Crear_modelo(IA, modelo, api_modelo, content_type string) error {

	/*
		intenta crear comportamiento que va a tener la ia
		por el momento se crea cada vez que se inicia la CLI
	*/

	fmt.Println("iniciando modelo...")

	data := map[string]string{
		"from":   IA,
		"model":  modelo,
		"system": utilidades.Instruccion,
	}

	instrucciones, marsherr := json.Marshal(&data)

	if marsherr != nil {
		return marsherr
	}

	if err := enviar_instruccion(instrucciones, api_modelo, content_type); err != nil {
		return err
	}

	return nil

}
