package store

import (
	"encoding/json"
	"errors"
	"reflect"
	"strings"

	"../crypt"
	"../types"
	"../utils"
)

// Documents are stored in a map, for quick retrieval
var documents = map[string]types.DocumentIndex{}

// Lookups map a field's value against its document
var lookups = map[string][]string{}

// List of all document IDs
var allIds = map[string]string{}

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
		invertedKeys := []string{}

		for fieldDotKey, fieldValue := range flattenedObject {

			sanitisedFieldKey := utils.RemoveNumericIndicesFromFlattenedKey(fieldDotKey)
			keyHash, error := storeKeyHash(id, sanitisedFieldKey, fieldValue, "full")

			if error == nil {
				invertedKeys = append(invertedKeys, keyHash)
			}

			// Now do the same but with words within the value if it's a string
			if valueString, ok := fieldValue.(string); ok {

				valueWords := utils.GetWordsFromString(valueString)

				for _, valueWord := range valueWords {

					wordKeyHash, error := storeKeyHash(id, sanitisedFieldKey, valueWord, "partial")

					if error == nil {
						invertedKeys = append(invertedKeys, wordKeyHash)
					}

				}

			}

		}

		// Then add the new version in
		documents[id] = types.DocumentIndex{Document: document, InvertedKeys: invertedKeys}
		allIds[id] = id

		return true

	}

}

// Generate a key/value lookup hash
func generateKeyHash(key string, value interface{}, entryType string) (string, error) {

	// If the value is a string, lowercase it
	if valueString, ok := value.(string); ok {
		value = strings.ToLower(valueString)
	}

	// Hash a JSON representation of the key and value
	keyHashData, error := json.Marshal(types.JsonDocument{"key": key, "value": value, "type": entryType})
	keyHash := crypt.Sha256(keyHashData)

	if error != nil {
		return "", error
	}

	return keyHash, nil

}

// If a document ID has not yet been stored against a lookup of a key/value
// hash, inert it into the lookup map
func storeKeyHash(id string, key string, value interface{}, entryType string) (string, error) {

	keyHash, error := generateKeyHash(key, value, entryType)

	if error != nil {
		return "", error
	}

	if isDocumentInLookup(keyHash, id) == true {
		return "", errors.New("The key hash has already been stored")
	}

	lookups[keyHash] = append(lookups[keyHash], id)

	return keyHash, nil

}

// Search for documents matching a single criterion
func searchCriterion(criterion map[string]interface{}) []string {

	result := []string{}

	for searchType, searchCriterion := range criterion {

		if remappedSearchCriterion, ok := searchCriterion.(map[string]interface{}); ok {

			for searchKey, searchValue := range remappedSearchCriterion {

				// Figure out what kind of search to do
				searchTypeName := "full"

				if searchType == "contains" || searchType == "not_contains" {
					searchTypeName = "partial"
				}

				// If the value is a string, lowercase it
				if valueString, ok := searchValue.(string); ok {
					searchValue = strings.ToLower(valueString)
				}

				// Generate a key hash for the criterion and return any document
				// IDs that have been stored against it
				keyHash, error := generateKeyHash(searchKey, searchValue, searchTypeName)

				if error == nil {

					if documentIds, ok := lookups[keyHash]; ok {

						// If the match is exclusive, build up a list of IDs not
						// found by the lookup
						if searchType == "not_equals" || searchType == "not_contains" {

							exclusiveIds := []string{}

							for _, singleId := range allIds {

								if utils.StringInSlice(singleId, documentIds) == false {
									exclusiveIds = append(exclusiveIds, singleId)
								}

							}

							return exclusiveIds

						}

						// If the match is inclusive, just return the IDs as they
						// are
						return documentIds

					}

				}

			}

		}

	}

	return result

}

// Search for document IDs by evaluating a set of JSON criteria
func SearchDocumentIds(criteria map[string][]interface{}) []string {

	result := []string{}
	ids := [][]string{}

	for groupType, groupCriteria := range criteria {

		for _, criterion := range groupCriteria {

			// Figure out what kind of criterion is being dealt with
			nestedCriterion := criterion.(map[string]interface{})
			isNested := false

			for nestedKey, nestedValue := range nestedCriterion {

				// Nested AND/OR criterion
				if strings.ToLower(nestedKey) == "and" || strings.ToLower(nestedKey) == "or" {

					isNested = true

					switch reflect.TypeOf(nestedValue).Kind() {

					case reflect.Slice:

						remappedAndOrCriteria := map[string][]interface{}{}

						for _, criteriaSlice := range reflect.ValueOf(nestedValue).Interface().([]interface{}) {
							remappedAndOrCriteria[nestedKey] = append(remappedAndOrCriteria[nestedKey], criteriaSlice)
						}

						ids = append(ids, SearchDocumentIds(remappedAndOrCriteria))

					}

					break

				}

			}

			// Regular criterion
			if isNested == false {

				regularCriterion := criterion.(map[string]interface{})
				ids = append(ids, searchCriterion(regularCriterion))

			}

		}

		// OR -- combine the IDs, deduplicating where necessary
		if strings.ToLower(groupType) == "or" {

			for _, idGroup := range ids {

				for _, id := range idGroup {

					if utils.StringInSlice(id, result) == false {
						result = append(result, id)
					}

				}

			}

			// AND -- compile a list of IDs appearing in all ID lists
		} else if strings.ToLower(groupType) == "and" {

			result = utils.StringSliceIntersection(ids)

		}

	}

	return result

}

// Search for documents by evaluating a set of JSON criteria
func SearchDocuments(criteria map[string][]interface{}) map[string]types.JsonDocument {

	ids := []string{}

	// If no criteria, retrieve everything
	if len(criteria) == 0 {

		for _, id := range allIds {
			ids = append(ids, id)
		}

		// Otherwise filter by the actual criteria
	} else {

		ids = SearchDocumentIds(criteria)

	}

	// Convert document IDs to actual documents
	results := map[string]types.JsonDocument{}

	for _, id := range ids {

		document, error := GetDocument(id)

		if error == nil {
			results[id] = document
		}

	}

	return results

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
	delete(allIds, id)

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

	stats["total_documents"] = len(documents)
	stats["total_inverted_indices"] = len(lookups)

	return stats

}
