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
	router.Post("/auth/signup", handlers.AuthHandler.Register)
	router.Post("/auth/signin-with-google", handlers.AuthHandler.SignInWithGoogle)

	// add private routes
	router.Group(func(grt httproute.Router) {
		grt.Use(handlers.AuthHandler.Authenticator())

		grt.Route("/auth/logout", func(r httproute.Router) {
			r.Post("/", handlers.AuthHandler.Logout)
		})

		grt.Route("/v1/users", func(r httproute.Router) {
			r.Get("/", handlers.UserHandler.List)
		})
	})

	return router
}
