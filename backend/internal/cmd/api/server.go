package api

import (
	"context"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

const (
	// DefaultReadTimeout is the default read timeout for the HTTP server.
	DefaultReadTimeout = 5 * time.Second

	// DefaultWriteTimeout is the default write timeout for the HTTP server.
	DefaultWriteTimeout = 10 * time.Second

	// DefaultIdleTimeout is the default idle timeout for the HTTP server.
	DefaultIdleTimeout = 120 * time.Second

	// DefaultShutdownTimeout is the default shutdown timeout for the HTTP server.
	DefaultShutdownTimeout = 30 * time.Second
)

// HTTPServerConfig is the configuration for the HTTP server.
type HTTPServerConfig struct {
	Addr string
}

// NewHTTPServer returns a new HTTP server.
func NewHTTPServer(config *HTTPServerConfig, handler http.Handler) *http.Server {
	server := &http.Server{
		Addr:         config.Addr,
		ReadTimeout:  DefaultReadTimeout,
		WriteTimeout: DefaultWriteTimeout,
		IdleTimeout:  DefaultIdleTimeout,
		Handler:      handler,
	}

	return server
}

// GracefulShutdown gracefully shuts down the application.
func GracefulShutdown(ctx context.Context, closers ...func(context.Context) error) error {
	sigCtx, sigCtxCancel := signal.NotifyContext(
		ctx,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGHUP,
		syscall.SIGQUIT,
	)
	defer sigCtxCancel()

	// wait for signal
	<-sigCtx.Done()

	newCtx, cancel := context.WithTimeout(ctx, DefaultShutdownTimeout)
	defer cancel()

	errs := make([]error, len(closers))

	for i, closer := range closers {
		if err := closer(newCtx); err != nil {
			errs[i] = err
		}
	}

	// return the first non-nil error
	for _, err := range errs {
		if err != nil {
			return err
		}
	}

	return nil
}
