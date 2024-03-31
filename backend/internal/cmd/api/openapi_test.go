package api

import (
	"bytes"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestNewOpenAPIValidator(t *testing.T) {
	t.Parallel()

	type args struct {
		oaiPath string
		logger  *slog.Logger
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				oaiPath: "../../../openapi.yaml",
				logger:  slog.New(slog.NewTextHandler(os.Stderr, nil)),
			},
		},
		{
			name: "empty path",
			args: args{
				oaiPath: "",
				logger:  slog.New(slog.NewTextHandler(os.Stderr, nil)),
			},
			wantErr: true,
		},
		{
			name: "error NewDocument",
			args: args{
				oaiPath: "testdata/empty_oas.yml",
				logger:  slog.New(slog.NewTextHandler(os.Stderr, nil)),
			},
			wantErr: true,
		},
		{
			name: "error NewValidator",
			args: args{
				oaiPath: "testdata/invalid_oas.yml",
				logger:  slog.New(slog.NewTextHandler(os.Stderr, nil)),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := NewOpenAPIValidator(tt.args.oaiPath, tt.args.logger)

			if (err != nil) != tt.wantErr {
				t.Errorf("NewOpenAPIValidator() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got == nil && !tt.wantErr {
				t.Errorf("NewOpenAPIValidator() = nil, want not nil")
			}
		})
	}
}

func newRequest(method, path string, body []byte) *http.Request {
	req := httptest.NewRequest(method, path, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	return req
}

func TestOpenAPIValidator_Middleware(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		next       http.Handler
		req        *http.Request
		wantStatus int
		wantBody   string
	}{
		{
			name: "/health",
			next: http.HandlerFunc(func(res http.ResponseWriter, _ *http.Request) {
				res.WriteHeader(http.StatusOK)
				res.Write([]byte("OK"))
			}),
			req:        httptest.NewRequest(http.MethodGet, "/health", nil),
			wantStatus: http.StatusOK,
			wantBody:   "OK",
		},
		{
			name: "error /auth/login",
			next: http.HandlerFunc(func(res http.ResponseWriter, _ *http.Request) {
				res.WriteHeader(http.StatusOK)
				res.Write([]byte("OK"))
			}),
			req:        httptest.NewRequest(http.MethodPost, "/auth/login", nil),
			wantStatus: http.StatusUnprocessableEntity,
			wantBody:   `{"type":"/errors/validation-error","title":"Validation Error","status":422,"detail":"Error: POST operation request content type '' does not exist, Reason: The content type '' of the POST request submitted has not been defined, it's an unknown type, Line: 30, Column: 9","instance":"/auth/login","errors":[]}` + "\n",
		},
		{
			name: "error /auth/login 2",
			next: http.HandlerFunc(func(res http.ResponseWriter, _ *http.Request) {
				res.WriteHeader(http.StatusOK)
				res.Write([]byte("OK"))
			}),
			req:        newRequest(http.MethodPost, "/auth/login", []byte(`{"username": "test"}`)),
			wantStatus: http.StatusUnprocessableEntity,
			wantBody:   `{"type":"/errors/validation-error","title":"Validation Error","status":422,"detail":"Error: POST request body for '/auth/login' failed to validate schema, Reason: The request body is defined as an object. However, it does not meet the schema requirements of the specification, Validation Errors: [Reason: missing properties: 'password', Location: /required], Line: 33, Column: 15","instance":"/auth/login","errors":[{"detail":"missing properties: 'password'","pointer":"/required"}]}` + "\n",
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			logger := slog.New(slog.NewTextHandler(os.Stderr, nil))

			v, err := NewOpenAPIValidator("../../../openapi.yaml", logger)
			if err != nil {
				t.Errorf("NewOpenAPIValidator() error = %v", err)
				return
			}

			handler := v.Middleware(tt.next)

			res := httptest.NewRecorder()

			handler.ServeHTTP(res, tt.req)

			if got := res.Code; got != tt.wantStatus {
				t.Errorf("Middleware() Status = %v, want %v", got, tt.wantStatus)
			}

			if got := res.Body.String(); got != tt.wantBody {
				t.Errorf("Middleware() Body = %v, want %v", got, tt.wantBody)
			}
		})
	}
}
