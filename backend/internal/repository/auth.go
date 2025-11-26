package repository

import (
	"context"
	"fmt"

	"personal-web-platform/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// AuthRepository defines methods for authentication data access
type AuthRepository interface {
	CreateUser(ctx context.Context, email string, role domain.Role) (*domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	GetUserByID(ctx context.Context, id int) (*domain.User, error)
	LinkOAuthProvider(ctx context.Context, provider *domain.OAuthProvider) error
	GetUserByProviderID(ctx context.Context, providerName, providerUserID string) (*domain.User, error)
	GetOAuthProvider(ctx context.Context, userID int, providerName string) (*domain.OAuthProvider, error)
	UpdateOAuthProvider(ctx context.Context, provider *domain.OAuthProvider) error
}

type authRepo struct {
	db *pgxpool.Pool
}

// NewAuthRepo creates a new auth repository implementation
func NewAuthRepo(db *pgxpool.Pool) AuthRepository {
	return &authRepo{db: db}
}

func (r *authRepo) CreateUser(ctx context.Context, email string, role domain.Role) (*domain.User, error) {
	var user domain.User

	query := `
		INSERT INTO users (email, role, created_at)
		VALUES ($1, $2, NOW())
		RETURNING id, email, role, created_at
	`

	err := r.db.QueryRow(ctx, query, email, role).Scan(
		&user.ID,
		&user.Email,
		&user.Role,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &user, nil
}

func (r *authRepo) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User

	query := `
		SELECT id, email, role, created_at
		FROM users
		WHERE email = $1
	`

	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Role,
		&user.CreatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil // User not found
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return &user, nil
}

func (r *authRepo) GetUserByID(ctx context.Context, id int) (*domain.User, error) {
	var user domain.User

	query := `
		SELECT id, email, role, created_at
		FROM users
		WHERE id = $1
	`

	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.Role,
		&user.CreatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil // User not found
		}
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}

	return &user, nil
}

func (r *authRepo) LinkOAuthProvider(ctx context.Context, provider *domain.OAuthProvider) error {
	query := `
		INSERT INTO oauth_providers (user_id, provider, provider_user_id, access_token, refresh_token, expires_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
		ON CONFLICT (provider, provider_user_id)
		DO UPDATE SET
			access_token = EXCLUDED.access_token,
			refresh_token = EXCLUDED.refresh_token,
			expires_at = EXCLUDED.expires_at,
			updated_at = NOW()
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(ctx, query,
		provider.UserID,
		provider.Provider,
		provider.ProviderUserID,
		provider.AccessToken,
		provider.RefreshToken,
		provider.ExpiresAt,
	).Scan(&provider.ID, &provider.CreatedAt, &provider.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to link oauth provider: %w", err)
	}

	return nil
}

func (r *authRepo) GetUserByProviderID(ctx context.Context, providerName, providerUserID string) (*domain.User, error) {
	var user domain.User

	query := `
		SELECT u.id, u.email, u.role, u.created_at
		FROM users u
		INNER JOIN oauth_providers op ON u.id = op.user_id
		WHERE op.provider = $1 AND op.provider_user_id = $2
	`

	err := r.db.QueryRow(ctx, query, providerName, providerUserID).Scan(
		&user.ID,
		&user.Email,
		&user.Role,
		&user.CreatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil // User not found
		}
		return nil, fmt.Errorf("failed to get user by provider id: %w", err)
	}

	return &user, nil
}

func (r *authRepo) GetOAuthProvider(ctx context.Context, userID int, providerName string) (*domain.OAuthProvider, error) {
	var provider domain.OAuthProvider

	query := `
		SELECT id, user_id, provider, provider_user_id, access_token, refresh_token, expires_at, created_at, updated_at
		FROM oauth_providers
		WHERE user_id = $1 AND provider = $2
	`

	err := r.db.QueryRow(ctx, query, userID, providerName).Scan(
		&provider.ID,
		&provider.UserID,
		&provider.Provider,
		&provider.ProviderUserID,
		&provider.AccessToken,
		&provider.RefreshToken,
		&provider.ExpiresAt,
		&provider.CreatedAt,
		&provider.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil // Provider not found
		}
		return nil, fmt.Errorf("failed to get oauth provider: %w", err)
	}

	return &provider, nil
}

func (r *authRepo) UpdateOAuthProvider(ctx context.Context, provider *domain.OAuthProvider) error {
	query := `
		UPDATE oauth_providers
		SET access_token = $1, refresh_token = $2, expires_at = $3, updated_at = NOW()
		WHERE id = $4
	`

	_, err := r.db.Exec(ctx, query,
		provider.AccessToken,
		provider.RefreshToken,
		provider.ExpiresAt,
		provider.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update oauth provider: %w", err)
	}

	return nil
}
