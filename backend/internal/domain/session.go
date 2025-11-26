package domain

import "time"

// Session represents user session
type Session struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id" validate:"required"`
	Token     string    `json:"token" validate:"required"`
	ExpiresAt time.Time `json:"expires_at" validate:"required"`
	CreatedAt time.Time `json:"created_at"`
}
