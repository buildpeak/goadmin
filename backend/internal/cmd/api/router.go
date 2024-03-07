package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

// NewRouter returns a new router with all routes defined
func NewRouter(validator *OpenAPIValidator, handlers *Handlers) http.Handler {
	router := mux.NewRouter()

	authHandler := handlers.AuthHandler

	// public routes
	router.HandleFunc("/health", handlers.HealthHandler).Methods(http.MethodGet)

	router.Use(validator.Middleware)
	router.HandleFunc("/auth/login", authHandler.Login).Methods(http.MethodPost)
	router.HandleFunc("/auth/register", authHandler.Register).Methods(http.MethodPost)

	// private routes
	privateRouter := router.PathPrefix("/v1").Subrouter()
	privateRouter.Use(authHandler.Authenticator())

	privateRouter.HandleFunc("/users", handlers.UserHandler.List).Methods(http.MethodGet)

	return router
}
