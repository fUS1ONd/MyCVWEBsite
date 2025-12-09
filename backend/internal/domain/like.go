package domain

import "time"

// PostLike represents a like on a post.
type PostLike struct {
	UserID    int       `json:"user_id" db:"user_id"`
	PostID    int       `json:"post_id" db:"post_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// CommentLike represents a like on a comment.
type CommentLike struct {
	UserID    int       `json:"user_id" db:"user_id"`
	CommentID int       `json:"comment_id" db:"comment_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// LikeStatus represents the like status and count for an entity.
type LikeStatus struct {
	IsLiked    bool `json:"is_liked"`
	LikesCount int  `json:"likes_count"`
}
