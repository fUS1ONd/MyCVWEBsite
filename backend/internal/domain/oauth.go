package domain

import "time"

// OAuthProvider represents OAuth provider information for a user
type OAuthProvider struct {
	ID             int       `json:"id"`
	UserID         int       `json:"user_id"`
	Provider       string    `json:"provider" validate:"required,oneof=vk google github"`
	ProviderUserID string    `json:"provider_user_id" validate:"required"`
	AccessToken    string    `json:"access_token,omitempty"`
	RefreshToken   string    `json:"refresh_token,omitempty"`
	ExpiresAt      time.Time `json:"expires_at,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
