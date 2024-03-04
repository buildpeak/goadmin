package auth

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"

	"goadmin-backend/internal/domain"
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

	// validate request
	validate := validator.New()
	if err := validate.Struct(credentials); err != nil {
		http.Error(res, "invalid request", http.StatusBadRequest)

		return
	}

	token, err := h.authService.Login(req.Context(), credentials)
	if err != nil {
		http.Error(res, "invalid credentials", http.StatusUnauthorized)

		return
	}

	if _, err := res.Write([]byte(token)); err != nil {
		http.Error(res, "internal server error", http.StatusInternalServerError)

		return
	}
}
