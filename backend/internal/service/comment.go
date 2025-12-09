package service

import (
	"context"
	"fmt"

	"personal-web-platform/internal/domain"
	"personal-web-platform/internal/pkg/validator"
	"personal-web-platform/internal/repository"
)

// CommentService defines methods for comment business logic
type CommentService interface {
	CreateComment(ctx context.Context, postID int, req *domain.CreateCommentRequest, userID int) (*domain.Comment, error)
	UpdateComment(ctx context.Context, commentID int, req *domain.UpdateCommentRequest, userID int, isAdmin bool) (*domain.Comment, error)
	DeleteComment(ctx context.Context, commentID int, userID int, isAdmin bool) error
	GetCommentByID(ctx context.Context, id int) (*domain.Comment, error)
	GetCommentsByPostSlug(ctx context.Context, slug string, userID int) ([]domain.Comment, error)
}

type commentService struct {
	commentRepo repository.CommentRepository
	postRepo    repository.PostRepository
}

// NewCommentService creates a new comment service implementation
func NewCommentService(commentRepo repository.CommentRepository, postRepo repository.PostRepository) CommentService {
	return &commentService{
		commentRepo: commentRepo,
		postRepo:    postRepo,
	}
}

func (s *commentService) CreateComment(ctx context.Context, postID int, req *domain.CreateCommentRequest, userID int) (*domain.Comment, error) {
	// Validate request
	if err := validator.Validate(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Check if post exists
	post, err := s.postRepo.GetByID(ctx, postID, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to get post: %w", err)
	}
	if post == nil {
		return nil, fmt.Errorf("post not found")
	}

	// If parent_id is specified, check if parent comment exists
	if req.ParentID != nil {
		parent, err := s.commentRepo.GetByID(ctx, *req.ParentID)
		if err != nil {
			return nil, fmt.Errorf("failed to get parent comment: %w", err)
		}
		if parent == nil {
			return nil, fmt.Errorf("parent comment not found")
		}
		// Parent comment must belong to the same post
		if parent.PostID != postID {
			return nil, fmt.Errorf("parent comment does not belong to this post")
		}
		// Cannot reply to deleted comments
		if parent.DeletedAt != nil {
			return nil, fmt.Errorf("cannot reply to deleted comment")
		}
	}

	// Create comment
	comment := &domain.Comment{
		PostID:   postID,
		UserID:   userID,
		Content:  req.Content,
		ParentID: req.ParentID,
	}

	createdComment, err := s.commentRepo.Create(ctx, comment)
	if err != nil {
		return nil, fmt.Errorf("failed to create comment: %w", err)
	}

	return createdComment, nil
}

func (s *commentService) UpdateComment(ctx context.Context, commentID int, req *domain.UpdateCommentRequest, userID int, _ bool) (*domain.Comment, error) { //nolint:revive // isAdmin reserved for future permission checks
	// Validate request
	if err := validator.Validate(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Get existing comment
	comment, err := s.commentRepo.GetByID(ctx, commentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get comment: %w", err)
	}
	if comment == nil {
		return nil, fmt.Errorf("comment not found")
	}

	// Check if comment is deleted
	if comment.DeletedAt != nil {
		return nil, fmt.Errorf("cannot update deleted comment")
	}

	// Check permission (only comment author can update, admin cannot)
	if comment.UserID != userID {
		return nil, fmt.Errorf("permission denied: you can only edit your own comments")
	}

	// Update comment
	comment.Content = req.Content

	updatedComment, err := s.commentRepo.Update(ctx, comment)
	if err != nil {
		return nil, fmt.Errorf("failed to update comment: %w", err)
	}

	return updatedComment, nil
}

func (s *commentService) DeleteComment(ctx context.Context, commentID int, userID int, isAdmin bool) error {
	// Get existing comment
	comment, err := s.commentRepo.GetByID(ctx, commentID)
	if err != nil {
		return fmt.Errorf("failed to get comment: %w", err)
	}
	if comment == nil {
		return fmt.Errorf("comment not found")
	}

	// Check if comment is already deleted
	if comment.DeletedAt != nil {
		return fmt.Errorf("comment already deleted")
	}

	// Check permission (only author or admin can delete)
	if !isAdmin && comment.UserID != userID {
		return fmt.Errorf("permission denied: you can only delete your own comments")
	}

	if err := s.commentRepo.Delete(ctx, commentID); err != nil {
		return fmt.Errorf("failed to delete comment: %w", err)
	}

	return nil
}

func (s *commentService) GetCommentByID(ctx context.Context, id int) (*domain.Comment, error) {
	comment, err := s.commentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get comment: %w", err)
	}
	if comment == nil {
		return nil, fmt.Errorf("comment not found")
	}
	return comment, nil
}

func (s *commentService) GetCommentsByPostSlug(ctx context.Context, slug string, userID int) ([]domain.Comment, error) {
	// Get post by slug
	post, err := s.postRepo.GetBySlug(ctx, slug, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to get post: %w", err)
	}
	if post == nil {
		return nil, fmt.Errorf("post not found")
	}

	// Get comments
	comments, err := s.commentRepo.GetByPostID(ctx, post.ID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get comments: %w", err)
	}

	return comments, nil
}
