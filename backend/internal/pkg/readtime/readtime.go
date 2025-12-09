// Package readtime provides utilities for estimating reading time of text content.
package readtime

import (
	"math"
	"strings"
	"unicode"
)

const (
	// Average reading speed in words per minute
	defaultWordsPerMinute = 200
)

// Calculate estimates the reading time for a given text in minutes.
// It counts words (sequences of letters/numbers separated by spaces)
// and divides by average reading speed (200 words/minute by default).
func Calculate(text string) int {
	if text == "" {
		return 0
	}

	wordCount := countWords(text)
	minutes := float64(wordCount) / float64(defaultWordsPerMinute)

	// Round up to at least 1 minute
	result := int(math.Ceil(minutes))
	if result < 1 {
		return 1
	}

	return result
}

// countWords counts the number of words in the text.
// A word is defined as a sequence of letters or numbers.
func countWords(text string) int {
	if text == "" {
		return 0
	}

	words := 0
	inWord := false

	for _, r := range text {
		if unicode.IsLetter(r) || unicode.IsNumber(r) {
			if !inWord {
				words++
				inWord = true
			}
		} else {
			inWord = false
		}
	}

	return words
}

// CalculateWithSpeed calculates reading time with a custom reading speed.
// wordsPerMinute should be positive, otherwise defaults to 200.
func CalculateWithSpeed(text string, wordsPerMinute int) int {
	if text == "" {
		return 0
	}

	if wordsPerMinute <= 0 {
		wordsPerMinute = defaultWordsPerMinute
	}

	wordCount := countWords(text)
	minutes := float64(wordCount) / float64(wordsPerMinute)

	result := int(math.Ceil(minutes))
	if result < 1 {
		return 1
	}

	return result
}

// EstimateMarkdown estimates reading time for markdown content.
// It strips markdown syntax before counting words.
func EstimateMarkdown(markdown string) int {
	// Simple markdown stripping - remove common syntax
	text := markdown

	// Remove code blocks
	text = removeCodeBlocks(text)

	// Remove inline code
	text = strings.ReplaceAll(text, "`", "")

	// Remove links [text](url) -> text
	text = removeLinks(text)

	// Remove images ![alt](url)
	text = removeImages(text)

	// Remove headers #
	text = strings.ReplaceAll(text, "#", "")

	// Remove bold/italic **text** or *text*
	text = strings.ReplaceAll(text, "**", "")
	text = strings.ReplaceAll(text, "*", "")

	return Calculate(text)
}

// removeCodeBlocks removes ```code``` blocks from text
func removeCodeBlocks(text string) string {
	result := strings.Builder{}
	inCodeBlock := false

	lines := strings.Split(text, "\n")
	for _, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "```") {
			inCodeBlock = !inCodeBlock
			continue
		}
		if !inCodeBlock {
			result.WriteString(line)
			result.WriteString("\n")
		}
	}

	return result.String()
}

// removeLinks removes [text](url) and keeps only text
func removeLinks(text string) string {
	result := strings.Builder{}
	i := 0

	for i < len(text) {
		if text[i] == '[' {
			// Find closing ]
			j := strings.IndexByte(text[i:], ']')
			if j != -1 && i+j+1 < len(text) && text[i+j+1] == '(' {
				// Find closing )
				k := strings.IndexByte(text[i+j+1:], ')')
				if k != -1 {
					// Extract text between [ and ]
					linkText := text[i+1 : i+j]
					result.WriteString(linkText)
					i = i + j + 1 + k + 1
					continue
				}
			}
		}
		result.WriteByte(text[i])
		i++
	}

	return result.String()
}

// removeImages removes ![alt](url) from text
func removeImages(text string) string {
	result := strings.Builder{}
	i := 0

	for i < len(text) {
		if i+1 < len(text) && text[i] == '!' && text[i+1] == '[' {
			// Find closing ]
			j := strings.IndexByte(text[i+1:], ']')
			if j != -1 && i+1+j+1 < len(text) && text[i+1+j+1] == '(' {
				// Find closing )
				k := strings.IndexByte(text[i+1+j+1:], ')')
				if k != -1 {
					// Skip the entire image syntax
					i = i + 1 + j + 1 + k + 1
					continue
				}
			}
		}
		result.WriteByte(text[i])
		i++
	}

	return result.String()
}
