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
