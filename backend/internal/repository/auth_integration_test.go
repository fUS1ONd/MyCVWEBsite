//go:build integration
// +build integration

package repository

import (
	"context"
	"testing"
	"time"

	"personal-web-platform/internal/domain"
	"personal-web-platform/internal/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthRepository_Integration(t *testing.T) {
	// Setup test database
	testDB := testutil.SetupTestDatabase(t)
	defer testDB.Cleanup(t)

	repo := NewAuthRepo(testDB.Pool)
	sessionRepo := NewSessionRepo(testDB.Pool)
	ctx := context.Background()

	// Clean up tables at the start
	err := testDB.TruncateTables(ctx, "oauth_providers", "sessions", "users")
	require.NoError(t, err)

	t.Run("CreateUser and GetUserByID", func(t *testing.T) {
		user, err := repo.CreateUser(ctx, "test@example.com", domain.RoleUser)
		require.NoError(t, err)
		require.NotNil(t, user)
		assert.Equal(t, "test@example.com", user.Email)
		assert.Equal(t, domain.RoleUser, user.Role)
		assert.NotZero(t, user.ID)

		// Get user by ID
		retrieved, err := repo.GetUserByID(ctx, user.ID)
		require.NoError(t, err)
		require.NotNil(t, retrieved)
		assert.Equal(t, user.ID, retrieved.ID)
		assert.Equal(t, user.Email, retrieved.Email)
	})

	t.Run("GetUserByEmail", func(t *testing.T) {
		email := "findme@example.com"
		created, err := repo.CreateUser(ctx, email, domain.RoleUser)
		require.NoError(t, err)

		found, err := repo.GetUserByEmail(ctx, email)
		require.NoError(t, err)
		require.NotNil(t, found)
		assert.Equal(t, created.ID, found.ID)
		assert.Equal(t, email, found.Email)
	})

	t.Run("LinkOAuthProvider and GetUserByProviderID", func(t *testing.T) {
		// Create user first
		user, err := repo.CreateUser(ctx, "oauth@example.com", domain.RoleUser)
		require.NoError(t, err)

		// Link OAuth provider
		provider := &domain.OAuthProvider{
			UserID:         user.ID,
			Provider:       "google",
			ProviderUserID: "google-123",
			AccessToken:    "access-token-123",
			RefreshToken:   "refresh-token-123",
			ExpiresAt:      time.Now().Add(1 * time.Hour),
		}

		err = repo.LinkOAuthProvider(ctx, provider)
		require.NoError(t, err)

		// Get user by provider ID
		found, err := repo.GetUserByProviderID(ctx, "google", "google-123")
		require.NoError(t, err)
		require.NotNil(t, found)
		assert.Equal(t, user.ID, found.ID)
		assert.Equal(t, user.Email, found.Email)
	})

	t.Run("GetOAuthProvider and UpdateOAuthProvider", func(t *testing.T) {
		// Create user
		user, err := repo.CreateUser(ctx, "provider@example.com", domain.RoleUser)
		require.NoError(t, err)

		// Link provider
		provider := &domain.OAuthProvider{
			UserID:         user.ID,
			Provider:       "github",
			ProviderUserID: "github-456",
			AccessToken:    "old-token",
			RefreshToken:   "old-refresh",
			ExpiresAt:      time.Now().Add(1 * time.Hour),
		}
		err = repo.LinkOAuthProvider(ctx, provider)
		require.NoError(t, err)

		// Get provider
		retrieved, err := repo.GetOAuthProvider(ctx, user.ID, "github")
		require.NoError(t, err)
		require.NotNil(t, retrieved)
		assert.Equal(t, "old-token", retrieved.AccessToken)

		// Update provider
		retrieved.AccessToken = "new-token"
		retrieved.RefreshToken = "new-refresh"
		err = repo.UpdateOAuthProvider(ctx, retrieved)
		require.NoError(t, err)

		// Verify update
		updated, err := repo.GetOAuthProvider(ctx, user.ID, "github")
		require.NoError(t, err)
		assert.Equal(t, "new-token", updated.AccessToken)
		assert.Equal(t, "new-refresh", updated.RefreshToken)
	})

	t.Run("Session CRUD operations", func(t *testing.T) {
		// Create user
		user, err := repo.CreateUser(ctx, "session@example.com", domain.RoleUser)
		require.NoError(t, err)

		// Create session
		session := &domain.Session{
			UserID:    user.ID,
			Token:     "test-session-token",
			ExpiresAt: time.Now().Add(24 * time.Hour),
		}
		err = sessionRepo.CreateSession(ctx, session)
		require.NoError(t, err)

		// Get session
		retrieved, err := sessionRepo.GetSession(ctx, "test-session-token")
		require.NoError(t, err)
		require.NotNil(t, retrieved)
		assert.Equal(t, user.ID, retrieved.UserID)
		assert.Equal(t, "test-session-token", retrieved.Token)

		// Delete session
		err = sessionRepo.DeleteSession(ctx, "test-session-token")
		require.NoError(t, err)

		// Verify deletion
		deleted, err := sessionRepo.GetSession(ctx, "test-session-token")
		require.NoError(t, err)
		assert.Nil(t, deleted)
	})

	t.Run("CleanupExpiredSessions", func(t *testing.T) {
		// Create user
		user, err := repo.CreateUser(ctx, "cleanup@example.com", domain.RoleUser)
		require.NoError(t, err)

		// Create expired session
		expiredSession := &domain.Session{
			UserID:    user.ID,
			Token:     "expired-token",
			ExpiresAt: time.Now().Add(-1 * time.Hour),
		}
		err = sessionRepo.CreateSession(ctx, expiredSession)
		require.NoError(t, err)

		// Create valid session
		validSession := &domain.Session{
			UserID:    user.ID,
			Token:     "valid-token",
			ExpiresAt: time.Now().Add(24 * time.Hour),
		}
		err = sessionRepo.CreateSession(ctx, validSession)
		require.NoError(t, err)

		// Cleanup expired sessions
		count, err := sessionRepo.CleanupExpiredSessions(ctx)
		require.NoError(t, err)
		assert.Equal(t, int64(1), count)

		// Verify expired session is gone
		expired, err := sessionRepo.GetSession(ctx, "expired-token")
		require.NoError(t, err)
		assert.Nil(t, expired)

		// Verify valid session still exists
		valid, err := sessionRepo.GetSession(ctx, "valid-token")
		require.NoError(t, err)
		assert.NotNil(t, valid)
	})

	t.Run("DeleteUserSessions", func(t *testing.T) {
		// Create user
		user, err := repo.CreateUser(ctx, "multisession@example.com", domain.RoleUser)
		require.NoError(t, err)

		// Create multiple sessions
		for i := 1; i <= 3; i++ {
			session := &domain.Session{
				UserID:    user.ID,
				Token:     "token-" + string(rune(i+'0')),
				ExpiresAt: time.Now().Add(24 * time.Hour),
			}
			err = sessionRepo.CreateSession(ctx, session)
			require.NoError(t, err)
		}

		// Delete all user sessions
		err = sessionRepo.DeleteUserSessions(ctx, user.ID)
		require.NoError(t, err)

		// Verify all sessions are deleted
		for i := 1; i <= 3; i++ {
			session, err := sessionRepo.GetSession(ctx, "token-"+string(rune(i+'0')))
			require.NoError(t, err)
			assert.Nil(t, session)
		}
	})
}
