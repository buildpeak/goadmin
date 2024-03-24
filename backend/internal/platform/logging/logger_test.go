package logging

import (
	"log/slog"
	"reflect"
	"testing"
)

func TestWithLevel(t *testing.T) {
	t.Parallel()

	type args struct {
		level string
	}

	tests := []struct {
		name string
		args args
		want Option
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := WithLevel(tt.args.level); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithLevel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithPretty(t *testing.T) {
	t.Parallel()

	type args struct {
		pretty bool
	}

	tests := []struct {
		name string
		args args
		want Option
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := WithPretty(tt.args.pretty); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithPretty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewLogger(t *testing.T) {
	t.Parallel()

	type args struct {
		opts []Option
	}

	tests := []struct {
		name string
		args args
		want *slog.Logger
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := NewLogger(tt.args.opts...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewLogger() = %v, want %v", got, tt.want)
			}
		})
	}
}
