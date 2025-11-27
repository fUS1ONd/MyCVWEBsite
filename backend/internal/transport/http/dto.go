package http

import "personal-web-platform/internal/domain"

// Swagger models for documentation purposes

// ProfileResponse represents profile data response
type ProfileResponse struct {
	ID          int              `json:"id" example:"1"`
	Name        string           `json:"name" example:"John Doe"`
	Description string           `json:"description" example:"Senior Software Engineer"`
	PhotoURL    string           `json:"photo_url,omitempty" example:"https://example.com/photo.jpg"`
	Activity    string           `json:"activity" example:"Building awesome products"`
	Contacts    ContactsResponse `json:"contacts"`
}

// ContactsResponse represents contact information
type ContactsResponse struct {
	Email    string `json:"email" example:"john@example.com"`
	GitHub   string `json:"github" example:"https://github.com/johndoe"`
	LinkedIn string `json:"linkedin" example:"https://linkedin.com/in/johndoe"`
}

// PostResponse represents post data response
type PostResponse struct {
	ID          int64  `json:"id" example:"1"`
	Title       string `json:"title" example:"Introduction to AI"`
	Slug        string `json:"slug" example:"introduction-to-ai"`
	Content     string `json:"content" example:"# Introduction\n\nThis is a post about AI..."`
	Preview     string `json:"preview" example:"Learn the basics of artificial intelligence"`
	AuthorID    int64  `json:"author_id" example:"1"`
	Published   bool   `json:"published" example:"true"`
	PublishedAt string `json:"published_at,omitempty" example:"2025-01-15T10:00:00Z"`
	CreatedAt   string `json:"created_at" example:"2025-01-15T09:00:00Z"`
	UpdatedAt   string `json:"updated_at" example:"2025-01-15T09:30:00Z"`
}

// PostListResponse represents paginated posts response
type PostListResponse struct {
	Posts []PostResponse `json:"posts"`
	Meta  MetaData       `json:"meta"`
}

// CommentResponse represents comment data response
type CommentResponse struct {
	ID        int64             `json:"id" example:"1"`
	PostID    int64             `json:"post_id" example:"1"`
	UserID    int64             `json:"user_id" example:"2"`
	Content   string            `json:"content" example:"Great post!"`
	ParentID  *int64            `json:"parent_id,omitempty" example:"null"`
	CreatedAt string            `json:"created_at" example:"2025-01-15T10:30:00Z"`
	UpdatedAt string            `json:"updated_at" example:"2025-01-15T10:30:00Z"`
	DeletedAt *string           `json:"deleted_at,omitempty" example:"null"`
	Replies   []CommentResponse `json:"replies,omitempty"`
}

// UserResponse represents user data response
type UserResponse struct {
	ID        int    `json:"id" example:"1"`
	Email     string `json:"email" example:"user@example.com"`
	Role      string `json:"role" example:"user"`
	CreatedAt string `json:"created_at" example:"2025-01-15T08:00:00Z"`
}

// CreatePostRequest represents post creation request
type CreatePostRequest struct {
	Title     string `json:"title" binding:"required" example:"My New Post"`
	Content   string `json:"content" binding:"required" example:"# Content\n\nPost content here..."`
	Preview   string `json:"preview" example:"Short preview of the post"`
	Published bool   `json:"published" example:"false"`
}

// UpdatePostRequest represents post update request
type UpdatePostRequest struct {
	Title     *string `json:"title,omitempty" example:"Updated Title"`
	Slug      *string `json:"slug,omitempty" example:"updated-slug"`
	Content   *string `json:"content,omitempty" example:"# Updated Content"`
	Preview   *string `json:"preview,omitempty" example:"Updated preview"`
	Published *bool   `json:"published,omitempty" example:"true"`
}

// UpdateProfileRequest represents profile update request
type UpdateProfileRequest struct {
	Name        *string `json:"name,omitempty" example:"John Doe"`
	Description *string `json:"description,omitempty" example:"Senior Engineer"`
	PhotoURL    *string `json:"photo_url,omitempty" example:"https://example.com/photo.jpg"`
	Activity    *string `json:"activity,omitempty" example:"Building products"`
	Email       *string `json:"email,omitempty" example:"john@example.com"`
	GitHub      *string `json:"github,omitempty" example:"https://github.com/johndoe"`
	LinkedIn    *string `json:"linkedin,omitempty" example:"https://linkedin.com/in/johndoe"`
}

// CreateCommentRequest represents comment creation request
type CreateCommentRequest struct {
	Content  string `json:"content" binding:"required" example:"This is a great post!"`
	ParentID *int64 `json:"parent_id,omitempty" example:"null"`
}

// UpdateCommentRequest represents comment update request
type UpdateCommentRequest struct {
	Content string `json:"content" binding:"required" example:"Updated comment text"`
}

// SuccessResponse represents a successful API response
type SuccessResponse struct {
	Success bool `json:"success" example:"true"`
	Data    any  `json:"data,omitempty"`
	Meta    any  `json:"meta,omitempty"`
}

// ErrorResponse represents an error API response
type ErrorResponse struct {
	Success bool       `json:"success" example:"false"`
	Error   *ErrorData `json:"error"`
}

// Convert domain models to response models

// ToProfileResponse converts domain.Profile to ProfileResponse
func ToProfileResponse(p *domain.Profile) ProfileResponse {
	return ProfileResponse{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		PhotoURL:    p.PhotoURL,
		Activity:    p.Activity,
		Contacts: ContactsResponse{
			Email:    p.Contacts.Email,
			GitHub:   p.Contacts.GitHub,
			LinkedIn: p.Contacts.LinkedIn,
		},
	}
}

// ToUserResponse converts domain.User to UserResponse
func ToUserResponse(u *domain.User) UserResponse {
	return UserResponse{
		ID:        u.ID,
		Email:     u.Email,
		Role:      string(u.Role),
		CreatedAt: u.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
