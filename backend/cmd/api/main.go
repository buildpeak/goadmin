package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/joho/godotenv/autoload"
	"google.golang.org/api/idtoken"

	"goadmin-backend/internal/auth"
	"goadmin-backend/internal/cmd/api"
	"goadmin-backend/internal/platform/logging"
	"goadmin-backend/internal/repository/postgres"
	"goadmin-backend/internal/user"
)

const (
	// Version is the version of the application.
	Version = "v0.0.1"
)

func main() {
	exitCode := 1
	defer func() {
		os.Exit(exitCode)
	}()

	apiCtx, cancelAPICtx := context.WithCancel(context.Background())
	defer cancelAPICtx()

	// config
	cfg, err := api.NewConfig()
	if err != nil {
		slog.Error("failed to load config: %v", slog.Any("err", err))

		return
	}

	// logger
	logger := logging.NewLogger(
		logging.WithLevel(cfg.Log.Level),
		logging.WithPretty(cfg.Log.Pretty),
	)

	// dbpool
	dbpool, err := pgxpool.New(apiCtx, cfg.DatabaseURL)
	if err != nil {
		logger.Error("failed to connect to database", slog.Any("err", err))

		return
	}

	// repositories
	userRepo := postgres.NewUserRepo(dbpool)
	revokedTokenRepo := postgres.NewRevokedTokenRepo(dbpool)

	// services
	idTknValidator, err := idtoken.NewValidator(apiCtx)
	if err != nil {
		logger.Error("failed to create id token validator", slog.Any("err", err))
	}

	authService := auth.NewAuthService(
		userRepo,
		revokedTokenRepo,
		[]byte(cfg.API.Auth.JWTSecret),
		idTknValidator,
		cfg.Google.ClientID,
	)
	userService := user.NewUserService(userRepo)

	// openapi-validator
	openapiValidator, err := api.NewOpenAPIValidator("", logger)
	if err != nil {
		logger.Error("failed to create openapi validator", slog.Any("err", err))

		return
	}

	// router
	apiHandler := api.NewRouter(openapiValidator, &api.Handlers{
		AuthHandler:   auth.NewHandler(authService, logger),
		UserHandler:   user.NewHandler(userService),
		HealthHandler: api.NewHealthHandler(logger),
	})

	// // otel
	// shutdownOtel, err := api.StartOtel(
	// 	api.ServiceInfo{
	// 		Name:    "goadmin",
	// 		Version: Version,
	// 		Env:     cfg.Env,
	// 	},
	// 	api.ObservabilityConfig{
	// 		Collector: api.Collector{
	// 			Host:               cfg.Observability.Collector.Host,
	// 			Port:               cfg.Observability.Collector.Port,
	// 			Headers:            cfg.Observability.Collector.Headers,
	// 			IsInsecure:         cfg.Observability.Collector.IsInsecure,
	// 			WithMetricsEnabled: cfg.Observability.Collector.WithMetricsEnabled,
	// 		},
	// 	},
	// )
	// if err != nil {
	// 	logger.Error("failed to start otel", slog.Any("err", err))
	//
	// 	return
	// }

	// web server
	mainAPIServer := api.NewHTTPServer(&api.HTTPServerConfig{
		Addr: fmt.Sprintf("0.0.0.0:%d", cfg.API.Port),
	}, apiHandler)

	go func() {
		if err := mainAPIServer.ListenAndServe(); err != nil &&
			!errors.Is(err, http.ErrServerClosed) {
			logger.Error("failed to start server", slog.Any("err", err))
			cancelAPICtx()
		}
	}()
	logger.Info("server started", slog.Any("addr", mainAPIServer.Addr))

	// graceful shutdown
	if err := api.GracefulShutdown(apiCtx, func(shutdownCtx context.Context) error {
		logger.Info("shutting down server")

		if err := mainAPIServer.Shutdown(shutdownCtx); err != nil {
			logger.Error("failed to shutdown server", slog.Any("err", err))
		}

		return nil
	}, func(_ context.Context) error {
		logger.Info("shutting down otel")
		// shutdownOtel()

		return nil
	}, func(_ context.Context) error {
		logger.Info("closing database connection")
		dbpool.Close()

		return nil
	}); err != nil {
		logger.Error("failed to shutdown server", slog.Any("err", err))
	}

	exitCode = 0
}
