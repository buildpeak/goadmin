package routers

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"goadmin-backend/internal/platform/httproute"
)

func NewGorillaMuxRouter() httproute.Router {
	router := mux.NewRouter()

	router.Use(handlers.RecoveryHandler()) // first install first run
	router.Use(mux.CORSMethodMiddleware(router))

	// other middlewares

	return &httproute.GorillaRouterWrapper{Router: router}
}
