// Package imageproc provides utilities for validating and processing image files.
package imageproc

import (
	"fmt"
	"image"
	// Register standard image format decoders
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"

	// Register WebP format decoder
	_ "golang.org/x/image/webp"
)

// SupportedFormats lists all supported image formats
var SupportedFormats = []string{".jpg", ".jpeg", ".png", ".gif", ".webp"}

// MaxFileSize is the maximum allowed file size (10MB)
const MaxFileSize = 10 * 1024 * 1024

// ImageInfo contains basic information about an image
type ImageInfo struct {
	Width     int
	Height    int
	Format    string
	SizeBytes int64
}

// ValidateImage checks if the uploaded file is a valid image
// and returns basic information about it.
func ValidateImage(file multipart.File, header *multipart.FileHeader) (*ImageInfo, error) {
	// Check file size
	if header.Size > MaxFileSize {
		return nil, fmt.Errorf("file size exceeds maximum allowed size of %d bytes", MaxFileSize)
	}

	// Check file extension
	ext := strings.ToLower(filepath.Ext(header.Filename))
	if !isValidExtension(ext) {
		return nil, fmt.Errorf("unsupported file format: %s (allowed: %v)", ext, SupportedFormats)
	}

	// Decode image to get dimensions and verify it's actually an image
	img, format, err := image.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	// Reset file pointer for potential further reads
	if seeker, ok := file.(io.Seeker); ok {
		_, err = seeker.Seek(0, io.SeekStart)
		if err != nil {
			return nil, fmt.Errorf("failed to reset file pointer: %w", err)
		}
	}

	bounds := img.Bounds()
	info := &ImageInfo{
		Width:     bounds.Dx(),
		Height:    bounds.Dy(),
		Format:    format,
		SizeBytes: header.Size,
	}

	return info, nil
}

// isValidExtension checks if the file extension is supported
func isValidExtension(ext string) bool {
	for _, supported := range SupportedFormats {
		if ext == supported {
			return true
		}
	}
	return false
}

// GenerateThumbnailFilename generates a thumbnail filename from the original
// Example: "image.jpg" -> "image_thumb.jpg"
func GenerateThumbnailFilename(originalFilename string) string {
	ext := filepath.Ext(originalFilename)
	nameWithoutExt := strings.TrimSuffix(originalFilename, ext)
	return fmt.Sprintf("%s_thumb%s", nameWithoutExt, ext)
}

// GetMimeType returns the MIME type for a given file extension
func GetMimeType(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".webp":
		return "image/webp"
	default:
		return "application/octet-stream"
	}
}

// ValidateDimensions checks if image dimensions are within acceptable limits
// minWidth, minHeight - minimum dimensions (0 means no minimum)
// maxWidth, maxHeight - maximum dimensions (0 means no maximum)
func ValidateDimensions(info *ImageInfo, minWidth, minHeight, maxWidth, maxHeight int) error {
	if minWidth > 0 && info.Width < minWidth {
		return fmt.Errorf("image width %d is less than minimum %d", info.Width, minWidth)
	}
	if minHeight > 0 && info.Height < minHeight {
		return fmt.Errorf("image height %d is less than minimum %d", info.Height, minHeight)
	}
	if maxWidth > 0 && info.Width > maxWidth {
		return fmt.Errorf("image width %d exceeds maximum %d", info.Width, maxWidth)
	}
	if maxHeight > 0 && info.Height > maxHeight {
		return fmt.Errorf("image height %d exceeds maximum %d", info.Height, maxHeight)
	}
	return nil
}

// IsSquare checks if the image has square dimensions (1:1 aspect ratio)
func IsSquare(info *ImageInfo) bool {
	return info.Width == info.Height
}

// GetAspectRatio calculates the aspect ratio of the image
func GetAspectRatio(info *ImageInfo) float64 {
	if info.Height == 0 {
		return 0
	}
	return float64(info.Width) / float64(info.Height)
}
