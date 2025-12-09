package service

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"personal-web-platform/internal/domain"
	"personal-web-platform/internal/domain/derr"
	"personal-web-platform/internal/pkg/imageproc"
	"personal-web-platform/internal/repository"
)

// MediaService defines the interface for media file operations.
type MediaService interface {
	Upload(ctx context.Context, file multipart.File, header *multipart.FileHeader, uploaderID int) (*domain.MediaFile, error)
	GetByID(ctx context.Context, id int) (*domain.MediaFile, error)
	Delete(ctx context.Context, id int, uploaderID int) error
	ListByUploader(ctx context.Context, uploaderID int) ([]domain.MediaFile, error)
}

type mediaService struct {
	repo       repository.MediaRepository
	uploadPath string
	baseURL    string
}

// NewMediaService creates a new MediaService instance.
func NewMediaService(repo repository.MediaRepository, uploadPath, baseURL string) MediaService {
	return &mediaService{
		repo:       repo,
		uploadPath: uploadPath,
		baseURL:    baseURL,
	}
}

// Upload handles file upload, validation, and storage.
func (s *mediaService) Upload(ctx context.Context, file multipart.File, header *multipart.FileHeader, uploaderID int) (*domain.MediaFile, error) {
	// Validate image
	info, err := imageproc.ValidateImage(file, header)
	if err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Generate unique filename
	timestamp := time.Now().Unix()
	ext := filepath.Ext(header.Filename)
	uniqueFilename := fmt.Sprintf("%d_%d%s", uploaderID, timestamp, ext)
	storagePath := filepath.Join(s.uploadPath, uniqueFilename)

	// Create upload directory if it doesn't exist
	if err := os.MkdirAll(s.uploadPath, 0750); err != nil {
		return nil, fmt.Errorf("failed to create upload directory: %w", err)
	}

	// Save file to disk
	outFile, err := os.Create(storagePath) // #nosec G304 - storagePath is constructed from controlled inputs
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %w", err)
	}
	defer func() {
		_ = outFile.Close() // Error ignored intentionally - file is already written
	}()

	// Reset file pointer before copying
	if seeker, ok := file.(io.Seeker); ok {
		if _, err := seeker.Seek(0, io.SeekStart); err != nil {
			return nil, fmt.Errorf("failed to reset file pointer: %w", err)
		}
	}

	if _, err := io.Copy(outFile, file); err != nil {
		_ = os.Remove(storagePath) // Cleanup on error
		return nil, fmt.Errorf("failed to save file: %w", err)
	}

	// Create database record
	media := &domain.MediaFile{
		Filename:    header.Filename,
		MimeType:    info.Format,
		Size:        info.SizeBytes,
		UploaderID:  uploaderID,
		StoragePath: storagePath,
	}

	if err := s.repo.Create(ctx, media); err != nil {
		_ = os.Remove(storagePath) // Cleanup on error
		return nil, fmt.Errorf("failed to save media record: %w", err)
	}

	// Set URL for response
	media.URL = s.generateURL(uniqueFilename)

	return media, nil
}

// GetByID retrieves a media file by ID.
func (s *mediaService) GetByID(ctx context.Context, id int) (*domain.MediaFile, error) {
	media, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get media: %w", err)
	}

	// Set URL for response
	filename := filepath.Base(media.StoragePath)
	media.URL = s.generateURL(filename)

	return media, nil
}

// Delete removes a media file (both from storage and database).
// Only the uploader can delete their files.
func (s *mediaService) Delete(ctx context.Context, id int, uploaderID int) error {
	// Get media to verify ownership
	media, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get media: %w", err)
	}

	// Verify ownership
	if media.UploaderID != uploaderID {
		return derr.ErrUnauthorized
	}

	// Delete from storage
	if err := os.Remove(media.StoragePath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	// Delete from database
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete media record: %w", err)
	}

	return nil
}

// ListByUploader retrieves all media files uploaded by a user.
func (s *mediaService) ListByUploader(ctx context.Context, uploaderID int) ([]domain.MediaFile, error) {
	files, err := s.repo.ListByUploader(ctx, uploaderID)
	if err != nil {
		return nil, fmt.Errorf("failed to list media: %w", err)
	}

	// Set URLs for all files
	for i := range files {
		filename := filepath.Base(files[i].StoragePath)
		files[i].URL = s.generateURL(filename)
	}

	return files, nil
}

// generateURL creates a public URL for a media file.
func (s *mediaService) generateURL(filename string) string {
	return fmt.Sprintf("%s/media/%s", s.baseURL, filename)
}
