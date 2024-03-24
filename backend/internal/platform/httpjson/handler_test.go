package httpjson

import (
	"log/slog"
	"net/http"
	"testing"
)

func TestHandler_ParseJSON(t *testing.T) {
	t.Parallel()

	type fields struct {
		Logger *slog.Logger
	}

	type args struct {
		res     http.ResponseWriter
		req     *http.Request
		dataPtr any
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			h := &Handler{
				Logger: tt.fields.Logger,
			}

			if err := h.ParseJSON(tt.args.res, tt.args.req, tt.args.dataPtr); (err != nil) != tt.wantErr {
				t.Errorf("Handler.ParseJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHandler_RespondJSON(t *testing.T) {
	t.Parallel()

	type fields struct {
		Logger *slog.Logger
	}

	type args struct {
		res    http.ResponseWriter
		data   any
		status int
	}

	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			h := &Handler{
				Logger: tt.fields.Logger,
			}

			h.RespondJSON(tt.args.res, tt.args.data, tt.args.status)
		})
	}
}
