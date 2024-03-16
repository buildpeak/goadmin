package api

import (
	"net/http"

	"goadmin-backend/internal/cmd/api/routers"
	"goadmin-backend/internal/platform/httproute"
)

// NewRouter returns a new router with all routes defined
func NewRouter(validator *OpenAPIValidator, handlers *Handlers) http.Handler {
	var router httproute.Router

	logger := handlers.HealthHandler.Logger

	// use chi router
	router = routers.NewChiRouter(logger)

	router.Use(validator.Middleware)

	// public routes
	router.Get("/health", handlers.HealthHandler.healthCheck)
	router.Post("/auth/login", handlers.AuthHandler.Login)
	router.Post("/auth/register", handlers.AuthHandler.Register)
	router.Post("/auth/signin-with-google", handlers.AuthHandler.SignInWithGoogle)

	// add private routes
	router.Group(func(r httproute.Router) {
		r.Use(handlers.AuthHandler.Authenticator())

		r.Route("/v1/users", func(r httproute.Router) {
			r.Get("/", handlers.UserHandler.List)
		})
	})

	return router
}
