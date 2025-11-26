package domain

import "time"

// Comment represents a comment on a post
type Comment struct {
	ID        int        `json:"id"`
	PostID    int        `json:"post_id" validate:"required"`
	UserID    int        `json:"user_id" validate:"required"`
	Content   string     `json:"content" validate:"required,min=1,max=5000"`
	ParentID  *int       `json:"parent_id,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
	User      *User      `json:"user,omitempty"`
	Replies   []Comment  `json:"replies,omitempty"`
}

// CreateCommentRequest represents the request to create a comment
type CreateCommentRequest struct {
	Content  string `json:"content" validate:"required,min=1,max=5000"`
	ParentID *int   `json:"parent_id,omitempty"`
}

// UpdateCommentRequest represents the request to update a comment
type UpdateCommentRequest struct {
	Content string `json:"content" validate:"required,min=1,max=5000"`
}
