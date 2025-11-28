package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"personal-web-platform/config"
	"personal-web-platform/internal/domain"
	"personal-web-platform/internal/repository"

	"github.com/markbates/goth"
)

// AuthService defines methods for authentication business logic
type AuthService interface {
	LoginWithOAuth(ctx context.Context, gothUser goth.User) (*domain.User, *domain.Session, error)
	Logout(ctx context.Context, sessionToken string) error
	ValidateSession(ctx context.Context, sessionToken string) (*domain.User, error)
	GetUserByID(ctx context.Context, userID int) (*domain.User, error)
}

type authService struct {
	authRepo    repository.AuthRepository
	sessionRepo repository.SessionRepository
	cfg         *config.Config
}

// NewAuthService creates a new auth service implementation
func NewAuthService(authRepo repository.AuthRepository, sessionRepo repository.SessionRepository, cfg *config.Config) AuthService {
	return &authService{
		authRepo:    authRepo,
		sessionRepo: sessionRepo,
		cfg:         cfg,
	}
}

func (s *authService) LoginWithOAuth(ctx context.Context, gothUser goth.User) (*domain.User, *domain.Session, error) {
	// Try to find existing user by provider ID
	user, err := s.authRepo.GetUserByProviderID(ctx, gothUser.Provider, gothUser.UserID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get user by provider id: %w", err)
	}

	// If user doesn't exist, create new user
	if user == nil {
		// Validate that email is provided (required for VK ID)
		if gothUser.Email == "" {
			return nil, nil, fmt.Errorf("email permission required for authentication")
		}

		user, err = s.authRepo.CreateUser(ctx, gothUser.Email, domain.RoleUser)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to create user: %w", err)
		}
	}

	// Link or update OAuth provider
	oauthProvider := &domain.OAuthProvider{
		UserID:         user.ID,
		Provider:       gothUser.Provider,
		ProviderUserID: gothUser.UserID,
		AccessToken:    gothUser.AccessToken,
		RefreshToken:   gothUser.RefreshToken,
		ExpiresAt:      gothUser.ExpiresAt,
	}

	err = s.authRepo.LinkOAuthProvider(ctx, oauthProvider)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to link oauth provider: %w", err)
	}

	// Create session
	session, err := s.createSession(ctx, user.ID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create session: %w", err)
	}

	return user, session, nil
}

func (s *authService) Logout(ctx context.Context, sessionToken string) error {
	if err := s.sessionRepo.DeleteSession(ctx, sessionToken); err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}
	return nil
}

func (s *authService) ValidateSession(ctx context.Context, sessionToken string) (*domain.User, error) {
	// Get session
	session, err := s.sessionRepo.GetSession(ctx, sessionToken)
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	if session == nil {
		return nil, nil // Session not found or expired
	}

	// Get user
	user, err := s.authRepo.GetUserByID(ctx, session.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

func (s *authService) GetUserByID(ctx context.Context, userID int) (*domain.User, error) {
	user, err := s.authRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}
	return user, nil
}

func (s *authService) createSession(ctx context.Context, userID int) (*domain.Session, error) {
	// Generate random session token
	token, err := generateSecureToken(32)
	if err != nil {
		return nil, fmt.Errorf("failed to generate session token: %w", err)
	}

	session := &domain.Session{
		UserID:    userID,
		Token:     token,
		ExpiresAt: time.Now().Add(s.cfg.Auth.SessionMaxAge),
	}

	if err := s.sessionRepo.CreateSession(ctx, session); err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return session, nil
}

// generateSecureToken generates a cryptographically secure random token
func generateSecureToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
