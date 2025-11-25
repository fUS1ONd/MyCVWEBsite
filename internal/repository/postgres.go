package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"personal-web-platform/config"
	"personal-web-platform/internal/domain"
)

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

type Repositories struct {
	Profile ProfileRepository
}

func NewRepositories(db *pgxpool.Pool, cfg *config.Config) *Repositories {
	return &Repositories{
		Profile: NewProfileRepo(cfg),
	}
}

type ProfileRepository interface {
	Get(ctx context.Context) (domain.Profile, error)
}
