package service

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"testing"
	"time"

	"personal-web-platform/config"
	"personal-web-platform/internal/domain"

	"github.com/markbates/goth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAuthRepository is a mock implementation of AuthRepository
type MockAuthRepository struct {
	mock.Mock
}

func (m *MockAuthRepository) GetUserByProviderID(ctx context.Context, provider, providerUserID string) (*domain.User, error) {
	args := m.Called(ctx, provider, providerUserID)
	if args.Get(0) == nil {
		return nil, args.Error(1) //nolint:errcheck // mock method
	}
	return args.Get(0).(*domain.User), args.Error(1) //nolint:errcheck // mock method
}

func (m *MockAuthRepository) CreateUser(ctx context.Context, email string, name string, avatarURL string, role domain.Role) (*domain.User, error) {
	args := m.Called(ctx, email, name, avatarURL, role)
	if args.Get(0) == nil {
		return nil, args.Error(1) //nolint:errcheck // mock method
	}
	return args.Get(0).(*domain.User), args.Error(1) //nolint:errcheck // mock method
}

func (m *MockAuthRepository) GetUserByID(ctx context.Context, userID int) (*domain.User, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1) //nolint:errcheck // mock method
	}
	return args.Get(0).(*domain.User), args.Error(1) //nolint:errcheck // mock method
}

func (m *MockAuthRepository) LinkOAuthProvider(ctx context.Context, provider *domain.OAuthProvider) error {
	args := m.Called(ctx, provider)
	return args.Error(0)
}

func (m *MockAuthRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1) //nolint:errcheck // mock method
	}
	return args.Get(0).(*domain.User), args.Error(1) //nolint:errcheck // mock method
}

func (m *MockAuthRepository) GetOAuthProvider(ctx context.Context, userID int, providerName string) (*domain.OAuthProvider, error) {
	args := m.Called(ctx, userID, providerName)
	if args.Get(0) == nil {
		return nil, args.Error(1) //nolint:errcheck // mock method
	}
	return args.Get(0).(*domain.OAuthProvider), args.Error(1) //nolint:errcheck // mock method
}

func (m *MockAuthRepository) UpdateOAuthProvider(ctx context.Context, provider *domain.OAuthProvider) error {
	args := m.Called(ctx, provider)
	return args.Error(0)
}

// MockSessionRepository is a mock implementation of SessionRepository
type MockSessionRepository struct {
	mock.Mock
}

func (m *MockSessionRepository) CreateSession(ctx context.Context, session *domain.Session) error {
	args := m.Called(ctx, session)
	return args.Error(0)
}

func (m *MockSessionRepository) GetSession(ctx context.Context, token string) (*domain.Session, error) {
	args := m.Called(ctx, token)
	if args.Get(0) == nil {
		return nil, args.Error(1) //nolint:errcheck // mock method
	}
	return args.Get(0).(*domain.Session), args.Error(1) //nolint:errcheck // mock method
}

func (m *MockSessionRepository) DeleteSession(ctx context.Context, token string) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}

func (m *MockSessionRepository) CleanupExpiredSessions(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1) //nolint:errcheck // mock method
}

func (m *MockSessionRepository) DeleteUserSessions(ctx context.Context, userID int) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func getTestConfig() *config.Config {
	return &config.Config{
		Auth: config.Auth{
			SessionMaxAge: 24 * time.Hour,
		},
	}
}

func TestAuthService_LoginWithOAuth_ExistingUser(t *testing.T) {
	mockAuthRepo := new(MockAuthRepository)
	mockSessionRepo := new(MockSessionRepository)
	cfg := getTestConfig()

	gothUser := goth.User{
		Provider:     "google",
		UserID:       "google-user-123",
		Email:        "test@example.com",
		AccessToken:  "access-token",
		RefreshToken: "refresh-token",
		ExpiresAt:    time.Now().Add(1 * time.Hour),
	}

	existingUser := &domain.User{
		ID:    1,
		Email: "test@example.com",
		Role:  domain.RoleUser,
	}

	// Setup mocks
	mockAuthRepo.On("GetUserByProviderID", mock.Anything, "google", "google-user-123").
		Return(existingUser, nil)
	mockAuthRepo.On("LinkOAuthProvider", mock.Anything, mock.AnythingOfType("*domain.OAuthProvider")).
		Return(nil)
	mockSessionRepo.On("CreateSession", mock.Anything, mock.AnythingOfType("*domain.Session")).
		Return(nil)

	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
	service := NewAuthService(mockAuthRepo, mockSessionRepo, cfg, log)
	user, session, err := service.LoginWithOAuth(context.Background(), gothUser)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.NotNil(t, session)
	assert.Equal(t, existingUser.ID, user.ID)
	assert.Equal(t, existingUser.Email, user.Email)
	assert.NotEmpty(t, session.Token)

	mockAuthRepo.AssertExpectations(t)
	mockSessionRepo.AssertExpectations(t)
}

func TestAuthService_LoginWithOAuth_NewUser(t *testing.T) {
	mockAuthRepo := new(MockAuthRepository)
	mockSessionRepo := new(MockSessionRepository)
	cfg := getTestConfig()

	gothUser := goth.User{
		Provider:     "github",
		UserID:       "github-user-456",
		Email:        "newuser@example.com",
		AccessToken:  "access-token",
		RefreshToken: "refresh-token",
		ExpiresAt:    time.Now().Add(1 * time.Hour),
	}

	newUser := &domain.User{
		ID:    2,
		Email: "newuser@example.com",
		Role:  domain.RoleUser,
	}

	// Setup mocks - user not found, create new one
	mockAuthRepo.On("GetUserByProviderID", mock.Anything, "github", "github-user-456").
		Return(nil, nil)
	mockAuthRepo.On("CreateUser", mock.Anything, "newuser@example.com", "", "", domain.RoleUser).
		Return(newUser, nil)
	mockAuthRepo.On("LinkOAuthProvider", mock.Anything, mock.AnythingOfType("*domain.OAuthProvider")).
		Return(nil)
	mockSessionRepo.On("CreateSession", mock.Anything, mock.AnythingOfType("*domain.Session")).
		Return(nil)

	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
	service := NewAuthService(mockAuthRepo, mockSessionRepo, cfg, log)
	user, session, err := service.LoginWithOAuth(context.Background(), gothUser)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.NotNil(t, session)
	assert.Equal(t, newUser.ID, user.ID)
	assert.Equal(t, newUser.Email, user.Email)

	mockAuthRepo.AssertExpectations(t)
	mockSessionRepo.AssertExpectations(t)
}

func TestAuthService_LoginWithOAuth_CreateUserError(t *testing.T) {
	mockAuthRepo := new(MockAuthRepository)
	mockSessionRepo := new(MockSessionRepository)
	cfg := getTestConfig()

	gothUser := goth.User{
		Provider: "vk",
		UserID:   "vk-user-789",
		Email:    "error@example.com",
	}

	mockAuthRepo.On("GetUserByProviderID", mock.Anything, "vk", "vk-user-789").
		Return(nil, nil)
	mockAuthRepo.On("CreateUser", mock.Anything, "error@example.com", "", "", domain.RoleUser).
		Return(nil, errors.New("database error"))

	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
	service := NewAuthService(mockAuthRepo, mockSessionRepo, cfg, log)
	user, session, err := service.LoginWithOAuth(context.Background(), gothUser)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create user")
	assert.Nil(t, user)
	assert.Nil(t, session)

	mockAuthRepo.AssertExpectations(t)
}

func TestAuthService_ValidateSession_Success(t *testing.T) {
	mockAuthRepo := new(MockAuthRepository)
	mockSessionRepo := new(MockSessionRepository)
	cfg := getTestConfig()

	sessionToken := "valid-session-token"
	storedSession := &domain.Session{
		ID:        1,
		UserID:    1,
		Token:     sessionToken,
		ExpiresAt: time.Now().Add(1 * time.Hour),
	}

	user := &domain.User{
		ID:    1,
		Email: "test@example.com",
		Role:  domain.RoleUser,
	}

	mockSessionRepo.On("GetSession", mock.Anything, sessionToken).
		Return(storedSession, nil)
	mockAuthRepo.On("GetUserByID", mock.Anything, 1).
		Return(user, nil)

	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
	service := NewAuthService(mockAuthRepo, mockSessionRepo, cfg, log)
	result, err := service.ValidateSession(context.Background(), sessionToken)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, user.ID, result.ID)
	assert.Equal(t, user.Email, result.Email)

	mockAuthRepo.AssertExpectations(t)
	mockSessionRepo.AssertExpectations(t)
}

func TestAuthService_ValidateSession_NotFound(t *testing.T) {
	mockAuthRepo := new(MockAuthRepository)
	mockSessionRepo := new(MockSessionRepository)
	cfg := getTestConfig()

	sessionToken := "invalid-session-token"

	mockSessionRepo.On("GetSession", mock.Anything, sessionToken).
		Return(nil, nil)

	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
	service := NewAuthService(mockAuthRepo, mockSessionRepo, cfg, log)
	result, err := service.ValidateSession(context.Background(), sessionToken)

	assert.NoError(t, err)
	assert.Nil(t, result)

	mockSessionRepo.AssertExpectations(t)
}

func TestAuthService_Logout_Success(t *testing.T) {
	mockAuthRepo := new(MockAuthRepository)
	mockSessionRepo := new(MockSessionRepository)
	cfg := getTestConfig()

	sessionToken := "session-to-delete"

	mockSessionRepo.On("DeleteSession", mock.Anything, sessionToken).
		Return(nil)

	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
	service := NewAuthService(mockAuthRepo, mockSessionRepo, cfg, log)
	err := service.Logout(context.Background(), sessionToken)

	assert.NoError(t, err)
	mockSessionRepo.AssertExpectations(t)
}

func TestAuthService_Logout_Error(t *testing.T) {
	mockAuthRepo := new(MockAuthRepository)
	mockSessionRepo := new(MockSessionRepository)
	cfg := getTestConfig()

	sessionToken := "session-to-delete"

	mockSessionRepo.On("DeleteSession", mock.Anything, sessionToken).
		Return(errors.New("database error"))

	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
	service := NewAuthService(mockAuthRepo, mockSessionRepo, cfg, log)
	err := service.Logout(context.Background(), sessionToken)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to delete session")
	mockSessionRepo.AssertExpectations(t)
}

func TestAuthService_GetUserByID(t *testing.T) {
	tests := []struct {
		name      string
		userID    int
		setupMock func(*MockAuthRepository)
		wantErr   bool
	}{
		{
			name:   "success",
			userID: 1,
			setupMock: func(m *MockAuthRepository) {
				m.On("GetUserByID", mock.Anything, 1).Return(&domain.User{
					ID:    1,
					Email: "test@example.com",
					Role:  domain.RoleUser,
				}, nil)
			},
			wantErr: false,
		},
		{
			name:   "error - repository error",
			userID: 999,
			setupMock: func(m *MockAuthRepository) {
				m.On("GetUserByID", mock.Anything, 999).Return(nil, errors.New("user not found"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAuthRepo := new(MockAuthRepository)
			mockSessionRepo := new(MockSessionRepository)
			cfg := getTestConfig()

			tt.setupMock(mockAuthRepo)

			log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
			service := NewAuthService(mockAuthRepo, mockSessionRepo, cfg, log)
			user, err := service.GetUserByID(context.Background(), tt.userID)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, tt.userID, user.ID)
			}

			mockAuthRepo.AssertExpectations(t)
		})
	}
}
