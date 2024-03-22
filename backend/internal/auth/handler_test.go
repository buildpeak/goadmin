package auth

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"goadmin-backend/internal/domain"
	"goadmin-backend/internal/platform/logging"
)

func TestNewHandler(t *testing.T) {
	t.Parallel()

	type args struct {
		authService Service
		logger      *slog.Logger
	}

	tests := []struct {
		name string
		args args
		want *Handler
	}{
		{
			name: "Success",
			args: args{
				authService: &ServiceMock{},
				logger:      logging.NewLogger(),
			},
			want: &Handler{
				authService: &ServiceMock{},
				logger:      logging.NewLogger(),
			},
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := NewHandler(tt.args.authService, tt.args.logger); !reflect.DeepEqual(
				got,
				tt.want,
			) {
				t.Errorf("NewHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func newRequest(method, urlPath string, body any) *http.Request {
	if body == nil {
		return httptest.NewRequest(method, urlPath, nil)
	}

	b, _ := json.Marshal(body)

	req := httptest.NewRequest(method, urlPath, bytes.NewReader(b))

	req.Header.Set("Content-Type", "application/json")

	return req
}

func TestHandler_Login(t *testing.T) {
	t.Parallel()

	type fields struct {
		authService Service
		logger      *slog.Logger
	}

	type args struct {
		res http.ResponseWriter
		req *http.Request
	}

	type want struct {
		code int
		body string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    want
		wantErr bool
	}{
		{
			name: "Success",
			fields: fields{
				authService: &ServiceMock{},
				logger:      logging.NewLogger(),
			},
			args: args{
				res: httptest.NewRecorder(),
				req: newRequest(http.MethodPost, "/login", domain.Credentials{
					Username: "username",
					Password: "password",
				}),
			},
			want: want{
				code: http.StatusOK,
				body: `{"access_token":"good_token","refresh_token":""}` + "\n",
			},
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			h := &Handler{
				authService: tt.fields.authService,
				logger:      tt.fields.logger,
			}

			h.Login(tt.args.res, tt.args.req)

			rr, ok := tt.args.res.(*httptest.ResponseRecorder)

			if ok && rr.Code != tt.want.code {
				t.Errorf(
					"Handler.Login() = %v, want %v",
					rr.Code,
					tt.want.code,
				)
			}

			if ok && rr.Body.String() != tt.want.body {
				t.Errorf(
					"Handler.Login() = %v, want %v",
					rr.Body.String(),
					tt.want.body,
				)
			}
		})
	}
}

func TestHandler_Register(t *testing.T) {
	t.Parallel()

	type fields struct {
		authService Service
		logger      *slog.Logger
	}

	type args struct {
		res http.ResponseWriter
		req *http.Request
	}

	type want struct {
		code int
		body string
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "Success",
			fields: fields{
				authService: &ServiceMock{},
				logger:      logging.NewLogger(),
			},
			args: args{
				res: httptest.NewRecorder(),
				req: newRequest(http.MethodPost, "/register", domain.User{
					Username: "username",
					Password: "password",
				}),
			},
			want: want{
				code: http.StatusCreated,
				body: `{"id":"1","username":"","first_name":"","last_name":"","email":"","active":false,"deleted_at":null}` + "\n",
			},
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			h := &Handler{
				authService: tt.fields.authService,
				logger:      tt.fields.logger,
			}

			h.Register(tt.args.res, tt.args.req)

			rr, ok := tt.args.res.(*httptest.ResponseRecorder)

			if ok && rr.Code != tt.want.code {
				t.Errorf(
					"Handler.Register() = %v, want %v",
					rr.Code,
					tt.want.code,
				)
			}

			if ok && rr.Body.String() != tt.want.body {
				t.Errorf(
					"Handler.Register() = %v, want %v",
					rr.Body.String(),
					tt.want.body,
				)
			}
		})
	}
}

func TestHandler_SignInWithGoogle(t *testing.T) {
	t.Parallel()

	type fields struct {
		authService Service
		logger      *slog.Logger
	}

	type args struct {
		res http.ResponseWriter
		req *http.Request
	}

	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Success",
			fields: fields{
				authService: &ServiceMock{},
				logger:      logging.NewLogger(),
			},
			args: args{
				res: httptest.NewRecorder(),
				req: newRequest(http.MethodPost, "/google", GoogleIDTokenVerifyRequest{
					IDToken: "good_token",
				}),
			},
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			h := &Handler{
				authService: tt.fields.authService,
				logger:      tt.fields.logger,
			}

			h.SignInWithGoogle(tt.args.res, tt.args.req)
		})
	}
}
