package prompts

import (
	"LLM-Chat/menu"
	"LLM-Chat/utilidades"
	"bufio"
	"bytes"
	"encoding/json"
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

func Crear_tool_scraping() Herramienta {

	prop_scraping := map[string]Prop_parametro{ // en ingles porque es leido por el llm (lenguaje base del llm)

		"url": {Tipo: "string", Descripcion: "URL of the website to scrape"},
	}

	scraping := Crear_tool([]string{"url"},
		"function to extract data from a website",
		"Scraping_web",
		prop_scraping)

	return *scraping

}

func Crear_tool(requerido []string, descripcion, nombre_tool string, props map[string]Prop_parametro) *Herramienta {

	/*
	   requerido: parametros de la funcion requeridos
	   descripcion: descripcion breve de lo que hace la funcion
	   props: propiedades de los parametros que se les provee

	*/

	param := parametros{

		Tipo:        "object",
		Requerido:   requerido,
		Propiedades: props,
	}

	f := funcion{
		Nombre:      nombre_tool,
		Descripcion: descripcion,
		Parametros:  param,
	}

	tool := Herramienta{

		Tipo:    "function",
		Funcion: f,
	}

	return &tool

}

// quita el modelo de la carga (no tiene nada que ver con la instalacion de un nuevo modelo)
func Descargar_modelo(modelo, content_type, endpoint string) {

	type Payload struct {
		Model      string `json:"model"`
		Keep_alive int    `json:"keep_alive"`
	}

	payload := Payload{Model: modelo, Keep_alive: 0}

	b, _ := json.Marshal(&payload)

	data := bytes.NewReader(b)

	http.Post(endpoint, content_type, data)

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
func recibir_prompt(resp *http.Response, carga *menu.Carga, wg *sync.WaitGroup, chat bool, prompt string, box utilidades.Box_info) utilidades.Respuesta_LLM {

	var cuerpo string
	respuesta := utilidades.Respuesta_LLM{Modelo: box.Modelo}

	escaner := bufio.NewScanner(resp.Body)

	defer resp.Body.Close()

	carga.Detener(wg)

	fmt.Print("\n\n" + utilidades.GRIS_AZUL)

	for escaner.Scan() {

		json_respuesta := Info{}

		if marsherr := json.Unmarshal(escaner.Bytes(), &json_respuesta); marsherr != nil {

			fmt.Print("\n\n")
			utilidades.Error(marsherr)

			return respuesta
		}

		respuesta.Tokens += json_respuesta.Num_tokens_prompt + json_respuesta.Num_tokens_resp

		if json_respuesta.Done_reason == "length" {

			fmt.Print("\n\n")

			utilidades.Advertencia("se agoto el contexto disponible para la generacion de nuevas respuestas")
			return respuesta

		}

		if !slices.Contains([]string{"", "stop"}, json_respuesta.Done_reason) {

			fmt.Print("\n\n")
			utilidades.Advertencia(fmt.Sprintf("se interrumpio la generacion de tokens desde el servidor, razon: %s", json_respuesta.Done_reason))
			return respuesta

		}

		if chat { // para chat

			cuerpo += json_respuesta.Message.Content
			fmt.Print(json_respuesta.Message.Thinking) //depende del modelo que se use

		} else { // para generate

			cuerpo += json_respuesta.Response
			fmt.Print(json_respuesta.Thinking)
		}

	}

	respuesta.Respuesta_raw = strings.TrimSpace(cuerpo)
	respuesta.Prompt = prompt

	return respuesta

}

func procesar_respuesta(r utilidades.Respuesta_LLM) {

	defer fmt.Printf("\n\n%stokens totales generados: %d%s\n", utilidades.FONDO_VERDE, r.Tokens, utilidades.RESET)

	if r.Respuesta_raw == "" {

		fmt.Print("\n\n")
		utilidades.Error("la respuesta llego vacia")
		return
	}

	Guardar_en_memoria(r.Respuesta_raw, "LLM (IA)")

	if err := utilidades.Imprimir_markdown(r); err != nil {

		fmt.Print("\n\n")
		utilidades.Error(err)
		return
	}

}

// envio el prompt desde el usuario al LLM
func enviar_prompt(prompt string, box utilidades.Box_info, endpoint, Content_type string, chat bool, imagenes []string) (*http.Response, error) {
	// TODO : meter una struct para ordenar, cambiar los inputs de la funcion

	var json_prompt_usuario any

	Guardar_en_memoria(prompt, "user")

	opciones := Opciones{
		Num_ctx:     box.Ctx,
		Num_predict: -1,
		Temperature: box.Temperatura,
	}

	//herramientas
	scraping := Crear_tool_scraping()

	if chat {

		json_prompt_usuario = Mensaje_usuario_chat{

			Model:    box.Modelo,
			Messages: Memoria,
			Stream:   true,
			Options:  opciones,
			Tools:    []Herramienta{scraping},
		}

	} else { //generate, para el procesamiento de imagenes

		json_prompt_usuario = Mensaje_usuario_generate{

			Model:   box.Modelo,
			Prompt:  prompt,
			Images:  imagenes,
			Options: opciones,
		}

	}

	return struct_a_respuesta(json_prompt_usuario, endpoint, Content_type)

}

// esta funcion se ocupa del envio y recepcion de los mensajes
func Comunicacion(prompt_archivo, prompt string, box utilidades.Box_info, endpoint, content_type string, carga *menu.Carga, wg *sync.WaitGroup, chat bool, imagenes []string) error {

	p := utilidades.Estructura_prompt{
		// prompt archivo se formatea por
		Fecha:  utilidades.Fecha_hora(),
		Prompt: prompt,
	}

	prompt_total := prompt_archivo + p.Formatear_prompt()
	// ver de reorganizar esto (quiza crear una struct para encapsular algunas cosas)
	resp, prompterr := enviar_prompt(prompt_total, box, endpoint, content_type, chat, imagenes)

	defer carga.Detener(wg)

	if prompterr != nil {

		return prompterr
	}

	defer menu.Esperar_tecla()

	respuesta := recibir_prompt(resp, carga, wg, chat, prompt, box)

	procesar_respuesta(respuesta)

	if !box.Archivo {

		return nil
	}

	err := respuesta.Exportar()

	if err != nil {

		return err
	}

	return nil

}
