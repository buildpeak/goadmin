package api

import (
	"log/slog"
	"net/http"
	"os"
	"testing"
)

func TestNewRouter(t *testing.T) {
	t.Parallel()

	type args struct {
		validator *OpenAPIValidator
		handlers  *Handlers
		logger    *slog.Logger
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	tests := []struct {
		name string
		args args
		want http.Handler
	}{
		{
			name: "success",
			args: args{
				validator: &OpenAPIValidator{},
				handlers:  &Handlers{},
				logger:    logger,
			},
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := NewRouter(tt.args.validator, tt.args.handlers, tt.args.logger)
			if got == nil {
				t.Errorf("NewRouter() = nil")
			}
		})
	}
}
