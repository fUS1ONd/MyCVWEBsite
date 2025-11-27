package service

import (
	"context"
	"testing"
	"time"

	"personal-web-platform/internal/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockCommentRepository is a mock implementation of CommentRepository
type MockCommentRepository struct {
	mock.Mock
}

func (m *MockCommentRepository) Create(ctx context.Context, comment *domain.Comment) (*domain.Comment, error) {
	args := m.Called(ctx, comment)
	if args.Get(0) == nil {
		return nil, args.Error(1) //nolint:errcheck // mock method
	}
	return args.Get(0).(*domain.Comment), args.Error(1) //nolint:errcheck // mock method
}

func (m *MockCommentRepository) Update(ctx context.Context, comment *domain.Comment) (*domain.Comment, error) {
	args := m.Called(ctx, comment)
	if args.Get(0) == nil {
		return nil, args.Error(1) //nolint:errcheck // mock method
	}
	return args.Get(0).(*domain.Comment), args.Error(1) //nolint:errcheck // mock method
}

func (m *MockCommentRepository) Delete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockCommentRepository) GetByID(ctx context.Context, id int) (*domain.Comment, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1) //nolint:errcheck // mock method
	}
	return args.Get(0).(*domain.Comment), args.Error(1) //nolint:errcheck // mock method
}

func (m *MockCommentRepository) GetByPostID(ctx context.Context, postID int) ([]domain.Comment, error) {
	args := m.Called(ctx, postID)
	return args.Get(0).([]domain.Comment), args.Error(1) //nolint:errcheck // mock method
}

func TestCommentService_CreateComment(t *testing.T) {
	tests := []struct {
		name             string
		postID           int
		request          *domain.CreateCommentRequest
		userID           int
		setupPostMock    func(*MockPostRepository)
		setupCommentMock func(*MockCommentRepository)
		wantErr          bool
		errContains      string
	}{
		{
			name:   "success - create root comment",
			postID: 1,
			request: &domain.CreateCommentRequest{
				Content: "Great post!",
			},
			userID: 1,
			setupPostMock: func(m *MockPostRepository) {
				m.On("GetByID", mock.Anything, 1).Return(&domain.Post{
					ID: 1,
				}, nil)
			},
			setupCommentMock: func(m *MockCommentRepository) {
				m.On("Create", mock.Anything, mock.Anything).Return(&domain.Comment{
					ID:      1,
					PostID:  1,
					UserID:  1,
					Content: "Great post!",
				}, nil)
			},
			wantErr: false,
		},
		{
			name:   "success - create reply comment",
			postID: 1,
			request: &domain.CreateCommentRequest{
				Content:  "I agree!",
				ParentID: intPtr(10),
			},
			userID: 2,
			setupPostMock: func(m *MockPostRepository) {
				m.On("GetByID", mock.Anything, 1).Return(&domain.Post{
					ID: 1,
				}, nil)
			},
			setupCommentMock: func(m *MockCommentRepository) {
				m.On("GetByID", mock.Anything, 10).Return(&domain.Comment{
					ID:     10,
					PostID: 1,
				}, nil)
				m.On("Create", mock.Anything, mock.Anything).Return(&domain.Comment{
					ID:       2,
					PostID:   1,
					UserID:   2,
					Content:  "I agree!",
					ParentID: intPtr(10),
				}, nil)
			},
			wantErr: false,
		},
		{
			name:   "error - post not found",
			postID: 999,
			request: &domain.CreateCommentRequest{
				Content: "Comment",
			},
			userID: 1,
			setupPostMock: func(m *MockPostRepository) {
				m.On("GetByID", mock.Anything, 999).Return(nil, nil)
			},
			setupCommentMock: func(m *MockCommentRepository) {},
			wantErr:          true,
			errContains:      "post not found",
		},
		{
			name:   "error - parent comment not found",
			postID: 1,
			request: &domain.CreateCommentRequest{
				Content:  "Reply",
				ParentID: intPtr(999),
			},
			userID: 1,
			setupPostMock: func(m *MockPostRepository) {
				m.On("GetByID", mock.Anything, 1).Return(&domain.Post{
					ID: 1,
				}, nil)
			},
			setupCommentMock: func(m *MockCommentRepository) {
				m.On("GetByID", mock.Anything, 999).Return(nil, nil)
			},
			wantErr:     true,
			errContains: "parent comment not found",
		},
		{
			name:   "error - parent comment deleted",
			postID: 1,
			request: &domain.CreateCommentRequest{
				Content:  "Reply",
				ParentID: intPtr(10),
			},
			userID: 1,
			setupPostMock: func(m *MockPostRepository) {
				m.On("GetByID", mock.Anything, 1).Return(&domain.Post{
					ID: 1,
				}, nil)
			},
			setupCommentMock: func(m *MockCommentRepository) {
				deletedAt := time.Now()
				m.On("GetByID", mock.Anything, 10).Return(&domain.Comment{
					ID:        10,
					PostID:    1,
					DeletedAt: &deletedAt,
				}, nil)
			},
			wantErr:     true,
			errContains: "cannot reply to deleted comment",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockPostRepo := new(MockPostRepository)
			mockCommentRepo := new(MockCommentRepository)
			tt.setupPostMock(mockPostRepo)
			tt.setupCommentMock(mockCommentRepo)

			service := NewCommentService(mockCommentRepo, mockPostRepo)
			comment, err := service.CreateComment(context.Background(), tt.postID, tt.request, tt.userID)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
				assert.Nil(t, comment)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, comment)
			}

			mockPostRepo.AssertExpectations(t)
			mockCommentRepo.AssertExpectations(t)
		})
	}
}

func TestCommentService_UpdateComment(t *testing.T) {
	tests := []struct {
		name        string
		commentID   int
		request     *domain.UpdateCommentRequest
		userID      int
		isAdmin     bool
		setupMock   func(*MockCommentRepository)
		wantErr     bool
		errContains string
	}{
		{
			name:      "success - author updates own comment",
			commentID: 1,
			request: &domain.UpdateCommentRequest{
				Content: "Updated comment",
			},
			userID:  1,
			isAdmin: false,
			setupMock: func(m *MockCommentRepository) {
				m.On("GetByID", mock.Anything, 1).Return(&domain.Comment{
					ID:      1,
					UserID:  1,
					Content: "Old content",
				}, nil)
				m.On("Update", mock.Anything, mock.Anything).Return(&domain.Comment{
					ID:      1,
					UserID:  1,
					Content: "Updated comment",
				}, nil)
			},
			wantErr: false,
		},
		{
			name:      "error - user tries to update others comment",
			commentID: 1,
			request: &domain.UpdateCommentRequest{
				Content: "Hacked",
			},
			userID:  999,
			isAdmin: false,
			setupMock: func(m *MockCommentRepository) {
				m.On("GetByID", mock.Anything, 1).Return(&domain.Comment{
					ID:     1,
					UserID: 1,
				}, nil)
			},
			wantErr:     true,
			errContains: "permission denied",
		},
		{
			name:      "error - cannot update deleted comment",
			commentID: 1,
			request: &domain.UpdateCommentRequest{
				Content: "Update",
			},
			userID:  1,
			isAdmin: false,
			setupMock: func(m *MockCommentRepository) {
				deletedAt := time.Now()
				m.On("GetByID", mock.Anything, 1).Return(&domain.Comment{
					ID:        1,
					UserID:    1,
					DeletedAt: &deletedAt,
				}, nil)
			},
			wantErr:     true,
			errContains: "cannot update deleted comment",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockCommentRepository)
			mockPostRepo := new(MockPostRepository)
			tt.setupMock(mockRepo)

			service := NewCommentService(mockRepo, mockPostRepo)
			comment, err := service.UpdateComment(context.Background(), tt.commentID, tt.request, tt.userID, tt.isAdmin)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
				assert.Nil(t, comment)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, comment)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestCommentService_DeleteComment(t *testing.T) {
	tests := []struct {
		name        string
		commentID   int
		userID      int
		isAdmin     bool
		setupMock   func(*MockCommentRepository)
		wantErr     bool
		errContains string
	}{
		{
			name:      "success - author deletes own comment",
			commentID: 1,
			userID:    1,
			isAdmin:   false,
			setupMock: func(m *MockCommentRepository) {
				m.On("GetByID", mock.Anything, 1).Return(&domain.Comment{
					ID:     1,
					UserID: 1,
				}, nil)
				m.On("Delete", mock.Anything, 1).Return(nil)
			},
			wantErr: false,
		},
		{
			name:      "success - admin deletes any comment",
			commentID: 1,
			userID:    999,
			isAdmin:   true,
			setupMock: func(m *MockCommentRepository) {
				m.On("GetByID", mock.Anything, 1).Return(&domain.Comment{
					ID:     1,
					UserID: 1,
				}, nil)
				m.On("Delete", mock.Anything, 1).Return(nil)
			},
			wantErr: false,
		},
		{
			name:      "error - non-author tries to delete",
			commentID: 1,
			userID:    999,
			isAdmin:   false,
			setupMock: func(m *MockCommentRepository) {
				m.On("GetByID", mock.Anything, 1).Return(&domain.Comment{
					ID:     1,
					UserID: 1,
				}, nil)
			},
			wantErr:     true,
			errContains: "permission denied",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockCommentRepository)
			mockPostRepo := new(MockPostRepository)
			tt.setupMock(mockRepo)

			service := NewCommentService(mockRepo, mockPostRepo)
			err := service.DeleteComment(context.Background(), tt.commentID, tt.userID, tt.isAdmin)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

// intPtr is a helper to create int pointer
func intPtr(i int) *int {
	return &i
}
