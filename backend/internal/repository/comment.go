package repository

import (
	"context"
	"fmt"
	"time"

	"personal-web-platform/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// CommentRepository defines methods for comment data access
type CommentRepository interface {
	Create(ctx context.Context, comment *domain.Comment) (*domain.Comment, error)
	Update(ctx context.Context, comment *domain.Comment) (*domain.Comment, error)
	Delete(ctx context.Context, id int) error
	GetByID(ctx context.Context, id int) (*domain.Comment, error)
	GetByPostID(ctx context.Context, postID int) ([]domain.Comment, error)
}

type commentRepo struct {
	db *pgxpool.Pool
}

// NewCommentRepo creates a new comment repository implementation
func NewCommentRepo(db *pgxpool.Pool) CommentRepository {
	return &commentRepo{db: db}
}

func (r *commentRepo) Create(ctx context.Context, comment *domain.Comment) (*domain.Comment, error) {
	var createdComment domain.Comment
	db := GetQueryEngine(ctx, r.db)

	query := `
		INSERT INTO comments (post_id, user_id, content, parent_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
		RETURNING id, post_id, user_id, content, parent_id, likes_count, created_at, updated_at, deleted_at
	`

	err := db.QueryRow(ctx, query,
		comment.PostID,
		comment.UserID,
		comment.Content,
		comment.ParentID,
	).Scan(
		&createdComment.ID,
		&createdComment.PostID,
		&createdComment.UserID,
		&createdComment.Content,
		&createdComment.ParentID,
		&createdComment.LikesCount,
		&createdComment.CreatedAt,
		&createdComment.UpdatedAt,
		&createdComment.DeletedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create comment: %w", err)
	}

	return &createdComment, nil
}

func (r *commentRepo) Update(ctx context.Context, comment *domain.Comment) (*domain.Comment, error) {
	var updatedComment domain.Comment
	db := GetQueryEngine(ctx, r.db)

	query := `
		UPDATE comments
		SET content = $1, updated_at = NOW()
		WHERE id = $2 AND deleted_at IS NULL
		RETURNING id, post_id, user_id, content, parent_id, likes_count, created_at, updated_at, deleted_at
	`

	err := db.QueryRow(ctx, query,
		comment.Content,
		comment.ID,
	).Scan(
		&updatedComment.ID,
		&updatedComment.PostID,
		&updatedComment.UserID,
		&updatedComment.Content,
		&updatedComment.ParentID,
		&updatedComment.LikesCount,
		&updatedComment.CreatedAt,
		&updatedComment.UpdatedAt,
		&updatedComment.DeletedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("comment not found or already deleted")
		}
		return nil, fmt.Errorf("failed to update comment: %w", err)
	}

	return &updatedComment, nil
}

func (r *commentRepo) Delete(ctx context.Context, id int) error {
	db := GetQueryEngine(ctx, r.db)

	query := `
		UPDATE comments
		SET deleted_at = $1
		WHERE id = $2 AND deleted_at IS NULL
	`

	result, err := db.Exec(ctx, query, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to delete comment: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("comment not found or already deleted")
	}

	return nil
}

func (r *commentRepo) GetByID(ctx context.Context, id int) (*domain.Comment, error) {
	var comment domain.Comment
	db := GetQueryEngine(ctx, r.db)

	query := `
		SELECT c.id, c.post_id, c.user_id, c.content, c.parent_id, c.likes_count, c.created_at, c.updated_at, c.deleted_at,
		       u.email, u.role
		FROM comments c
		LEFT JOIN users u ON c.user_id = u.id
		WHERE c.id = $1
	`

	var userEmail string
	var userRole domain.Role
	err := db.QueryRow(ctx, query, id).Scan(
		&comment.ID,
		&comment.PostID,
		&comment.UserID,
		&comment.Content,
		&comment.ParentID,
		&comment.LikesCount,
		&comment.CreatedAt,
		&comment.UpdatedAt,
		&comment.DeletedAt,
		&userEmail,
		&userRole,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil // Comment not found
		}
		return nil, fmt.Errorf("failed to get comment by id: %w", err)
	}

	// Set user if exists
	if userEmail != "" {
		comment.User = &domain.User{
			ID:    comment.UserID,
			Email: userEmail,
			Role:  userRole,
		}
	}

	return &comment, nil
}

func (r *commentRepo) GetByPostID(ctx context.Context, postID int) ([]domain.Comment, error) {
	db := GetQueryEngine(ctx, r.db)
	query := `
		SELECT c.id, c.post_id, c.user_id, c.content, c.parent_id, c.likes_count, c.created_at, c.updated_at, c.deleted_at,
		       u.email, u.role
		FROM comments c
		LEFT JOIN users u ON c.user_id = u.id
		WHERE c.post_id = $1
		ORDER BY c.created_at ASC
	`

	rows, err := db.Query(ctx, query, postID)
	if err != nil {
		return nil, fmt.Errorf("failed to get comments by post id: %w", err)
	}
	defer rows.Close()

	var comments []domain.Comment
	commentsMap := make(map[int]*domain.Comment)

	// First pass: scan all comments
	for rows.Next() {
		var comment domain.Comment
		var userEmail string
		var userRole domain.Role

		err := rows.Scan(
			&comment.ID,
			&comment.PostID,
			&comment.UserID,
			&comment.Content,
			&comment.ParentID,
			&comment.LikesCount,
			&comment.CreatedAt,
			&comment.UpdatedAt,
			&comment.DeletedAt,
			&userEmail,
			&userRole,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan comment: %w", err)
		}

		// Set user if exists
		if userEmail != "" {
			comment.User = &domain.User{
				ID:    comment.UserID,
				Email: userEmail,
				Role:  userRole,
			}
		}

		// Initialize replies slice
		comment.Replies = []domain.Comment{}

		commentsMap[comment.ID] = &comment
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating comments: %w", err)
	}

	// Second pass: build tree structure
	for _, comment := range commentsMap {
		if comment.ParentID == nil {
			// Top-level comment
			comments = append(comments, *comment)
		} else {
			// Nested comment - add to parent's replies
			if parent, ok := commentsMap[*comment.ParentID]; ok {
				parent.Replies = append(parent.Replies, *comment)
			}
		}
	}

	return comments, nil
}
