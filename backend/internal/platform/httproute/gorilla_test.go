package httproute

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gorilla/mux"
)

func ExampleGorillaRouterWrapper() {
	// mux.NewRouter() returns a new gorilla router
	router := mux.NewRouter()

	// mux.Router is a mux.Router
	rwrp := &GorillaRouterWrapper{Router: router}

	rwrp.Group(func(r Router) {
		// r is a mux.Router
		r.Route("/v1/users", func(r Router) {
			// r is a mux.Router
			r.Get("/", func(w http.ResponseWriter, _ *http.Request) {
				// handle get request
				w.Write([]byte("GET /v1/users"))
			})
		})
	})

	srv := httptest.NewServer(rwrp)
	defer srv.Close()

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, srv.URL+"/v1/users", nil)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%d %s", res.StatusCode, string(body))
	// Output: 200 GET /v1/users
}

func TestGorillaRouterWrapper_Use(t *testing.T) {
	t.Parallel()

	type fields struct {
		Router *mux.Router
	}

	type args struct {
		middleware []func(http.Handler) http.Handler
	}

	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			g := &GorillaRouterWrapper{
				Router: tt.fields.Router,
			}

			g.Use(tt.args.middleware...)
		})
	}
}

func TestGorillaRouterWrapper_Group(t *testing.T) {
	t.Parallel()

	type fields struct {
		Router *mux.Router
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			g := &GorillaRouterWrapper{
				Router: tt.fields.Router,
			}

			if got := g.Group(tt.args.grpHandler); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GorillaRouterWrapper.Group() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGorillaRouterWrapper_Route(t *testing.T) {
	t.Parallel()

	type fields struct {
		Router *mux.Router
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			g := &GorillaRouterWrapper{
				Router: tt.fields.Router,
			}

			if got := g.Route(tt.args.pattern, tt.args.rtHandler); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GorillaRouterWrapper.Route() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGorillaRouterWrapper_Connect(t *testing.T) {
	t.Parallel()

	type fields struct {
		Router *mux.Router
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			g := &GorillaRouterWrapper{
				Router: tt.fields.Router,
			}

			g.Connect(tt.args.pattern, tt.args.h)
		})
	}
}

func TestGorillaRouterWrapper_Delete(t *testing.T) {
	t.Parallel()

	type fields struct {
		Router *mux.Router
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			g := &GorillaRouterWrapper{
				Router: tt.fields.Router,
			}

			g.Delete(tt.args.pattern, tt.args.h)
		})
	}
}

func TestGorillaRouterWrapper_Get(t *testing.T) {
	t.Parallel()

	type fields struct {
		Router *mux.Router
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			g := &GorillaRouterWrapper{
				Router: tt.fields.Router,
			}

			g.Get(tt.args.pattern, tt.args.h)
		})
	}
}

func TestGorillaRouterWrapper_Head(t *testing.T) {
	t.Parallel()

	type fields struct {
		Router *mux.Router
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			g := &GorillaRouterWrapper{
				Router: tt.fields.Router,
			}

			g.Head(tt.args.pattern, tt.args.h)
		})
	}
}

func TestGorillaRouterWrapper_Options(t *testing.T) {
	t.Parallel()

	type fields struct {
		Router *mux.Router
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			g := &GorillaRouterWrapper{
				Router: tt.fields.Router,
			}

			g.Options(tt.args.pattern, tt.args.h)
		})
	}
}

func TestGorillaRouterWrapper_Patch(t *testing.T) {
	t.Parallel()

	type fields struct {
		Router *mux.Router
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			g := &GorillaRouterWrapper{
				Router: tt.fields.Router,
			}

			g.Patch(tt.args.pattern, tt.args.h)
		})
	}
}

func TestGorillaRouterWrapper_Post(t *testing.T) {
	t.Parallel()

	type fields struct {
		Router *mux.Router
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			g := &GorillaRouterWrapper{
				Router: tt.fields.Router,
			}

			g.Post(tt.args.pattern, tt.args.h)
		})
	}
}

func TestGorillaRouterWrapper_Put(t *testing.T) {
	t.Parallel()

	type fields struct {
		Router *mux.Router
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			g := &GorillaRouterWrapper{
				Router: tt.fields.Router,
			}

			g.Put(tt.args.pattern, tt.args.h)
		})
	}
}

func TestGorillaRouterWrapper_Trace(t *testing.T) {
	t.Parallel()

	type fields struct {
		Router *mux.Router
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			g := &GorillaRouterWrapper{
				Router: tt.fields.Router,
			}

			g.Trace(tt.args.pattern, tt.args.h)
		})
	}
}
