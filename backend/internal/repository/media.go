package repository

import (
	"context"
	"fmt"

	"personal-web-platform/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

// MediaRepository defines the interface for media file operations.
type MediaRepository interface {
	Create(ctx context.Context, media *domain.MediaFile) error
	GetByID(ctx context.Context, id int) (*domain.MediaFile, error)
	Delete(ctx context.Context, id int) error
	ListByUploader(ctx context.Context, uploaderID int) ([]domain.MediaFile, error)
}

type mediaRepository struct {
	db *pgxpool.Pool
}

// NewMediaRepository creates a new instance of MediaRepository.
func NewMediaRepository(db *pgxpool.Pool) MediaRepository {
	return &mediaRepository{db: db}
}

// Create inserts a new media file record into the database.
func (r *mediaRepository) Create(ctx context.Context, media *domain.MediaFile) error {
	query := `
		INSERT INTO media_files (filename, mime_type, size, uploader_id, storage_path, uploaded_at)
		VALUES ($1, $2, $3, $4, $5, NOW())
		RETURNING id, uploaded_at
	`
	err := r.db.QueryRow(
		ctx, query,
		media.Filename, media.MimeType, media.Size, media.UploaderID, media.StoragePath,
	).Scan(&media.ID, &media.UploadedAt)

	if err != nil {
		return fmt.Errorf("failed to create media file: %w", err)
	}
	return nil
}

// GetByID retrieves a media file by its ID.
func (r *mediaRepository) GetByID(ctx context.Context, id int) (*domain.MediaFile, error) {
	query := `
		SELECT id, filename, mime_type, size, uploader_id, storage_path, uploaded_at
		FROM media_files
		WHERE id = $1
	`
	var media domain.MediaFile
	err := r.db.QueryRow(ctx, query, id).Scan(
		&media.ID, &media.Filename, &media.MimeType, &media.Size,
		&media.UploaderID, &media.StoragePath, &media.UploadedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get media file: %w", err)
	}
	return &media, nil
}

// Delete removes a media file record from the database.
func (r *mediaRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM media_files WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete media file: %w", err)
	}
	return nil
}

// ListByUploader retrieves all media files uploaded by a specific user.
func (r *mediaRepository) ListByUploader(ctx context.Context, uploaderID int) ([]domain.MediaFile, error) {
	query := `
		SELECT id, filename, mime_type, size, uploader_id, storage_path, uploaded_at
		FROM media_files
		WHERE uploader_id = $1
		ORDER BY uploaded_at DESC
	`
	rows, err := r.db.Query(ctx, query, uploaderID)
	if err != nil {
		return nil, fmt.Errorf("failed to list media files: %w", err)
	}
	defer rows.Close()

	var files []domain.MediaFile
	for rows.Next() {
		var media domain.MediaFile
		err := rows.Scan(
			&media.ID, &media.Filename, &media.MimeType, &media.Size,
			&media.UploaderID, &media.StoragePath, &media.UploadedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan media file: %w", err)
		}
		files = append(files, media)
	}
	return files, nil
}
