package httproute

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/mux"
)

func TestURLParam(t *testing.T) {
	t.Parallel()

	type args struct {
		req                 *http.Request
		key                 string
		routerFrameworkName string
	}

	tests := []struct {
		name      string
		router    Router
		args      args
		want      string
		wantPanic bool
	}{
		{
			name: "gorilla/mux",
			router: &GorillaRouterWrapper{
				Router: mux.NewRouter(),
			},
			args: args{
				req:                 httptest.NewRequest(http.MethodGet, "/users/1", nil),
				key:                 "id",
				routerFrameworkName: "gorilla/mux",
			},
			want: "1",
		},
		{
			name: "chi",
			router: &ChiRouterWrapper{
				Mux: chi.NewRouter(),
			},
			args: args{
				req:                 httptest.NewRequest(http.MethodGet, "/users/1", nil),
				key:                 "id",
				routerFrameworkName: "chi",
			},
			want: "1",
		},
		{
			name: "unsupported router framework",
			router: &GorillaRouterWrapper{
				Router: mux.NewRouter(),
			},
			args: args{
				req:                 httptest.NewRequest(http.MethodGet, "/users/1", nil),
				key:                 "id",
				routerFrameworkName: "unsupported",
			},
			wantPanic: true,
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if tt.wantPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("URLParam() did not panic")
					}
				}()
			}

			tt.router.Get("/users/{id}", func(w http.ResponseWriter, req *http.Request) {
				if got := URLParam(req, tt.args.key, tt.args.routerFrameworkName); got != tt.want {
					t.Errorf("URLParam() = %v, want %v", got, tt.want)
				}

				w.WriteHeader(http.StatusOK)
			})

			tt.router.ServeHTTP(httptest.NewRecorder(), tt.args.req)
		})
	}
}
