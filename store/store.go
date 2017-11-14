package store


import (
	"encoding/json"
	"../output"
	"github.com/17twenty/flatter"
	"crypto/sha256"
	"encoding/base64"
)


// Documents could be any format once they're parsed, so make a generic object
// type to define them
type Document map[string]interface {}


// Documents are stored in a map, for quick retrieval
var documents = map[string]Document {}


// Lookups map a field's value against its document
var lookups = map[string][]string {}


// Parse a document (represented by a JSON string) and store it in the document
// map by its ID
func IndexDocument(id string, document []byte) bool {

	var parsed Document

	err := json.Unmarshal(document, &parsed)

	// Document is not valid JSON
	if err != nil {
		
		output.Write("Request body is not valid JSON")

	// Store the document
	} else {

		documents[id] = parsed

		// Flatten the document using dot-notation
		flattenedObject := flatten.Flatten(parsed)

		for fieldDotKey, fieldValue := range flattenedObject {

			hasher         := sha256.New()
			keyHashData, _ := json.Marshal(Document{"key": fieldDotKey, "value": fieldValue})
			
			hasher.Write(keyHashData)

			keyHash := base64.URLEncoding.EncodeToString(hasher.Sum(nil))

			// If the document ID has not yet been stored against a lookup of
			// its hashed key and value, store it now
			if isDocumentInLookup(keyHash, id) == false {
				lookups[keyHash] = append(lookups[keyHash], id)
			}

		}

		return true

	}

	return false

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