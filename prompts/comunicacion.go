package prompts

import (
	"LLM-Chat/menu"
	"LLM-Chat/utilidades"
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strings"
	"sync"
)

var Memoria []message_chat

func Borrar_memoria() {

	Memoria = []message_chat{}

}

func Guardar_en_memoria(prompt, rol string) {

	mensaje_usuario := message_chat{Role: rol, Content: prompt}

	Memoria = append(Memoria, mensaje_usuario)

}

// recibe una struct y la envia por POST al servidor
func struct_a_respuesta(info any, endpoint, content_type string) (*http.Response, error) {

	msg_byte, jsonerr := json.Marshal(info)

	if jsonerr != nil {
		return &http.Response{}, jsonerr
	}

	data := bytes.NewReader(msg_byte)

	resp, resperr := http.Post(endpoint, content_type, data)

	if resperr != nil {

		return resp, resperr
	}

	if resp.StatusCode != http.StatusOK {

		return resp, fmt.Errorf("hubo un problema con la solicitud post, codigo de estado: %d", resp.StatusCode)
	}

	return resp, nil

}

// recibo el prompt desde el LLM al usuario
func recibir_prompt(resp *http.Response, carga *menu.Carga, wg *sync.WaitGroup, chat bool) error {

	var cuerpo string

	escaner := bufio.NewScanner(resp.Body)
	defer resp.Body.Close()

	carga.Detener(wg)

	fmt.Print("\n\n" + utilidades.GRIS_AZUL)
	for escaner.Scan() {

		json_respuesta := Info{}

		if marsherr := json.Unmarshal(escaner.Bytes(), &json_respuesta); marsherr != nil {

			return marsherr
		}

		if json_respuesta.Done_reason == "length" {

			return errors.New("se agoto el contexto disponible para la generacion de nuevas respuestas")

		}

		if !slices.Contains([]string{"", "stop"}, json_respuesta.Done_reason) {

			return fmt.Errorf("se interrumpio la generacion de tokens desde el servidor, razon: %s", json_respuesta.Done_reason)

		}

		if chat { // para chat

			cuerpo += json_respuesta.Message.Content
			fmt.Print(json_respuesta.Message.Thinking) //depende del modelo que se use

		} else { // para generate

			cuerpo += json_respuesta.Response
			fmt.Print(json_respuesta.Thinking)
		}

	}

	cuerpo = strings.TrimSpace(cuerpo)

	if cuerpo == "" {

		return errors.New("la respuesta llego vacia")
	}

	Guardar_en_memoria(cuerpo, "LLM (IA)")

	utilidades.Limpieza_rapida()

	if markerr := utilidades.Imprimir_markdown(cuerpo); markerr != nil {

		return markerr
	}

	return nil
}

// envio el prompt desde el usuario al LLM
func enviar_prompt(prompt, Modelo, endpoint, Content_type string, ctx int, temp float64, chat bool, imagenes []string) (*http.Response, error) {

	var json_prompt_usuario any

	Guardar_en_memoria(prompt, "user")

	opciones := Opciones{
		num_ctx:     ctx,
		num_predict: -1,
		temperature: temp,
	}

	if chat {

		json_prompt_usuario = Mensaje_usuario_chat{

			Model:    Modelo,
			Messages: Memoria,
			Stream:   true,
			Options:  opciones,
		}

	} else { //generate

		json_prompt_usuario = Mensaje_usuario_generate{

			Model:  Modelo,
			Prompt: prompt,
			Images: imagenes,
		}

	}

	return struct_a_respuesta(json_prompt_usuario, endpoint, Content_type)

}

// esta funcion se ocupa del envio y recepcion de los mensajes
func Comunicacion(prompt, modelo, endpoint, content_type string, ctx int, temp float64, carga *menu.Carga, wg *sync.WaitGroup, chat bool, imagenes []string) error {

	resp, prompterr := enviar_prompt(prompt, modelo, endpoint, content_type, ctx, temp, chat, imagenes)

	defer carga.Detener(wg)

	if prompterr != nil {

		return prompterr
	}

	if recerr := recibir_prompt(resp, carga, wg, chat); recerr != nil {

		return recerr
	}

	return nil

}
