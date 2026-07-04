package utilidades

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/pterm/pterm"
	"golang.org/x/term"
)

const (
	AMARILLO      = "\033[0;33m"
	RESET         = "\033[0m"
	CELESTE_CLARO = "\033[38;2;125;255;216m"
	GRIS_AZUL     = "\033[38;2;90;112;176m"
	BLANCO        = "\033[47m"
)

func separador() {

	x, _, _ := term.GetSize(int(os.Stdout.Fd()))
	fmt.Println(BLANCO)
	fmt.Println(strings.Repeat(" ", x))
	fmt.Println(RESET)

}

func Ollama_instalado() bool {

	/*
		con esta funcion miro si Ollama esta en las variables de entorno del sistema

	*/

	ollama, err := exec.LookPath("Ollama")

	return err == nil || ollama != ""

}

func Imprimir_markdown(txt string) error {

	render, _ := glamour.NewTermRenderer(glamour.WithStylesFromJSONBytes([]byte(estilos)))

	separador()
	md, err := render.Render(txt)

	if err != nil {
		return err
	}
	fmt.Print(md)
	separador()

	return nil
}

func Box(msgs ...string) {

	superficie := pterm.DefaultBox.WithHorizontalPadding(5).WithBottomPadding(1)

	superficie.Println(strings.Join(msgs, "\n"))

}

func Limpieza_rapida() {

	fmt.Print("\033[2J")
	fmt.Print("\033[H")

}

func Formato_string_box(cuerpo map[string]string) []string {

	var retorno []string

	for clave, valor := range cuerpo {

		elemento := CELESTE_CLARO + clave + RESET + " : " + valor

		retorno = append(retorno, elemento)

	}

	return retorno

}

var estilos = `{
  "document": {
    "block_prefix": "\n",
    "block_suffix": "\n",
    "color": "#CED2F5",
    "margin": 2
  },
  "block_quote": {
    "indent": 1,
    "indent_token": "│ "
  },
  "paragraph": {},
  "list": {
    "level_indent": 2
  },
  "heading": {
    "block_suffix": "\n",
    "color": "39",
    "bold": true
  },
  "h1": {
    "prefix": " ",
    "suffix": " ",
    "color": "228",
    "background_color": "#42C77A",
    "bold": true
  },
  "h2": {
    "prefix": "> ",
    "background_color": "#3AB56E"
  },
  "h3": {
    "prefix": ">>> ",
    "background_color": "#35A664"
  },
  "h4": {
    "prefix": ">>>> ",
    "background_color": "#2E8F56"
  },
  "h5": {
    "prefix": ">>>>> ",
    "background_color": "#277546"
  },
  "h6": {
    "prefix": ">>>>>> ",
    "background_color": "#087832",
    "color": "#0DDB58",
    "bold": false
  },
  "text": {},
  "strikethrough": {
    "crossed_out": true
  },
  "emph": {
    "italic": true
  },
  "strong": {
    "bold": true
  },
  "hr": {
    "color": "#0DDB58",
    "format": "\n--------\n"
  },
  "item": {
    "block_prefix": "► "
  },
  "enumeration": {
    "block_prefix": ". "
  },
  "task": {
    "ticked": "[✓] ",
    "unticked": "[ ] "
  },
  "link": {
    "color": "#C0A2F2",
    "underline": true
  },
  "link_text": {
    "color": "#5dce51",
    "bold": true
  },
  "image": {
    "color": "#C0A2F2",
    "underline": true
  },
  "image_text": {
    "color": "243",
    "format": "Image: {{.text}} →"
  },
  "code": {
    "prefix": " ",
    "suffix": " ",
    "color": "#C3BEF7",
    "background_color": "#555087"
  },
  "code_block": {
    "color": "244",
    "margin": 2,
    "chroma": {
      "text": {
        "color": "#ffffff"
      },
      "error": {
        "color": "#F1F1F1",
        "background_color": "#e04c4c"
      },
      "comment": {
        "color": "#534c68"
      },
      "comment_preproc": {
        "color": "#f46738"
      },
      "keyword": {
        "color": "#04FF00"
      },
      "keyword_reserved": {
        "color": "#FF5FD2"
      },
      "keyword_namespace": {
        "color": "#ff5fa7"
      },
      "keyword_type": {
        "color": "#6E6ED8"
      },
      "operator": {
        "color": "#efe680"
      },
      "punctuation": {
        "color": "#E8E8A8"
      },
      "name": {
        "color": "#C4C4C4"
      },
      "name_builtin": {
        "color": "#f260e3"
      },
      "name_tag": {
        "color": "#B083EA"
      },
      "name_attribute": {
        "color": "#7A7AE6"
      },
      "name_class": {
        "color": "#F1F1F1",
        "underline": true,
        "bold": true
      },
      "name_constant": {},
      "name_decorator": {
        "color": "#7bdc63"
      },
      "name_exception": {},
      "name_function": {
        "color": "#00d773"
      },
      "name_other": {},
      "literal": {},
      "literal_number": {
        "color": "#45ecaf"
      },
      "literal_date": {},
      "literal_string": {
        "color": "#9564E3"
      },
      "literal_string_escape": {
        "color": "#AFFFD7"
      },
      "generic_deleted": {
        "color": "#9bf75a"
      },
      "generic_emph": {
        "italic": true
      },
      "generic_inserted": {
        "color": "#00D787"
      },
      "generic_strong": {
        "bold": true
      },
      "generic_subheading": {
        "color": "#777777"
      },
      "background": {
        "background_color": "#515173"
      }
    }
  },
  "table": {},
  "definition_list": {},
  "definition_term": {},
  "definition_description": {
    "block_prefix": "\n🠶 "
  },
  "html_block": {},
  "html_span": {}
}`
