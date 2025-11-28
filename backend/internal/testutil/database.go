// Package testutil provides utilities for integration testing with testcontainers
package testutil

import (
	"context"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // Import postgres driver for migrations
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	testDBName = "testdb"
	testUser   = "testuser"
	testPass   = "testpass"
)

// TestDatabase holds the test database container and connection pool
type TestDatabase struct {
	Container *postgres.PostgresContainer
	Pool      *pgxpool.Pool
	ConnStr   string
}

// SetupTestDatabase creates a PostgreSQL container, runs migrations, and returns a connection pool
func SetupTestDatabase(t *testing.T) *TestDatabase {
	t.Helper()
	ctx := context.Background()

	// Create PostgreSQL container
	pgContainer, err := postgres.Run(ctx,
		"postgres:15-alpine",
		postgres.WithDatabase(testDBName),
		postgres.WithUsername(testUser),
		postgres.WithPassword(testPass),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(60*time.Second),
		),
	)
	if err != nil {
		t.Fatalf("Failed to start PostgreSQL container: %v", err)
	}

	// Get connection string
	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("Failed to get connection string: %v", err)
	}

	// Run migrations
	if err := runMigrations(connStr); err != nil {
		_ = pgContainer.Terminate(ctx) //nolint:errcheck,gosec // cleanup, error can be ignored
		t.Fatalf("Failed to run migrations: %v", err)
	}

	// Create connection pool
	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		_ = pgContainer.Terminate(ctx) //nolint:errcheck,gosec // cleanup, error can be ignored
		t.Fatalf("Failed to create connection pool: %v", err)
	}

	// Verify connection
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		_ = pgContainer.Terminate(ctx) //nolint:errcheck,gosec // cleanup, error can be ignored
		t.Fatalf("Failed to ping database: %v", err)
	}

	return &TestDatabase{
		Container: pgContainer,
		Pool:      pool,
		ConnStr:   connStr,
	}
}

// Cleanup closes the connection pool and terminates the container
func (td *TestDatabase) Cleanup(t *testing.T) {
	t.Helper()
	ctx := context.Background()

	if td.Pool != nil {
		td.Pool.Close()
	}

	if td.Container != nil {
		if err := td.Container.Terminate(ctx); err != nil {
			t.Logf("Failed to terminate container: %v", err)
		}
	}
}

// runMigrations runs database migrations
func runMigrations(connStr string) error {
	// Get the migrations directory path relative to the project root
	migrationsPath := filepath.Join("..", "..", "migrations")

	m, err := migrate.New(
		fmt.Sprintf("file://%s", migrationsPath),
		connStr,
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}
	defer m.Close() //nolint:errcheck // cleanup, error not important

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}

// TruncateTables truncates all tables for test isolation
func (td *TestDatabase) TruncateTables(ctx context.Context, tables ...string) error {
	for _, table := range tables {
		_, err := td.Pool.Exec(ctx, fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table))
		if err != nil {
			return fmt.Errorf("failed to truncate table %s: %w", table, err)
		}
	}
	return nil
}
