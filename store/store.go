package store


import (
	"encoding/json"
	maputils "../utils"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"../types"
)


// Documents are stored in a map, for quick retrieval
var documents = map[string]types.JsonDocument{}


// Lookups map a field's value against its document
var lookups = map[string][]string{}


// Parse a document (represented by a JSON string) and store it in the document
// map by its ID
func IndexDocument(id string, document []byte) bool {

	var parsedDocument types.JsonDocument

	err := json.Unmarshal(document, &parsedDocument)

	// Document is not valid JSON
	if err != nil {
		
		return false

	// Store the document
	} else {

		documents[id] = parsedDocument

		// Flatten the document using dot-notation
		flattenedObject := maputils.FlattenDocumentToDotNotation(parsedDocument)

		for fieldDotKey, fieldValue := range flattenedObject {

			hasher         := sha256.New()
			keyHashData, _ := json.Marshal(types.JsonDocument{"key": fieldDotKey, "value": fieldValue})
			
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

}


// Get a document by its ID
func GetDocument(id string) (types.JsonDocument, error) {

	if document, ok := documents[id]; ok {
		return document, nil
	}

	return nil, errors.New("Document does not exist")

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