package httpjson

import (
	"bytes"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
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
		{
			name: "Success",
			fields: fields{
				Logger: slog.New(slog.NewTextHandler(os.Stderr, nil)),
			},
			args: args{
				res: httptest.NewRecorder(),
				req: httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte(`{"name":"John"}`))),
				dataPtr: &struct {
					Name string `json:"name"`
				}{},
			},
		},
		{
			name: "Success with nil data",
			fields: fields{
				Logger: slog.New(slog.NewTextHandler(os.Stderr, nil)),
			},
			args: args{
				res:     httptest.NewRecorder(),
				req:     httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte(`{"name":"John"}`))),
				dataPtr: nil,
			},
			wantErr: true,
		},
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

	type Person struct {
		Name   string
		Friend *Person // Circular reference
	}

	circular := Person{Name: "John"}
	circular.Friend = &circular

	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Success",
			fields: fields{
				Logger: slog.New(slog.NewTextHandler(os.Stderr, nil)),
			},
			args: args{
				res:    httptest.NewRecorder(),
				data:   struct{ Name string }{"John"},
				status: http.StatusOK,
			},
		},
		{
			name: "Success with nil data",
			fields: fields{
				Logger: slog.New(slog.NewTextHandler(os.Stderr, nil)),
			},
			args: args{
				res:    httptest.NewRecorder(),
				data:   nil,
				status: http.StatusOK,
			},
		},
		{
			name: "Fail with circular reference",
			fields: fields{
				Logger: slog.New(slog.NewTextHandler(os.Stderr, nil)),
			},
			args: args{
				res:    httptest.NewRecorder(),
				data:   circular,
				status: http.StatusOK,
			},
		},
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
