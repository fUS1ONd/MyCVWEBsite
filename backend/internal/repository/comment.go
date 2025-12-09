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
	SoftDelete(ctx context.Context, id int, placeholder string) error
	HardDelete(ctx context.Context, id int) error
	HasReplies(ctx context.Context, id int) (bool, error)
	GetByID(ctx context.Context, id int) (*domain.Comment, error)
	GetByPostID(ctx context.Context, postID, userID int) ([]domain.Comment, error)
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

	// Increment post comments count
	_, err = db.Exec(ctx, "UPDATE posts SET comments_count = comments_count + 1 WHERE id = $1", comment.PostID)
	if err != nil {
		return nil, fmt.Errorf("failed to increment comments count: %w", err)
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

func (r *commentRepo) SoftDelete(ctx context.Context, id int, placeholder string) error {
	db := GetQueryEngine(ctx, r.db)

	query := `
		UPDATE comments
		SET deleted_at = $1, content = $2
		WHERE id = $3 AND deleted_at IS NULL
	`

	result, err := db.Exec(ctx, query, time.Now(), placeholder, id)
	if err != nil {
		return fmt.Errorf("failed to soft delete comment: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("comment not found or already deleted")
	}

	// Decrement post comments count
	var postID int
	if err := db.QueryRow(ctx, "SELECT post_id FROM comments WHERE id = $1", id).Scan(&postID); err == nil {
		_, _ = db.Exec(ctx, "UPDATE posts SET comments_count = GREATEST(comments_count - 1, 0) WHERE id = $1", postID)
	}

	return nil
}

func (r *commentRepo) HardDelete(ctx context.Context, id int) error {
	db := GetQueryEngine(ctx, r.db)

	// Get post_id to decrement count before deletion
	var postID int
	// We ignore error here, if comment doesn't exist, DELETE will affect 0 rows
	_ = db.QueryRow(ctx, "SELECT post_id FROM comments WHERE id = $1", id).Scan(&postID)

	query := `DELETE FROM comments WHERE id = $1`
	result, err := db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to hard delete comment: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("comment not found")
	}

	// Get post_id to decrement count
	if postID != 0 {
		_, _ = db.Exec(ctx, "UPDATE posts SET comments_count = GREATEST(comments_count - 1, 0) WHERE id = $1", postID)
	}

	return nil
}

func (r *commentRepo) HasReplies(ctx context.Context, id int) (bool, error) {
	db := GetQueryEngine(ctx, r.db)
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM comments WHERE parent_id = $1)`
	err := db.QueryRow(ctx, query, id).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check replies: %w", err)
	}
	return exists, nil
}

func (r *commentRepo) GetByID(ctx context.Context, id int) (*domain.Comment, error) {
	var comment domain.Comment
	db := GetQueryEngine(ctx, r.db)

	query := `
		SELECT c.id, c.post_id, c.user_id, c.content, c.parent_id, c.likes_count, c.created_at, c.updated_at, c.deleted_at,
		       u.email, u.name, u.avatar_url, u.role,
		       (SELECT photo_url FROM profile_info LIMIT 1) as profile_photo
		FROM comments c
		LEFT JOIN users u ON c.user_id = u.id
		WHERE c.id = $1
	`

	var userEmail string
	var userName string
	var userAvatar string
	var userRole string
	var profilePhoto *string

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
		&userName,
		&userAvatar,
		&userRole,
		&profilePhoto,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil // Comment not found
		}
		return nil, fmt.Errorf("failed to get comment by id: %w", err)
	}

	// Set user if exists
	if userEmail != "" {
		avatar := userAvatar
		if userRole == string(domain.RoleAdmin) && profilePhoto != nil && *profilePhoto != "" {
			avatar = *profilePhoto
		}

		comment.User = &domain.User{
			ID:        comment.UserID,
			Email:     userEmail,
			Name:      userName,
			AvatarURL: avatar,
			Role:      domain.Role(userRole),
		}
	}

	return &comment, nil
}

func (r *commentRepo) GetByPostID(ctx context.Context, postID, userID int) ([]domain.Comment, error) {
	db := GetQueryEngine(ctx, r.db)
	query := `
		SELECT c.id, c.post_id, c.user_id, c.content, c.parent_id, c.likes_count, c.created_at, c.updated_at, c.deleted_at,
		       u.email, u.name, u.avatar_url, u.role,
		       EXISTS(SELECT 1 FROM comment_likes cl WHERE cl.comment_id = c.id AND cl.user_id = $2) as is_liked,
		       (SELECT photo_url FROM profile_info LIMIT 1) as profile_photo
		FROM comments c
		LEFT JOIN users u ON c.user_id = u.id
		WHERE c.post_id = $1
		ORDER BY c.created_at ASC
	`

	rows, err := db.Query(ctx, query, postID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get comments by post id: %w", err)
	}
	defer rows.Close()

	var comments []domain.Comment
	commentsMap := make(map[int]*domain.Comment)
	var rawComments []*domain.Comment
	// First pass: scan all comments
	for rows.Next() {
		var comment domain.Comment
		var userEmail string
		var userName string
		var userAvatar string
		var userRole string
		var profilePhoto *string

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
			&userName,
			&userAvatar,
			&userRole,
			&comment.IsLiked,
			&profilePhoto,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan comment: %w", err)
		}

		// Set user if exists
		if userEmail != "" {
			avatar := userAvatar
			if userRole == string(domain.RoleAdmin) && profilePhoto != nil && *profilePhoto != "" {
				avatar = *profilePhoto
			}

			comment.User = &domain.User{
				ID:        comment.UserID,
				Email:     userEmail,
				Name:      userName,
				AvatarURL: avatar,
				Role:      domain.Role(userRole),
			}
		}

		// Initialize replies slice
		comment.Replies = []*domain.Comment{}

		ptr := &comment
		commentsMap[comment.ID] = ptr
		rawComments = append(rawComments, ptr)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating comments: %w", err)
	}

	// Second pass: build tree structure
	// We iterate over rawComments to preserve order
	for _, comment := range rawComments {
		if comment.ParentID != nil {
			// Nested comment - add to parent's replies
			if parent, ok := commentsMap[*comment.ParentID]; ok {
				parent.Replies = append(parent.Replies, comment)
			}
		}
	}

	// Third pass: collect roots
	for _, comment := range rawComments {
		if comment.ParentID == nil {
			comments = append(comments, *comment)
		}
	}

	return comments, nil
}
