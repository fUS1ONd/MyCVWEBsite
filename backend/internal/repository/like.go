package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// LikeRepository defines the interface for like operations.
type LikeRepository interface {
	// Post likes
	TogglePostLike(ctx context.Context, userID, postID int) (bool, error)
	GetPostLikesCount(ctx context.Context, postID int) (int, error)
	IsPostLikedByUser(ctx context.Context, userID, postID int) (bool, error)

	// Comment likes
	ToggleCommentLike(ctx context.Context, userID, commentID int) (bool, error)
	GetCommentLikesCount(ctx context.Context, commentID int) (int, error)
	IsCommentLikedByUser(ctx context.Context, userID, commentID int) (bool, error)
}

type likeRepository struct {
	db *pgxpool.Pool
}

// NewLikeRepository creates a new instance of LikeRepository.
func NewLikeRepository(db *pgxpool.Pool) LikeRepository {
	return &likeRepository{db: db}
}

// TogglePostLike adds or removes a like for a post. Returns true if liked, false if unliked.
func (r *likeRepository) TogglePostLike(ctx context.Context, userID, postID int) (bool, error) {
	// Check if like exists
	isLiked, err := r.IsPostLikedByUser(ctx, userID, postID)
	if err != nil {
		return false, fmt.Errorf("failed to check post like status: %w", err)
	}

	if isLiked {
		// Unlike: delete the like
		query := `DELETE FROM post_likes WHERE user_id = $1 AND post_id = $2`
		_, err := r.db.Exec(ctx, query, userID, postID)
		if err != nil {
			return false, fmt.Errorf("failed to unlike post: %w", err)
		}
		return false, nil
	}

	// Like: insert a new like
	query := `INSERT INTO post_likes (user_id, post_id) VALUES ($1, $2)`
	_, err = r.db.Exec(ctx, query, userID, postID)
	if err != nil {
		return false, fmt.Errorf("failed to like post: %w", err)
	}
	return true, nil
}

// GetPostLikesCount returns the total number of likes for a post.
func (r *likeRepository) GetPostLikesCount(ctx context.Context, postID int) (int, error) {
	query := `SELECT COUNT(*) FROM post_likes WHERE post_id = $1`
	var count int
	err := r.db.QueryRow(ctx, query, postID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get post likes count: %w", err)
	}
	return count, nil
}

// IsPostLikedByUser checks if a user has liked a specific post.
func (r *likeRepository) IsPostLikedByUser(ctx context.Context, userID, postID int) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM post_likes WHERE user_id = $1 AND post_id = $2)`
	var exists bool
	err := r.db.QueryRow(ctx, query, userID, postID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check if post is liked: %w", err)
	}
	return exists, nil
}

// ToggleCommentLike adds or removes a like for a comment. Returns true if liked, false if unliked.
func (r *likeRepository) ToggleCommentLike(ctx context.Context, userID, commentID int) (bool, error) {
	// Check if like exists
	isLiked, err := r.IsCommentLikedByUser(ctx, userID, commentID)
	if err != nil {
		return false, fmt.Errorf("failed to check comment like status: %w", err)
	}

	if isLiked {
		// Unlike: delete the like
		query := `DELETE FROM comment_likes WHERE user_id = $1 AND comment_id = $2`
		_, err := r.db.Exec(ctx, query, userID, commentID)
		if err != nil {
			return false, fmt.Errorf("failed to unlike comment: %w", err)
		}
		return false, nil
	}

	// Like: insert a new like
	query := `INSERT INTO comment_likes (user_id, comment_id) VALUES ($1, $2)`
	_, err = r.db.Exec(ctx, query, userID, commentID)
	if err != nil {
		return false, fmt.Errorf("failed to like comment: %w", err)
	}
	return true, nil
}

// GetCommentLikesCount returns the total number of likes for a comment.
func (r *likeRepository) GetCommentLikesCount(ctx context.Context, commentID int) (int, error) {
	query := `SELECT COUNT(*) FROM comment_likes WHERE comment_id = $1`
	var count int
	err := r.db.QueryRow(ctx, query, commentID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get comment likes count: %w", err)
	}
	return count, nil
}

// IsCommentLikedByUser checks if a user has liked a specific comment.
func (r *likeRepository) IsCommentLikedByUser(ctx context.Context, userID, commentID int) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM comment_likes WHERE user_id = $1 AND comment_id = $2)`
	var exists bool
	err := r.db.QueryRow(ctx, query, userID, commentID).Scan(&exists)
	if err != nil && err != pgx.ErrNoRows {
		return false, fmt.Errorf("failed to check if comment is liked: %w", err)
	}
	return exists, nil
}
