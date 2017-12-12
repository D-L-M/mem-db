package routing

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"../crypt"
	"../data"
	"../messaging"
	"../output"
	"../store"
	"../types"
	"../utils"
)

// RegisterRoutes registers all HTTP routes
func RegisterRoutes() {

	// Welcome message
	Register("GET", "/", false, func(request *http.Request, response *http.ResponseWriter, body *[]byte, id string) {

		output.WriteJSONResponse(response, data.GetWelcomeMessage(), http.StatusOK)

	})

	// Database stats
	Register("GET", "/_stats", false, func(request *http.Request, response *http.ResponseWriter, body *[]byte, id string) {

		output.WriteJSONResponse(response, store.GetStats(), http.StatusOK)

	})

	// Create or update a user
	Register("POST|PUT", "/_user", true, func(request *http.Request, response *http.ResponseWriter, body *[]byte, id string) {

		var credentials map[string]interface{}

		err := json.Unmarshal(*body, &credentials)

		if err != nil || utils.MapHasKey(&credentials, "username") == false || utils.MapHasKey(&credentials, "password") == false {

			output.WriteJSONResponse(response, types.JSONDocument{"success": false, "message": "Malformed request"}, http.StatusBadRequest)

		} else {

			go messaging.AddUser(credentials["username"].(string), credentials["password"].(string))

			output.WriteJSONResponse(response, types.JSONDocument{"success": true, "message": "User will be created or updated"}, http.StatusAccepted)

		}

	})

	// Store a document
	Register("PUT", "/*", false, func(request *http.Request, response *http.ResponseWriter, body *[]byte, id string) {

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

			_, err := store.ParseDocument(*body)

			if err != nil {

				output.WriteJSONErrorMessage(response, id, "Document is not valid JSON", http.StatusBadRequest)

			} else {

				go messaging.AddDocument(id, body)

				output.WriteJSONSuccessMessage(response, id, "Document will be stored", http.StatusAccepted)

			}

		}

	})

	// Truncate the database
	Register("DELETE", "/_all", false, func(request *http.Request, response *http.ResponseWriter, body *[]byte, id string) {

		go messaging.RemoveAllDocuments()

		output.WriteJSONResponse(response, types.JSONDocument{"success": true, "message": "All documents will be removed"}, http.StatusAccepted)

	})

	// Remove a document
	Register("DELETE", "/*", false, func(request *http.Request, response *http.ResponseWriter, body *[]byte, id string) {

		_, err := store.GetRawDocument(id)

		if err != nil {

			output.WriteJSONErrorMessage(response, id, "Document does not exist", http.StatusNotFound)

		} else {

			go messaging.RemoveDocument(id)

			output.WriteJSONSuccessMessage(response, id, "Document will be removed", http.StatusAccepted)

		}

	})

	// Search for documents by criteria
	Register("GET|POST", "/_search", false, func(request *http.Request, response *http.ResponseWriter, body *[]byte, id string) {

		// If no body sent, assume an empty criteria
		if string((*body)[:]) == "" {
			emptyBody := []byte("{}")
			body = &emptyBody
		}

		// Get the actual JSON criteria
		var criteria map[string][]interface{}

		err := json.Unmarshal(*body, &criteria)

		if err != nil {

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
	Register("GET|POST", "/_delete", false, func(request *http.Request, response *http.ResponseWriter, body *[]byte, id string) {

		// If no body sent, assume an empty criteria
		if string((*body)[:]) == "" {
			emptyBody := []byte("{}")
			body = &emptyBody
		}

		// Get the actual JSON criteria
		var criteria map[string][]interface{}

		err := json.Unmarshal(*body, &criteria)

		if err != nil {

			output.WriteJSONErrorMessage(response, "", "Search criteria is not valid JSON", http.StatusBadRequest)

			// Retrieve documents matching the search criteria
		} else {

			criteria := map[string][]interface{}(criteria)

			documentIds := store.SearchDocumentIds(criteria)

			for _, documentID := range documentIds {
				go messaging.RemoveDocument(documentID)
			}

			output.WriteJSONResponse(response, types.JSONDocument{"success": true, "message": strconv.Itoa(len(documentIds)) + " document(s) will be removed"}, http.StatusAccepted)

		}

	})

	// Get a document
	Register("GET", "/*", false, func(request *http.Request, response *http.ResponseWriter, body *[]byte, id string) {

		document, err := store.GetDocument(id)

		if err != nil {

			output.WriteJSONErrorMessage(response, id, "Document does not exist", http.StatusNotFound)

		} else {

			output.WriteJSONResponse(response, document, http.StatusOK)

		}

	})

}
