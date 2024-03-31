package api

import (
	"context"
	"errors"
	"net/http"
	"os"
	"testing"
	"time"
)

func TestNewHTTPServer(t *testing.T) {
	t.Parallel()

	type args struct {
		config  *HTTPServerConfig
		handler http.Handler
	}

	tests := []struct {
		name string
		args args
		want *http.Server
	}{
		{
			name: "success",
			args: args{
				config: &HTTPServerConfig{
					Addr: "localhost:8080",
				},
				handler: http.DefaultServeMux,
			},
			want: &http.Server{
				Addr:         "localhost:8080",
				ReadTimeout:  DefaultReadTimeout,
				WriteTimeout: DefaultWriteTimeout,
				IdleTimeout:  DefaultIdleTimeout,
				Handler:      http.DefaultServeMux,
			},
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := NewHTTPServer(tt.args.config, tt.args.handler)
			if got.Addr != tt.want.Addr {
				t.Errorf("NewHTTPServer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGracefulShutdown(t *testing.T) {
	t.Parallel()

	type args struct {
		closers []func(context.Context) error
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				closers: []func(context.Context) error{
					func(_ context.Context) error {
						return nil
					},
				},
			},
		},
		{
			name: "error",
			args: args{
				closers: []func(context.Context) error{
					func(_ context.Context) error {
						return errors.New("error")
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			go func() {
				if err := GracefulShutdown(context.Background(), tt.args.closers...); (err != nil) != tt.wantErr {
					t.Errorf("GracefulShutdown() error = %v, wantErr %v", err, tt.wantErr)
				}
			}()

			time.Sleep(200 * time.Millisecond)

			process, err := os.FindProcess(os.Getpid())
			if err != nil {
				t.Fatal(err)
			}

			err = process.Signal(os.Interrupt)
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}
