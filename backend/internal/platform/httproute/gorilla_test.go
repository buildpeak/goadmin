package httproute

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
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
		{
			name: "use middleware",
			fields: fields{
				Router: mux.NewRouter(),
			},
			args: args{
				middleware: []func(http.Handler) http.Handler{
					func(next http.Handler) http.Handler {
						return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
							w.Header().Set("X-Test", "test")
							next.ServeHTTP(w, r)
						})
					},
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			g := &GorillaRouterWrapper{
				Router: tt.fields.Router,
			}

			g.Use(tt.args.middleware...)

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			w := httptest.NewRecorder()

			g.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("OK"))
			})

			g.ServeHTTP(w, req)

			if got := w.Header().Get("X-Test"); got != "test" {
				t.Errorf("GorillaRouterWrapper.Use() = %v, want %v", got, "test")
			}

			if got := w.Code; got != http.StatusOK {
				t.Errorf("GorillaRouterWrapper.Use() = %v, want %v", got, http.StatusOK)
			}

			if got := w.Body.String(); got != "OK" {
				t.Errorf("GorillaRouterWrapper.Use() = %v, want %v", got, "OK")
			}
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
		{
			name: "group",
			fields: fields{
				Router: mux.NewRouter(),
			},
			args: args{
				grpHandler: func(r Router) {
					r.Get("/", func(w http.ResponseWriter, _ *http.Request) {
						w.WriteHeader(http.StatusOK)
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

			g := &GorillaRouterWrapper{
				Router: tt.fields.Router,
			}

			g.Group(tt.args.grpHandler)

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			w := httptest.NewRecorder()

			g.ServeHTTP(w, req)

			if got := w.Code; got != http.StatusOK {
				t.Errorf("GorillaRouterWrapper.Route() = %v, want %v", got, http.StatusOK)
			}

			if got := w.Body.String(); got != "OK" {
				t.Errorf("GorillaRouterWrapper.Route() = %v, want %v", got, "OK")
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
		{
			name: "route",
			fields: fields{
				Router: mux.NewRouter(),
			},
			args: args{
				pattern: "/test",
				rtHandler: func(r Router) {
					r.Get("/", func(w http.ResponseWriter, _ *http.Request) {
						w.WriteHeader(http.StatusOK)
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

			g := &GorillaRouterWrapper{
				Router: tt.fields.Router,
			}

			g.Route(tt.args.pattern, tt.args.rtHandler)

			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			w := httptest.NewRecorder()

			g.ServeHTTP(w, req)

			if got := w.Code; got != http.StatusOK {
				t.Errorf("GorillaRouterWrapper.Route() = %v, want %v", got, http.StatusOK)
			}

			if got := w.Body.String(); got != "OK" {
				t.Errorf("GorillaRouterWrapper.Route() = %v, want %v", got, "OK")
			}
		})
	}
}

func checkHTTPMethodFunc(t *testing.T, g *GorillaRouterWrapper, method string) {
	t.Helper()

	req := httptest.NewRequest(method, "/", nil)
	w := httptest.NewRecorder()

	g.ServeHTTP(w, req)

	if got := w.Code; got != http.StatusOK {
		t.Errorf("GorillaRouterWrapper.Trace() = %v, want %v", got, http.StatusOK)
	}

	if got := w.Body.String(); got != "OK" {
		t.Errorf("GorillaRouterWrapper.Trace() = %v, want %v", got, "OK")
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
		{
			name: "connect",
			fields: fields{
				Router: mux.NewRouter(),
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

			g := &GorillaRouterWrapper{
				Router: tt.fields.Router,
			}

			g.Connect(tt.args.pattern, tt.args.h)

			checkHTTPMethodFunc(t, g, http.MethodConnect)
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
		{
			name: "delete",
			fields: fields{
				Router: mux.NewRouter(),
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

			g := &GorillaRouterWrapper{
				Router: tt.fields.Router,
			}

			g.Delete(tt.args.pattern, tt.args.h)

			checkHTTPMethodFunc(t, g, http.MethodDelete)
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
		{
			name: "get",
			fields: fields{
				Router: mux.NewRouter(),
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

			g := &GorillaRouterWrapper{
				Router: tt.fields.Router,
			}

			g.Get(tt.args.pattern, tt.args.h)

			checkHTTPMethodFunc(t, g, http.MethodGet)
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
		{
			name: "head",
			fields: fields{
				Router: mux.NewRouter(),
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

			g := &GorillaRouterWrapper{
				Router: tt.fields.Router,
			}

			g.Head(tt.args.pattern, tt.args.h)

			checkHTTPMethodFunc(t, g, http.MethodHead)
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
		{
			name: "options",
			fields: fields{
				Router: mux.NewRouter(),
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

			g := &GorillaRouterWrapper{
				Router: tt.fields.Router,
			}

			g.Options(tt.args.pattern, tt.args.h)

			checkHTTPMethodFunc(t, g, http.MethodOptions)
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
		{
			name: "patch",
			fields: fields{
				Router: mux.NewRouter(),
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

			g := &GorillaRouterWrapper{
				Router: tt.fields.Router,
			}

			g.Patch(tt.args.pattern, tt.args.h)

			checkHTTPMethodFunc(t, g, http.MethodPatch)
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
		{
			name: "post",
			fields: fields{
				Router: mux.NewRouter(),
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

			g := &GorillaRouterWrapper{
				Router: tt.fields.Router,
			}

			g.Post(tt.args.pattern, tt.args.h)

			checkHTTPMethodFunc(t, g, http.MethodPost)
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
		{
			name: "put",
			fields: fields{
				Router: mux.NewRouter(),
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

			g := &GorillaRouterWrapper{
				Router: tt.fields.Router,
			}

			g.Put(tt.args.pattern, tt.args.h)

			checkHTTPMethodFunc(t, g, http.MethodPut)
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
		{
			name: "trace",
			fields: fields{
				Router: mux.NewRouter(),
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

			g := &GorillaRouterWrapper{
				Router: tt.fields.Router,
			}

			g.Trace(tt.args.pattern, tt.args.h)

			checkHTTPMethodFunc(t, g, http.MethodTrace)
		})
	}
}
