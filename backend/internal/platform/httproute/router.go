package httproute

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/mux"
)

// Router consisting of the core routing methods used by chi's Mux,
//
//nolint:interfacebloat // it's a wrapper of chi.Router
type Router interface {
	http.Handler
	Use(...func(http.Handler) http.Handler)
	Group(func(r Router)) Router
	Route(pattern string, fn func(r Router)) Router

	// HTTP-method routing along `pattern`
	Connect(pattern string, h http.HandlerFunc)
	Delete(pattern string, h http.HandlerFunc)
	Get(pattern string, h http.HandlerFunc)
	Head(pattern string, h http.HandlerFunc)
	Options(pattern string, h http.HandlerFunc)
	Patch(pattern string, h http.HandlerFunc)
	Post(pattern string, h http.HandlerFunc)
	Put(pattern string, h http.HandlerFunc)
	Trace(pattern string, h http.HandlerFunc)
}

// URLParam returns the URL path parameter value for a given key
func URLParam(req *http.Request, key, routerFrameworkName string) string {
	switch {
	case routerFrameworkName == "gorilla/mux":
		vars := mux.Vars(req)

		return vars[key]
	case routerFrameworkName == "chi":
		return chi.URLParam(req, key)
	default:
		panic("unsupported router framework")
	}
}
