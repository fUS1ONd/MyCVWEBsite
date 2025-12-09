package domain

import "time"

// Post represents a blog post
type Post struct {
	ID              int       `json:"id"`
	Title           string    `json:"title" validate:"required,min=3,max=255"`
	Slug            string    `json:"slug" validate:"required,min=3,max=255"`
	Content         string    `json:"content" validate:"required,min=10"`
	Preview         string    `json:"preview"`
	AuthorID        int       `json:"author_id" validate:"required"`
	Published       bool      `json:"published"`
	PublishedAt     time.Time `json:"published_at,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	CoverImage      string    `json:"cover_image,omitempty" db:"cover_image"`
	ReadTimeMinutes int       `json:"read_time_minutes" db:"read_time_minutes"`
	LikesCount      int       `json:"likes_count" db:"likes_count"`
	CommentsCount   int       `json:"comments_count" db:"comments_count"`
	IsLiked         bool      `json:"is_liked" db:"-"`
	Author          *User     `json:"author,omitempty"`
}

// Media represents media attached to a post
type Media struct {
	ID        int    `json:"id"`
	Filename  string `json:"filename"`
	Path      string `json:"path"`
	MimeType  string `json:"mime_type"`
	Size      int64  `json:"size"`
	SortOrder int    `json:"sort_order"`
}

// CreatePostRequest represents the request to create a post
type CreatePostRequest struct {
	Title      string `json:"title" validate:"required,min=3,max=255"`
	Content    string `json:"content" validate:"required,min=10"`
	Preview    string `json:"preview"`
	Published  bool   `json:"published"`
	CoverImage string `json:"cover_image" validate:"omitempty,url"`
}

// UpdatePostRequest represents the request to update a post
type UpdatePostRequest struct {
	Title      string `json:"title" validate:"required,min=3,max=255"`
	Content    string `json:"content" validate:"required,min=10"`
	Preview    string `json:"preview"`
	Published  bool   `json:"published"`
	CoverImage string `json:"cover_image" validate:"omitempty,url"`
}

// ListPostsRequest represents query parameters for listing posts
type ListPostsRequest struct {
	Page      int   `json:"page" validate:"omitempty,min=1"`
	Limit     int   `json:"limit" validate:"omitempty,min=1,max=100"`
	Published *bool `json:"published,omitempty"`
	UserID    int   `json:"user_id,omitempty"`
}

// PostsListResponse represents paginated posts response
type PostsListResponse struct {
	Posts      []Post `json:"posts"`
	TotalCount int    `json:"total_count"`
	Page       int    `json:"page"`
	Limit      int    `json:"limit"`
	TotalPages int    `json:"total_pages"`
}
