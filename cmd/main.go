package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"urlshorter/internal/config"
	"urlshorter/internal/http-server/handlers/url/delete"
	"urlshorter/internal/http-server/handlers/url/redirect"
	"urlshorter/internal/http-server/handlers/url/save"
	"urlshorter/internal/storage/sqlite"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoadConfig()

	log := setupLogger(cfg.Env)
	log.Info(
		"Starting server",
		slog.String("env", cfg.Env),
	)
	log.Debug("debug messages are enabled")

	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error(
			"failed to init storage",
			slog.String("error", err.Error()),
		)
		os.Exit(1)
	}

	router := chi.NewRouter()

	//middleware
	router.Use(middleware.RequestID) // Adds request ID for each request, useful for debugging
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	router.Use(middleware.Logger)

	router.Route("/url", func(r chi.Router) {
		r.Use(middleware.BasicAuth("url-shortener", map[string]string{
			cfg.HTTPServer.User: cfg.HTTPServer.Password,
		}))
		r.Post("/", save.New(log, storage))
		r.Delete("/{alias}", delete.New(log, storage))
	})
	router.Get("/{alias}", redirect.New(log, storage))

	log.Info("Server listening", slog.String("address", cfg.Address))

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.Timeout,
	}

	// Channel to listen for interrupt signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	// Run server in a goroutine
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Error("failed to start server", slog.String("error", err.Error()))
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
	<-stop
	log.Info("shutting down server gracefully")

	// Give 5 seconds for ongoing requests to finish
	timeout := 5 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Error("failed to gracefully shutdown server", slog.String("error", err.Error()))
	}

	log.Info("server stopped")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal, envDev:
		log = slog.New(
			slog.NewTextHandler(
				os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug},
			),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}
	return log
}
