//go:build integration
// +build integration

package repository

import (
	"context"
	"testing"

	"personal-web-platform/internal/domain"
	"personal-web-platform/internal/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProfileRepository_Integration(t *testing.T) {
	// Setup test database
	testDB := testutil.SetupTestDatabase(t)
	defer testDB.Cleanup(t)

	repo := NewProfileRepo(testDB.Pool)
	ctx := context.Background()

	t.Run("GetProfile - initial profile exists", func(t *testing.T) {
		profile, err := repo.GetProfile(ctx)
		require.NoError(t, err)
		require.NotNil(t, profile)
		assert.Equal(t, 1, profile.ID)
	})

	t.Run("UpdateProfile - success", func(t *testing.T) {
		req := &domain.UpdateProfileRequest{
			Name:        "John Doe",
			Description: "Senior Software Engineer",
			Activity:    "Building awesome applications",
			Contacts: domain.Contacts{
				Email:    "john.doe@example.com",
				GitHub:   "https://github.com/johndoe",
				LinkedIn: "https://linkedin.com/in/johndoe",
				VK:       "https://vk.com/johndoe",
			},
		}

		updated, err := repo.UpdateProfile(ctx, req)
		require.NoError(t, err)
		require.NotNil(t, updated)

		assert.Equal(t, "John Doe", updated.Name)
		assert.Equal(t, "Senior Software Engineer", updated.Description)
		assert.Equal(t, "Building awesome applications", updated.Activity)
		assert.Equal(t, "john.doe@example.com", updated.Contacts.Email)
		assert.Equal(t, "https://github.com/johndoe", updated.Contacts.GitHub)
	})

	t.Run("GetProfile - returns updated data", func(t *testing.T) {
		profile, err := repo.GetProfile(ctx)
		require.NoError(t, err)
		require.NotNil(t, profile)

		assert.Equal(t, "John Doe", profile.Name)
		assert.Equal(t, "Senior Software Engineer", profile.Description)
	})
}
