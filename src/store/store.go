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
var documents = map[string][]byte{}


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

		// First remove any old version that might exist
		RemoveDocument(id)

		// Then add the new version in
		documents[id] = document

		// Flatten the document using dot-notation so the inverted index can be
		// created
		flattenedObject := maputils.FlattenDocumentToDotNotation(parsedDocument)

		for fieldDotKey, fieldValue := range flattenedObject {

			hasher         := sha256.New()
			keyHashData, _ := json.Marshal(types.JsonDocument{"key": fieldDotKey, "value": fieldValue})
			
			hasher.Write(keyHashData)

			keyHash := base64.URLEncoding.EncodeToString(hasher.Sum(nil))

			// If the document ID has not yet been stored against a lookup of
			// its hashed key and value, store it now
			//
			// TODO Also store in another lookup where the ID is the key, so on
			// reindex these lookups can be removed first
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
		
		var parsedDocument types.JsonDocument

		err := json.Unmarshal(document, &parsedDocument)
		
		if err != nil {			

			return nil, errors.New("Document is corrupted")

		} else {
		
			return parsedDocument, nil

		}
	}

	return nil, errors.New("Document does not exist")

}


// Remove a document by its ID
// TODO Also remove from lookups
func RemoveDocument(id string) {

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

	stats["total_documents"]        = len(documents)
	stats["total_inverted_indices"] = len(lookups)

	return stats

}