// Package imageproc provides utilities for validating and processing image files.
package imageproc

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"

	"golang.org/x/image/draw"

	// Register additional image format decoders
	_ "golang.org/x/image/webp"
	_ "image/gif"
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

// DecodeImage decodes an image from a reader
func DecodeImage(r io.Reader) (image.Image, string, error) {
	img, format, err := image.Decode(r)
	if err != nil {
		return nil, "", fmt.Errorf("failed to decode image: %w", err)
	}
	return img, format, nil
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

// MaxDimensions defines maximum allowed dimensions for uploaded images
const (
	MaxImageWidth  = 2000
	MaxImageHeight = 2000
)

// OptimizeImage resizes an image if it exceeds maximum dimensions
// and returns the optimized image as bytes with the appropriate format
func OptimizeImage(img image.Image, format string) ([]byte, error) {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// Check if resizing is needed
	needsResize := width > MaxImageWidth || height > MaxImageHeight

	var resizedImg image.Image
	if needsResize {
		// Calculate new dimensions maintaining aspect ratio
		newWidth, newHeight := calculateResizedDimensions(width, height, MaxImageWidth, MaxImageHeight)

		// Create a new image with the target dimensions
		dst := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))

		// Use high-quality scaling algorithm
		draw.CatmullRom.Scale(dst, dst.Bounds(), img, img.Bounds(), draw.Over, nil)
		resizedImg = dst
	} else {
		resizedImg = img
	}

	// Encode the image to bytes
	var buf bytes.Buffer
	switch format {
	case "jpeg", "jpg":
		// Use quality 85 for good balance between size and quality
		if err := jpeg.Encode(&buf, resizedImg, &jpeg.Options{Quality: 85}); err != nil {
			return nil, fmt.Errorf("failed to encode JPEG: %w", err)
		}
	case "png":
		if err := png.Encode(&buf, resizedImg); err != nil {
			return nil, fmt.Errorf("failed to encode PNG: %w", err)
		}
	default:
		// For other formats (gif, webp), we can't easily re-encode, so return error
		return nil, fmt.Errorf("optimization not supported for format: %s", format)
	}

	return buf.Bytes(), nil
}

// calculateResizedDimensions calculates new dimensions maintaining aspect ratio
func calculateResizedDimensions(width, height, maxWidth, maxHeight int) (int, int) {
	if width <= maxWidth && height <= maxHeight {
		return width, height
	}

	widthRatio := float64(maxWidth) / float64(width)
	heightRatio := float64(maxHeight) / float64(height)

	// Use the smaller ratio to ensure both dimensions fit
	ratio := widthRatio
	if heightRatio < widthRatio {
		ratio = heightRatio
	}

	newWidth := int(float64(width) * ratio)
	newHeight := int(float64(height) * ratio)

	return newWidth, newHeight
}
