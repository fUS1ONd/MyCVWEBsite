package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"personal-web-platform/config"
	"personal-web-platform/internal/domain"
	"personal-web-platform/internal/repository"

	"github.com/jackc/pgx/v5/pgconn"
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
	transactor  repository.Transactor
	cfg         *config.Config
	log         *slog.Logger
}

// NewAuthService creates a new auth service implementation
func NewAuthService(authRepo repository.AuthRepository, sessionRepo repository.SessionRepository, transactor repository.Transactor, cfg *config.Config, log *slog.Logger) AuthService {
	return &authService{
		authRepo:    authRepo,
		sessionRepo: sessionRepo,
		transactor:  transactor,
		cfg:         cfg,
		log:         log,
	}
}

func (s *authService) LoginWithOAuth(ctx context.Context, gothUser goth.User) (*domain.User, *domain.Session, error) {
	s.log.Info("auth_service: processing OAuth login",
		"provider", gothUser.Provider,
		"provider_user_id", gothUser.UserID,
		"email", gothUser.Email,
	)

	var user *domain.User

	// Execute user lookup/creation + provider linking in transaction
	err := s.transactor.RunInTransaction(ctx, func(txCtx context.Context) error {
		// Step 1: Try to find by provider ID first
		foundUser, err := s.authRepo.GetUserByProviderID(txCtx, gothUser.Provider, gothUser.UserID)
		if err != nil {
			s.log.Error("auth_service: failed to get user by provider id",
				"error", err,
				"provider", gothUser.Provider,
				"provider_user_id", gothUser.UserID,
			)
			return fmt.Errorf("failed to get user by provider id: %w", err)
		}

		// Step 2: If not found by provider, try to find by email
		if foundUser == nil && gothUser.Email != "" {
			foundUser, err = s.authRepo.GetUserByEmail(txCtx, gothUser.Email)
			if err != nil {
				s.log.Error("auth_service: failed to get user by email",
					"error", err,
					"email", gothUser.Email,
				)
				return fmt.Errorf("failed to get user by email: %w", err)
			}

			if foundUser != nil {
				s.log.Info("auth_service: existing user found by email",
					"user_id", foundUser.ID,
					"email", foundUser.Email,
					"new_provider", gothUser.Provider,
				)
			}
		}

		// Step 3: If still not found, create new user
		if foundUser == nil {
			if gothUser.Email == "" {
				s.log.Error("auth_service: email required but not provided",
					"provider", gothUser.Provider,
					"provider_user_id", gothUser.UserID,
				)
				return fmt.Errorf("email permission required for authentication")
			}

			s.log.Info("auth_service: user not found, creating new user",
				"provider", gothUser.Provider,
				"email", gothUser.Email,
			)

			foundUser, err = s.authRepo.CreateUser(txCtx, gothUser.Email, gothUser.Name, gothUser.AvatarURL, domain.RoleUser)
			if err != nil {
				// Check if it's a duplicate key error (race condition)
				if isDuplicateKeyError(err) {
					s.log.Warn("auth_service: duplicate key error, retrying lookup",
						"email", gothUser.Email,
					)
					// Another goroutine created the user, retry lookup
					foundUser, err = s.authRepo.GetUserByEmail(txCtx, gothUser.Email)
					if err != nil || foundUser == nil {
						s.log.Error("auth_service: failed to get user after duplicate key",
							"error", err,
							"email", gothUser.Email,
						)
						return fmt.Errorf("failed to get user after duplicate key: %w", err)
					}
					s.log.Info("auth_service: user found after retry", "user_id", foundUser.ID)
				} else {
					s.log.Error("auth_service: failed to create user",
						"error", err,
						"email", gothUser.Email,
					)
					return fmt.Errorf("failed to create user: %w", err)
				}
			} else {
				s.log.Info("auth_service: new user created",
					"user_id", foundUser.ID,
					"email", foundUser.Email,
				)
			}
		} else {
			s.log.Info("auth_service: existing user found",
				"user_id", foundUser.ID,
				"email", foundUser.Email,
			)
		}

		// Step 4: Link OAuth provider (UPSERT)
		oauthProvider := &domain.OAuthProvider{
			UserID:         foundUser.ID,
			Provider:       gothUser.Provider,
			ProviderUserID: gothUser.UserID,
			AccessToken:    gothUser.AccessToken,
			RefreshToken:   gothUser.RefreshToken,
			ExpiresAt:      gothUser.ExpiresAt,
		}

		err = s.authRepo.LinkOAuthProvider(txCtx, oauthProvider)
		if err != nil {
			s.log.Error("auth_service: failed to link OAuth provider",
				"error", err,
				"user_id", foundUser.ID,
				"provider", gothUser.Provider,
			)
			return fmt.Errorf("failed to link oauth provider: %w", err)
		}

		s.log.Debug("auth_service: OAuth provider linked",
			"user_id", foundUser.ID,
			"provider", gothUser.Provider,
		)

		user = foundUser
		return nil
	})

	if err != nil {
		return nil, nil, err
	}

	// Step 5: Create session (outside transaction)
	session, err := s.createSession(ctx, user.ID)
	if err != nil {
		s.log.Error("auth_service: failed to create session",
			"error", err,
			"user_id", user.ID,
		)
		return nil, nil, fmt.Errorf("failed to create session: %w", err)
	}

	s.log.Info("auth_service: OAuth login successful",
		"user_id", user.ID,
		"provider", gothUser.Provider,
		"session_expires", session.ExpiresAt,
	)

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

// isDuplicateKeyError checks if error is PostgreSQL unique constraint violation
func isDuplicateKeyError(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505" // unique_violation
	}
	return false
}
