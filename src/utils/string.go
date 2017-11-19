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
