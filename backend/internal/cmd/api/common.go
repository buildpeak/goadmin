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
	HealthHandler http.HandlerFunc
}

type healthHandler struct {
	logger *slog.Logger
}

// NewHealthHandler returns a new health handler.
func NewHealthHandler(logger *slog.Logger) http.HandlerFunc {
	h := &healthHandler{logger: logger}

	return h.health
}

func (h *healthHandler) health(res http.ResponseWriter, _ *http.Request) {
	res.WriteHeader(http.StatusOK)

	if _, err := res.Write([]byte("OK")); err != nil {
		h.logger.Error("failed to write response", slog.Any("err", err))
	}
}
