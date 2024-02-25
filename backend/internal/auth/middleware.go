package auth

import (
	"context"
	"net/http"
	"strings"
)

func (h *Handler) Authenticator() func(http.Handler) http.Handler {
	// return middleware
	return func(next http.Handler) http.Handler {
		// return HandlerFunc
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenString := FindToken(r)
			if tokenString == "" {
				// h.logger.Error(fmt.Sprintf("Cannot find token"))
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			ctx := r.Context()

			user, err := h.authService.VerifyToken(ctx, tokenString)
			if err != nil {
				// h.logger.Error(fmt.Sprintf("Authenticate err: %v", err))
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// create new context with user value
			newCtx := context.WithValue(ctx, "user", *user)

			next.ServeHTTP(w, r.WithContext(newCtx))
		})
	}
}

func FindToken(r *http.Request) string {
	for _, f := range []func(*http.Request) string{
		TokenFromQuery,
		TokenFromHeader,
		TokenFromCookie,
	} {
		if token := f(r); token != "" {
			return token
		}
	}

	return ""
}

// TokenFromHeader tries to retreive the token string from the
// "Authorization" reqeust header: "Authorization: BEARER T".
func TokenFromHeader(r *http.Request) string {
	// Get token from authorization header.
	bearer := r.Header.Get("Authorization")
	if len(bearer) > 7 && strings.ToUpper(bearer[0:6]) == "BEARER" {
		return bearer[7:]
	}
	return ""
}

// TokenFromQuery tries to retreive the token string from the "jwt" URI
// query parameter.
//
// To use it, build our own middleware handler, such as:
//
//	func Verifier(ja *JWTAuth) func(http.Handler) http.Handler {
//		return func(next http.Handler) http.Handler {
//			return Verify(ja, TokenFromQuery, TokenFromHeader, TokenFromCookie)(next)
//		}
//	}
func TokenFromQuery(r *http.Request) string {
	// Get token from query param named "jwt".
	return r.URL.Query().Get("jwt")
}

// TokenFromCookie tries to retreive the token string from a cookie named
// "jwt".
func TokenFromCookie(r *http.Request) string {
	cookie, err := r.Cookie("jwt")
	if err != nil {
		return ""
	}
	return cookie.Value
}
