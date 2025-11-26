package repository

import (
	"context"

	"personal-web-platform/config"
	"personal-web-platform/internal/domain"
)

type profileRepo struct {
	cfg *config.Config
}

// NewProfileRepo creates a new profile repository implementation
func NewProfileRepo(cfg *config.Config) ProfileRepository {
	return &profileRepo{cfg: cfg}
}

func (r *profileRepo) Get(ctx context.Context) (domain.Profile, error) {
	return domain.Profile{
		Name:        r.cfg.Profile.Name,
		Description: r.cfg.Profile.Description,
		Skills:      r.cfg.Profile.Skills,
		Contacts: domain.Contacts{
			Email:    r.cfg.Profile.Contacts.Email,
			Github:   r.cfg.Profile.Contacts.Github,
			LinkedIn: r.cfg.Profile.Contacts.LinkedIn,
		},
	}, nil
}
