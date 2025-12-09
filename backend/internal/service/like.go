package service

import (
	"context"
	"fmt"

	"personal-web-platform/internal/repository"
)

// LikeService defines the interface for like operations.
type LikeService interface {
	// Post likes
	TogglePostLike(ctx context.Context, userID, postID int) (bool, error)
	GetPostLikesCount(ctx context.Context, postID int) (int, error)
	IsPostLikedByUser(ctx context.Context, userID, postID int) (bool, error)

	// Comment likes
	ToggleCommentLike(ctx context.Context, userID, commentID int) (bool, error)
	GetCommentLikesCount(ctx context.Context, commentID int) (int, error)
	IsCommentLikedByUser(ctx context.Context, userID, commentID int) (bool, error)
}

type likeService struct {
	repo     repository.LikeRepository
	postRepo repository.PostRepository
	commRepo repository.CommentRepository
}

// NewLikeService creates a new LikeService instance.
func NewLikeService(repo repository.LikeRepository, postRepo repository.PostRepository, commRepo repository.CommentRepository) LikeService {
	return &likeService{
		repo:     repo,
		postRepo: postRepo,
		commRepo: commRepo,
	}
}

// TogglePostLike toggles a like for a post. Returns true if liked, false if unliked.
func (s *likeService) TogglePostLike(ctx context.Context, userID, postID int) (bool, error) {
	// Verify post exists
	_, err := s.postRepo.GetByID(ctx, postID)
	if err != nil {
		return false, fmt.Errorf("failed to get post: %w", err)
	}

	liked, err := s.repo.TogglePostLike(ctx, userID, postID)
	if err != nil {
		return false, fmt.Errorf("failed to toggle post like: %w", err)
	}

	return liked, nil
}

// GetPostLikesCount returns the total number of likes for a post.
func (s *likeService) GetPostLikesCount(ctx context.Context, postID int) (int, error) {
	count, err := s.repo.GetPostLikesCount(ctx, postID)
	if err != nil {
		return 0, fmt.Errorf("failed to get post likes count: %w", err)
	}
	return count, nil
}

// IsPostLikedByUser checks if a user has liked a specific post.
func (s *likeService) IsPostLikedByUser(ctx context.Context, userID, postID int) (bool, error) {
	liked, err := s.repo.IsPostLikedByUser(ctx, userID, postID)
	if err != nil {
		return false, fmt.Errorf("failed to check if post is liked: %w", err)
	}
	return liked, nil
}

// ToggleCommentLike toggles a like for a comment. Returns true if liked, false if unliked.
func (s *likeService) ToggleCommentLike(ctx context.Context, userID, commentID int) (bool, error) {
	// Verify comment exists
	_, err := s.commRepo.GetByID(ctx, commentID)
	if err != nil {
		return false, fmt.Errorf("failed to get comment: %w", err)
	}

	liked, err := s.repo.ToggleCommentLike(ctx, userID, commentID)
	if err != nil {
		return false, fmt.Errorf("failed to toggle comment like: %w", err)
	}

	return liked, nil
}

// GetCommentLikesCount returns the total number of likes for a comment.
func (s *likeService) GetCommentLikesCount(ctx context.Context, commentID int) (int, error) {
	count, err := s.repo.GetCommentLikesCount(ctx, commentID)
	if err != nil {
		return 0, fmt.Errorf("failed to get comment likes count: %w", err)
	}
	return count, nil
}

// IsCommentLikedByUser checks if a user has liked a specific comment.
func (s *likeService) IsCommentLikedByUser(ctx context.Context, userID, commentID int) (bool, error) {
	liked, err := s.repo.IsCommentLikedByUser(ctx, userID, commentID)
	if err != nil {
		return false, fmt.Errorf("failed to check if comment is liked: %w", err)
	}
	return liked, nil
}
