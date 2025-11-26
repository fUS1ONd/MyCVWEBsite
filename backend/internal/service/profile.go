// Package service provides business logic layer
package service

import (
	"context"
	"fmt"

	"personal-web-platform/internal/domain"
	"personal-web-platform/internal/pkg/validator"
	"personal-web-platform/internal/repository"
)

// ProfileService defines methods for profile business logic
type ProfileService interface {
	GetProfile(ctx context.Context) (*domain.Profile, error)
	UpdateProfile(ctx context.Context, req *domain.UpdateProfileRequest) (*domain.Profile, error)
}

type profileService struct {
	profileRepo repository.ProfileRepository
}

// NewProfileService creates a new profile service implementation
func NewProfileService(profileRepo repository.ProfileRepository) ProfileService {
	return &profileService{
		profileRepo: profileRepo,
	}
}

func (s *profileService) GetProfile(ctx context.Context) (*domain.Profile, error) {
	profile, err := s.profileRepo.GetProfile(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get profile: %w", err)
	}

	if profile == nil {
		return nil, fmt.Errorf("profile not found")
	}

	return profile, nil
}

func (s *profileService) UpdateProfile(ctx context.Context, req *domain.UpdateProfileRequest) (*domain.Profile, error) {
	// Validate request
	if err := validator.Validate(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Update profile
	profile, err := s.profileRepo.UpdateProfile(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to update profile: %w", err)
	}

	return profile, nil
}
