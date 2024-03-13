package api

import (
	"log/slog"
	"net/http"

	"goadmin-backend/internal/auth"
	"goadmin-backend/internal/user"
)

// Handlers contains all the HTTP handlers for the API.
type Handlers struct {
	AuthHandler   *auth.Handler
	UserHandler   *user.Handler
	HealthHandler *HealthHandler
}

type HealthHandler struct {
	Logger *slog.Logger
}

// NewHealthHandler returns a new health handler.
func NewHealthHandler(logger *slog.Logger) *HealthHandler {
	return &HealthHandler{Logger: logger}
}

func (h *HealthHandler) healthCheck(res http.ResponseWriter, _ *http.Request) {
	res.WriteHeader(http.StatusOK)

	if _, err := res.Write([]byte("OK")); err != nil {
		h.Logger.Error("failed to write response", slog.Any("err", err))
	}
}
