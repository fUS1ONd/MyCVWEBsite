// Package http provides HTTP transport layer handlers
package http

import (
	"log/slog"
	"net/http"

	"personal-web-platform/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// Handler aggregates all HTTP handlers
type Handler struct {
	services *service.Services
	log      *slog.Logger
}

// NewHandler creates a new HTTP handler
func NewHandler(services *service.Services, log *slog.Logger) *Handler {
	return &Handler{
		services: services,
		log:      log,
	}
}

// InitRoutes initializes all HTTP routes
func (h *Handler) InitRoutes() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Health checks (outside /api/v1 for easier access)
	r.Get("/health", h.health)
	r.Get("/ready", h.ready)

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/profile", h.getProfile)
	})

	return r
}
