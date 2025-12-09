package imageproc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValidExtension(t *testing.T) {
	tests := []struct {
		name     string
		ext      string
		expected bool
	}{
		{
			name:     "valid jpg",
			ext:      ".jpg",
			expected: true,
		},
		{
			name:     "valid jpeg",
			ext:      ".jpeg",
			expected: true,
		},
		{
			name:     "valid png",
			ext:      ".png",
			expected: true,
		},
		{
			name:     "valid gif",
			ext:      ".gif",
			expected: true,
		},
		{
			name:     "valid webp",
			ext:      ".webp",
			expected: true,
		},
		{
			name:     "invalid bmp",
			ext:      ".bmp",
			expected: false,
		},
		{
			name:     "invalid txt",
			ext:      ".txt",
			expected: false,
		},
		{
			name:     "uppercase JPG (should fail - expects lowercase)",
			ext:      ".JPG",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isValidExtension(tt.ext)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGenerateThumbnailFilename(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		expected string
	}{
		{
			name:     "simple jpg",
			filename: "image.jpg",
			expected: "image_thumb.jpg",
		},
		{
			name:     "png file",
			filename: "photo.png",
			expected: "photo_thumb.png",
		},
		{
			name:     "file with multiple dots",
			filename: "my.image.file.jpg",
			expected: "my.image.file_thumb.jpg",
		},
		{
			name:     "file with path",
			filename: "path/to/image.jpg",
			expected: "path/to/image_thumb.jpg",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GenerateThumbnailFilename(tt.filename)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetMimeType(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		expected string
	}{
		{
			name:     "jpg file",
			filename: "image.jpg",
			expected: "image/jpeg",
		},
		{
			name:     "jpeg file",
			filename: "image.jpeg",
			expected: "image/jpeg",
		},
		{
			name:     "png file",
			filename: "image.png",
			expected: "image/png",
		},
		{
			name:     "gif file",
			filename: "image.gif",
			expected: "image/gif",
		},
		{
			name:     "webp file",
			filename: "image.webp",
			expected: "image/webp",
		},
		{
			name:     "unknown extension",
			filename: "file.unknown",
			expected: "application/octet-stream",
		},
		{
			name:     "uppercase extension",
			filename: "IMAGE.JPG",
			expected: "image/jpeg",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetMimeType(tt.filename)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestValidateDimensions(t *testing.T) {
	tests := []struct {
		name      string
		info      *ImageInfo
		minWidth  int
		minHeight int
		maxWidth  int
		maxHeight int
		wantError bool
	}{
		{
			name: "valid dimensions",
			info: &ImageInfo{
				Width:  800,
				Height: 600,
			},
			minWidth:  100,
			minHeight: 100,
			maxWidth:  1920,
			maxHeight: 1080,
			wantError: false,
		},
		{
			name: "width too small",
			info: &ImageInfo{
				Width:  50,
				Height: 600,
			},
			minWidth:  100,
			minHeight: 100,
			maxWidth:  1920,
			maxHeight: 1080,
			wantError: true,
		},
		{
			name: "height too small",
			info: &ImageInfo{
				Width:  800,
				Height: 50,
			},
			minWidth:  100,
			minHeight: 100,
			maxWidth:  1920,
			maxHeight: 1080,
			wantError: true,
		},
		{
			name: "width too large",
			info: &ImageInfo{
				Width:  2000,
				Height: 600,
			},
			minWidth:  100,
			minHeight: 100,
			maxWidth:  1920,
			maxHeight: 1080,
			wantError: true,
		},
		{
			name: "height too large",
			info: &ImageInfo{
				Width:  800,
				Height: 1200,
			},
			minWidth:  100,
			minHeight: 100,
			maxWidth:  1920,
			maxHeight: 1080,
			wantError: true,
		},
		{
			name: "no minimum constraints",
			info: &ImageInfo{
				Width:  10,
				Height: 10,
			},
			minWidth:  0,
			minHeight: 0,
			maxWidth:  1920,
			maxHeight: 1080,
			wantError: false,
		},
		{
			name: "no maximum constraints",
			info: &ImageInfo{
				Width:  5000,
				Height: 5000,
			},
			minWidth:  100,
			minHeight: 100,
			maxWidth:  0,
			maxHeight: 0,
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateDimensions(tt.info, tt.minWidth, tt.minHeight, tt.maxWidth, tt.maxHeight)
			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestIsSquare(t *testing.T) {
	tests := []struct {
		name     string
		info     *ImageInfo
		expected bool
	}{
		{
			name: "square image",
			info: &ImageInfo{
				Width:  100,
				Height: 100,
			},
			expected: true,
		},
		{
			name: "rectangular image",
			info: &ImageInfo{
				Width:  100,
				Height: 200,
			},
			expected: false,
		},
		{
			name: "wide image",
			info: &ImageInfo{
				Width:  200,
				Height: 100,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsSquare(tt.info)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetAspectRatio(t *testing.T) {
	tests := []struct {
		name     string
		info     *ImageInfo
		expected float64
	}{
		{
			name: "square image",
			info: &ImageInfo{
				Width:  100,
				Height: 100,
			},
			expected: 1.0,
		},
		{
			name: "16:9 image",
			info: &ImageInfo{
				Width:  1920,
				Height: 1080,
			},
			expected: 1.7777777777777777,
		},
		{
			name: "4:3 image",
			info: &ImageInfo{
				Width:  800,
				Height: 600,
			},
			expected: 1.3333333333333333,
		},
		{
			name: "portrait image",
			info: &ImageInfo{
				Width:  600,
				Height: 800,
			},
			expected: 0.75,
		},
		{
			name: "zero height",
			info: &ImageInfo{
				Width:  100,
				Height: 0,
			},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetAspectRatio(tt.info)
			assert.Equal(t, tt.expected, result)
		})
	}
}
