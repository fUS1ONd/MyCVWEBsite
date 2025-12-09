// Package repository provides data access layer implementations
package repository

import (
	"context"
	"fmt"

	"personal-web-platform/config"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// txKey is a key for context storage of transaction
type txKey struct{}

// QueryEngine defines the interface for database operations (satisfied by *pgxpool.Pool and pgx.Tx)
type QueryEngine interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

// Transactor defines the interface for executing code within a transaction
type Transactor interface {
	RunInTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}

// GetQueryEngine returns the transaction from context if available, otherwise returns the pool
func GetQueryEngine(ctx context.Context, pool *pgxpool.Pool) QueryEngine {
	if tx, ok := ctx.Value(txKey{}).(pgx.Tx); ok {
		return tx
	}
	return pool
}

type txManager struct {
	pool *pgxpool.Pool
}

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

// RunInTransaction executes the given function within a database transaction
func (tm *txManager) RunInTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	// If already in transaction, just run the function
	if _, ok := ctx.Value(txKey{}).(pgx.Tx); ok {
		return fn(ctx)
	}

	tx, err := tm.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		// Rollback is safe to call even if transaction is already committed
		_ = tx.Rollback(ctx) //nolint:errcheck // rollback error is expected when already committed
	}()

	ctxWithTx := context.WithValue(ctx, txKey{}, tx)
	if err := fn(ctxWithTx); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// Repositories aggregates all repository interfaces
type Repositories struct {
	Profile    ProfileRepository
	Auth       AuthRepository
	Session    SessionRepository
	Post       PostRepository
	Comment    CommentRepository
	Media      MediaRepository
	Transactor Transactor
	db         *pgxpool.Pool
}

// NewRepositories creates a new Repositories instance with all implementations
func NewRepositories(db *pgxpool.Pool, _ *config.Config) *Repositories { //nolint:revive // cfg reserved for future use
	return &Repositories{
		Profile:    NewProfileRepo(db),
		Auth:       NewAuthRepo(db),
		Session:    NewSessionRepo(db),
		Post:       NewPostRepo(db),
		Comment:    NewCommentRepo(db),
		Media:      NewMediaRepository(db),
		Transactor: &txManager{pool: db},
		db:         db,
	}
}

// Ping checks database connectivity
func (r *Repositories) Ping(ctx context.Context) error {
	return r.db.Ping(ctx)
}
