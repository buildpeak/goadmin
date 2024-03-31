package api

import (
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
)

func TestNewHealthHandler(t *testing.T) {
	t.Parallel()

	type args struct {
		logger *slog.Logger
	}

	tests := []struct {
		name string
		args args
		want *HealthHandler
	}{
		{
			name: "success",
			args: args{
				logger: slog.New(slog.NewTextHandler(os.Stderr, nil)),
			},
			want: &HealthHandler{
				Logger: slog.New(slog.NewTextHandler(os.Stderr, nil)),
			},
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := NewHealthHandler(tt.args.logger); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewHealthHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

type errRecorder struct {
	httptest.ResponseRecorder
}

func (r *errRecorder) Write(_ []byte) (int, error) {
	return 0, errors.New("error")
}

func TestHealthHandler_healthCheck(t *testing.T) {
	t.Parallel()

	type fields struct {
		Logger *slog.Logger
	}

	type args struct {
		res http.ResponseWriter
		in1 *http.Request
	}

	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "success",
			fields: fields{
				Logger: slog.New(slog.NewTextHandler(os.Stderr, nil)),
			},
			args: args{
				res: httptest.NewRecorder(),
				in1: httptest.NewRequest(http.MethodGet, "/health", nil),
			},
		},
		{
			name: "error",
			fields: fields{
				Logger: slog.New(slog.NewTextHandler(os.Stderr, nil)),
			},
			args: args{
				res: &errRecorder{},
				in1: httptest.NewRequest(http.MethodPost, "/health", nil),
			},
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			h := &HealthHandler{
				Logger: tt.fields.Logger,
			}

			h.healthCheck(tt.args.res, tt.args.in1)
		})
	}
}
