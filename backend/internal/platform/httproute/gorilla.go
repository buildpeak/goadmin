package httproute

import (
	"net/http"

	"github.com/gorilla/mux"

	"goadmin-backend/internal/platform/slices"
)

type GorillaRouterWrapper struct {
	*mux.Router

	pathPrefixes []string
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
	wrp := &GorillaRouterWrapper{Router: sr}
	grpHandler(wrp)

	return wrp
}

// Route is a method that returns a new router
func (g *GorillaRouterWrapper) Route(
	pattern string,
	rtHandler func(r Router),
) Router {
	sr := g.PathPrefix(pattern).Subrouter()
	pathPrefixes := make([]string, len(g.pathPrefixes)+1)
	copy(pathPrefixes, g.pathPrefixes)
	pathPrefixes[len(pathPrefixes)-1] = pattern
	wrp := &GorillaRouterWrapper{Router: sr, pathPrefixes: pathPrefixes}
	rtHandler(wrp)

	return wrp
}

func (g *GorillaRouterWrapper) installHander(
	method, pattern string,
	handler http.HandlerFunc,
) {
	if len(g.pathPrefixes) > 0 && pattern == "/" {
		pattern = ""
	}

	g.Router.SkipClean(true)

	g.HandleFunc(pattern, handler).Methods(method)
}

// HTTP-method routing along `pattern`
func (g *GorillaRouterWrapper) Connect(pattern string, h http.HandlerFunc) {
	g.installHander(http.MethodConnect, pattern, h)
}

func (g *GorillaRouterWrapper) Delete(pattern string, h http.HandlerFunc) {
	g.installHander(http.MethodDelete, pattern, h)
}

func (g *GorillaRouterWrapper) Get(pattern string, h http.HandlerFunc) {
	g.installHander(http.MethodGet, pattern, h)
}

func (g *GorillaRouterWrapper) Head(pattern string, h http.HandlerFunc) {
	g.installHander(http.MethodHead, pattern, h)
}

func (g *GorillaRouterWrapper) Options(pattern string, h http.HandlerFunc) {
	g.installHander(http.MethodOptions, pattern, h)
}

func (g *GorillaRouterWrapper) Patch(pattern string, h http.HandlerFunc) {
	g.installHander(http.MethodPatch, pattern, h)
}

func (g *GorillaRouterWrapper) Post(pattern string, h http.HandlerFunc) {
	g.installHander(http.MethodPost, pattern, h)
}

func (g *GorillaRouterWrapper) Put(pattern string, h http.HandlerFunc) {
	g.installHander(http.MethodPut, pattern, h)
}

func (g *GorillaRouterWrapper) Trace(pattern string, h http.HandlerFunc) {
	g.installHander(http.MethodTrace, pattern, h)
}
