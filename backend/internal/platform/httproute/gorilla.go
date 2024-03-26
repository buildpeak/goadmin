package httproute

import (
	"net/http"

	"github.com/gorilla/mux"

	"goadmin-backend/internal/platform/slices"
)

type GorillaRouterWrapper struct {
	*mux.Router
}

// Use is a wrapper method for mux.Router.Use
func (g *GorillaRouterWrapper) Use(
	middleware ...func(http.Handler) http.Handler,
) {
	mwf := slices.Map(
		middleware,
		func(f func(http.Handler) http.Handler) mux.MiddlewareFunc {
			return mux.MiddlewareFunc(f)
		},
	)
	g.Router.Use(mwf...)
}

// Group is a wrapper method for mux.Router.PathPrefix
func (g *GorillaRouterWrapper) Group(
	grpHandler func(r Router),
) Router {
	sr := g.NewRoute().Subrouter()
	wrp := &GorillaRouterWrapper{sr}
	grpHandler(wrp)

	return wrp
}

// Route is a method that returns a new router
func (g *GorillaRouterWrapper) Route(
	pattern string,
	rtHandler func(r Router),
) Router {
	sr := g.NewRoute().Path(pattern).Subrouter()
	wrp := &GorillaRouterWrapper{sr}
	rtHandler(wrp)

	return wrp
}

// HTTP-method routing along `pattern`
func (g *GorillaRouterWrapper) Connect(pattern string, h http.HandlerFunc) {
	g.HandleFunc(pattern, h).Methods(http.MethodConnect)
}

func (g *GorillaRouterWrapper) Delete(pattern string, h http.HandlerFunc) {
	g.HandleFunc(pattern, h).Methods(http.MethodDelete)
}

func (g *GorillaRouterWrapper) Get(pattern string, h http.HandlerFunc) {
	g.HandleFunc(pattern, h).Methods(http.MethodGet)
}

func (g *GorillaRouterWrapper) Head(pattern string, h http.HandlerFunc) {
	g.HandleFunc(pattern, h).Methods(http.MethodHead)
}

func (g *GorillaRouterWrapper) Options(pattern string, h http.HandlerFunc) {
	g.HandleFunc(pattern, h).Methods(http.MethodOptions)
}

func (g *GorillaRouterWrapper) Patch(pattern string, h http.HandlerFunc) {
	g.HandleFunc(pattern, h).Methods(http.MethodPatch)
}

func (g *GorillaRouterWrapper) Post(pattern string, h http.HandlerFunc) {
	g.HandleFunc(pattern, h).Methods(http.MethodPost)
}

func (g *GorillaRouterWrapper) Put(pattern string, h http.HandlerFunc) {
	g.HandleFunc(pattern, h).Methods(http.MethodPut)
}

func (g *GorillaRouterWrapper) Trace(pattern string, h http.HandlerFunc) {
	g.HandleFunc(pattern, h).Methods(http.MethodTrace)
}
