package routing

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"../crypt"
	"../data"
	"../output"
	"../store"
	"../types"
)

// RegisterRoutes registers all HTTP routes
func RegisterRoutes(documentMessageQueue chan types.DocumentMessage) {

	// Welcome message
	Register("GET", "/", func(response http.ResponseWriter, body *[]byte, id string) {

		output.WriteJSONResponse(response, data.GetWelcomeMessage(), http.StatusOK)

	})

	// Database stats
	Register("GET", "/_stats", func(response http.ResponseWriter, body *[]byte, id string) {

		output.WriteJSONResponse(response, store.GetStats(), http.StatusOK)

	})

	// Store a document
	Register("PUT", "/*", func(response http.ResponseWriter, body *[]byte, id string) {

		// If an ID was not provided, create one
		if id == "" {

			id, _ = crypt.GenerateUUID()

			if id == "" {

				output.WriteJSONErrorMessage(response, "", "An error occurred whilst generating a document ID", http.StatusInternalServerError)

			}

		}

		// If an ID was provided or has been generated, attempt to store the
		// document under it
		if id != "" {

			_, error := store.ParseDocument(*body)

			if error != nil {

				output.WriteJSONErrorMessage(response, id, "Document is not valid JSON", http.StatusBadRequest)

			} else {

				documentMessageQueue <- types.DocumentMessage{ID: id, Document: *body, Action: "add"}

				output.WriteJSONSuccessMessage(response, id, "Document will be stored", http.StatusAccepted)

			}

		}

	})

	// Truncate the database
	Register("DELETE", "/_all", func(response http.ResponseWriter, body *[]byte, id string) {

		documentMessageQueue <- types.DocumentMessage{ID: "_all", Document: []byte{}, Action: "remove"}

		output.WriteJSONResponse(response, types.JSONDocument{"success": true, "message": "All documents will be removed"}, http.StatusAccepted)

	})

	// Remove a document
	Register("DELETE", "/*", func(response http.ResponseWriter, body *[]byte, id string) {

		_, error := store.GetRawDocument(id)

		if error != nil {

			output.WriteJSONErrorMessage(response, id, "Document does not exist", http.StatusNotFound)

		} else {

			documentMessageQueue <- types.DocumentMessage{ID: id, Document: []byte{}, Action: "remove"}

			output.WriteJSONSuccessMessage(response, id, "Document will be removed", http.StatusAccepted)

		}

	})

	// Search for documents by criteria
	Register("GET|POST", "/_search", func(response http.ResponseWriter, body *[]byte, id string) {

		// If no body sent, assume an empty criteria
		if string((*body)[:]) == "" {
			emptyBody := []byte("{}")
			body = &emptyBody
		}

		// Get the actual JSON criteria
		var criteria map[string][]interface{}

		error := json.Unmarshal(*body, &criteria)

		if error != nil {

			output.WriteJSONErrorMessage(response, "", "Search criteria is not valid JSON", http.StatusBadRequest)

			// Retrieve documents matching the search criteria
		} else {

			criteria := map[string][]interface{}(criteria)
			startTime := time.Now()
			documents := store.SearchDocuments(criteria)
			timeTaken := (time.Since(startTime).Nanoseconds() / int64(time.Millisecond))
			info := map[string]interface{}{"total_matches": len(documents), "time_taken": timeTaken}
			searchResults := map[string]interface{}{"criteria": criteria, "information": info, "results": documents}

			output.WriteJSONResponse(response, searchResults, http.StatusOK)

		}

	})

	// Delete documents by criteria
	Register("GET|POST", "/_delete", func(response http.ResponseWriter, body *[]byte, id string) {

		// If no body sent, assume an empty criteria
		if string((*body)[:]) == "" {
			emptyBody := []byte("{}")
			body = &emptyBody
		}

		// Get the actual JSON criteria
		var criteria map[string][]interface{}

		error := json.Unmarshal(*body, &criteria)

		if error != nil {

			output.WriteJSONErrorMessage(response, "", "Search criteria is not valid JSON", http.StatusBadRequest)

			// Retrieve documents matching the search criteria
		} else {

			criteria := map[string][]interface{}(criteria)

			documentIds := store.SearchDocumentIds(criteria)

			for _, documentID := range documentIds {
				documentMessageQueue <- types.DocumentMessage{ID: documentID, Document: []byte{}, Action: "remove"}
			}

			output.WriteJSONResponse(response, types.JSONDocument{"success": true, "message": strconv.Itoa(len(documentIds)) + " document(s) will be removed"}, http.StatusAccepted)

		}

	})

	// Get a document
	Register("GET", "/*", func(response http.ResponseWriter, body *[]byte, id string) {

		document, error := store.GetDocument(id)

		if error != nil {

			output.WriteJSONErrorMessage(response, id, "Document does not exist", http.StatusNotFound)

		} else {

			output.WriteJSONResponse(response, document, http.StatusOK)

		}

	})

}
