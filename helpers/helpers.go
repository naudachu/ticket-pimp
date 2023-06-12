package helpers

import (
	"regexp"
	"strings"
)

func GitNaming(input string) string {
	// Remove leading and trailing whitespace
	input = strings.TrimSpace(input)

	// Replace non-Latin letters with spaces
	reg := regexp.MustCompile("[^a-zA-Z0-9]+")
	input = strings.TrimSpace(reg.ReplaceAllString(input, " "))

	// Split into words
	words := strings.Fields(input)
	for i, word := range words {
		words[i] = strings.ToLower(word)
	}

	// Join words and return
	return strings.Join(words, "-")
}
