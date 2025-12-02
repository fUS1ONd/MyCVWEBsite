// Package main is the entry point of the application
package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"personal-web-platform/config"
	"personal-web-platform/internal/pkg/logger"
	"personal-web-platform/internal/pkg/oauth"
	"personal-web-platform/internal/repository"
	"personal-web-platform/internal/service"
	transport "personal-web-platform/internal/transport/http"
)

func main() {
	cfg := config.MustLoad()

	log := logger.SetupLogger(cfg.Env)
	log.Info("starting application", slog.String("env", cfg.Env))

	// Database initialization
	db, err := repository.NewPostgresDB(context.Background(), cfg.Database.URL)
	if err != nil {
		log.Error("failed to init db", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer db.Close()
	log.Info("connected to database")

	// Initialize OAuth providers
	oauth.InitProviders(cfg, log)

	// Layers initialization
	repo := repository.NewRepositories(db, cfg)
	services := service.NewServices(repo, cfg, log)
	handlers := transport.NewHandler(services, log, cfg)

	// Start background session cleanup
	go startSessionCleanup(log, repo.Session)
	log.Info("session cleanup started")

	// HTTP Server
	srv := &http.Server{
		Addr:         cfg.HTTPServer.Address,
		Handler:      handlers.InitRoutes(),
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("failed to start server", slog.String("error", err.Error()))
		}
	}()

	log.Info("server started", slog.String("address", cfg.HTTPServer.Address))

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Info("shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("failed to shutdown server", slog.String("error", err.Error()))
	}
	log.Info("server stopped")
}

// startSessionCleanup runs periodic cleanup of expired sessions
func startSessionCleanup(log *slog.Logger, sessionRepo repository.SessionRepository) {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		deleted, err := sessionRepo.CleanupExpiredSessions(ctx)
		if err != nil {
			log.Error("failed to cleanup expired sessions", slog.String("error", err.Error()))
		} else if deleted > 0 {
			log.Info("cleaned up expired sessions", slog.Int64("count", deleted))
		}
		cancel()
	}
}
