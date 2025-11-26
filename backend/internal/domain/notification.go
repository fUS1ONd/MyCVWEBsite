package domain

import "time"

// NotificationSettings represents user notification preferences
type NotificationSettings struct {
	UserID          int       `json:"user_id"`
	EmailEnabled    bool      `json:"email_enabled"`
	PushEnabled     bool      `json:"push_enabled"`
	NewPostsEnabled bool      `json:"new_posts_enabled"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// Notification represents a notification sent to a user
type Notification struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id" validate:"required"`
	Type      string    `json:"type" validate:"required"`
	Title     string    `json:"title" validate:"required,max=255"`
	Message   string    `json:"message" validate:"required"`
	Read      bool      `json:"read"`
	CreatedAt time.Time `json:"created_at"`
}

// NotificationType constants
const (
	NotificationTypeNewPost    = "new_post"
	NotificationTypeNewComment = "new_comment"
)
