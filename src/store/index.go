package store


import (
	"encoding/json"
	"../utils"
	"../crypt"
	"errors"
	"../types"
	"strings"
)


// Documents are stored in a map, for quick retrieval
var documents = map[string]types.DocumentIndex{}


// Lookups map a field's value against its document
var lookups = map[string][]string{}


// Parse a raw JSON document into an object
func ParseDocument(document []byte) (map[string]interface{}, error) {

	var parsedDocument types.JsonDocument

	error := json.Unmarshal(document, &parsedDocument)

	// Document is not valid JSON
	if error != nil {

		return nil, errors.New("Document is not valid JSON")

	// Store the document
	} else {

		return parsedDocument, nil

	}

}


// Parse a document (represented by a JSON string) and store it in the document
// map by its ID
func IndexDocument(id string, document []byte) bool {

	parsedDocument, error := ParseDocument(document)

	// Document is not valid JSON
	if error != nil {

		return false

	// Store the document
	} else {

		// First remove any old version that might exist
		RemoveDocument(id)

		// Flatten the document using dot-notation so the inverted index can be
		// created
		flattenedObject := utils.FlattenDocumentToDotNotation(parsedDocument)
		invertedKeys	:= []string{}

		for fieldDotKey, fieldValue := range flattenedObject {

			sanitisedFieldKey := utils.RemoveNumericIndicesFromFlattenedKey(fieldDotKey)
			keyHash	          := storeKeyHash(id, sanitisedFieldKey, fieldValue, "full")
			invertedKeys       = append(invertedKeys, keyHash)

			// Now do the same but with words within the value if it's a string
			if valueString, ok := fieldValue.(string); ok {

				valueWords := strings.Split(string(valueString), " ")

				for _, valueWord := range valueWords {

					if valueWord != "" && valueWord != " " {
						wordKeyHash  := storeKeyHash(id, sanitisedFieldKey, valueWord, "partial")
						invertedKeys  = append(invertedKeys, wordKeyHash)
					}

				}

			}

		}

		// Then add the new version in
		documents[id] = types.DocumentIndex{Document: document, InvertedKeys: invertedKeys}

		return true

	}

}


// If a document ID has not yet been stored against a lookup of a key/value
// hash, inert it into the lookup map
func storeKeyHash(id string, key string, value interface{}, entryType string) string {

	keyHashData, error := json.Marshal(types.JsonDocument{"key": key, "value": value, "type": entryType})
	keyHash			:= crypt.Sha256(keyHashData)

	if error == nil && isDocumentInLookup(keyHash, id) == false {
		lookups[keyHash] = append(lookups[keyHash], id)
	}

	return keyHash

}


// Get a raw document by its ID
func GetRawDocument(id string) ([]byte, error) {

	if document, ok := documents[id]; ok {
		return document.Document, nil
	}

	return nil, errors.New("Document does not exist")

}


// Get a document by its ID
func GetDocument(id string) (types.JsonDocument, error) {

	document, error := GetRawDocument(id)

	if error == nil {

		var parsedDocument types.JsonDocument

		error := json.Unmarshal(document, &parsedDocument)

		if error != nil {

			return nil, errors.New("Document is corrupted")

		} else {

			return parsedDocument, nil

		}

	}

	return nil, errors.New("Document does not exist")

}


// Remove a document by its ID
func RemoveDocument(id string) {

	// Remove it from any inverted indices using its own inverted lookup
	for _, lookupKey := range documents[id].InvertedKeys {

		// Iterate through all document IDs for the lookup
		for i, lookupValue := range lookups[lookupKey] {

			// If the ID matches the document that's being removed, take
			// that ID out of the lookup slice
			if lookupValue == id {

				lookups[lookupKey] = append(lookups[lookupKey][:i], lookups[lookupKey][i+1:]...)

				// Also remove the whole inverted index if it's now empty
				if len(lookups[lookupKey]) == 0 {
					delete(lookups, lookupKey)
				}

			}

		}

	}

	// Remove the document itself
	delete(documents, id)

}


// Check whether a document ID exists within a given key hash lookup
func isDocumentInLookup(keyHash string, documentId string) bool {

	for _, lookupValue := range lookups[keyHash] {

		if lookupValue == documentId {
			return true
		}
	}

	return false

}


// Get lookup map
func GetLookups() map[string][]string {

	return lookups

}


// Get stats about the index
func GetStats() map[string]interface{} {

	stats := make(map[string]interface{})

	stats["total_documents"]		= len(documents)
	stats["total_inverted_indices"] = len(lookups)

	return stats

}