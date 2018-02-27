package store

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/D-L-M/mem-db/src/crypt"
	"github.com/D-L-M/mem-db/src/data"
	"github.com/D-L-M/mem-db/src/types"
	"github.com/D-L-M/mem-db/src/utils"
)

// Documents are stored in a map, for quick retrieval
var documents = map[string]types.DocumentIndex{}

// Lookups map a field's value against its document
var lookups = map[string][]string{}

// List of all document IDs
var allIds = map[string]string{}

// documentsLock allows locking of the documents map during reads/writes
var documentsLock = sync.RWMutex{}

// lookupsLock allows locking of the lookups map during reads/writes
var lookupsLock = sync.RWMutex{}

// allIdsLock allows locking of the allIds map during reads/writes
var allIdsLock = sync.RWMutex{}

// ParseDocument parses a raw JSON document into an object
func ParseDocument(document []byte) (map[string]interface{}, error) {

	var parsedDocument types.JSONDocument

	err := json.Unmarshal(document, &parsedDocument)

	// Document is not valid JSON
	if err != nil {
		return nil, err
	}

	// Store the document
	return parsedDocument, nil

}

// IndexDocument parses a document (represented by a JSON string) and store it in the document
// map by its ID
func IndexDocument(id string, document []byte, removeFromDiskBeforehand bool) bool {

	parsedDocument, err := ParseDocument(document)

	// If document is not valid JSON
	if err != nil {
		return false
	}

	// First remove any old version that might exist
	RemoveDocument(id, "", removeFromDiskBeforehand)

	// Flatten the document using dot-notation so the inverted index can be
	// created
	flattenedObject := utils.FlattenDocumentToDotNotation(parsedDocument)
	invertedKeys := []string{}

	for fieldDotKey, fieldValue := range flattenedObject {

		sanitisedFieldKey := utils.RemoveNumericIndicesFromFlattenedKey(fieldDotKey)
		keyHash, err := storeKeyHash(id, sanitisedFieldKey, fieldValue, "full")

		if err == nil {
			invertedKeys = append(invertedKeys, keyHash)
		}

		// Now do the same but with words within the value if it's a string
		if valueString, ok := fieldValue.(string); ok {

			_, valueWords := utils.GetPhrasesFromString(valueString)

			for _, valueWord := range valueWords {

				wordKeyHash, err := storeKeyHash(id, sanitisedFieldKey, valueWord, "partial")

				if err == nil {
					invertedKeys = append(invertedKeys, wordKeyHash)
				}

			}

		}

	}

	// Then add the new version in
	documentsLock.Lock()
	allIdsLock.Lock()

	documents[id] = types.DocumentIndex{Document: document, InvertedKeys: invertedKeys}
	allIds[id] = id

	documentsLock.Unlock()
	allIdsLock.Unlock()

	return true

}

// Generate a key/value lookup hash
func generateKeyHash(key string, value interface{}, entryType string) (string, error) {

	// If the value is a string, lowercase it
	if valueString, ok := value.(string); ok {
		value = strings.ToLower(valueString)
	}

	// Hash a JSON representation of the key and value
	keyHashData, err := json.Marshal(types.JSONDocument{"key": key, "value": value, "type": entryType})
	keyHash := crypt.Sha512(keyHashData)

	if err != nil {
		return "", err
	}

	return keyHash, nil

}

// If a document ID has not yet been stored against a lookup of a key/value
// hash, insert it into the lookup map
func storeKeyHash(id string, key string, value interface{}, entryType string) (string, error) {

	keyHash, err := generateKeyHash(key, value, entryType)

	if err != nil {
		return "", err
	}

	if isDocumentInLookup(keyHash, id) == true {
		return "", errors.New("The key hash has already been stored")
	}

	lookupsLock.Lock()
	lookups[keyHash] = append(lookups[keyHash], id)
	lookupsLock.Unlock()

	return keyHash, nil

}

// GetRawDocument gets a raw document by its ID
func GetRawDocument(id string) ([]byte, error) {

	documentsLock.RLock()
	defer documentsLock.RUnlock()

	if document, ok := documents[id]; ok {
		return document.Document, nil
	}

	return nil, errors.New("Document does not exist")

}

// GetDocument gets a document by its ID
func GetDocument(id string) (types.JSONDocument, error) {

	document, err := GetRawDocument(id)

	if err == nil {

		var parsedDocument types.JSONDocument

		err := json.Unmarshal(document, &parsedDocument)

		if err != nil {
			return nil, err
		}

		return parsedDocument, nil

	}

	return nil, errors.New("Document does not exist")

}

// RemoveAllDocuments removes all documents
func RemoveAllDocuments(removeFromDisk bool) {

	data.SetState("truncating")

	documentsLock.Lock()
	lookupsLock.Lock()
	allIdsLock.Lock()

	documents = map[string]types.DocumentIndex{}
	lookups = map[string][]string{}
	allIds = map[string]string{}

	documentsLock.Unlock()
	lookupsLock.Unlock()
	allIdsLock.Unlock()

	if removeFromDisk {

		storageDirectory, err := data.GetStorageDirectory()

		if err != nil {
			log.Fatal(err)
		}

		// Iterate through and delete all flushed JSON files
		files, err := filepath.Glob(storageDirectory + "/*.json")

		if err != nil {
			log.Fatal("Cannot read from storage directory")
		}

		for _, filename := range files {
			os.Remove(filename)
		}

	}

	data.SetState("active")

}

// RemoveDocument removes a document by its ID
func RemoveDocument(id string, filepath string, removeFromDisk bool) {

	// Remove it from any inverted indices using its own inverted lookup
	documentsLock.Lock()

	for _, lookupKey := range documents[id].InvertedKeys {

		// Iterate through all document IDs for the lookup
		lookupsLock.Lock()

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

		lookupsLock.Unlock()

	}

	// Remove the document itself
	allIdsLock.Lock()

	delete(documents, id)
	delete(allIds, id)

	allIdsLock.Unlock()
	documentsLock.Unlock()

	// Optionally also remove the flushed file from disk
	if removeFromDisk && filepath != "" {
		os.Remove(filepath)
	}

}

// Check whether a document ID exists within a given key hash lookup
func isDocumentInLookup(keyHash string, documentID string) bool {

	lookupsLock.RLock()
	defer lookupsLock.RUnlock()

	for _, lookupValue := range lookups[keyHash] {

		if lookupValue == documentID {
			return true
		}

	}

	return false

}

// GetLookups gets the lookup map
func GetLookups() map[string][]string {

	lookupsLock.RLock()
	defer lookupsLock.RUnlock()

	return lookups

}

// GetStats gets stats about the index
func GetStats() map[string]interface{} {

	documentsLock.RLock()
	lookupsLock.RLock()

	stats := map[string]interface{}{
		"totals": map[string]int{
			"documents":        len(documents),
			"inverted_indices": len(lookups)}}

	documentsLock.RUnlock()
	lookupsLock.RUnlock()

	return stats

}
