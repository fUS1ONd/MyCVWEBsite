package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"personal-web-platform/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ProfileRepository defines methods for profile data access
type ProfileRepository interface {
	GetProfile(ctx context.Context) (*domain.Profile, error)
	UpdateProfile(ctx context.Context, profile *domain.UpdateProfileRequest) (*domain.Profile, error)
}

type profileRepo struct {
	db *pgxpool.Pool
}

// NewProfileRepo creates a new profile repository implementation
func NewProfileRepo(db *pgxpool.Pool) ProfileRepository {
	return &profileRepo{db: db}
}

func (r *profileRepo) GetProfile(ctx context.Context) (*domain.Profile, error) {
	var profile domain.Profile
	var contactsJSON []byte
	db := GetQueryEngine(ctx, r.db)

	query := `
		SELECT id, name, description, photo_url, activity, contacts, created_at, updated_at
		FROM profile_info
		ORDER BY id ASC
		LIMIT 1
	`

	err := db.QueryRow(ctx, query).Scan(
		&profile.ID,
		&profile.Name,
		&profile.Description,
		&profile.PhotoURL,
		&profile.Activity,
		&contactsJSON,
		&profile.CreatedAt,
		&profile.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil // Profile not found
		}
		return nil, fmt.Errorf("failed to get profile: %w", err)
	}

	// Unmarshal contacts JSON
	if err := json.Unmarshal(contactsJSON, &profile.Contacts); err != nil {
		return nil, fmt.Errorf("failed to unmarshal contacts: %w", err)
	}

	return &profile, nil
}

func (r *profileRepo) UpdateProfile(ctx context.Context, req *domain.UpdateProfileRequest) (*domain.Profile, error) {
	var profile domain.Profile
	var contactsJSON []byte
	db := GetQueryEngine(ctx, r.db)

	// Marshal contacts to JSON
	contacts, err := json.Marshal(req.Contacts)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal contacts: %w", err)
	}

	query := `
		UPDATE profile_info
		SET name = $1, description = $2, photo_url = $3, activity = $4, contacts = $5, updated_at = NOW()
		WHERE id = (SELECT id FROM profile_info ORDER BY id ASC LIMIT 1)
		RETURNING id, name, description, photo_url, activity, contacts, created_at, updated_at
	`

	err = db.QueryRow(ctx, query,
		req.Name,
		req.Description,
		req.PhotoURL,
		req.Activity,
		contacts,
	).Scan(
		&profile.ID,
		&profile.Name,
		&profile.Description,
		&profile.PhotoURL,
		&profile.Activity,
		&contactsJSON,
		&profile.CreatedAt,
		&profile.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update profile: %w", err)
	}

	// Unmarshal contacts JSON
	if err := json.Unmarshal(contactsJSON, &profile.Contacts); err != nil {
		return nil, fmt.Errorf("failed to unmarshal contacts: %w", err)
	}

	return &profile, nil
}
