package auth

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"goadmin-backend/internal/domain"
)

func newRequestWithToken(token, whereIsToken string) *http.Request {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	if whereIsToken == "header" || whereIsToken == "" {
		req.Header = http.Header{
			"Authorization": []string{"Bearer " + token},
		}
	} else if whereIsToken == "query" {
		req.URL.RawQuery = "jwt=" + token
	} else if whereIsToken == "cookie" {
		req.AddCookie(&http.Cookie{
			Name:  "jwt",
			Value: token,
		})
	}

	return req
}

func TestHandler_Authenticator(t *testing.T) {
	t.Parallel()

	type fields struct {
		authService Service
	}

	tests := []struct {
		name       string
		fields     fields
		req        *http.Request
		wantStatus int
		want       string
	}{
		{
			name: "Test Authenticator() with valid token in header",
			fields: fields{
				authService: &ServiceMock{},
			},
			req:        newRequestWithToken("good_token", "header"),
			wantStatus: http.StatusOK,
			want:       "OK\n",
		},
		{
			name: "Test Authenticator() with valid token in query",
			fields: fields{
				authService: &ServiceMock{},
			},
			req:        newRequestWithToken("good_token", "query"),
			wantStatus: http.StatusOK,
			want:       "OK\n",
		},
		{
			name: "Test Authenticator() with valid token in cookie",
			fields: fields{
				authService: &ServiceMock{},
			},
			req:        newRequestWithToken("good_token", "cookie"),
			wantStatus: http.StatusOK,
			want:       "OK\n",
		},
		{
			name: "Test Authenticator() with invalid token in header",
			fields: fields{
				authService: &ServiceMock{
					hasError: true,
				},
			},
			req:        newRequestWithToken("invalid_token", "header"),
			wantStatus: http.StatusUnauthorized,
			want:       "Unauthorized\n",
		},
		{
			name: "Test Authenticator() with no token",
			fields: fields{
				authService: &ServiceMock{},
			},
			req:        httptest.NewRequest(http.MethodGet, "/", nil),
			wantStatus: http.StatusUnauthorized,
			want:       "Unauthorized\n",
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			h := &Handler{
				authService: tt.fields.authService,
			}

			handler := h.Authenticator()(
				http.HandlerFunc(
					func(res http.ResponseWriter, _ *http.Request) {
						// Do nothing
						t.Log("Auth passed")
						res.Write([]byte("OK\n"))
					},
				),
			)

			res := httptest.NewRecorder()

			handler.ServeHTTP(res, tt.req)

			if got := res.Code; got != tt.wantStatus {
				t.Errorf(
					"Handler.Authenticator()(...) Status = %v, want %v",
					got,
					tt.wantStatus,
				)
			}

			if got := res.Body.String(); got != tt.want {
				t.Errorf(
					"Handler.Authenticator()(...) Body = %v, want %v",
					got,
					tt.want,
				)
			}
		})
	}
}

func TestFindToken(t *testing.T) {
	t.Parallel()

	type args struct {
		req *http.Request
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Success",
			args: args{
				req: newRequestWithToken("good_token", "header"),
			},
			want: "good_token",
		},
		{
			name: "Fail",
			args: args{
				req: httptest.NewRequest(http.MethodGet, "/", nil),
			},
			want: "",
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := FindToken(tt.args.req); got != tt.want {
				t.Errorf("FindToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTokenFromHeader(t *testing.T) {
	t.Parallel()

	type args struct {
		req *http.Request
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Success",
			args: args{
				req: newRequestWithToken("good_token", "header"),
			},
			want: "good_token",
		},
		{
			name: "Fail",
			args: args{
				req: httptest.NewRequest(http.MethodGet, "/", nil),
			},
			want: "",
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := TokenFromHeader(tt.args.req); got != tt.want {
				t.Errorf("TokenFromHeader() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTokenFromQuery(t *testing.T) {
	t.Parallel()

	type args struct {
		req *http.Request
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Success",
			args: args{
				req: newRequestWithToken("good_token", "query"),
			},
			want: "good_token",
		},
		{
			name: "Fail",
			args: args{
				req: httptest.NewRequest(http.MethodGet, "/", nil),
			},
			want: "",
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := TokenFromQuery(tt.args.req); got != tt.want {
				t.Errorf("TokenFromQuery() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTokenFromCookie(t *testing.T) {
	t.Parallel()

	type args struct {
		req *http.Request
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Success",
			args: args{
				req: newRequestWithToken("good_token", "cookie"),
			},
			want: "good_token",
		},
		{
			name: "Fail",
			args: args{
				req: httptest.NewRequest(http.MethodGet, "/", nil),
			},
			want: "",
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := TokenFromCookie(tt.args.req); got != tt.want {
				t.Errorf("TokenFromCookie() = %v, want %v", got, tt.want)
			}
		})
	}
}

var _ Service = &ServiceMock{}

type ServiceMock struct {
	hasError bool
	err      error
}

func (s *ServiceMock) Login(
	_ context.Context,
	_ domain.Credentials,
) (*domain.JWTToken, error) {
	if s.hasError {
		if s.err != nil {
			return nil, s.err
		}

		return nil, ErrInvalidCredentials
	}

	return &domain.JWTToken{
		AccessToken: "good_token",
	}, nil
}

func (s *ServiceMock) VerifyToken(
	_ context.Context,
	_ string,
) (*domain.User, error) {
	if s.hasError {
		return nil, ErrInvalidToken
	}

	return &domain.User{
		ID: "1",
	}, nil
}

func (s *ServiceMock) Register(
	_ context.Context,
	_ *domain.User,
) (*domain.User, error) {
	if s.hasError {
		return nil, ErrInvalidCredentials
	}

	return &domain.User{
		ID: "1",
	}, nil
}

func (s *ServiceMock) ValidateGoogleIDToken(
	_ context.Context,
	_ string,
	_ string,
) (*domain.JWTToken, error) {
	if s.hasError {
		return nil, ErrInvalidToken
	}

	return &domain.JWTToken{
		AccessToken: "good_token",
	}, nil
}
