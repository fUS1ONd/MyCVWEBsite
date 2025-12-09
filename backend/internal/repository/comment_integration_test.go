//go:build integration

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

func TestCommentRepository_Integration(t *testing.T) {
	// Setup test database
	testDB := testutil.SetupTestDatabase(t)
	defer testDB.Cleanup(t)

	authRepo := NewAuthRepo(testDB.Pool)
	postRepo := NewPostRepo(testDB.Pool)
	commentRepo := NewCommentRepo(testDB.Pool)
	ctx := context.Background()

	// Clean up tables at the start
	err := testDB.TruncateTables(ctx, "comments", "posts", "users")
	require.NoError(t, err)

	// Create test user and post
	user, err := authRepo.CreateUser(ctx, "commenter@example.com", "", "", domain.RoleUser)
	require.NoError(t, err)

	post := &domain.Post{
		Title:     "Post for Comments",
		Slug:      "post-for-comments",
		Content:   "This is a post to test comments functionality",
		Preview:   "Comment test preview",
		AuthorID:  user.ID,
		Published: true,
	}
	post, err = postRepo.Create(ctx, post)
	require.NoError(t, err)

	t.Run("Create and GetByID", func(t *testing.T) {
		comment := &domain.Comment{
			PostID:  post.ID,
			UserID:  user.ID,
			Content: "This is a test comment",
		}

		created, err := commentRepo.Create(ctx, comment)
		require.NoError(t, err)
		require.NotNil(t, created)
		assert.NotZero(t, created.ID)
		assert.Equal(t, post.ID, created.PostID)
		assert.Equal(t, user.ID, created.UserID)
		assert.Equal(t, "This is a test comment", created.Content)
		assert.NotZero(t, created.CreatedAt)

		// Get by ID
		retrieved, err := commentRepo.GetByID(ctx, created.ID)
		require.NoError(t, err)
		require.NotNil(t, retrieved)
		assert.Equal(t, created.ID, retrieved.ID)
		assert.Equal(t, created.Content, retrieved.Content)
	})

	t.Run("Create nested comment", func(t *testing.T) {
		// Create parent comment
		parent := &domain.Comment{
			PostID:  post.ID,
			UserID:  user.ID,
			Content: "Parent comment",
		}
		parent, err := commentRepo.Create(ctx, parent)
		require.NoError(t, err)

		// Create reply
		parentID := parent.ID
		reply := &domain.Comment{
			PostID:   post.ID,
			UserID:   user.ID,
			Content:  "Reply to parent",
			ParentID: &parentID,
		}
		reply, err = commentRepo.Create(ctx, reply)
		require.NoError(t, err)
		require.NotNil(t, reply.ParentID)
		assert.Equal(t, parent.ID, *reply.ParentID)
	})

	t.Run("Update comment", func(t *testing.T) {
		comment := &domain.Comment{
			PostID:  post.ID,
			UserID:  user.ID,
			Content: "Original content",
		}
		created, err := commentRepo.Create(ctx, comment)
		require.NoError(t, err)

		// Update
		created.Content = "Updated content"
		updated, err := commentRepo.Update(ctx, created)
		require.NoError(t, err)
		require.NotNil(t, updated)
		assert.Equal(t, "Updated content", updated.Content)
		assert.True(t, updated.UpdatedAt.After(updated.CreatedAt))
	})

	t.Run("Delete comment (soft delete)", func(t *testing.T) {
		comment := &domain.Comment{
			PostID:  post.ID,
			UserID:  user.ID,
			Content: "To be deleted",
		}
		created, err := commentRepo.Create(ctx, comment)
		require.NoError(t, err)

		// Delete
		err = commentRepo.Delete(ctx, created.ID)
		require.NoError(t, err)

		// Get deleted comment - should still exist but have DeletedAt set
		deleted, err := commentRepo.GetByID(ctx, created.ID)
		require.NoError(t, err)
		require.NotNil(t, deleted)
		assert.NotNil(t, deleted.DeletedAt)
		assert.True(t, deleted.DeletedAt.After(time.Time{}))
	})

	t.Run("GetByPostID", func(t *testing.T) {
		// Create multiple comments for the post
		for i := 1; i <= 5; i++ {
			comment := &domain.Comment{
				PostID:  post.ID,
				UserID:  user.ID,
				Content: "Comment number " + string(rune(i+'0')),
			}
			_, err := commentRepo.Create(ctx, comment)
			require.NoError(t, err)
		}

		// Get all comments for the post
		comments, err := commentRepo.GetByPostID(ctx, post.ID, 0)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(comments), 5)

		// Verify all comments belong to the post
		for _, comment := range comments {
			assert.Equal(t, post.ID, comment.PostID)
		}
	})

	t.Run("GetByPostID includes deleted comments with DeletedAt set", func(t *testing.T) {
		// Create comment and delete it
		comment := &domain.Comment{
			PostID:  post.ID,
			UserID:  user.ID,
			Content: "Will be deleted",
		}
		created, err := commentRepo.Create(ctx, comment)
		require.NoError(t, err)

		err = commentRepo.Delete(ctx, created.ID)
		require.NoError(t, err)

		// Get comments - repository returns all comments (filtering is done at service layer)
		comments, err := commentRepo.GetByPostID(ctx, post.ID, 0)
		require.NoError(t, err)

		// Find the deleted comment and verify DeletedAt is set
		foundDeleted := false
		for _, c := range comments {
			if c.ID == created.ID {
				foundDeleted = true
				assert.NotNil(t, c.DeletedAt, "Deleted comment should have DeletedAt set")
			}
		}
		assert.True(t, foundDeleted, "Deleted comment should be in the result set")
	})

	t.Run("Nested comments structure", func(t *testing.T) {
		// Create root comment
		root := &domain.Comment{
			PostID:  post.ID,
			UserID:  user.ID,
			Content: "Root comment",
		}
		root, err := commentRepo.Create(ctx, root)
		require.NoError(t, err)

		// Create first-level reply
		rootID := root.ID
		reply1 := &domain.Comment{
			PostID:   post.ID,
			UserID:   user.ID,
			Content:  "First level reply",
			ParentID: &rootID,
		}
		reply1, err = commentRepo.Create(ctx, reply1)
		require.NoError(t, err)

		// Create second-level reply
		reply1ID := reply1.ID
		reply2 := &domain.Comment{
			PostID:   post.ID,
			UserID:   user.ID,
			Content:  "Second level reply",
			ParentID: &reply1ID,
		}
		reply2, err = commentRepo.Create(ctx, reply2)
		require.NoError(t, err)

		// Verify structure
		assert.Nil(t, root.ParentID)
		assert.NotNil(t, reply1.ParentID)
		assert.Equal(t, root.ID, *reply1.ParentID)
		assert.NotNil(t, reply2.ParentID)
		assert.Equal(t, reply1.ID, *reply2.ParentID)
	})

	t.Run("GetByID returns nil for non-existent comment", func(t *testing.T) {
		comment, err := commentRepo.GetByID(ctx, 99999)
		require.NoError(t, err)
		assert.Nil(t, comment)
	})
}
