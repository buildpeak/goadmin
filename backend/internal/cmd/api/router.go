package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

// NewRouter returns a new router with all routes defined
func NewRouter(handlers *Handlers) http.Handler {
	router := mux.NewRouter()

	authHandler := handlers.AuthHandler

	// public routes
	router.HandleFunc("/health", handlers.HealthHandler).Methods(http.MethodGet)

	// private routes
	privateRouter := router.PathPrefix("/v1").Subrouter()
	privateRouter.Use(authHandler.Authenticator())

	privateRouter.HandleFunc("/users", handlers.UserHandler.List).Methods(http.MethodGet)

	return router
}
