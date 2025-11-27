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

func TestPostRepository_Integration(t *testing.T) {
	// Setup test database
	testDB := testutil.SetupTestDatabase(t)
	defer testDB.Cleanup(t)

	authRepo := NewAuthRepo(testDB.Pool)
	postRepo := NewPostRepo(testDB.Pool)
	ctx := context.Background()

	// Clean up tables at the start
	err := testDB.TruncateTables(ctx, "posts", "users")
	require.NoError(t, err)

	// Create test user
	author, err := authRepo.CreateUser(ctx, "author@example.com", domain.RoleAdmin)
	require.NoError(t, err)

	t.Run("Create and GetByID", func(t *testing.T) {
		post := &domain.Post{
			Title:     "Test Post",
			Slug:      "test-post",
			Content:   "This is test content with more than 10 characters",
			Preview:   "Test preview",
			AuthorID:  author.ID,
			Published: true,
		}

		created, err := postRepo.Create(ctx, post)
		require.NoError(t, err)
		require.NotNil(t, created)
		assert.NotZero(t, created.ID)
		assert.Equal(t, "Test Post", created.Title)
		assert.Equal(t, "test-post", created.Slug)
		assert.NotZero(t, created.CreatedAt)

		// Get by ID
		retrieved, err := postRepo.GetByID(ctx, created.ID)
		require.NoError(t, err)
		require.NotNil(t, retrieved)
		assert.Equal(t, created.ID, retrieved.ID)
		assert.Equal(t, created.Title, retrieved.Title)
		assert.Equal(t, created.Slug, retrieved.Slug)
	})

	t.Run("GetBySlug", func(t *testing.T) {
		post := &domain.Post{
			Title:     "Slug Test",
			Slug:      "unique-slug",
			Content:   "Content for slug test with more than 10 characters",
			Preview:   "Preview",
			AuthorID:  author.ID,
			Published: true,
		}

		created, err := postRepo.Create(ctx, post)
		require.NoError(t, err)

		// Get by slug
		retrieved, err := postRepo.GetBySlug(ctx, "unique-slug")
		require.NoError(t, err)
		require.NotNil(t, retrieved)
		assert.Equal(t, created.ID, retrieved.ID)
		assert.Equal(t, "unique-slug", retrieved.Slug)
	})

	t.Run("Update post", func(t *testing.T) {
		post := &domain.Post{
			Title:     "Original Title",
			Slug:      "original-slug",
			Content:   "Original content with more than 10 characters",
			Preview:   "Original preview",
			AuthorID:  author.ID,
			Published: true,
		}

		created, err := postRepo.Create(ctx, post)
		require.NoError(t, err)

		// Update
		created.Title = "Updated Title"
		created.Slug = "updated-slug"
		created.Content = "Updated content with more than 10 characters"

		updated, err := postRepo.Update(ctx, created)
		require.NoError(t, err)
		require.NotNil(t, updated)
		assert.Equal(t, "Updated Title", updated.Title)
		assert.Equal(t, "updated-slug", updated.Slug)
		assert.True(t, updated.Published)
		assert.True(t, updated.UpdatedAt.After(updated.CreatedAt))
	})

	t.Run("Delete post", func(t *testing.T) {
		post := &domain.Post{
			Title:     "To Delete",
			Slug:      "to-delete",
			Content:   "This post will be deleted, content is long enough",
			Preview:   "Delete preview",
			AuthorID:  author.ID,
			Published: true,
		}

		created, err := postRepo.Create(ctx, post)
		require.NoError(t, err)

		// Delete
		err = postRepo.Delete(ctx, created.ID)
		require.NoError(t, err)

		// Verify deletion
		deleted, err := postRepo.GetByID(ctx, created.ID)
		require.NoError(t, err)
		assert.Nil(t, deleted)
	})

	t.Run("List posts with pagination", func(t *testing.T) {
		// Clean up posts from previous subtests
		err := testDB.TruncateTables(ctx, "posts")
		require.NoError(t, err)

		// Create multiple posts (all published to avoid NULL published_at issues)
		for i := 1; i <= 15; i++ {
			post := &domain.Post{
				Title:     "Post " + string(rune(i+'0')),
				Slug:      "post-" + string(rune(i+'0')),
				Content:   "Content for post " + string(rune(i+'0')) + " with enough characters",
				Preview:   "Preview " + string(rune(i+'0')),
				AuthorID:  author.ID,
				Published: true,
			}
			_, err := postRepo.Create(ctx, post)
			require.NoError(t, err)
		}

		// List all posts (page 1, limit 10)
		req := &domain.ListPostsRequest{
			Page:  1,
			Limit: 10,
		}
		posts, total, err := postRepo.List(ctx, req)
		require.NoError(t, err)
		assert.Len(t, posts, 10)
		assert.Equal(t, 15, total)

		// List page 2
		req.Page = 2
		posts, total, err = postRepo.List(ctx, req)
		require.NoError(t, err)
		assert.Len(t, posts, 5)
		assert.Equal(t, 15, total)

		// List only published posts
		publishedOnly := true
		req = &domain.ListPostsRequest{
			Page:      1,
			Limit:     20,
			Published: &publishedOnly,
		}
		posts, total, err = postRepo.List(ctx, req)
		require.NoError(t, err)
		assert.Equal(t, 15, len(posts)) // All 15 posts are published
		assert.Equal(t, 15, total)
		for _, post := range posts {
			assert.True(t, post.Published)
		}
	})

	t.Run("GetByID returns nil for non-existent post", func(t *testing.T) {
		post, err := postRepo.GetByID(ctx, 99999)
		require.NoError(t, err)
		assert.Nil(t, post)
	})

	t.Run("GetBySlug returns nil for non-existent slug", func(t *testing.T) {
		post, err := postRepo.GetBySlug(ctx, "non-existent-slug")
		require.NoError(t, err)
		assert.Nil(t, post)
	})
}
