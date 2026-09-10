package utilidades_test

import (
	"LLM-Chat/utilidades"
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

func TestString_vacio(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		str  string
		want bool
	}{

		{str: "", want: true},
		{str: "    ", want: true},
		{str: "   test", want: false},
		{str: "  \n", want: true},
		{str: "\t \n", want: true},
		{str: "test", want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := utilidades.String_vacio(tt.str)

			if got != tt.want {
				t.Errorf("String_vacio() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSlices_Iguales(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		sl1  [][]string
		sl2  [][]string
		want bool
	}{

		{name: "caso normal", sl1: [][]string{{"1", "2", "3"}, {"4"}}, sl2: [][]string{{"1", "2", "3"}, {"4"}}, want: true},
		{name: "caso vacia", sl1: [][]string{}, sl2: [][]string{}, want: true},
		{name: "caso falso", sl1: [][]string{{"1", "2", "3"}, {"4"}}, sl2: [][]string{{"1", "2"}}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := utilidades.Slices_Iguales(tt.sl1, tt.sl2)

			if got != tt.want {
				t.Errorf("Slices_Iguales() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSeparar_lista(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		lista []string
		want  [][]string
	}{

		{name: "caso normal", lista: []string{"1", "2", "3", "4", "5", "6"}, want: [][]string{{"1", "2", "3", "4", "5"}, {"6"}}},
		{name: "vacia", lista: []string{}, want: [][]string{{}}},
		{name: "caso iguales", lista: []string{"1", "2", "3"}, want: [][]string{{"1", "2", "3"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := utilidades.Separar_lista(tt.lista)

			if !utilidades.Slices_Iguales(got, tt.want) {

				t.Errorf("Separar_lista() = %v, want %v", got, tt.want)
			}
		})
	}
}
