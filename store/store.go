package store


import (
	"encoding/json"
	"../output"
)


// Documents could be any format once they're parsed, so make a generic object
// type to define them
type Document interface {}


// Documents are stored in a map, for quick retrieval
var documents = map[string]Document {}


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

		return true

	}

	return false

}