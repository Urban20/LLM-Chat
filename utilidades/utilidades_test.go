package utilidades_test

import (
	"LLM-Chat/utilidades"
	_ "embed"
	"slices"
	"testing"
)

func TestEliminar_repetidos(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		elementos []string
		want      []string
	}{
		{elementos: []string{"1", "1", "2", "2", "3", "3"}, want: []string{"1", "2", "3"}},
		{elementos: []string{}, want: []string{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := utilidades.Eliminar_repetidos(tt.elementos)

			if slices.Compare(got, tt.want) != 0 {

				t.Errorf("Eliminar_repetidos() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVacio(t *testing.T) {
	tests := []struct {
		name  string
		linea string
		want  bool
	}{

		{linea: "   ", want: true},
		{linea: "  test ", want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := utilidades.Vacio(tt.linea)

			if got != tt.want {
				t.Errorf("Vacio() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEliminar_espacios(t *testing.T) {
	tests := []struct {
		name string

		texto string
		want  string
	}{

		{name: "caso simple", texto: "esto es un\n  \ntest", want: "esto es un\ntest"},
		{texto: "esto es un test", want: "esto es un test"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := utilidades.Eliminar_espacios(tt.texto)

			if got != tt.want {
				t.Errorf("Eliminar_espacios() = %v, want %v", got, tt.want)
			}
		})
	}
}
