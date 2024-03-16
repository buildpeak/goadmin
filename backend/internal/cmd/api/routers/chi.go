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

type ChiRouterWrapper struct {
	*chi.Mux
}

// Group is a wrapper method for chi.Router.Group
func (c *ChiRouterWrapper) Group(
	grpHandler func(r httproute.Router),
) httproute.Router {
	c.Mux.Group(func(r chi.Router) {
		m, ok := r.(*chi.Mux)
		if !ok {
			panic("chi.Router is not chi.Mux")
		}

		grpHandler(&ChiRouterWrapper{m})
	})

	return c
}

// Route is a method that returns a new router
func (c *ChiRouterWrapper) Route(
	pattern string,
	rtHandler func(r httproute.Router),
) httproute.Router {
	c.Mux.Route(pattern, func(r chi.Router) {
		m, ok := r.(*chi.Mux)
		if !ok {
			panic("chi.Router is not chi.Mux")
		}

		rtHandler(&ChiRouterWrapper{m})
	})

	return c
}

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

	return &ChiRouterWrapper{router}
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
