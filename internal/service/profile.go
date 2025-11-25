package service

import (
	"context"

	"personal-web-platform/internal/domain"
	"personal-web-platform/internal/repository"
)

type ProfileService interface {
	GetProfile(ctx context.Context) (domain.Profile, error)
}

type profileService struct {
	repo repository.ProfileRepository
}

func NewProfileService(repo repository.ProfileRepository) ProfileService {
	return &profileService{repo: repo}
}

func (s *profileService) GetProfile(ctx context.Context) (domain.Profile, error) {
	// Logic/transformations can be added here if needed
	return s.repo.Get(ctx)
}
