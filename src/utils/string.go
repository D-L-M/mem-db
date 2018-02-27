package utils

import (
	"regexp"
	"strings"

	porterstemmer "github.com/agonopol/go-stem"
)

// RemoveNumericIndicesFromFlattenedKey strips numeric indices from a
// dot-notation-flattened map key
func RemoveNumericIndicesFromFlattenedKey(dotNotationKey string) string {

	pattern := regexp.MustCompile(`(^[0-9]+\.)|(\.[0-9]+\.)|(\.[0-9]+$)`)
	sanitised := pattern.ReplaceAllString(dotNotationKey, `.`)
	trimmed := strings.Trim(sanitised, ".")

	return trimmed

}

// ContainsPunctuation checks whether an input string contains punctuation
func ContainsPunctuation(inputString string) bool {

	pattern := regexp.MustCompile(`([^\w\s])`)

	return pattern.MatchString(inputString)

}

// PadPunctuationWithSpaces pads any punctuation in a string with spaces so
// words are clearly defined
func PadPunctuationWithSpaces(inputString string) string {

	pattern := regexp.MustCompile(`(\w\S+\w)|(\w+)|(\s*\.{3}\s*)|(\s*[^\w\s]\s*)|\s+`)
	padded := pattern.ReplaceAllString(inputString, ` $0 `)

	return padded

}

// GetPhrasesFromString splits a string into a slice of individual words
func GetPhrasesFromString(inputString string) ([]string, []string) {

	phraseWordLimit := 3 // TODO: Move into config
	words := strings.Split(PadPunctuationWithSpaces(inputString), " ")
	validPlainWords := []string{}
	validStemmedWords := []string{}
	plainResult := []string{}
	stemmedResult := []string{}

	// Remove spaces
	for _, word := range words {

		if word != "" {
			validPlainWords = append(validPlainWords, word)
			stemmedWord := string(porterstemmer.Stem([]byte(word)))
			validStemmedWords = append(validStemmedWords, stemmedWord)
		}

	}

	// Build up a list of phrases, starting at one word each and building to
	// the phrase word limit
	for i := 1; i <= phraseWordLimit; i++ {

		for j := 0; j <= (len(validPlainWords) - i); j++ {

			plainPhrase := strings.Join(validPlainWords[j:(j+i)], " ")
			plainResult = append(plainResult, plainPhrase)
			stemmedPhrase := strings.Join(validStemmedWords[j:(j+i)], " ")
			stemmedResult = append(stemmedResult, stemmedPhrase)

		}

	}

	return plainResult, stemmedResult

}
