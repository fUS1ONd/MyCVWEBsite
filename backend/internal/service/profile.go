// Package service provides business logic layer
package service

import (
	"context"

	"personal-web-platform/internal/domain"
	"personal-web-platform/internal/repository"
)

// ProfileService defines methods for profile business logic
type ProfileService interface {
	GetProfile(ctx context.Context) (domain.Profile, error)
}

type profileService struct {
	repo repository.ProfileRepository
}

// NewProfileService creates a new profile service implementation
func NewProfileService(repo repository.ProfileRepository) ProfileService {
	return &profileService{repo: repo}
}

func (s *profileService) GetProfile(ctx context.Context) (domain.Profile, error) {
	// Logic/transformations can be added here if needed
	return s.repo.Get(ctx)
}
