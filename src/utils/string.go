package utils

import (
	"regexp"
	"strings"
)

// Strip numeric indices from a dot-notation-flattened map key
func RemoveNumericIndicesFromFlattenedKey(dotNotationKey string) string {

	pattern := regexp.MustCompile(`(^[0-9]+\.)|(\.[0-9]+\.)|(\.[0-9]+$)`)
	sanitised := pattern.ReplaceAllString(dotNotationKey, `.`)
	trimmed := strings.Trim(sanitised, ".")

	return trimmed

}

// Pad any punctuation in a string with spaces so words are clearly defined
func PadPunctuationWithSpaces(inputString string) string {

	pattern := regexp.MustCompile(`(\w\S+\w)|(\w+)|(\s*\.{3}\s*)|(\s*[^\w\s]\s*)|\s+`)
	padded := pattern.ReplaceAllString(inputString, ` $0 `)

	return padded

}

// Split a string into a slice of individual words
func GetWordsFromString(inputString string) []string {

	words := strings.Split(PadPunctuationWithSpaces(inputString), " ")
	result := []string{}

	for _, word := range words {

		if word != "" {
			result = append(result, word)
		}

	}

	return result

}
