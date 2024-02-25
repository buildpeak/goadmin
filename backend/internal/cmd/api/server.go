package api

import (
	"context"
	"log/slog"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	HttpServer *http.Server
	Config     *Config
	Logger     *slog.Logger
	DB         *pgxpool.Pool
}

type HTTPServerConfig struct {
	Addr string
}

func NewHTTPServer(config *HTTPServerConfig, handler http.Handler) *http.Server {
	server := &http.Server{
		Addr:         config.Addr,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      handler,
	}

	return server
}

func GracefulShutdown(ctx context.Context, closers ...func() error) error {
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

	for _, closer := range closers {
		if err := closer(); err != nil {
			return err
		}
	}

	return nil
}
