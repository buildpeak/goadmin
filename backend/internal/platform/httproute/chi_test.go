package httproute

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
)

func ExampleChiRouterWrapper() {
	// chi.NewRouter() returns a new chi router
	router := chi.NewRouter()

	// chi.Router is a chi.Mux
	rwrp := &ChiRouterWrapper{Mux: router}

	rwrp.Group(func(r Router) {
		// r is a chi.Mux
		r.Route("/v1/users", func(r Router) {
			// r is a chi.Mux
			r.Get("/", func(w http.ResponseWriter, _ *http.Request) {
				// handle get request
				w.Write([]byte("GET /v1/users"))
			})
		})
	})

	srv := httptest.NewServer(rwrp)
	defer srv.Close()

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, srv.URL+"/v1/users", nil)

	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)
	fmt.Printf("%d %s", res.StatusCode, string(body))
	// Output: 200 GET /v1/users
}

func TestChiRouterWrapper_Group(t *testing.T) {
	t.Parallel()

	type fields struct {
		Mux *chi.Mux
	}

	type args struct {
		grpHandler func(r Router)
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   Router
	}{
		{
			name: "group",
			fields: fields{
				Mux: chi.NewRouter(),
			},
			args: args{
				grpHandler: func(r Router) {
					r.Get("/", func(w http.ResponseWriter, _ *http.Request) {
						w.Write([]byte("OK"))
					})
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c := &ChiRouterWrapper{
				Mux: tt.fields.Mux,
			}

			c.Group(tt.args.grpHandler)

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			w := httptest.NewRecorder()

			c.ServeHTTP(w, req)

			if got := w.Code; got != http.StatusOK {
				t.Errorf("GorillaRouterWrapper.Route() = %v, want %v", got, http.StatusOK)
			}

			if got := w.Body.String(); got != "OK" {
				t.Errorf("GorillaRouterWrapper.Route() = %v, want %v", got, "OK")
			}
		})
	}
}

func TestChiRouterWrapper_Route(t *testing.T) {
	t.Parallel()

	type fields struct {
		Mux *chi.Mux
	}

	type args struct {
		pattern   string
		rtHandler func(r Router)
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   Router
	}{
		{
			name: "route",
			fields: fields{
				Mux: chi.NewRouter(),
			},
			args: args{
				pattern: "/test",
				rtHandler: func(r Router) {
					r.Get("/", func(w http.ResponseWriter, _ *http.Request) {
						w.Write([]byte("OK"))
					})
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c := &ChiRouterWrapper{
				Mux: tt.fields.Mux,
			}

			c.Route(tt.args.pattern, tt.args.rtHandler)

			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			w := httptest.NewRecorder()

			c.ServeHTTP(w, req)

			if got := w.Code; got != http.StatusOK {
				t.Errorf("GorillaRouterWrapper.Route() = %v, want %v", got, http.StatusOK)
			}

			if got := w.Body.String(); got != "OK" {
				t.Errorf("GorillaRouterWrapper.Route() = %v, want %v", got, "OK")
			}
		})
	}
}

func TestChiRouterWrapper_Get(t *testing.T) {
	t.Parallel()

	type fields struct {
		Mux *chi.Mux
	}

	type args struct {
		pattern string
		h       http.HandlerFunc
	}

	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "get",
			fields: fields{
				Mux: chi.NewRouter(),
			},
			args: args{
				pattern: "/",
				h: func(w http.ResponseWriter, _ *http.Request) {
					w.WriteHeader(http.StatusOK)
					w.Write([]byte("OK"))
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c := &ChiRouterWrapper{
				Mux: tt.fields.Mux,
			}

			c.Get(tt.args.pattern, tt.args.h)

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			w := httptest.NewRecorder()

			c.ServeHTTP(w, req)

			if got := w.Code; got != http.StatusOK {
				t.Errorf("GorillaRouterWrapper.Trace() = %v, want %v", got, http.StatusOK)
			}

			if got := w.Body.String(); got != "OK" {
				t.Errorf("GorillaRouterWrapper.Trace() = %v, want %v", got, "OK")
			}
		})
	}
}
