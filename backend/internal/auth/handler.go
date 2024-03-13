package auth

import (
	"encoding/json"
	"net/http"

	"goadmin-backend/internal/domain"
	"goadmin-backend/internal/platform/httperr"
)

type Handler struct {
	authService Service
}

func NewHandler(authService Service) *Handler {
	return &Handler{authService: authService}
}

// Login handler signs in a user.
func (h *Handler) Login(res http.ResponseWriter, req *http.Request) {
	var credentials domain.Credentials

	if err := json.NewDecoder(req.Body).Decode(&credentials); err != nil {
		http.Error(res, "invalid request", http.StatusBadRequest)

		return
	}

	token, err := h.authService.Login(req.Context(), credentials)
	if err != nil {
		httperr.JSONError(res, err, http.StatusInternalServerError)

		return
	}

	if _, err := res.Write([]byte(token)); err != nil {
		httperr.JSONError(res, err, http.StatusInternalServerError)

		return
	}
}

// register handler signs up a user
func (h *Handler) Register(res http.ResponseWriter, req *http.Request) {
	var user domain.User

	if err := json.NewDecoder(req.Body).Decode(&user); err != nil {
		httperr.JSONError(res, err, http.StatusInternalServerError)

		return
	}

	newUser, err := h.authService.Register(req.Context(), &user)
	if err != nil {
		httperr.JSONError(res, err, http.StatusInternalServerError)

		return
	}

	res.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(res).Encode(newUser); err != nil {
		httperr.JSONError(res, err, http.StatusInternalServerError)

		return
	}
}

// SignInWithGoogle handler signs in a user using Google.
func (h *Handler) SignInWithGoogle(res http.ResponseWriter, req *http.Request) {
	var idTokenReq GoogleIDTokenVerifyRequest

	if err := json.NewDecoder(req.Body).Decode(&idTokenReq); err != nil {
		httperr.JSONError(res, err, http.StatusInternalServerError)

		return
	}

	cookie, err := req.Cookie("g_csrf_token")
	if err != nil {
		httperr.JSONError(res, err, http.StatusUnauthorized)

		return
	}

	if cookie.Value != idTokenReq.GCSRFToken {
		httperr.JSONError(res, httperr.ErrUnauthorized, http.StatusUnauthorized)

		return
	}

	user, err := h.authService.VerifyGoogleIDToken(req.Context(), idTokenReq.IDToken)
	if err != nil {
		httperr.JSONError(res, err, http.StatusInternalServerError)

		return
	}

	if err := json.NewEncoder(res).Encode(user); err != nil {
		httperr.JSONError(res, err, http.StatusInternalServerError)

		return
	}
}
