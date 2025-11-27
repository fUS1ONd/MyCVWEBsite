// Package service provides business logic layer
package service

import (
	"context"
	"fmt"
	"sync"
	"time"

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
	cache       *domain.Profile
	cacheMu     sync.RWMutex
	lastUpdate  time.Time
}

// NewProfileService creates a new profile service implementation
func NewProfileService(profileRepo repository.ProfileRepository) ProfileService {
	return &profileService{
		profileRepo: profileRepo,
	}
}

func (s *profileService) GetProfile(ctx context.Context) (*domain.Profile, error) {
	// Check cache
	s.cacheMu.RLock()
	if s.cache != nil && time.Since(s.lastUpdate) < 5*time.Minute {
		defer s.cacheMu.RUnlock()
		return s.cache, nil
	}
	s.cacheMu.RUnlock()

	profile, err := s.profileRepo.GetProfile(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get profile: %w", err)
	}

	if profile == nil {
		return nil, fmt.Errorf("profile not found")
	}

	// Update cache
	s.cacheMu.Lock()
	s.cache = profile
	s.lastUpdate = time.Now()
	s.cacheMu.Unlock()

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

	// Update cache
	s.cacheMu.Lock()
	s.cache = profile
	s.lastUpdate = time.Now()
	s.cacheMu.Unlock()

	return profile, nil
}
