package auth

import (
	"net/http"
	"reflect"
	"testing"
)

func TestHandler_Authenticator(t *testing.T) {
	t.Parallel()

	type fields struct {
		authService Service
	}

	tests := []struct {
		name   string
		fields fields
		want   func(http.Handler) http.Handler
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			h := &Handler{
				authService: tt.fields.authService,
			}

			if got := h.Authenticator(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Handler.Authenticator() = %v, want %v", got, tt.want)
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
