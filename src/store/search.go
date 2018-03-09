package store

import (
	"reflect"
	"sort"
	"strings"

	"github.com/D-L-M/jsonserver"
	"github.com/D-L-M/mem-db/src/data"
	"github.com/D-L-M/mem-db/src/utils"
	"github.com/kljensen/snowball"
)

// significantTermsSort is a custom sorting algorithm for significant terms --
// sort by document count (highest first), then term (alphabetically)
type significantTermsSort []map[string]interface{}

func (significantTerms significantTermsSort) Len() int {

	return len(significantTerms)

}

func (significantTerms significantTermsSort) Swap(i, j int) {

	significantTerms[i], significantTerms[j] = significantTerms[j], significantTerms[i]

}

func (significantTerms significantTermsSort) Less(i, j int) bool {

	if significantTerms[i]["doc_count"].(int) > significantTerms[j]["doc_count"].(int) {
		return true
	}

	if significantTerms[i]["doc_count"].(int) < significantTerms[j]["doc_count"].(int) {
		return false
	}

	return strings.ToLower(significantTerms[i]["term"].(string)) < strings.ToLower(significantTerms[j]["term"].(string))

}

// DiscoverSignificantTerms returns a slice of significant terms discovered in
// a specific field of a slice of documents, compared to the rest of the index
func DiscoverSignificantTerms(targetedDocuments *[]jsonserver.JSON, field string, percentageThreshold int, minimumOccurrences float64) []map[string]interface{} {

	collectedFragmentHashes := map[string]string{}
	fragmentHashCounts := map[string]int{}
	result := []map[string]interface{}{}

	// Get counts from the documents provided
	for _, document := range *targetedDocuments {

		termFragments, err := getTermFragmentHashesForDocumentField(document["document"].(jsonserver.JSON), field, true, false)

		if err != nil {
			continue
		}

		for hashedTerm, plainTerm := range termFragments {
			collectedFragmentHashes[hashedTerm] = plainTerm
			fragmentHashCounts[hashedTerm]++
		}

	}

	// Compare against the rest of the index
	for hashedTerm, hashTermCount := range fragmentHashCounts {

		if utils.StringInSlice(collectedFragmentHashes[hashedTerm], data.StopWords) {
			continue
		}

		if ((float64(hashTermCount) / float64(len(*targetedDocuments))) * 100) < minimumOccurrences {
			continue
		}

		if utils.ContainsPunctuation(collectedFragmentHashes[hashedTerm]) {
			continue
		}

		targetedFrequencyPerDocument := (float64(hashTermCount) / float64(len(*targetedDocuments)))
		comparisonFrequencyPerDocument := (float64(len(lookups[hashedTerm])) / float64(len(documents)))

		if ((targetedFrequencyPerDocument / comparisonFrequencyPerDocument) * 100) >= float64(percentageThreshold) {
			result = append(result, map[string]interface{}{"term": collectedFragmentHashes[hashedTerm], "doc_count": hashTermCount})
		}

	}

	sort.Sort(significantTermsSort(result))

	return result

}

// Get the hashes and plain forms of all terms for a specific field in a
// document
func getTermFragmentHashesForDocumentField(document jsonserver.JSON, field string, stemHash bool, stemValue bool) (map[string]string, error) {

	result := map[string]string{}
	flattenedObject := utils.FlattenDocumentToDotNotation(document)

	for fieldKey, fieldValue := range flattenedObject {

		sanitisedFieldKey := utils.RemoveNumericIndicesFromFlattenedKey(fieldKey)

		if field == sanitisedFieldKey {

			if valueString, ok := fieldValue.(string); ok {

				valueWords, stemmedValueWords := utils.GetPhrasesFromString(valueString)

				for i, valueWord := range valueWords {

					// Decide which version of the word to use for the hash and
					// the stored value
					hashWordValue := valueWord
					storedWordValue := valueWord

					if stemHash {
						hashWordValue = stemmedValueWords[i]
					}

					if stemValue {
						storedWordValue = stemmedValueWords[i]
					}

					// Generate a hash of the field value
					wordKeyHash, err := generateKeyHash(sanitisedFieldKey, hashWordValue, "partial")

					if err == nil {
						result[wordKeyHash] = strings.ToLower(storedWordValue)
					}

				}

			}

		}

	}

	return result, nil

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

				// Stem words for partial matches
				if searchType == "contains" || searchType == "not_contains" {

					partialWords := strings.Split(utils.PadPunctuationWithSpaces(searchValue.(string)), " ")
					stemmedPhrase := []string{}

					for _, partialWord := range partialWords {

						stemmedWord, err := snowball.Stem(partialWord, "english", true)

						if err == nil && stemmedWord != "" {
							stemmedPhrase = append(stemmedPhrase, stemmedWord)
						}

					}

					searchValue = strings.Join(stemmedPhrase, " ")

				}

				// Generate a key hash for the criterion and return any document
				// IDs that have been stored against it
				keyHash, err := generateKeyHash(searchKey, searchValue, searchTypeName)

				if err == nil {

					lookupsLock.RLock()

					if documentIds, ok := lookups[keyHash]; ok {

						lookupsLock.RUnlock()

						// If the match is exclusive, build up a list of IDs not
						// found by the lookup
						if searchType == "not_equals" || searchType == "not_contains" {

							exclusiveIds := []string{}

							allIdsLock.RLock()

							for _, singleID := range allIds {

								if utils.StringInSlice(singleID, documentIds) == false {
									exclusiveIds = append(exclusiveIds, singleID)
								}

							}

							allIdsLock.RUnlock()

							return exclusiveIds

						}

						// If the match is inclusive, just return the IDs as they
						// are
						return documentIds

					}

					lookupsLock.RUnlock()

				}

			}

		}

	}

	return result

}

// SearchDocumentIds searches for document IDs by evaluating a set of JSON criteria
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

// SearchDocuments searches for documents by evaluating a set of JSON criteria
func SearchDocuments(criteria map[string][]interface{}, from int, size int, alsoReturnAll bool) (int, []jsonserver.JSON, []jsonserver.JSON) {

	ids := []string{}

	// If no criteria, retrieve everything
	if len(criteria) == 0 {

		allIdsLock.RLock()

		for _, id := range allIds {
			ids = append(ids, id)
		}

		allIdsLock.RUnlock()

		// Otherwise filter by the actual criteria
	} else {
		ids = SearchDocumentIds(criteria)
	}

	// Sort IDs (later we will allow sorting by custom fields)
	sort.Strings(ids)

	// Convert document IDs to actual documents
	filtered := []jsonserver.JSON{}
	all := []jsonserver.JSON{}

	for sliceKey, id := range ids {

		// Use only the required IDs (pagination)
		if sliceKey >= from && sliceKey < from+size {

			document, err := GetDocument(id)

			if err == nil {

				filtered = append(filtered, map[string]interface{}{"id": id, "document": document})

				if alsoReturnAll {
					all = append(all, map[string]interface{}{"id": id, "document": document})
				}

			}

		} else if alsoReturnAll {

			document, err := GetDocument(id)

			if err == nil {
				all = append(all, map[string]interface{}{"id": id, "document": document})
			}

		}

	}

	if alsoReturnAll {
		return len(ids), filtered, all
	}

	return len(ids), filtered, nil

}
