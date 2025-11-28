package slugify

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerate(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Basic Latin",
			input:    "Hello World",
			expected: "hello-world",
		},
		{
			name:     "Cyrillic to Latin",
			input:    "Привет Мир",
			expected: "privet-mir",
		},
		{
			name:     "Mixed and Special Characters",
			input:    "Go 1.25 & Rust!",
			expected: "go-125-rust",
		},
		{
			name:     "Multiple Spaces and Hyphens",
			input:    "Hello   World---Go",
			expected: "hello-world-go",
		},
		{
			name:     "Trim Hyphens",
			input:    "-Start End-",
			expected: "start-end",
		},
		{
			name:     "Long String Truncation",
			input:    strings.Repeat("a", 150),
			expected: strings.Repeat("a", 100),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Generate(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
