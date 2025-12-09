// Package service provides business logic layer
package service

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"
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
	log         *slog.Logger
}

// NewProfileService creates a new profile service implementation
func NewProfileService(profileRepo repository.ProfileRepository, log *slog.Logger) ProfileService {
	return &profileService{
		profileRepo: profileRepo,
		log:         log,
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

	// Get old profile for cleanup
	oldProfile, _ := s.profileRepo.GetProfile(ctx)

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

	// Cleanup old file if it was a local upload and changed
	if oldProfile != nil && oldProfile.PhotoURL != "" && oldProfile.PhotoURL != profile.PhotoURL {
		if strings.HasPrefix(oldProfile.PhotoURL, "/uploads/") {
			// Assuming relative path from app root. path starts with /uploads/, so "." + ... works
			if err := os.Remove("." + oldProfile.PhotoURL); err != nil {
				s.log.Warn("failed to remove old profile photo", "path", oldProfile.PhotoURL, "error", err)
			}
		}
	}

	return profile, nil
}
