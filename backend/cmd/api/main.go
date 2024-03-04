package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/joho/godotenv/autoload"

	"goadmin-backend/internal/auth"
	"goadmin-backend/internal/cmd/api"
	"goadmin-backend/internal/platform/logging"
	"goadmin-backend/internal/repository/postgres"
	"goadmin-backend/internal/user"
)

func main() {
	exitCode := 0
	defer func() {
		os.Exit(exitCode)
	}()

	apiCtx, cancelAPICtx := context.WithCancel(context.Background())
	defer cancelAPICtx()

	// config
	cfg, err := api.NewConfig()
	if err != nil {
		slog.Error("failed to load config: %v", slog.Any("err", err))

		exitCode = 1

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

		exitCode = 1

		return
	}

	// repositories
	userRepo := postgres.NewUserRepo(dbpool)
	revokedTokenRepo := postgres.NewRevokedTokenRepo(dbpool)

	// services
	authService := auth.NewAuthService(
		userRepo,
		revokedTokenRepo,
		[]byte(cfg.API.Auth.JWTSecret),
	)
	userService := user.NewUserService(userRepo)

	// openapi-validator
	openapiValidator, err := api.NewOpenAPIValidator("")
	if err != nil {
		logger.Error("failed to create openapi validator", slog.Any("err", err))

		exitCode = 1

		return
	}

	// router
	apiHandler := api.NewRouter(openapiValidator, &api.Handlers{
		AuthHandler:   auth.NewHandler(authService),
		UserHandler:   user.NewHandler(userService),
		HealthHandler: api.NewHealthHandler(logger),
	})

	// otel
	shutdownOtel, err := api.StartOtel(
		api.ServiceInfo{},
		api.ObservabilityConfig{},
	)
	if err != nil {
		logger.Error("failed to start otel", slog.Any("err", err))

		exitCode = 1

		return
	}

	// web server
	mainAPIServer := api.NewHTTPServer(&api.HTTPServerConfig{
		Addr: fmt.Sprintf("0.0.0.0:%d", cfg.API.Port),
	}, apiHandler)

	go func() {
		if err := mainAPIServer.ListenAndServe(); err != nil {
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
		shutdownOtel()

		return nil
	}, func(_ context.Context) error {
		logger.Info("closing database connection")
		dbpool.Close()

		return nil
	}); err != nil {
		logger.Error("failed to shutdown server", slog.Any("err", err))
	}
}
