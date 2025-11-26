// Package slugify provides utilities for generating URL-friendly slugs
package slugify

import (
	"regexp"
	"strings"
)

var (
	// Transliteration map for Cyrillic to Latin
	translitMap = map[rune]string{
		'а': "a", 'б': "b", 'в': "v", 'г': "g", 'д': "d", 'е': "e", 'ё': "yo",
		'ж': "zh", 'з': "z", 'и': "i", 'й': "y", 'к': "k", 'л': "l", 'м': "m",
		'н': "n", 'о': "o", 'п': "p", 'р': "r", 'с': "s", 'т': "t", 'у': "u",
		'ф': "f", 'х': "h", 'ц': "ts", 'ч': "ch", 'ш': "sh", 'щ': "sch",
		'ъ': "", 'ы': "y", 'ь': "", 'э': "e", 'ю': "yu", 'я': "ya",
		'А': "A", 'Б': "B", 'В': "V", 'Г': "G", 'Д': "D", 'Е': "E", 'Ё': "Yo",
		'Ж': "Zh", 'З': "Z", 'И': "I", 'Й': "Y", 'К': "K", 'Л': "L", 'М': "M",
		'Н': "N", 'О': "O", 'П': "P", 'Р': "R", 'С': "S", 'Т': "T", 'У': "U",
		'Ф': "F", 'Х': "H", 'Ц': "Ts", 'Ч': "Ch", 'Ш': "Sh", 'Щ': "Sch",
		'Ъ': "", 'Ы': "Y", 'Ь': "", 'Э': "E", 'Ю': "Yu", 'Я': "Ya",
	}

	// Regex to match non-alphanumeric characters (except hyphens)
	nonAlphanumeric = regexp.MustCompile(`[^a-zA-Z0-9\-]+`)
	// Regex to match multiple consecutive hyphens
	multipleHyphens = regexp.MustCompile(`-+`)
)

// Generate creates a URL-friendly slug from the input string
// It transliterates Cyrillic characters, converts to lowercase,
// removes special characters, and replaces spaces with hyphens
func Generate(s string) string {
	// Transliterate Cyrillic to Latin
	var result strings.Builder
	for _, char := range s {
		if latin, ok := translitMap[char]; ok {
			result.WriteString(latin)
		} else {
			result.WriteRune(char)
		}
	}

	slug := result.String()

	// Convert to lowercase
	slug = strings.ToLower(slug)

	// Replace spaces with hyphens
	slug = strings.ReplaceAll(slug, " ", "-")

	// Remove non-alphanumeric characters (except hyphens)
	slug = nonAlphanumeric.ReplaceAllString(slug, "")

	// Replace multiple consecutive hyphens with a single hyphen
	slug = multipleHyphens.ReplaceAllString(slug, "-")

	// Trim hyphens from start and end
	slug = strings.Trim(slug, "-")

	// Ensure slug is not too long (max 100 characters)
	if len(slug) > 100 {
		slug = slug[:100]
		// Trim trailing hyphen if cut in the middle
		slug = strings.TrimRight(slug, "-")
	}

	return slug
}
