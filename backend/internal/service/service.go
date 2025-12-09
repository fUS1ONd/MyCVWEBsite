// Package service provides business logic layer
package service

import (
	"context"
	"log/slog"

	"personal-web-platform/config"
	"personal-web-platform/internal/repository"
)

// Services aggregates all service interfaces
type Services struct {
	Profile ProfileService
	Auth    AuthService
	Post    PostService
	Comment CommentService
	repos   *repository.Repositories
	cfg     *config.Config
}

// NewServices creates a new Services instance with all implementations
func NewServices(repos *repository.Repositories, cfg *config.Config, log *slog.Logger) *Services {
	return &Services{
		Profile: NewProfileService(repos.Profile),
		Auth:    NewAuthService(repos.Auth, repos.Session, repos.Transactor, cfg, log),
		Post:    NewPostService(repos.Post),
		Comment: NewCommentService(repos.Comment, repos.Post),
		repos:   repos,
		cfg:     cfg,
	}
}

// HealthCheck performs health check for all dependencies
func (s *Services) HealthCheck(ctx context.Context) error {
	// Check database connectivity
	if err := s.repos.Ping(ctx); err != nil {
		return err
	}

	return nil
}
