package auth

import (
	"context"
	"net/http"
	"strings"
)

type contextKey string

const userKey = contextKey("user")

func (h *Handler) Authenticator() func(http.Handler) http.Handler {
	// returns middleware
	return func(next http.Handler) http.Handler {
		// returns HandlerFunc
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			tokenString := FindToken(req)
			if tokenString == "" {
				http.Error(res, "Unauthorized", http.StatusUnauthorized)

				return
			}

			ctx := req.Context()

			user, err := h.authService.VerifyToken(ctx, tokenString)
			if err != nil {
				http.Error(res, "Unauthorized", http.StatusUnauthorized)

				return
			}

			// create new context with user value
			newCtx := context.WithValue(ctx, userKey, *user)

			next.ServeHTTP(res, req.WithContext(newCtx))
		})
	}
}

func FindToken(req *http.Request) string {
	for _, f := range []func(*http.Request) string{
		TokenFromQuery,
		TokenFromHeader,
		TokenFromCookie,
	} {
		if token := f(req); token != "" {
			return token
		}
	}

	return ""
}

// TokenFromHeader tries to retrieve the token string from the
// "Authorization" reqeust header: "Authorization: BEARER T".
func TokenFromHeader(req *http.Request) string {
	// Get token from authorization header.
	bearer := req.Header.Get("Authorization")
	if len(bearer) > 7 && strings.ToUpper(bearer[0:6]) == "BEARER" {
		return bearer[7:]
	}

	return ""
}

// TokenFromQuery tries to retrieve the token string from the "jwt" URI
// query parameter.
//
// To use it, build our own middleware handler, such as:
//
//	func Verifier(ja *JWTAuth) func(http.Handler) http.Handler {
//		return func(next http.Handler) http.Handler {
//			return Verify(ja, TokenFromQuery, TokenFromHeader, TokenFromCookie)(next)
//		}
//	}
func TokenFromQuery(req *http.Request) string {
	// Get token from query param named "jwt".
	return req.URL.Query().Get("jwt")
}

// TokenFromCookie tries to retrieve the token string from a cookie named
// "jwt".
func TokenFromCookie(req *http.Request) string {
	cookie, err := req.Cookie("jwt")
	if err != nil {
		return ""
	}

	return cookie.Value
}
