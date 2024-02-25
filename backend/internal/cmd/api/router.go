package api

import (
	"net/http"

	"github.com/gorilla/mux"

	"goadmin-backend/internal/auth"
	"goadmin-backend/internal/user"
)

type Handlers struct {
	AuthHandler *auth.Handler
	UserHandler *user.Handler
}

func NewRouter(handlers *Handlers) http.Handler {
	router := mux.NewRouter()

	authHandler := handlers.AuthHandler

	// public routes
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// private routes
	privateRouter := router.PathPrefix("/v1").Subrouter()
	privateRouter.Use(authHandler.Authenticator())

	privateRouter.HandleFunc("/users", handlers.UserHandler.List).Methods(http.MethodGet)

	return router
}
