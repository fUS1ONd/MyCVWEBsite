package service

import (
	"bytes"
	"context"
	"image"
	"image/png"
	"mime/multipart"
	"os"
	"path/filepath"
	"testing"
	"time"

	"personal-web-platform/internal/domain"
	"personal-web-platform/internal/domain/derr"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockMediaRepository is a mock implementation of MediaRepository
type MockMediaRepository struct {
	mock.Mock
}

func (m *MockMediaRepository) Create(ctx context.Context, media *domain.MediaFile) error {
	args := m.Called(ctx, media)
	// Set ID and timestamp for testing
	media.ID = 1
	media.UploadedAt = time.Now()
	return args.Error(0)
}

func (m *MockMediaRepository) GetByID(ctx context.Context, id int) (*domain.MediaFile, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.MediaFile), args.Error(1)
}

func (m *MockMediaRepository) Delete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockMediaRepository) ListByUploader(ctx context.Context, uploaderID int) ([]domain.MediaFile, error) {
	args := m.Called(ctx, uploaderID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.MediaFile), args.Error(1)
}

// createTestImage creates a simple test image
func createTestImage(t *testing.T, width, height int) *bytes.Buffer {
	t.Helper()

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	buf := new(bytes.Buffer)
	err := png.Encode(buf, img)
	assert.NoError(t, err)
	return buf
}

// createMultipartFile creates a multipart file for testing
func createMultipartFile(t *testing.T, filename string, content []byte) (multipart.File, *multipart.FileHeader) {
	t.Helper()

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filename)
	assert.NoError(t, err)

	_, err = part.Write(content)
	assert.NoError(t, err)
	assert.NoError(t, writer.Close())

	reader := multipart.NewReader(body, writer.Boundary())
	form, err := reader.ReadForm(10 << 20)
	assert.NoError(t, err)

	files := form.File["file"]
	assert.Len(t, files, 1)

	file, err := files[0].Open()
	assert.NoError(t, err)

	return file, files[0]
}

func TestMediaService_Upload(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockMediaRepository)

	// Create temporary directory for uploads
	tmpDir, err := os.MkdirTemp("", "media_test")
	assert.NoError(t, err)
	defer func() { _ = os.RemoveAll(tmpDir) }()

	service := NewMediaService(mockRepo, tmpDir, "http://localhost:8080")

	t.Run("successful upload", func(t *testing.T) {
		// Create test image
		imgBuf := createTestImage(t, 100, 100)
		file, header := createMultipartFile(t, "test.png", imgBuf.Bytes())
		defer func() { _ = file.Close() }()

		mockRepo.On("Create", ctx, mock.AnythingOfType("*domain.MediaFile")).Return(nil).Once()

		result, err := service.Upload(ctx, file, header, 1)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "test.png", result.Filename)
		assert.Equal(t, "png", result.MimeType)
		assert.Equal(t, 1, result.UploaderID)
		assert.Contains(t, result.URL, "http://localhost:8080/media/")

		// Verify file was created on disk
		assert.FileExists(t, result.StoragePath)

		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid image format", func(t *testing.T) {
		// Create non-image file
		content := []byte("not an image")
		file, header := createMultipartFile(t, "test.txt", content)
		defer func() { _ = file.Close() }()

		result, err := service.Upload(ctx, file, header, 1)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "validation failed")
	})

	t.Run("repository error", func(t *testing.T) {
		imgBuf := createTestImage(t, 100, 100)
		file, header := createMultipartFile(t, "test.png", imgBuf.Bytes())
		defer func() { _ = file.Close() }()

		mockRepo.On("Create", ctx, mock.AnythingOfType("*domain.MediaFile")).Return(assert.AnError).Once()

		result, err := service.Upload(ctx, file, header, 1)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to save media record")

		mockRepo.AssertExpectations(t)
	})
}

func TestMediaService_GetByID(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockMediaRepository)
	service := NewMediaService(mockRepo, "/uploads", "http://localhost:8080")

	t.Run("successful get", func(t *testing.T) {
		expectedMedia := &domain.MediaFile{
			ID:          1,
			Filename:    "test.png",
			MimeType:    "png",
			Size:        1024,
			UploaderID:  1,
			StoragePath: "/uploads/1_1234567890.png",
			UploadedAt:  time.Now(),
		}

		mockRepo.On("GetByID", ctx, 1).Return(expectedMedia, nil).Once()

		result, err := service.GetByID(ctx, 1)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedMedia.ID, result.ID)
		assert.Equal(t, expectedMedia.Filename, result.Filename)
		assert.Contains(t, result.URL, "http://localhost:8080/media/")

		mockRepo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo.On("GetByID", ctx, 999).Return(nil, assert.AnError).Once()

		result, err := service.GetByID(ctx, 999)

		assert.Error(t, err)
		assert.Nil(t, result)

		mockRepo.AssertExpectations(t)
	})
}

func TestMediaService_Delete(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockMediaRepository)

	// Create temporary directory and test file
	tmpDir, err := os.MkdirTemp("", "media_test")
	assert.NoError(t, err)
	defer func() { _ = os.RemoveAll(tmpDir) }()

	testFile := filepath.Join(tmpDir, "test.png")
	err = os.WriteFile(testFile, []byte("test"), 0600)
	assert.NoError(t, err)

	service := NewMediaService(mockRepo, tmpDir, "http://localhost:8080")

	t.Run("successful delete", func(t *testing.T) {
		media := &domain.MediaFile{
			ID:          1,
			UploaderID:  1,
			StoragePath: testFile,
		}

		mockRepo.On("GetByID", ctx, 1).Return(media, nil).Once()
		mockRepo.On("Delete", ctx, 1).Return(nil).Once()

		err := service.Delete(ctx, 1, 1)

		assert.NoError(t, err)
		assert.NoFileExists(t, testFile)

		mockRepo.AssertExpectations(t)
	})

	t.Run("unauthorized delete", func(t *testing.T) {
		// Recreate test file
		err = os.WriteFile(testFile, []byte("test"), 0600)
		assert.NoError(t, err)

		media := &domain.MediaFile{
			ID:          1,
			UploaderID:  1,
			StoragePath: testFile,
		}

		mockRepo.On("GetByID", ctx, 1).Return(media, nil).Once()

		err := service.Delete(ctx, 1, 2) // Different user ID

		assert.Error(t, err)
		assert.Equal(t, derr.ErrUnauthorized, err)
		assert.FileExists(t, testFile) // File should still exist

		mockRepo.AssertExpectations(t)
	})
}

func TestMediaService_ListByUploader(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockMediaRepository)
	service := NewMediaService(mockRepo, "/uploads", "http://localhost:8080")

	t.Run("successful list", func(t *testing.T) {
		expectedFiles := []domain.MediaFile{
			{
				ID:          1,
				Filename:    "test1.png",
				StoragePath: "/uploads/1_123.png",
				UploaderID:  1,
			},
			{
				ID:          2,
				Filename:    "test2.png",
				StoragePath: "/uploads/1_456.png",
				UploaderID:  1,
			},
		}

		mockRepo.On("ListByUploader", ctx, 1).Return(expectedFiles, nil).Once()

		result, err := service.ListByUploader(ctx, 1)

		assert.NoError(t, err)
		assert.Len(t, result, 2)
		for _, file := range result {
			assert.Contains(t, file.URL, "http://localhost:8080/media/")
		}

		mockRepo.AssertExpectations(t)
	})

	t.Run("empty list", func(t *testing.T) {
		mockRepo.On("ListByUploader", ctx, 1).Return([]domain.MediaFile{}, nil).Once()

		result, err := service.ListByUploader(ctx, 1)

		assert.NoError(t, err)
		assert.Empty(t, result)

		mockRepo.AssertExpectations(t)
	})
}
