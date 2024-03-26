package routers

import (
	"log/slog"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httplog/v2"

	"goadmin-backend/internal/platform/httproute"
)

const healthQuietDownPeriod = 5 * time.Second

// NewChiRouter is a function that returns a new router
func NewChiRouter(logger *slog.Logger) httproute.Router {
	router := chi.NewRouter()

	router.Use(httplog.RequestLogger(newHTTPLogger(logger, nil)))
	router.Use(middleware.Recoverer)
	router.Use(cors.Handler(cors.Options{
		AllowedHeaders: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		Debug:          true,
	}))

	return &httproute.ChiRouterWrapper{Mux: router}
}

func newHTTPLogger(
	logger *slog.Logger,
	quietDownRoutes []string,
) *httplog.Logger {
	if len(quietDownRoutes) == 0 {
		quietDownRoutes = []string{"/", "/health", "/liveness"}
	}

	return &httplog.Logger{
		Logger: logger,
		Options: httplog.Options{
			Concise:          true,
			RequestHeaders:   true,
			MessageFieldName: "msg",
			QuietDownRoutes:  quietDownRoutes,
			QuietDownPeriod:  healthQuietDownPeriod,
			SourceFieldName:  "caller",
		},
	}
}

// mux.Use(httplog.RequestLogger(logger))
// // add otel middleware
// mux.Use(otelHandler)
