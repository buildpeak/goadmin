package auth

import (
	"bytes"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"goadmin-backend/internal/domain"
	"goadmin-backend/internal/platform/httpjson"
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
				Handler: httpjson.Handler{
					Logger: logging.NewLogger(),
				},
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

//nolint:unparam // unit test
func newRequest(method, urlPath string, body any) *http.Request {
	if body == nil {
		return httptest.NewRequest(method, urlPath, nil)
	}

	//nolint:errchkjson // unit test ignore error
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
		{
			name: "Fail",
			fields: fields{
				authService: &ServiceMock{err: ErrInvalidCredentials},
				logger:      logging.NewLogger(),
			},
			args: args{
				res: httptest.NewRecorder(),
				req: newRequest(http.MethodPost, "/login", domain.Credentials{}),
			},
			want: want{
				code: http.StatusUnauthorized,
				body: `{"type":"/errors/unauthorized","title":"Unauthorized","status":401,"detail":"You are not authorized to perform this action","instance":"/login"}` + "\n",
			},
		},
		{
			name: "Fail Internal Server Error",
			fields: fields{
				authService: &ServiceMock{err: errors.New("db error")},
				logger:      logging.NewLogger(),
			},
			args: args{
				res: httptest.NewRecorder(),
				req: newRequest(http.MethodPost, "/login", domain.Credentials{}),
			},
			want: want{
				code: http.StatusInternalServerError,
				body: `{"type":"/errors/internal-server-error","title":"Internal Server Error","status":500,"detail":"An internal server error occurred","instance":"/login"}` + "\n",
			},
		},
		{
			name: "Fail Decode",
			fields: fields{
				authService: &ServiceMock{},
				logger:      logging.NewLogger(),
			},
			args: args{
				res: httptest.NewRecorder(),
				req: httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader([]byte("invalid"))),
			},
			want: want{
				code: http.StatusBadRequest,
				body: `{"type":"/errors/bad-request","title":"Bad Request","status":400,"detail":"The request was invalid or cannot be served","instance":"/login"}` + "\n",
			},
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			h := &Handler{
				authService: tt.fields.authService,
				Handler: httpjson.Handler{
					Logger: tt.fields.logger,
				},
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
		{
			name: "Fail Internal Server Error",
			fields: fields{
				authService: &ServiceMock{err: errors.New("db error")},
				logger:      logging.NewLogger(),
			},
			args: args{
				res: httptest.NewRecorder(),
				req: newRequest(http.MethodPost, "/register", domain.User{}),
			},
			want: want{
				code: http.StatusInternalServerError,
				body: `{"type":"/errors/internal-server-error","title":"Internal Server Error","status":500,"detail":"An internal server error occurred","instance":"/register"}` + "\n",
			},
		},
		{
			name: "Fail Decode",
			fields: fields{
				authService: &ServiceMock{},
				logger:      logging.NewLogger(),
			},
			args: args{
				res: httptest.NewRecorder(),
				req: httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader([]byte("invalid"))),
			},
			want: want{
				code: http.StatusBadRequest,
				body: `{"type":"/errors/bad-request","title":"Bad Request","status":400,"detail":"The request was invalid or cannot be served","instance":"/register"}` + "\n",
			},
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			h := &Handler{
				authService: tt.fields.authService,
				Handler: httpjson.Handler{
					Logger: tt.fields.logger,
				},
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
		httpjson.Handler
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
				Handler: httpjson.Handler{
					Logger: logging.NewLogger(),
				},
			},
			args: args{
				res: httptest.NewRecorder(),
				req: newRequest(http.MethodPost, "/signin-with-google", GoogleIDTokenVerifyRequest{
					IDToken: "good_token",
				}),
			},
		},
		{
			name: "Fail Invalid Id Token",
			fields: fields{
				authService: &ServiceMock{err: ErrInvalidIDToken},
				Handler: httpjson.Handler{
					Logger: logging.NewLogger(),
				},
			},
			args: args{
				res: httptest.NewRecorder(),
				req: newRequest(http.MethodPost, "/signin-with-google", GoogleIDTokenVerifyRequest{
					IDToken: "bad_token",
				}),
			},
		},
		{
			name: "Fail Decode",
			fields: fields{
				authService: &ServiceMock{},
				Handler: httpjson.Handler{
					Logger: logging.NewLogger(),
				},
			},
			args: args{
				res: httptest.NewRecorder(),
				req: httptest.NewRequest(http.MethodPost, "/signin-with-google", bytes.NewReader([]byte("invalid"))),
			},
		},
		{
			name: "Fail Internal Server Error",
			fields: fields{
				authService: &ServiceMock{err: errors.New("db error")},
				Handler: httpjson.Handler{
					Logger: logging.NewLogger(),
				},
			},
			args: args{
				res: httptest.NewRecorder(),
				req: newRequest(http.MethodPost, "/signin-with-google", GoogleIDTokenVerifyRequest{
					IDToken: "good_token",
				}),
			},
		},
		{
			name: "Fail User Not Found",
			fields: fields{
				authService: &ServiceMock{err: domain.NewResourceNotFoundError("user", "id=1")},
				Handler: httpjson.Handler{
					Logger: logging.NewLogger(),
				},
			},
			args: args{
				res: httptest.NewRecorder(),
				req: newRequest(http.MethodPost, "/signin-with-google", GoogleIDTokenVerifyRequest{
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
				Handler:     tt.fields.Handler,
			}

			h.SignInWithGoogle(tt.args.res, tt.args.req)
		})
	}
}
