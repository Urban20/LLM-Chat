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
		{name: "basico", elementos: []string{"1", "1", "2", "2", "3", "3"}, want: []string{"1", "2", "3"}},
		{name: "slice vacio", elementos: []string{}, want: []string{}},
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

		{name: "basico", str: "", want: true},
		{name: "espacio", str: "    ", want: true},
		{name: "texto con espacio", str: "   test", want: false},
		{name: "salto de linea", str: "  \n", want: true},
		{name: "salto de linea con tab", str: "\t \n", want: true},
		{name: "texto normal", str: "test", want: false},
		{name: "retorno tab", str: "\t\r", want: true},
		{name: "retorno tab con string", str: "\t\rtest", want: false},
		{name: "doble salto de linea", str: "\n\n", want: true},
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
		{name: "caso falso mismo largo", sl1: [][]string{{"1", "2", "3"}}, sl2: [][]string{{"1", "2", "4"}}, want: false},
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
		lista  []string
		limite int
		want   [][]string
	}{

		{name: "caso normal", lista: []string{"1", "2", "3", "4", "5", "6"}, limite: 5, want: [][]string{{"1", "2", "3", "4", "5"}, {"6"}}},
		{name: "vacia", lista: []string{}, limite: 5, want: [][]string{{}}},
		{name: "caso iguales", lista: []string{"1", "2", "3"}, limite: 5, want: [][]string{{"1", "2", "3"}}},
		{name: "caso invalido", lista: []string{"1", "2", "3"}, limite: -1, want: [][]string{}},
		{name: "caso iguales", lista: []string{"1", "2", "3"}, limite: 2, want: [][]string{{"1", "2"}, {"3"}}},
		{name: "caso cero", lista: []string{}, limite: 0, want: [][]string{}},
		{name: "limite cero", lista: []string{"1", "2", "3"}, limite: 0, want: [][]string{}},
		{name: "varias paginas", lista: []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"}, limite: 2, want: [][]string{{"1", "2"}, {"3", "4"}, {"5", "6"}, {"7", "8"}, {"9", "10"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := utilidades.Separar_lista(tt.lista, tt.limite)

			if !utilidades.Slices_Iguales(got, tt.want) {

				t.Errorf("Separar_lista() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMax(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		l    []string
		want int
	}{

		{name: "vacio", l: []string{}, want: 0},
		{name: "caso normal", l: []string{"test", "otorrinolaringologo", "testeando"}, want: len("otorrinolaringologo")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := utilidades.Max(tt.l)

			if got != tt.want {
				t.Errorf("Max() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLimpiar_listas(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		l    []string
		want []string
	}{

		{l: []string{"", "test", "\n", " ", "   "}, want: []string{"test"}},
		{l: []string{}, want: []string{}},
		{l: []string{"test", "testing"}, want: []string{"test", "testing"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := utilidades.Limpiar_listas(tt.l)

			if !slices.Equal(got, tt.want) {
				t.Errorf("Limpiar_listas() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSlice_vacio(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		sl   []string
		want bool
	}{

		{name: "basico", sl: []string{}, want: true},
		{name: "strings invalidos", sl: []string{" ", "", "\t", "\n"}, want: true},
		{name: "en blanco", sl: []string{"   "}, want: true},
		{name: "caso falso", sl: []string{"test", " ", ""}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := utilidades.Slice_vacio(tt.sl)

			if got != tt.want {
				t.Errorf("Slice_vacio() = %v, want %v", got, tt.want)
			}
		})
	}
}
