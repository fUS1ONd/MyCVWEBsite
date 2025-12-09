package service

import (
	"context"
	"fmt"

	"personal-web-platform/internal/domain"
	"personal-web-platform/internal/pkg/slugify"
	"personal-web-platform/internal/pkg/validator"
	"personal-web-platform/internal/repository"
)

// PostService defines methods for post business logic
type PostService interface {
	CreatePost(ctx context.Context, req *domain.CreatePostRequest, authorID int) (*domain.Post, error)
	UpdatePost(ctx context.Context, postID int, req *domain.UpdatePostRequest, userID int, isAdmin bool) (*domain.Post, error)
	DeletePost(ctx context.Context, postID int, userID int, isAdmin bool) error
	GetPostByID(ctx context.Context, id, userID int) (*domain.Post, error)
	GetPostBySlug(ctx context.Context, slug string, userID int) (*domain.Post, error)
	ListPosts(ctx context.Context, req *domain.ListPostsRequest) (*domain.PostsListResponse, error)
}

type postService struct {
	postRepo repository.PostRepository
}

// NewPostService creates a new post service implementation
func NewPostService(postRepo repository.PostRepository) PostService {
	return &postService{
		postRepo: postRepo,
	}
}

func (s *postService) CreatePost(ctx context.Context, req *domain.CreatePostRequest, authorID int) (*domain.Post, error) {
	// Validate request
	if err := validator.Validate(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Generate slug from title
	slug := slugify.Generate(req.Title)
	if slug == "" {
		return nil, fmt.Errorf("failed to generate slug from title")
	}

	// Check if slug already exists
	existingPost, err := s.postRepo.GetBySlug(ctx, slug, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to check slug uniqueness: %w", err)
	}
	if existingPost != nil {
		return nil, fmt.Errorf("post with this slug already exists")
	}

	// Create post
	post := &domain.Post{
		Title:     req.Title,
		Slug:      slug,
		Content:   req.Content,
		Preview:   req.Preview,
		AuthorID:  authorID,
		Published: req.Published,
	}

	createdPost, err := s.postRepo.Create(ctx, post)
	if err != nil {
		return nil, fmt.Errorf("failed to create post: %w", err)
	}

	return createdPost, nil
}

func (s *postService) UpdatePost(ctx context.Context, postID int, req *domain.UpdatePostRequest, userID int, isAdmin bool) (*domain.Post, error) {
	// Validate request
	if err := validator.Validate(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Get existing post
	post, err := s.postRepo.GetByID(ctx, postID, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to get post: %w", err)
	}
	if post == nil {
		return nil, fmt.Errorf("post not found")
	}

	// Check permission (only author or admin can update)
	if !isAdmin && post.AuthorID != userID {
		return nil, fmt.Errorf("permission denied: you can only edit your own posts")
	}

	// Generate new slug if title changed
	newSlug := slugify.Generate(req.Title)
	if newSlug == "" {
		return nil, fmt.Errorf("failed to generate slug from title")
	}

	// Check if new slug conflicts with other posts (excluding current post)
	if newSlug != post.Slug {
		existingPost, err := s.postRepo.GetBySlug(ctx, newSlug, 0)
		if err != nil {
			return nil, fmt.Errorf("failed to check slug uniqueness: %w", err)
		}
		if existingPost != nil && existingPost.ID != postID {
			return nil, fmt.Errorf("post with this slug already exists")
		}
	}

	// Update post fields
	post.Title = req.Title
	post.Slug = newSlug
	post.Content = req.Content
	post.Preview = req.Preview
	post.Published = req.Published

	updatedPost, err := s.postRepo.Update(ctx, post)
	if err != nil {
		return nil, fmt.Errorf("failed to update post: %w", err)
	}

	return updatedPost, nil
}

func (s *postService) DeletePost(ctx context.Context, postID int, userID int, isAdmin bool) error {
	// Get existing post
	post, err := s.postRepo.GetByID(ctx, postID, 0)
	if err != nil {
		return fmt.Errorf("failed to get post: %w", err)
	}
	if post == nil {
		return fmt.Errorf("post not found")
	}

	// Check permission (only author or admin can delete)
	if !isAdmin && post.AuthorID != userID {
		return fmt.Errorf("permission denied: you can only delete your own posts")
	}

	if err := s.postRepo.Delete(ctx, postID); err != nil {
		return fmt.Errorf("failed to delete post: %w", err)
	}

	return nil
}

func (s *postService) GetPostByID(ctx context.Context, id, userID int) (*domain.Post, error) {
	post, err := s.postRepo.GetByID(ctx, id, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get post: %w", err)
	}
	if post == nil {
		return nil, fmt.Errorf("post not found")
	}
	return post, nil
}

func (s *postService) GetPostBySlug(ctx context.Context, slug string, userID int) (*domain.Post, error) {
	post, err := s.postRepo.GetBySlug(ctx, slug, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get post: %w", err)
	}
	if post == nil {
		return nil, fmt.Errorf("post not found")
	}
	return post, nil
}

func (s *postService) ListPosts(ctx context.Context, req *domain.ListPostsRequest) (*domain.PostsListResponse, error) {
	// Validate request
	if err := validator.Validate(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Set defaults
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 1 {
		req.Limit = 10
	}

	posts, totalCount, err := s.postRepo.List(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to list posts: %w", err)
	}

	// Calculate total pages
	totalPages := (totalCount + req.Limit - 1) / req.Limit

	return &domain.PostsListResponse{
		Posts:      posts,
		TotalCount: totalCount,
		Page:       req.Page,
		Limit:      req.Limit,
		TotalPages: totalPages,
	}, nil
}
