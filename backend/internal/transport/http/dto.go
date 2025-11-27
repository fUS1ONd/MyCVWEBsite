package http

import "personal-web-platform/internal/domain"

// ProfileResponse represents profile data response
type ProfileResponse struct {
	ID          int              `json:"id"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	PhotoURL    string           `json:"photo_url,omitempty"`
	Activity    string           `json:"activity"`
	Contacts    ContactsResponse `json:"contacts"`
}

// ContactsResponse represents contact information
type ContactsResponse struct {
	Email    string `json:"email"`
	GitHub   string `json:"github"`
	LinkedIn string `json:"linkedin"`
}

// PostResponse represents post data response
type PostResponse struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Slug        string `json:"slug"`
	Content     string `json:"content"`
	Preview     string `json:"preview"`
	AuthorID    int64  `json:"author_id"`
	Published   bool   `json:"published"`
	PublishedAt string `json:"published_at,omitempty"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// PostListResponse represents paginated posts response
type PostListResponse struct {
	Posts []PostResponse `json:"posts"`
	Meta  MetaData       `json:"meta"`
}

// CommentResponse represents comment data response
type CommentResponse struct {
	ID        int64             `json:"id"`
	PostID    int64             `json:"post_id"`
	UserID    int64             `json:"user_id"`
	Content   string            `json:"content"`
	ParentID  *int64            `json:"parent_id,omitempty"`
	CreatedAt string            `json:"created_at"`
	UpdatedAt string            `json:"updated_at"`
	DeletedAt *string           `json:"deleted_at,omitempty"`
	Replies   []CommentResponse `json:"replies,omitempty"`
}

// UserResponse represents user data response
type UserResponse struct {
	ID        int    `json:"id"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	CreatedAt string `json:"created_at"`
}

// CreatePostRequest represents post creation request
type CreatePostRequest struct {
	Title     string `json:"title" binding:"required"`
	Content   string `json:"content" binding:"required"`
	Preview   string `json:"preview"`
	Published bool   `json:"published"`
}

// UpdatePostRequest represents post update request
type UpdatePostRequest struct {
	Title     *string `json:"title,omitempty"`
	Slug      *string `json:"slug,omitempty"`
	Content   *string `json:"content,omitempty"`
	Preview   *string `json:"preview,omitempty"`
	Published *bool   `json:"published,omitempty"`
}

// UpdateProfileRequest represents profile update request
type UpdateProfileRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	PhotoURL    *string `json:"photo_url,omitempty"`
	Activity    *string `json:"activity,omitempty"`
	Email       *string `json:"email,omitempty"`
	GitHub      *string `json:"github,omitempty"`
	LinkedIn    *string `json:"linkedin,omitempty"`
}

// CreateCommentRequest represents comment creation request
type CreateCommentRequest struct {
	Content  string `json:"content" binding:"required"`
	ParentID *int64 `json:"parent_id,omitempty"`
}

// UpdateCommentRequest represents comment update request
type UpdateCommentRequest struct {
	Content string `json:"content" binding:"required"`
}

// SuccessResponse represents a successful API response
type SuccessResponse struct {
	Success bool `json:"success"`
	Data    any  `json:"data,omitempty"`
	Meta    any  `json:"meta,omitempty"`
}

// ErrorResponse represents an error API response
type ErrorResponse struct {
	Success bool       `json:"success"`
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
