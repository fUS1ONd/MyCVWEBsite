package domain

import "time"

// Post represents a blog post
type Post struct {
	ID          int       `json:"id"`
	Title       string    `json:"title" validate:"required,min=3,max=255"`
	Slug        string    `json:"slug" validate:"required,min=3,max=255"`
	Content     string    `json:"content" validate:"required,min=10"`
	Preview     string    `json:"preview"`
	AuthorID    int       `json:"author_id" validate:"required"`
	Published   bool      `json:"published"`
	PublishedAt time.Time `json:"published_at,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Author      *User     `json:"author,omitempty"`
	Media       []Media   `json:"media,omitempty"`
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
