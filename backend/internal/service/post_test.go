package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"personal-web-platform/internal/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockPostRepository is a mock implementation of PostRepository
type MockPostRepository struct {
	mock.Mock
}

func (m *MockPostRepository) Create(ctx context.Context, post *domain.Post) (*domain.Post, error) {
	args := m.Called(ctx, post)
	if args.Get(0) == nil {
		return nil, args.Error(1) //nolint:errcheck // mock method
	}
	return args.Get(0).(*domain.Post), args.Error(1) //nolint:errcheck // mock method
}

func (m *MockPostRepository) Update(ctx context.Context, post *domain.Post) (*domain.Post, error) {
	args := m.Called(ctx, post)
	if args.Get(0) == nil {
		return nil, args.Error(1) //nolint:errcheck // mock method
	}
	return args.Get(0).(*domain.Post), args.Error(1) //nolint:errcheck // mock method
}

func (m *MockPostRepository) Delete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockPostRepository) GetByID(ctx context.Context, id int) (*domain.Post, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1) //nolint:errcheck // mock method
	}
	return args.Get(0).(*domain.Post), args.Error(1) //nolint:errcheck // mock method
}

func (m *MockPostRepository) GetBySlug(ctx context.Context, slug string) (*domain.Post, error) {
	args := m.Called(ctx, slug)
	if args.Get(0) == nil {
		return nil, args.Error(1) //nolint:errcheck // mock method
	}
	return args.Get(0).(*domain.Post), args.Error(1) //nolint:errcheck // mock method
}

func (m *MockPostRepository) List(ctx context.Context, req *domain.ListPostsRequest) ([]domain.Post, int, error) {
	args := m.Called(ctx, req)
	return args.Get(0).([]domain.Post), args.Get(1).(int), args.Error(2) //nolint:errcheck // mock method
}

func TestPostService_CreatePost(t *testing.T) {
	tests := []struct {
		name        string
		request     *domain.CreatePostRequest
		authorID    int
		setupMock   func(*MockPostRepository)
		wantErr     bool
		errContains string
	}{
		{
			name: "success - create published post",
			request: &domain.CreatePostRequest{
				Title:     "Introduction to Go",
				Content:   "# Go is awesome\n\nLet me tell you why...",
				Preview:   "Learn about Go programming language",
				Published: true,
			},
			authorID: 1,
			setupMock: func(m *MockPostRepository) {
				// Check slug doesn't exist
				m.On("GetBySlug", mock.Anything, "introduction-to-go").Return(nil, nil)
				// Create post
				m.On("Create", mock.Anything, mock.MatchedBy(func(p *domain.Post) bool {
					return p.Title == "Introduction to Go" && p.Slug == "introduction-to-go"
				})).Return(&domain.Post{
					ID:        1,
					Title:     "Introduction to Go",
					Slug:      "introduction-to-go",
					Content:   "# Go is awesome\n\nLet me tell you why...",
					Preview:   "Learn about Go programming language",
					AuthorID:  1,
					Published: true,
					CreatedAt: time.Now(),
				}, nil)
			},
			wantErr: false,
		},
		{
			name: "error - validation failed (empty title)",
			request: &domain.CreatePostRequest{
				Title:     "",
				Content:   "Content",
				Published: false,
			},
			authorID:    1,
			setupMock:   func(_ *MockPostRepository) {},
			wantErr:     true,
			errContains: "validation failed",
		},
		{
			name: "error - slug already exists",
			request: &domain.CreatePostRequest{
				Title:     "Existing Post",
				Content:   "This is content with more than 10 characters",
				Published: false,
			},
			authorID: 1,
			setupMock: func(m *MockPostRepository) {
				m.On("GetBySlug", mock.Anything, "existing-post").Return(&domain.Post{
					ID:   999,
					Slug: "existing-post",
				}, nil)
			},
			wantErr:     true,
			errContains: "post with this slug already exists",
		},
		{
			name: "error - repository error on create",
			request: &domain.CreatePostRequest{
				Title:     "New Post",
				Content:   "This is content with more than 10 characters",
				Published: false,
			},
			authorID: 1,
			setupMock: func(m *MockPostRepository) {
				m.On("GetBySlug", mock.Anything, "new-post").Return(nil, nil)
				m.On("Create", mock.Anything, mock.Anything).Return(nil, errors.New("database error"))
			},
			wantErr:     true,
			errContains: "failed to create post",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockPostRepository)
			tt.setupMock(mockRepo)

			service := NewPostService(mockRepo)
			post, err := service.CreatePost(context.Background(), tt.request, tt.authorID)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
				assert.Nil(t, post)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, post)
				assert.Equal(t, tt.request.Title, post.Title)
				assert.Equal(t, "introduction-to-go", post.Slug)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestPostService_UpdatePost(t *testing.T) {
	tests := []struct {
		name        string
		postID      int
		request     *domain.UpdatePostRequest
		userID      int
		isAdmin     bool
		setupMock   func(*MockPostRepository)
		wantErr     bool
		errContains string
	}{
		{
			name:   "success - author updates own post",
			postID: 1,
			request: &domain.UpdatePostRequest{
				Title:     "Updated Title",
				Content:   "Updated content",
				Preview:   "Updated preview",
				Published: true,
			},
			userID:  1,
			isAdmin: false,
			setupMock: func(m *MockPostRepository) {
				m.On("GetByID", mock.Anything, 1).Return(&domain.Post{
					ID:       1,
					Title:    "Old Title",
					Slug:     "old-title",
					AuthorID: 1,
				}, nil)
				m.On("GetBySlug", mock.Anything, "updated-title").Return(nil, nil)
				m.On("Update", mock.Anything, mock.Anything).Return(&domain.Post{
					ID:        1,
					Title:     "Updated Title",
					Slug:      "updated-title",
					Content:   "Updated content",
					Preview:   "Updated preview",
					Published: true,
				}, nil)
			},
			wantErr: false,
		},
		{
			name:   "success - admin updates any post",
			postID: 1,
			request: &domain.UpdatePostRequest{
				Title:     "Admin Update",
				Content:   "This is content with more than 10 characters",
				Preview:   "Preview",
				Published: true,
			},
			userID:  999,
			isAdmin: true,
			setupMock: func(m *MockPostRepository) {
				m.On("GetByID", mock.Anything, 1).Return(&domain.Post{
					ID:       1,
					Slug:     "old-slug",
					AuthorID: 1,
				}, nil)
				m.On("GetBySlug", mock.Anything, "admin-update").Return(nil, nil)
				m.On("Update", mock.Anything, mock.Anything).Return(&domain.Post{
					ID:    1,
					Title: "Admin Update",
					Slug:  "admin-update",
				}, nil)
			},
			wantErr: false,
		},
		{
			name:   "error - non-author tries to update",
			postID: 1,
			request: &domain.UpdatePostRequest{
				Title:   "Hacked",
				Content: "This is content with more than 10 characters",
				Preview: "Preview",
			},
			userID:  999,
			isAdmin: false,
			setupMock: func(m *MockPostRepository) {
				m.On("GetByID", mock.Anything, 1).Return(&domain.Post{
					ID:       1,
					AuthorID: 1,
				}, nil)
			},
			wantErr:     true,
			errContains: "permission denied",
		},
		{
			name:   "error - post not found",
			postID: 999,
			request: &domain.UpdatePostRequest{
				Title:   "Title",
				Content: "This is content with more than 10 characters",
				Preview: "Preview",
			},
			userID:  1,
			isAdmin: false,
			setupMock: func(m *MockPostRepository) {
				m.On("GetByID", mock.Anything, 999).Return(nil, nil)
			},
			wantErr:     true,
			errContains: "post not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockPostRepository)
			tt.setupMock(mockRepo)

			service := NewPostService(mockRepo)
			post, err := service.UpdatePost(context.Background(), tt.postID, tt.request, tt.userID, tt.isAdmin)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
				assert.Nil(t, post)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, post)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestPostService_DeletePost(t *testing.T) {
	tests := []struct {
		name        string
		postID      int
		userID      int
		isAdmin     bool
		setupMock   func(*MockPostRepository)
		wantErr     bool
		errContains string
	}{
		{
			name:    "success - author deletes own post",
			postID:  1,
			userID:  1,
			isAdmin: false,
			setupMock: func(m *MockPostRepository) {
				m.On("GetByID", mock.Anything, 1).Return(&domain.Post{
					ID:       1,
					AuthorID: 1,
				}, nil)
				m.On("Delete", mock.Anything, 1).Return(nil)
			},
			wantErr: false,
		},
		{
			name:    "success - admin deletes any post",
			postID:  1,
			userID:  999,
			isAdmin: true,
			setupMock: func(m *MockPostRepository) {
				m.On("GetByID", mock.Anything, 1).Return(&domain.Post{
					ID:       1,
					AuthorID: 1,
				}, nil)
				m.On("Delete", mock.Anything, 1).Return(nil)
			},
			wantErr: false,
		},
		{
			name:    "error - non-author tries to delete",
			postID:  1,
			userID:  999,
			isAdmin: false,
			setupMock: func(m *MockPostRepository) {
				m.On("GetByID", mock.Anything, 1).Return(&domain.Post{
					ID:       1,
					AuthorID: 1,
				}, nil)
			},
			wantErr:     true,
			errContains: "permission denied",
		},
		{
			name:    "error - post not found",
			postID:  999,
			userID:  1,
			isAdmin: false,
			setupMock: func(m *MockPostRepository) {
				m.On("GetByID", mock.Anything, 999).Return(nil, nil)
			},
			wantErr:     true,
			errContains: "post not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockPostRepository)
			tt.setupMock(mockRepo)

			service := NewPostService(mockRepo)
			err := service.DeletePost(context.Background(), tt.postID, tt.userID, tt.isAdmin)

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

func TestPostService_ListPosts(t *testing.T) {
	tests := []struct {
		name        string
		request     *domain.ListPostsRequest
		setupMock   func(*MockPostRepository)
		wantErr     bool
		checkResult func(*testing.T, *domain.PostsListResponse)
	}{
		{
			name: "success - default pagination",
			request: &domain.ListPostsRequest{
				Page:  1,
				Limit: 10,
			},
			setupMock: func(m *MockPostRepository) {
				m.On("List", mock.Anything, mock.Anything).Return([]domain.Post{
					{ID: 1, Title: "Post 1", Slug: "post-1"},
					{ID: 2, Title: "Post 2", Slug: "post-2"},
				}, 2, nil)
			},
			wantErr: false,
			checkResult: func(t *testing.T, result *domain.PostsListResponse) {
				assert.Len(t, result.Posts, 2)
				assert.Equal(t, 2, result.TotalCount)
				assert.Equal(t, 1, result.TotalPages)
			},
		},
		{
			name: "success - custom pagination values",
			request: &domain.ListPostsRequest{
				Page:  3,
				Limit: 20,
			},
			setupMock: func(m *MockPostRepository) {
				m.On("List", mock.Anything, mock.MatchedBy(func(req *domain.ListPostsRequest) bool {
					return req.Page == 3 && req.Limit == 20
				})).Return([]domain.Post{
					{ID: 41, Title: "Post 41", Slug: "post-41"},
				}, 50, nil)
			},
			wantErr: false,
			checkResult: func(t *testing.T, result *domain.PostsListResponse) {
				assert.Equal(t, 3, result.Page)
				assert.Equal(t, 20, result.Limit)
				assert.Equal(t, 50, result.TotalCount)
				assert.Equal(t, 3, result.TotalPages) // 50/20 = 3 pages
			},
		},
		{
			name: "success - calculate total pages correctly",
			request: &domain.ListPostsRequest{
				Page:  2,
				Limit: 5,
			},
			setupMock: func(m *MockPostRepository) {
				m.On("List", mock.Anything, mock.Anything).Return([]domain.Post{}, 23, nil)
			},
			wantErr: false,
			checkResult: func(t *testing.T, result *domain.PostsListResponse) {
				assert.Equal(t, 23, result.TotalCount)
				assert.Equal(t, 5, result.TotalPages) // 23/5 = 5 pages
			},
		},
		{
			name: "error - repository error",
			request: &domain.ListPostsRequest{
				Page:  1,
				Limit: 10,
			},
			setupMock: func(m *MockPostRepository) {
				m.On("List", mock.Anything, mock.Anything).Return([]domain.Post{}, 0, errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockPostRepository)
			tt.setupMock(mockRepo)

			service := NewPostService(mockRepo)
			result, err := service.ListPosts(context.Background(), tt.request)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				if tt.checkResult != nil {
					tt.checkResult(t, result)
				}
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestPostService_GetPostByID(t *testing.T) { //nolint:dupl // similar test pattern to GetPostBySlug
	tests := []struct {
		name        string
		postID      int
		setupMock   func(*MockPostRepository)
		wantErr     bool
		errContains string
	}{
		{
			name:   "success - post found",
			postID: 1,
			setupMock: func(m *MockPostRepository) {
				m.On("GetByID", mock.Anything, 1).Return(&domain.Post{
					ID:    1,
					Title: "Test Post",
					Slug:  "test-post",
				}, nil)
			},
			wantErr: false,
		},
		{
			name:   "error - post not found",
			postID: 999,
			setupMock: func(m *MockPostRepository) {
				m.On("GetByID", mock.Anything, 999).Return(nil, nil)
			},
			wantErr:     true,
			errContains: "post not found",
		},
		{
			name:   "error - repository error",
			postID: 1,
			setupMock: func(m *MockPostRepository) {
				m.On("GetByID", mock.Anything, 1).Return(nil, errors.New("database connection failed"))
			},
			wantErr:     true,
			errContains: "failed to get post",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockPostRepository)
			tt.setupMock(mockRepo)

			service := NewPostService(mockRepo)
			post, err := service.GetPostByID(context.Background(), tt.postID)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
				assert.Nil(t, post)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, post)
				assert.Equal(t, tt.postID, post.ID)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestPostService_GetPostBySlug(t *testing.T) { //nolint:dupl // similar test pattern to GetPostByID
	tests := []struct {
		name        string
		slug        string
		setupMock   func(*MockPostRepository)
		wantErr     bool
		errContains string
	}{
		{
			name: "success - post found by slug",
			slug: "my-awesome-post",
			setupMock: func(m *MockPostRepository) {
				m.On("GetBySlug", mock.Anything, "my-awesome-post").Return(&domain.Post{
					ID:    5,
					Title: "My Awesome Post",
					Slug:  "my-awesome-post",
				}, nil)
			},
			wantErr: false,
		},
		{
			name: "error - post not found",
			slug: "non-existent-slug",
			setupMock: func(m *MockPostRepository) {
				m.On("GetBySlug", mock.Anything, "non-existent-slug").Return(nil, nil)
			},
			wantErr:     true,
			errContains: "post not found",
		},
		{
			name: "error - repository error",
			slug: "test-slug",
			setupMock: func(m *MockPostRepository) {
				m.On("GetBySlug", mock.Anything, "test-slug").Return(nil, errors.New("db connection lost"))
			},
			wantErr:     true,
			errContains: "failed to get post",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockPostRepository)
			tt.setupMock(mockRepo)

			service := NewPostService(mockRepo)
			post, err := service.GetPostBySlug(context.Background(), tt.slug)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
				assert.Nil(t, post)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, post)
				assert.Equal(t, tt.slug, post.Slug)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
