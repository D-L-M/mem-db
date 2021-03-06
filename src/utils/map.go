package utils

import (
	"strconv"
)

// FlattenDocumentToDotNotation flattens a map to a key/value map using dot
// notation to represent nested layers of keys
func FlattenDocumentToDotNotation(document map[string]interface{}) map[string]interface{} {

	flattenedMap := make(map[string]interface{})

	for key, value := range document {

		switch child := value.(type) {

		// Nested maps an go straight back through
		case map[string]interface{}:

			subMap := FlattenDocumentToDotNotation(child)

			for subKey, subValue := range subMap {
				flattenedMap[key+"."+subKey] = subValue
			}

		// Slices need to first be converted to maps by casting their
		// numeric indices as strings
		case []interface{}:

			sliceMap := make(map[string]interface{})

			for subKey, subValue := range child {
				sliceMap[strconv.Itoa(subKey)] = subValue
			}

			// Then send through as normal
			subMap := FlattenDocumentToDotNotation(sliceMap)

			for subKey, subValue := range subMap {
				flattenedMap[key+"."+subKey] = subValue
			}

		// Any other value is a leaf node
		default:
			flattenedMap[key] = value

		}

	}

	return flattenedMap

}

// MapHasKey checks whether a map has a key
func MapHasKey(inputMap *map[string]interface{}, key string) bool {

	if _, ok := (*inputMap)[key]; ok {
		return true
	}

	return false

}
