package repository

import (
	"context"
	"fmt"
	"time"

	"personal-web-platform/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PostRepository defines methods for post data access
type PostRepository interface {
	Create(ctx context.Context, post *domain.Post) (*domain.Post, error)
	Update(ctx context.Context, post *domain.Post) (*domain.Post, error)
	Delete(ctx context.Context, id int) error
	GetByID(ctx context.Context, id int) (*domain.Post, error)
	GetBySlug(ctx context.Context, slug string) (*domain.Post, error)
	List(ctx context.Context, req *domain.ListPostsRequest) ([]domain.Post, int, error)
}

type postRepo struct {
	db *pgxpool.Pool
}

// NewPostRepo creates a new post repository implementation
func NewPostRepo(db *pgxpool.Pool) PostRepository {
	return &postRepo{db: db}
}

func (r *postRepo) Create(ctx context.Context, post *domain.Post) (*domain.Post, error) {
	var createdPost domain.Post
	var publishedAt *time.Time
	db := GetQueryEngine(ctx, r.db)

	if post.Published {
		now := time.Now()
		publishedAt = &now
	}

	query := `
		INSERT INTO posts (title, slug, content, preview, author_id, published, published_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())
		RETURNING id, title, slug, content, preview, author_id, published, published_at, created_at, updated_at
	`

	err := db.QueryRow(ctx, query,
		post.Title,
		post.Slug,
		post.Content,
		post.Preview,
		post.AuthorID,
		post.Published,
		publishedAt,
	).Scan(
		&createdPost.ID,
		&createdPost.Title,
		&createdPost.Slug,
		&createdPost.Content,
		&createdPost.Preview,
		&createdPost.AuthorID,
		&createdPost.Published,
		&createdPost.PublishedAt,
		&createdPost.CreatedAt,
		&createdPost.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create post: %w", err)
	}

	return &createdPost, nil
}

func (r *postRepo) Update(ctx context.Context, post *domain.Post) (*domain.Post, error) {
	var updatedPost domain.Post
	var publishedAt *time.Time
	db := GetQueryEngine(ctx, r.db)

	// If publishing for the first time, set published_at
	if post.Published && post.PublishedAt.IsZero() {
		now := time.Now()
		publishedAt = &now
	} else if post.Published {
		publishedAt = &post.PublishedAt
	}

	query := `
		UPDATE posts
		SET title = $1, slug = $2, content = $3, preview = $4, published = $5, published_at = $6, updated_at = NOW()
		WHERE id = $7
		RETURNING id, title, slug, content, preview, author_id, published, published_at, created_at, updated_at
	`

	err := db.QueryRow(ctx, query,
		post.Title,
		post.Slug,
		post.Content,
		post.Preview,
		post.Published,
		publishedAt,
		post.ID,
	).Scan(
		&updatedPost.ID,
		&updatedPost.Title,
		&updatedPost.Slug,
		&updatedPost.Content,
		&updatedPost.Preview,
		&updatedPost.AuthorID,
		&updatedPost.Published,
		&updatedPost.PublishedAt,
		&updatedPost.CreatedAt,
		&updatedPost.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update post: %w", err)
	}

	return &updatedPost, nil
}

func (r *postRepo) Delete(ctx context.Context, id int) error {
	db := GetQueryEngine(ctx, r.db)
	query := `DELETE FROM posts WHERE id = $1`

	result, err := db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete post: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("post not found")
	}

	return nil
}

func (r *postRepo) GetByID(ctx context.Context, id int) (*domain.Post, error) { //nolint:dupl // similar to GetBySlug but different query
	var post domain.Post
	db := GetQueryEngine(ctx, r.db)

	query := `
		SELECT p.id, p.title, p.slug, p.content, p.preview, p.author_id, p.published, p.published_at, p.created_at, p.updated_at,
		       u.email
		FROM posts p
		LEFT JOIN users u ON p.author_id = u.id
		WHERE p.id = $1
	`

	var authorEmail string
	err := db.QueryRow(ctx, query, id).Scan(
		&post.ID,
		&post.Title,
		&post.Slug,
		&post.Content,
		&post.Preview,
		&post.AuthorID,
		&post.Published,
		&post.PublishedAt,
		&post.CreatedAt,
		&post.UpdatedAt,
		&authorEmail,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil // Post not found
		}
		return nil, fmt.Errorf("failed to get post by id: %w", err)
	}

	// Set author if exists
	if authorEmail != "" {
		post.Author = &domain.User{
			ID:    post.AuthorID,
			Email: authorEmail,
		}
	}

	return &post, nil
}

func (r *postRepo) GetBySlug(ctx context.Context, slug string) (*domain.Post, error) { //nolint:dupl // similar to GetByID but different query
	var post domain.Post
	db := GetQueryEngine(ctx, r.db)

	query := `
		SELECT p.id, p.title, p.slug, p.content, p.preview, p.author_id, p.published, p.published_at, p.created_at, p.updated_at,
		       u.email
		FROM posts p
		LEFT JOIN users u ON p.author_id = u.id
		WHERE p.slug = $1
	`

	var authorEmail string
	err := db.QueryRow(ctx, query, slug).Scan(
		&post.ID,
		&post.Title,
		&post.Slug,
		&post.Content,
		&post.Preview,
		&post.AuthorID,
		&post.Published,
		&post.PublishedAt,
		&post.CreatedAt,
		&post.UpdatedAt,
		&authorEmail,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil // Post not found
		}
		return nil, fmt.Errorf("failed to get post by slug: %w", err)
	}

	// Set author if exists
	if authorEmail != "" {
		post.Author = &domain.User{
			ID:    post.AuthorID,
			Email: authorEmail,
		}
	}

	return &post, nil
}

func (r *postRepo) List(ctx context.Context, req *domain.ListPostsRequest) ([]domain.Post, int, error) {
	db := GetQueryEngine(ctx, r.db)

	// Set default values
	page := req.Page
	if page < 1 {
		page = 1
	}
	limit := req.Limit
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	offset := (page - 1) * limit

	// Build query with filters
	baseQuery := `
		SELECT p.id, p.title, p.slug, p.content, p.preview, p.author_id, p.published, p.published_at, p.created_at, p.updated_at,
		       u.email
		FROM posts p
		LEFT JOIN users u ON p.author_id = u.id
	`
	countQuery := `SELECT COUNT(*) FROM posts p`

	var whereClause string
	var args []interface{}
	argIndex := 1

	// Filter by published status if specified
	if req.Published != nil {
		whereClause = fmt.Sprintf(" WHERE p.published = $%d", argIndex)
		args = append(args, *req.Published)
		argIndex++
	}

	// Order by published_at desc for published posts, created_at desc for drafts
	orderClause := " ORDER BY p.published_at DESC NULLS LAST, p.created_at DESC"

	// Limit and offset
	limitClause := fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, limit, offset)

	// Execute count query
	var totalCount int
	countArgs := args[:len(args)-2] // Exclude limit and offset for count query
	err := db.QueryRow(ctx, countQuery+whereClause, countArgs...).Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count posts: %w", err)
	}

	// Execute list query
	rows, err := db.Query(ctx, baseQuery+whereClause+orderClause+limitClause, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list posts: %w", err)
	}
	defer rows.Close()

	var posts []domain.Post
	for rows.Next() {
		var post domain.Post
		var authorEmail string

		err := rows.Scan(
			&post.ID,
			&post.Title,
			&post.Slug,
			&post.Content,
			&post.Preview,
			&post.AuthorID,
			&post.Published,
			&post.PublishedAt,
			&post.CreatedAt,
			&post.UpdatedAt,
			&authorEmail,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan post: %w", err)
		}

		// Set author if exists
		if authorEmail != "" {
			post.Author = &domain.User{
				ID:    post.AuthorID,
				Email: authorEmail,
			}
		}

		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating posts: %w", err)
	}

	return posts, totalCount, nil
}
