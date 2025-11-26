// Package repository provides data access layer implementations
package repository

import (
	"context"
	"fmt"

	"personal-web-platform/config"
	"personal-web-platform/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

// NewPostgresDB creates a new PostgreSQL connection pool
func NewPostgresDB(ctx context.Context, url string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	return pool, nil
}

// Repositories aggregates all repository interfaces
type Repositories struct {
	Profile ProfileRepository
	Auth    AuthRepository
	Session SessionRepository
	db      *pgxpool.Pool
}

// NewRepositories creates a new Repositories instance with all implementations
func NewRepositories(db *pgxpool.Pool, cfg *config.Config) *Repositories {
	return &Repositories{
		Profile: NewProfileRepo(cfg),
		Auth:    NewAuthRepo(db),
		Session: NewSessionRepo(db),
		db:      db,
	}
}

// Ping checks database connectivity
func (r *Repositories) Ping(ctx context.Context) error {
	return r.db.Ping(ctx)
}

// ProfileRepository defines methods for profile data access
type ProfileRepository interface {
	Get(ctx context.Context) (domain.Profile, error)
}
