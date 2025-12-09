package http

import (
	"context"
	"io"
	"log/slog"
	"mime/multipart"
	"testing"

	"personal-web-platform/config"
	"personal-web-platform/internal/domain"
	"personal-web-platform/internal/service"

	"github.com/markbates/goth"
	"github.com/stretchr/testify/mock"
)

// MockServices holds all mocked services
type MockServices struct {
	Post    *MockPostService
	Comment *MockCommentService
	Profile *MockProfileService
	Auth    *MockAuthService
	Media   *MockMediaService
	Like    *MockLikeService
}

// setupHandler creates a handler with mocked services
func setupHandler(_ *testing.T) (*Handler, *MockServices) { //nolint:revive // t is kept for consistency
	mocks := &MockServices{
		Post:    new(MockPostService),
		Comment: new(MockCommentService),
		Profile: new(MockProfileService),
		Auth:    new(MockAuthService),
		Media:   new(MockMediaService),
		Like:    new(MockLikeService),
	}

	services := &service.Services{
		Post:    mocks.Post,
		Comment: mocks.Comment,
		Profile: mocks.Profile,
		Auth:    mocks.Auth,
		Media:   mocks.Media,
		Like:    mocks.Like,
	}

	cfg := &config.Config{
		Auth: config.Auth{CookieName: "session_id"},
	}
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	return NewHandler(services, logger, cfg), mocks
}

// --- Service Mocks ---

type MockPostService struct {
	mock.Mock
}

func (m *MockPostService) CreatePost(ctx context.Context, req *domain.CreatePostRequest, authorID int) (*domain.Post, error) {
	args := m.Called(ctx, req, authorID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Post), args.Error(1)
}

func (m *MockPostService) UpdatePost(ctx context.Context, postID int, req *domain.UpdatePostRequest, userID int, isAdmin bool) (*domain.Post, error) {
	args := m.Called(ctx, postID, req, userID, isAdmin)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Post), args.Error(1)
}

func (m *MockPostService) DeletePost(ctx context.Context, postID int, userID int, isAdmin bool) error {
	args := m.Called(ctx, postID, userID, isAdmin)
	return args.Error(0)
}

func (m *MockPostService) GetPostByID(ctx context.Context, id, userID int) (*domain.Post, error) {
	args := m.Called(ctx, id, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Post), args.Error(1)
}

func (m *MockPostService) GetPostBySlug(ctx context.Context, slug string, userID int) (*domain.Post, error) {
	args := m.Called(ctx, slug, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Post), args.Error(1)
}

func (m *MockPostService) ListPosts(ctx context.Context, req *domain.ListPostsRequest) (*domain.PostsListResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.PostsListResponse), args.Error(1)
}

type MockCommentService struct {
	mock.Mock
}

func (m *MockCommentService) CreateComment(ctx context.Context, postID int, req *domain.CreateCommentRequest, userID int) (*domain.Comment, error) {
	args := m.Called(ctx, postID, req, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Comment), args.Error(1)
}

func (m *MockCommentService) UpdateComment(ctx context.Context, commentID int, req *domain.UpdateCommentRequest, userID int, isAdmin bool) (*domain.Comment, error) {
	args := m.Called(ctx, commentID, req, userID, isAdmin)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Comment), args.Error(1)
}

func (m *MockCommentService) DeleteComment(ctx context.Context, commentID int, userID int, isAdmin bool) error {
	args := m.Called(ctx, commentID, userID, isAdmin)
	return args.Error(0)
}

func (m *MockCommentService) GetCommentByID(ctx context.Context, id int) (*domain.Comment, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Comment), args.Error(1)
}

func (m *MockCommentService) GetCommentsByPostSlug(ctx context.Context, slug string, userID int) ([]domain.Comment, error) {
	args := m.Called(ctx, slug, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Comment), args.Error(1)
}

type MockProfileService struct {
	mock.Mock
}

func (m *MockProfileService) GetProfile(ctx context.Context) (*domain.Profile, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Profile), args.Error(1)
}

func (m *MockProfileService) UpdateProfile(ctx context.Context, req *domain.UpdateProfileRequest) (*domain.Profile, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Profile), args.Error(1)
}

type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) LoginWithOAuth(ctx context.Context, gothUser goth.User) (*domain.User, *domain.Session, error) {
	args := m.Called(ctx, gothUser)
	if args.Get(0) == nil {
		return nil, nil, args.Error(2)
	}
	return args.Get(0).(*domain.User), args.Get(1).(*domain.Session), args.Error(2)
}

func (m *MockAuthService) Logout(ctx context.Context, sessionToken string) error {
	args := m.Called(ctx, sessionToken)
	return args.Error(0)
}

func (m *MockAuthService) ValidateSession(ctx context.Context, sessionToken string) (*domain.User, error) {
	args := m.Called(ctx, sessionToken)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockAuthService) GetUserByID(ctx context.Context, userID int) (*domain.User, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

type MockMediaService struct {
	mock.Mock
}

func (m *MockMediaService) Upload(ctx context.Context, file multipart.File, header *multipart.FileHeader, uploaderID int) (*domain.MediaFile, error) {
	args := m.Called(ctx, file, header, uploaderID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.MediaFile), args.Error(1)
}

func (m *MockMediaService) GetByID(ctx context.Context, id int) (*domain.MediaFile, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.MediaFile), args.Error(1)
}

func (m *MockMediaService) GetFileReader(ctx context.Context, filename string) (io.ReadSeekCloser, string, error) {
	args := m.Called(ctx, filename)
	if args.Get(0) == nil {
		return nil, "", args.Error(2)
	}
	return args.Get(0).(io.ReadSeekCloser), args.String(1), args.Error(2)
}

func (m *MockMediaService) Delete(ctx context.Context, id int, uploaderID int) error {
	args := m.Called(ctx, id, uploaderID)
	return args.Error(0)
}

func (m *MockMediaService) ListByUploader(ctx context.Context, uploaderID int) ([]domain.MediaFile, error) {
	args := m.Called(ctx, uploaderID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.MediaFile), args.Error(1)
}

type MockLikeService struct {
	mock.Mock
}

func (m *MockLikeService) TogglePostLike(ctx context.Context, userID, postID int) (bool, error) {
	args := m.Called(ctx, userID, postID)
	return args.Bool(0), args.Error(1)
}

func (m *MockLikeService) GetPostLikesCount(ctx context.Context, postID int) (int, error) {
	args := m.Called(ctx, postID)
	return args.Int(0), args.Error(1)
}

func (m *MockLikeService) IsPostLikedByUser(ctx context.Context, userID, postID int) (bool, error) {
	args := m.Called(ctx, userID, postID)
	return args.Bool(0), args.Error(1)
}

func (m *MockLikeService) ToggleCommentLike(ctx context.Context, userID, commentID int) (bool, error) {
	args := m.Called(ctx, userID, commentID)
	return args.Bool(0), args.Error(1)
}

func (m *MockLikeService) GetCommentLikesCount(ctx context.Context, commentID int) (int, error) {
	args := m.Called(ctx, commentID)
	return args.Int(0), args.Error(1)
}

func (m *MockLikeService) IsCommentLikedByUser(ctx context.Context, userID, commentID int) (bool, error) {
	args := m.Called(ctx, userID, commentID)
	return args.Bool(0), args.Error(1)
}
