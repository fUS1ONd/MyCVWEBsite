package service

import (
	"context"
	"testing"
	"time"

	"personal-web-platform/internal/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockLikeRepository is a mock implementation of LikeRepository
type MockLikeRepository struct {
	mock.Mock
}

func (m *MockLikeRepository) TogglePostLike(ctx context.Context, userID, postID int) (bool, error) {
	args := m.Called(ctx, userID, postID)
	return args.Bool(0), args.Error(1)
}

func (m *MockLikeRepository) GetPostLikesCount(ctx context.Context, postID int) (int, error) {
	args := m.Called(ctx, postID)
	return args.Int(0), args.Error(1)
}

func (m *MockLikeRepository) IsPostLikedByUser(ctx context.Context, userID, postID int) (bool, error) {
	args := m.Called(ctx, userID, postID)
	return args.Bool(0), args.Error(1)
}

func (m *MockLikeRepository) ToggleCommentLike(ctx context.Context, userID, commentID int) (bool, error) {
	args := m.Called(ctx, userID, commentID)
	return args.Bool(0), args.Error(1)
}

func (m *MockLikeRepository) GetCommentLikesCount(ctx context.Context, commentID int) (int, error) {
	args := m.Called(ctx, commentID)
	return args.Int(0), args.Error(1)
}

func (m *MockLikeRepository) IsCommentLikedByUser(ctx context.Context, userID, commentID int) (bool, error) {
	args := m.Called(ctx, userID, commentID)
	return args.Bool(0), args.Error(1)
}

func TestLikeService_TogglePostLike(t *testing.T) {
	ctx := context.Background()
	mockLikeRepo := new(MockLikeRepository)
	mockPostRepo := new(MockPostRepository)
	mockCommentRepo := new(MockCommentRepository)

	service := NewLikeService(mockLikeRepo, mockPostRepo, mockCommentRepo)

	t.Run("successfully like post", func(t *testing.T) {
		post := &domain.Post{
			ID:        1,
			Title:     "Test Post",
			CreatedAt: time.Now(),
		}

		mockPostRepo.On("GetByID", ctx, 1, 0).Return(post, nil).Once()
		mockLikeRepo.On("TogglePostLike", ctx, 1, 1).Return(true, nil).Once()

		liked, err := service.TogglePostLike(ctx, 1, 1)

		assert.NoError(t, err)
		assert.True(t, liked)

		mockPostRepo.AssertExpectations(t)
		mockLikeRepo.AssertExpectations(t)
	})

	t.Run("successfully unlike post", func(t *testing.T) {
		post := &domain.Post{
			ID:        1,
			Title:     "Test Post",
			CreatedAt: time.Now(),
		}

		mockPostRepo.On("GetByID", ctx, 1, 0).Return(post, nil).Once()
		mockLikeRepo.On("TogglePostLike", ctx, 1, 1).Return(false, nil).Once()

		liked, err := service.TogglePostLike(ctx, 1, 1)

		assert.NoError(t, err)
		assert.False(t, liked)

		mockPostRepo.AssertExpectations(t)
		mockLikeRepo.AssertExpectations(t)
	})

	t.Run("post not found", func(t *testing.T) {
		mockPostRepo.On("GetByID", ctx, 999, 0).Return(nil, assert.AnError).Once()

		liked, err := service.TogglePostLike(ctx, 1, 999)

		assert.Error(t, err)
		assert.False(t, liked)

		mockPostRepo.AssertExpectations(t)
	})

	t.Run("repository error", func(t *testing.T) {
		post := &domain.Post{
			ID:        1,
			Title:     "Test Post",
			CreatedAt: time.Now(),
		}

		mockPostRepo.On("GetByID", ctx, 1, 0).Return(post, nil).Once()
		mockLikeRepo.On("TogglePostLike", ctx, 1, 1).Return(false, assert.AnError).Once()

		liked, err := service.TogglePostLike(ctx, 1, 1)

		assert.Error(t, err)
		assert.False(t, liked)

		mockPostRepo.AssertExpectations(t)
		mockLikeRepo.AssertExpectations(t)
	})
}

func TestLikeService_GetPostLikesCount(t *testing.T) {
	ctx := context.Background()
	mockLikeRepo := new(MockLikeRepository)
	mockPostRepo := new(MockPostRepository)
	mockCommentRepo := new(MockCommentRepository)

	service := NewLikeService(mockLikeRepo, mockPostRepo, mockCommentRepo)

	t.Run("successful count", func(t *testing.T) {
		mockLikeRepo.On("GetPostLikesCount", ctx, 1).Return(5, nil).Once()

		count, err := service.GetPostLikesCount(ctx, 1)

		assert.NoError(t, err)
		assert.Equal(t, 5, count)

		mockLikeRepo.AssertExpectations(t)
	})

	t.Run("zero likes", func(t *testing.T) {
		mockLikeRepo.On("GetPostLikesCount", ctx, 1).Return(0, nil).Once()

		count, err := service.GetPostLikesCount(ctx, 1)

		assert.NoError(t, err)
		assert.Equal(t, 0, count)

		mockLikeRepo.AssertExpectations(t)
	})

	t.Run("repository error", func(t *testing.T) {
		mockLikeRepo.On("GetPostLikesCount", ctx, 1).Return(0, assert.AnError).Once()

		count, err := service.GetPostLikesCount(ctx, 1)

		assert.Error(t, err)
		assert.Equal(t, 0, count)

		mockLikeRepo.AssertExpectations(t)
	})
}

func TestLikeService_IsPostLikedByUser(t *testing.T) { //nolint:dupl // Similar test structure is acceptable for different entities
	ctx := context.Background()
	mockLikeRepo := new(MockLikeRepository)
	mockPostRepo := new(MockPostRepository)
	mockCommentRepo := new(MockCommentRepository)

	service := NewLikeService(mockLikeRepo, mockPostRepo, mockCommentRepo)

	t.Run("post is liked", func(t *testing.T) {
		mockLikeRepo.On("IsPostLikedByUser", ctx, 1, 1).Return(true, nil).Once()

		liked, err := service.IsPostLikedByUser(ctx, 1, 1)

		assert.NoError(t, err)
		assert.True(t, liked)

		mockLikeRepo.AssertExpectations(t)
	})

	t.Run("post is not liked", func(t *testing.T) {
		mockLikeRepo.On("IsPostLikedByUser", ctx, 1, 1).Return(false, nil).Once()

		liked, err := service.IsPostLikedByUser(ctx, 1, 1)

		assert.NoError(t, err)
		assert.False(t, liked)

		mockLikeRepo.AssertExpectations(t)
	})
}

func TestLikeService_ToggleCommentLike(t *testing.T) {
	ctx := context.Background()
	mockLikeRepo := new(MockLikeRepository)
	mockPostRepo := new(MockPostRepository)
	mockCommentRepo := new(MockCommentRepository)

	service := NewLikeService(mockLikeRepo, mockPostRepo, mockCommentRepo)

	t.Run("successfully like comment", func(t *testing.T) {
		comment := &domain.Comment{
			ID:        1,
			PostID:    1,
			UserID:    1,
			Content:   "Test comment",
			CreatedAt: time.Now(),
		}

		mockCommentRepo.On("GetByID", ctx, 1).Return(comment, nil).Once()
		mockLikeRepo.On("ToggleCommentLike", ctx, 1, 1).Return(true, nil).Once()

		liked, err := service.ToggleCommentLike(ctx, 1, 1)

		assert.NoError(t, err)
		assert.True(t, liked)

		mockCommentRepo.AssertExpectations(t)
		mockLikeRepo.AssertExpectations(t)
	})

	t.Run("successfully unlike comment", func(t *testing.T) {
		comment := &domain.Comment{
			ID:        1,
			PostID:    1,
			UserID:    1,
			Content:   "Test comment",
			CreatedAt: time.Now(),
		}

		mockCommentRepo.On("GetByID", ctx, 1).Return(comment, nil).Once()
		mockLikeRepo.On("ToggleCommentLike", ctx, 1, 1).Return(false, nil).Once()

		liked, err := service.ToggleCommentLike(ctx, 1, 1)

		assert.NoError(t, err)
		assert.False(t, liked)

		mockCommentRepo.AssertExpectations(t)
		mockLikeRepo.AssertExpectations(t)
	})

	t.Run("comment not found", func(t *testing.T) {
		mockCommentRepo.On("GetByID", ctx, 999).Return(nil, assert.AnError).Once()

		liked, err := service.ToggleCommentLike(ctx, 1, 999)

		assert.Error(t, err)
		assert.False(t, liked)

		mockCommentRepo.AssertExpectations(t)
	})
}

func TestLikeService_GetCommentLikesCount(t *testing.T) {
	ctx := context.Background()
	mockLikeRepo := new(MockLikeRepository)
	mockPostRepo := new(MockPostRepository)
	mockCommentRepo := new(MockCommentRepository)

	service := NewLikeService(mockLikeRepo, mockPostRepo, mockCommentRepo)

	t.Run("successful count", func(t *testing.T) {
		mockLikeRepo.On("GetCommentLikesCount", ctx, 1).Return(3, nil).Once()

		count, err := service.GetCommentLikesCount(ctx, 1)

		assert.NoError(t, err)
		assert.Equal(t, 3, count)

		mockLikeRepo.AssertExpectations(t)
	})

	t.Run("zero likes", func(t *testing.T) {
		mockLikeRepo.On("GetCommentLikesCount", ctx, 1).Return(0, nil).Once()

		count, err := service.GetCommentLikesCount(ctx, 1)

		assert.NoError(t, err)
		assert.Equal(t, 0, count)

		mockLikeRepo.AssertExpectations(t)
	})
}

func TestLikeService_IsCommentLikedByUser(t *testing.T) { //nolint:dupl // Similar test structure is acceptable for different entities
	ctx := context.Background()
	mockLikeRepo := new(MockLikeRepository)
	mockPostRepo := new(MockPostRepository)
	mockCommentRepo := new(MockCommentRepository)

	service := NewLikeService(mockLikeRepo, mockPostRepo, mockCommentRepo)

	t.Run("comment is liked", func(t *testing.T) {
		mockLikeRepo.On("IsCommentLikedByUser", ctx, 1, 1).Return(true, nil).Once()

		liked, err := service.IsCommentLikedByUser(ctx, 1, 1)

		assert.NoError(t, err)
		assert.True(t, liked)

		mockLikeRepo.AssertExpectations(t)
	})

	t.Run("comment is not liked", func(t *testing.T) {
		mockLikeRepo.On("IsCommentLikedByUser", ctx, 1, 1).Return(false, nil).Once()

		liked, err := service.IsCommentLikedByUser(ctx, 1, 1)

		assert.NoError(t, err)
		assert.False(t, liked)

		mockLikeRepo.AssertExpectations(t)
	})
}
