package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"

	"goadmin-backend/internal/auth"
	"goadmin-backend/internal/cmd/api"
	"goadmin-backend/internal/platform/logging"
	"goadmin-backend/internal/repository/postgres"
	"goadmin-backend/internal/user"
)

func main() {
	apiCtx, cancelAPICtx := context.WithCancel(context.Background())
	defer cancelAPICtx()

	// config
	cfg, err := api.NewConfig()
	if err != nil {
		slog.Error("failed to load config: %v", slog.Any("err", err))
		cancelAPICtx()
		os.Exit(1)
	}

	// logger
	logger := logging.NewLogger()

	// db
	db, err := pgxpool.New(apiCtx, cfg.DatabaseURL)
	if err != nil {
		logger.Error("failed to connect to database", slog.Any("err", err))
		cancelAPICtx()
		os.Exit(1)
	}

	// app := &api.App{
	// 	Config: cfg,
	// 	Logger: logger,
	// 	DB:     db,
	// }

	// repositories
	userRepo := postgres.NewUserRepo(db)
	revokedTokenRepo := postgres.NewRevokedTokenRepo(db)

	// services
	authService := auth.NewAuthService(
		userRepo,
		revokedTokenRepo,
		[]byte(cfg.APIServer.Auth.JWTSecret),
	)
	userService := user.NewUserService(userRepo)

	// router
	apiHandler := api.NewRouter(&api.Handlers{
		AuthHandler: auth.NewHandler(authService),
		UserHandler: user.NewHandler(userService),
	})

	// otel
	shutdownOtel, err := api.StartOtel(
		api.ServiceInfo{},
		api.ObservabilityConfig{},
	)
	if err != nil {
		logger.Error("failed to start otel", slog.Any("err", err))
		cancelAPICtx()
		os.Exit(1)
	}

	// web server
	mainAPIServer := api.NewHTTPServer(&api.HTTPServerConfig{
		Addr: fmt.Sprintf("0.0.0.0:%d", cfg.APIServer.Port),
	}, apiHandler)

	go func() {
		if err := mainAPIServer.ListenAndServe(); err != nil {
			logger.Error("failed to start server", slog.Any("err", err))
			cancelAPICtx()
		}
	}()
	logger.Info("server started", slog.Any("addr", mainAPIServer.Addr))

	// graceful shutdown
	if err := api.GracefulShutdown(apiCtx, func() error {
		logger.Info("shutting down server")
		return mainAPIServer.Shutdown(apiCtx)
	}, func() error {
		logger.Info("shutting down otel")
		shutdownOtel()
		return nil
	}, func() error {
		logger.Info("closing database connection")
		db.Close()
		return nil
	}); err != nil {
		logger.Error("failed to shutdown server", slog.Any("err", err))
	}
}
