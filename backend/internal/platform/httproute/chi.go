package httproute

import "github.com/go-chi/chi/v5"

var _ Router = &ChiRouterWrapper{}

type ChiRouterWrapper struct {
	*chi.Mux
}

// Group is a wrapper method for chi.Router.Group
func (c *ChiRouterWrapper) Group(
	grpHandler func(r Router),
) Router {
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
	rtHandler func(r Router),
) Router {
	c.Mux.Route(pattern, func(r chi.Router) {
		m, ok := r.(*chi.Mux)
		if !ok {
			panic("chi.Router is not chi.Mux")
		}

		rtHandler(&ChiRouterWrapper{m})
	})

	return c
}
