package readtime

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculate(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		expected int
	}{
		{
			name:     "empty text",
			text:     "",
			expected: 0,
		},
		{
			name:     "very short text (less than 1 minute)",
			text:     "Hello world",
			expected: 1,
		},
		{
			name:     "exactly 200 words",
			text:     strings.Repeat("word ", 200),
			expected: 1,
		},
		{
			name:     "201 words (should be 2 minutes)",
			text:     strings.Repeat("word ", 201),
			expected: 2,
		},
		{
			name:     "400 words",
			text:     strings.Repeat("word ", 400),
			expected: 2,
		},
		{
			name:     "401 words (should be 3 minutes)",
			text:     strings.Repeat("word ", 401),
			expected: 3,
		},
		{
			name:     "text with punctuation",
			text:     "Hello, world! How are you? I'm fine, thank you.",
			expected: 1,
		},
		{
			name:     "text with numbers",
			text:     "The year is 2024 and we have 365 days.",
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Calculate(tt.text)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCountWords(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		expected int
	}{
		{
			name:     "empty string",
			text:     "",
			expected: 0,
		},
		{
			name:     "single word",
			text:     "word",
			expected: 1,
		},
		{
			name:     "multiple words",
			text:     "hello world foo bar",
			expected: 4,
		},
		{
			name:     "words with punctuation",
			text:     "Hello, world! How are you?",
			expected: 5,
		},
		{
			name:     "words with numbers",
			text:     "I have 2 cats and 3 dogs",
			expected: 6,
		},
		{
			name:     "multiple spaces",
			text:     "hello    world",
			expected: 2,
		},
		{
			name:     "tabs and newlines",
			text:     "hello\tworld\nfoo\rbar",
			expected: 4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := countWords(tt.text)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCalculateWithSpeed(t *testing.T) {
	tests := []struct {
		name           string
		text           string
		wordsPerMinute int
		expected       int
	}{
		{
			name:           "empty text",
			text:           "",
			wordsPerMinute: 100,
			expected:       0,
		},
		{
			name:           "custom speed 100 wpm",
			text:           strings.Repeat("word ", 100),
			wordsPerMinute: 100,
			expected:       1,
		},
		{
			name:           "custom speed 100 wpm (101 words)",
			text:           strings.Repeat("word ", 101),
			wordsPerMinute: 100,
			expected:       2,
		},
		{
			name:           "invalid speed (defaults to 200)",
			text:           strings.Repeat("word ", 200),
			wordsPerMinute: 0,
			expected:       1,
		},
		{
			name:           "negative speed (defaults to 200)",
			text:           strings.Repeat("word ", 200),
			wordsPerMinute: -100,
			expected:       1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateWithSpeed(tt.text, tt.wordsPerMinute)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestEstimateMarkdown(t *testing.T) {
	tests := []struct {
		name     string
		markdown string
		expected int
	}{
		{
			name:     "empty markdown",
			markdown: "",
			expected: 0,
		},
		{
			name:     "plain text",
			markdown: "Hello world",
			expected: 1,
		},
		{
			name:     "with headers",
			markdown: "# Header\n\nSome text here",
			expected: 1,
		},
		{
			name:     "with bold",
			markdown: "This is **bold** text",
			expected: 1,
		},
		{
			name:     "with italic",
			markdown: "This is *italic* text",
			expected: 1,
		},
		{
			name:     "with inline code",
			markdown: "Use `console.log()` to print",
			expected: 1,
		},
		{
			name:     "with code block",
			markdown: "Some text\n```go\nfunc main() {\n}\n```\nMore text",
			expected: 1,
		},
		{
			name:     "with links",
			markdown: "Check [this link](https://example.com) out",
			expected: 1,
		},
		{
			name:     "with images",
			markdown: "![alt text](image.png)",
			expected: 0,
		},
		{
			name:     "complex markdown",
			markdown: "# Title\n\n**Bold** and *italic* text.\n\n```\ncode block\n```\n\n[Link](url) and ![image](img.png)\n\n" + strings.Repeat("word ", 200),
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := EstimateMarkdown(tt.markdown)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestRemoveCodeBlocks(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		expected string
	}{
		{
			name:     "no code blocks",
			text:     "plain text",
			expected: "plain text\n",
		},
		{
			name:     "single code block",
			text:     "before\n```\ncode\n```\nafter",
			expected: "before\nafter\n",
		},
		{
			name:     "multiple code blocks",
			text:     "text1\n```\ncode1\n```\ntext2\n```\ncode2\n```\ntext3",
			expected: "text1\ntext2\ntext3\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := removeCodeBlocks(tt.text)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestRemoveLinks(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		expected string
	}{
		{
			name:     "no links",
			text:     "plain text",
			expected: "plain text",
		},
		{
			name:     "single link",
			text:     "Check [this link](https://example.com) out",
			expected: "Check this link out",
		},
		{
			name:     "multiple links",
			text:     "[link1](url1) and [link2](url2)",
			expected: "link1 and link2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := removeLinks(tt.text)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestRemoveImages(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		expected string
	}{
		{
			name:     "no images",
			text:     "plain text",
			expected: "plain text",
		},
		{
			name:     "single image",
			text:     "Look at this ![alt](image.png) image",
			expected: "Look at this  image",
		},
		{
			name:     "multiple images",
			text:     "![img1](url1) and ![img2](url2)",
			expected: " and ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := removeImages(tt.text)
			assert.Equal(t, tt.expected, result)
		})
	}
}
