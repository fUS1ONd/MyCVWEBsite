// Package service provides business logic layer
package service

import (
	"context"

	"personal-web-platform/internal/repository"
)

// Services aggregates all service interfaces
type Services struct {
	Profile ProfileService
	repos   *repository.Repositories
}

// NewServices creates a new Services instance with all implementations
func NewServices(repos *repository.Repositories) *Services {
	return &Services{
		Profile: NewProfileService(repos.Profile),
		repos:   repos,
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
