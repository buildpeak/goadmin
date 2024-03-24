package slices

import (
	"reflect"
	"testing"
)

func TestMap(t *testing.T) {
	t.Parallel()

	type args struct {
		src []T
		fnc func(T) U
	}

	tests := []struct {
		name string
		args args
		want []U
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := Map(tt.args.src, tt.args.fnc); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Map() = %v, want %v", got, tt.want)
			}
		})
	}
}
