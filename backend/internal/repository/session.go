package repository

import (
	"context"
	"fmt"
	"time"

	"personal-web-platform/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// SessionRepository defines methods for session data access
type SessionRepository interface {
	CreateSession(ctx context.Context, session *domain.Session) error
	GetSession(ctx context.Context, token string) (*domain.Session, error)
	DeleteSession(ctx context.Context, token string) error
	DeleteUserSessions(ctx context.Context, userID int) error
	CleanupExpiredSessions(ctx context.Context) (int64, error)
}

type sessionRepo struct {
	db *pgxpool.Pool
}

// NewSessionRepo creates a new session repository implementation
func NewSessionRepo(db *pgxpool.Pool) SessionRepository {
	return &sessionRepo{db: db}
}

func (r *sessionRepo) CreateSession(ctx context.Context, session *domain.Session) error {
	db := GetQueryEngine(ctx, r.db)

	query := `
		INSERT INTO sessions (user_id, token, expires_at, created_at)
		VALUES ($1, $2, $3, NOW())
		RETURNING id, created_at
	`

	err := db.QueryRow(ctx, query,
		session.UserID,
		session.Token,
		session.ExpiresAt,
	).Scan(&session.ID, &session.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}

	return nil
}

func (r *sessionRepo) GetSession(ctx context.Context, token string) (*domain.Session, error) {
	var session domain.Session
	db := GetQueryEngine(ctx, r.db)

	query := `
		SELECT id, user_id, token, expires_at, created_at
		FROM sessions
		WHERE token = $1 AND expires_at > NOW()
	`

	err := db.QueryRow(ctx, query, token).Scan(
		&session.ID,
		&session.UserID,
		&session.Token,
		&session.ExpiresAt,
		&session.CreatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil // Session not found or expired
		}
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	return &session, nil
}

func (r *sessionRepo) DeleteSession(ctx context.Context, token string) error {
	db := GetQueryEngine(ctx, r.db)

	query := `DELETE FROM sessions WHERE token = $1`

	_, err := db.Exec(ctx, query, token)
	if err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}

	return nil
}

func (r *sessionRepo) DeleteUserSessions(ctx context.Context, userID int) error {
	db := GetQueryEngine(ctx, r.db)

	query := `DELETE FROM sessions WHERE user_id = $1`

	_, err := db.Exec(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to delete user sessions: %w", err)
	}

	return nil
}

func (r *sessionRepo) CleanupExpiredSessions(ctx context.Context) (int64, error) {
	db := GetQueryEngine(ctx, r.db)

	query := `DELETE FROM sessions WHERE expires_at < $1`

	result, err := db.Exec(ctx, query, time.Now())
	if err != nil {
		return 0, fmt.Errorf("failed to cleanup expired sessions: %w", err)
	}

	return result.RowsAffected(), nil
}
