// Package http provides HTTP transport layer handlers
package http

import (
	"log/slog"
	"net/http"

	"personal-web-platform/config"
	"personal-web-platform/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// Handler aggregates all HTTP handlers
type Handler struct {
	services *service.Services
	log      *slog.Logger
	cfg      *config.Config
}

// NewHandler creates a new HTTP handler
func NewHandler(services *service.Services, log *slog.Logger, cfg *config.Config) *Handler {
	return &Handler{
		services: services,
		log:      log,
		cfg:      cfg,
	}
}

// InitRoutes initializes all HTTP routes
func (h *Handler) InitRoutes() http.Handler {
	r := chi.NewRouter()

	// Standard middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)

	// Custom middleware
	r.Use(h.RequestLogger)
	r.Use(h.CORS(&h.cfg.CORS))
	r.Use(h.RateLimit(&h.cfg.RateLimit))

	// Health checks (outside /api/v1 for easier access)
	r.Get("/health", h.health)
	r.Get("/ready", h.ready)

	// Auth routes
	r.Route("/auth", func(r chi.Router) {
		r.Get("/{provider}", h.authLogin)
		r.Get("/{provider}/callback", h.authCallback)

		// Protected routes
		r.Group(func(r chi.Router) {
			r.Use(h.AuthRequired)
			r.Get("/me", h.authMe)
			r.Post("/logout", h.authLogout)
		})
	})

	r.Route("/api/v1", func(r chi.Router) {
		// Public endpoints
		r.Get("/profile", h.getProfile)

		// Protected endpoints (require authentication)
		r.Group(func(r chi.Router) {
			r.Use(h.AuthRequired)

			// Posts endpoints
			r.Get("/posts", h.listPosts)
			r.Get("/posts/{slug}", h.getPostBySlug)
			r.Get("/posts/{slug}/comments", h.getCommentsByPostSlug)

			// Comments (authenticated users)
			r.Post("/posts/{slug}/comments", h.createComment)
			r.Put("/comments/{id}", h.updateComment)
			r.Delete("/comments/{id}", h.deleteComment)
		})

		// Admin endpoints
		r.Group(func(r chi.Router) {
			r.Use(h.AuthRequired)
			r.Use(h.AdminRequired)
			r.Put("/admin/profile", h.updateProfile)
			r.Post("/admin/posts", h.createPost)
			r.Put("/admin/posts/{id}", h.updatePost)
			r.Delete("/admin/posts/{id}", h.deletePost)
		})
	})

	return r
}
