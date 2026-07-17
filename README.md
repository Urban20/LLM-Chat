# LLM-Chat

Cliente de terminal en Go para conversar con modelos alojados localmente a través de [Ollama](https://ollama.com). Ofrece un menú navegable, indicador de carga durante la generación, renderizado de Markdown en la respuesta y manejo de memoria conversacional por sesión.

## Requisitos

- Go 1.25 o superior (definido en `go.mod`).
- [Ollama](https://ollama.com/download) instalado y accesible en el `PATH` del sistema.
- Al menos un modelo descargado en Ollama.
- Una terminal que soporte secuencias ANSI (en Windows el propio binario habilita el modo virtual terminal al iniciar).

## 1. Levantar el servidor de Ollama

El cliente no inicia Ollama por sí mismo: solo se conecta a un servidor que ya debe estar corriendo.

1. Instalar Ollama siguiendo la guía oficial para tu sistema operativo.
2. Iniciar el servidor:

   ```bash
   ollama serve
   ```

   Por defecto Ollama escucha en `127.0.0.1:11434`. Si ya lo instalaste como servicio del sistema, es posible que ya esté corriendo en segundo plano y este paso no sea necesario.

3. Descargar al menos un modelo (en otra terminal, con el servidor ya corriendo):

   ```bash
   ollama pull llama3
   ```

   Podés reemplazar `llama3` por cualquier modelo disponible en [ollama.com/search](https://ollama.com/search).

4. Verificar que el servidor responde:

   ```bash
   curl http://localhost:11434/api/tags
   ```

   Si devuelve un JSON con la lista de modelos, el servidor está listo para que el cliente se conecte.

## 2. Compilar el cliente

Clonar el repositorio y compilar con Go:

```bash
git clone https://github.com/Urban20/LLM-Chat.git
cd LLM-Chat
go build .
```

Esto genera un binario `llm-chat` (o `llm-chat.exe` en Windows) en el directorio actual.

## 3. Ejecutar el cliente

Ejecución básica, asumiendo que Ollama corre en `localhost:11434`:

```bash
./llm-chat
```

### Flags disponibles

| Flag       | Valor por defecto | Descripción                                              |
|------------|--------------------|-----------------------------------------------------------|
| `-host`    | `localhost`        | Host/URL donde escucha el endpoint de Ollama              |
| `-puerto`  | `11434`             | Puerto del endpoint de Ollama                             |
| `-ctx`     | `16000`             | Cantidad de contexto (tokens) que usará el LLM             |
| `-temp`    | `0.5`               | Temperatura del modelo (creatividad de las respuestas)     |

Ejemplo apuntando a un servidor remoto con más contexto y menor temperatura:

```bash
./llm-chat -host 192.168.1.50 -puerto 11434 -ctx 32000 -temp 0.2
```

## 4. Uso del cliente

Al iniciar, el programa hace lo siguiente automáticamente:

1. Verifica si Ollama está en las variables de entorno del sistema (advierte si no lo encuentra, pero continúa).
2. Verifica que el servidor responda en la URL configurada.
3. Lista los modelos ya descargados en el servidor. Si no hay ninguno, te va a pedir que corras `ollama pull (modelo)` y termina la ejecución.

### Navegación del menú

- Se despliega un menú con las opciones disponibles (modelos instalados, o acciones dentro de una conversación).
- Navegación con las flechas `↑` `↓`.
- `Enter` confirma la selección.

### Selección de modelo

Se muestra la lista de modelos disponibles más la opción `[Salir]`. Al elegir un modelo se muestra un panel con:

- Modelo seleccionado
- Host y puerto
- Límite de memoria configurado
- Sistema operativo
- Temperatura y contexto del LLM

### Dentro de una conversación

Una vez elegido el modelo, el menú ofrece tres opciones:

- **Ingresar prompt**: abre el modo de escritura. Se escribe el mensaje y se envía con `TAB` seguido de `ENTER`. La respuesta se muestra en streaming y luego se renderiza como Markdown.
- **Borrar contexto**: limpia la memoria de la conversación actual (el historial que se envía al modelo en cada request), sin salir del modelo seleccionado.
- **Volver**: sale de la conversación actual, borra la memoria y regresa al selector de modelos.


