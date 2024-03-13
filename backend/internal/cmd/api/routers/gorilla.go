package routers

import (
	"net/http"

	"github.com/gorilla/mux"
)

func newMuxRouter() http.Handler {
	router := mux.NewRouter()

	// authHandler := handlers.AuthHandler
	//
	// // public routes
	// router.HandleFunc("/health", handlers.HealthHandler).Methods(http.MethodGet)
	//
	// router.Use(mux.CORSMethodMiddleware(router))
	// router.Use(validator.Middleware)
	//
	// router.Use(ghndlrs.RecoveryHandler()) // must be last middleware
	//
	// router.HandleFunc("/auth/login", authHandler.Login).Methods(http.MethodPost)
	// router.HandleFunc("/auth/register", authHandler.Register).Methods(http.MethodPost)
	// router.HandleFunc("/auth/signin-with-google", authHandler.SignInWithGoogle).Methods(http.MethodPost)
	//
	// // private routes
	// privateRouter := router.PathPrefix("/v1").Subrouter()
	// privateRouter.Use(authHandler.Authenticator())
	//
	// privateRouter.HandleFunc("/users", handlers.UserHandler.List).Methods(http.MethodGet)

	return router
}
