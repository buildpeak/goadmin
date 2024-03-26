package httproute

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
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
		// TODO
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c := &ChiRouterWrapper{
				Mux: tt.fields.Mux,
			}

			if got := c.Group(tt.args.grpHandler); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ChiRouterWrapper.Group() = %v, want %v", got, tt.want)
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
		// TODO
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c := &ChiRouterWrapper{
				Mux: tt.fields.Mux,
			}

			if got := c.Route(tt.args.pattern, tt.args.rtHandler); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ChiRouterWrapper.Route() = %v, want %v", got, tt.want)
			}
		})
	}
}
