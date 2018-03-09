package routing

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/D-L-M/jsonserver"
	"github.com/D-L-M/mem-db/src/auth"
	"github.com/D-L-M/mem-db/src/crypt"
	"github.com/D-L-M/mem-db/src/data"
	"github.com/D-L-M/mem-db/src/messaging"
	"github.com/D-L-M/mem-db/src/output"
	"github.com/D-L-M/mem-db/src/store"
	"github.com/D-L-M/mem-db/src/types"
	"github.com/D-L-M/mem-db/src/utils"
)

// RegisterRoutes registers all HTTP routes
func RegisterRoutes() {

	// Check that the user is logged in
	authMiddleware := func(request *http.Request, body *[]byte, queryParams url.Values, routeParams jsonserver.RouteParams) (bool, int) {
		return auth.CheckCredentials(request, body), 401
	}

	// Check that the user is the root user
	adminMiddleware := func(request *http.Request, body *[]byte, queryParams url.Values, routeParams jsonserver.RouteParams) (bool, int) {

		username, _, _ := auth.GetCredentials(request)
		isRootUser := username == "root"

		return isRootUser, 401

	}

	// Welcome message
	jsonserver.RegisterRoute("GET", "/", []jsonserver.Middleware{authMiddleware}, func(request *http.Request, response *http.ResponseWriter, body *[]byte, queryParams url.Values, routeParams jsonserver.RouteParams) {

		output.WriteJSONResponse(response, data.GetWelcomeMessage(), http.StatusOK)

	})

	// Database stats
	jsonserver.RegisterRoute("GET", "/_stats", []jsonserver.Middleware{authMiddleware}, func(request *http.Request, response *http.ResponseWriter, body *[]byte, queryParams url.Values, routeParams jsonserver.RouteParams) {

		stats := store.GetStats()
		stats["peers"] = messaging.GetPeers()

		output.WriteJSONResponse(response, stats, http.StatusOK)

	})

	// Create or update/delete a user
	jsonserver.RegisterRoute("POST|PUT", "/_user", []jsonserver.Middleware{authMiddleware, adminMiddleware}, func(request *http.Request, response *http.ResponseWriter, body *[]byte, queryParams url.Values, routeParams jsonserver.RouteParams) {

		var credentials map[string]interface{}

		err := json.Unmarshal(*body, &credentials)

		hasUsername := err == nil && utils.MapHasKey(&credentials, "username")
		hasPassword := err == nil && utils.MapHasKey(&credentials, "password")
		isCreateAction := err == nil && utils.MapHasKey(&credentials, "action") && credentials["action"] == "create"
		isUpdateAction := err == nil && utils.MapHasKey(&credentials, "action") && credentials["action"] == "update"
		isDeleteAction := err == nil && utils.MapHasKey(&credentials, "action") && credentials["action"] == "delete"
		isCreateOrUpdateAction := isCreateAction || isUpdateAction

		if isCreateOrUpdateAction && hasUsername && hasPassword {

			go messaging.AddUser(credentials["username"].(string), credentials["password"].(string))
			output.WriteJSONResponse(response, types.JSONDocument{"success": true, "message": "User will be created or updated"}, http.StatusAccepted)

		} else if isDeleteAction && hasUsername && credentials["username"].(string) != "root" {

			go messaging.DeleteUser(credentials["username"].(string))
			output.WriteJSONResponse(response, types.JSONDocument{"success": true, "message": "User will be deleted"}, http.StatusAccepted)

		} else {

			output.WriteJSONResponse(response, types.JSONDocument{"success": false, "message": "Malformed request"}, http.StatusBadRequest)

		}

	})

	// Receive an instructional message from a peer server
	jsonserver.RegisterRoute("POST", "/_peer-message", []jsonserver.Middleware{authMiddleware}, func(request *http.Request, response *http.ResponseWriter, body *[]byte, queryParams url.Values, routeParams jsonserver.RouteParams) {

		var message types.PeerMessage

		err := json.Unmarshal(*body, &message)

		if err != nil {

			output.WriteJSONResponse(response, types.JSONDocument{"success": false, "message": "Malformed request"}, http.StatusBadRequest)

		} else {

			messaging.PeerMessageQueue <- message

			output.WriteJSONResponse(response, types.JSONDocument{"success": true, "message": "Instructions will be acted upon"}, http.StatusAccepted)

		}

	})

	// Store a document
	putDocumentAction := func(request *http.Request, response *http.ResponseWriter, body *[]byte, queryParams url.Values, routeParams jsonserver.RouteParams) {

		id, ok := routeParams["id"]

		// If an ID was not provided, create one
		if ok == false {

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

				go messaging.AddDocument(id, body, true)

				output.WriteJSONSuccessMessage(response, id, "Document will be stored", http.StatusAccepted)

			}

		}

	}

	jsonserver.RegisterRoute("PUT", "/", []jsonserver.Middleware{authMiddleware}, putDocumentAction)
	jsonserver.RegisterRoute("PUT", "/{id}", []jsonserver.Middleware{authMiddleware}, putDocumentAction)

	// Truncate the database
	jsonserver.RegisterRoute("DELETE", "/_all", []jsonserver.Middleware{authMiddleware}, func(request *http.Request, response *http.ResponseWriter, body *[]byte, queryParams url.Values, routeParams jsonserver.RouteParams) {

		go messaging.RemoveAllDocuments(true)

		output.WriteJSONResponse(response, types.JSONDocument{"success": true, "message": "All documents will be removed"}, http.StatusAccepted)

	})

	// Remove a document
	jsonserver.RegisterRoute("DELETE", "/{id}", []jsonserver.Middleware{authMiddleware}, func(request *http.Request, response *http.ResponseWriter, body *[]byte, queryParams url.Values, routeParams jsonserver.RouteParams) {

		id := routeParams["id"]
		_, err := store.GetRawDocument(id)

		if err != nil {

			output.WriteJSONErrorMessage(response, id, "Document does not exist", http.StatusNotFound)

		} else {

			go messaging.RemoveDocument(id, true)

			output.WriteJSONSuccessMessage(response, id, "Document will be removed", http.StatusAccepted)

		}

	})

	// Search for documents by criteria
	jsonserver.RegisterRoute("GET|POST", "/_search", []jsonserver.Middleware{authMiddleware}, func(request *http.Request, response *http.ResponseWriter, body *[]byte, queryParams url.Values, routeParams jsonserver.RouteParams) {

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

			from, _ := strconv.Atoi(GetFirstParamValue(queryParams, "from", "0"))
			size, _ := strconv.Atoi(GetFirstParamValue(queryParams, "size", "25"))
			significantTermsField := GetFirstParamValue(queryParams, "significant_terms_field", "")
			significantTermsThreshold, _ := strconv.Atoi(GetFirstParamValue(queryParams, "significant_terms_threshold", "200"))
			significantTermsMinimumOccurrencePercentage, _ := strconv.ParseFloat(GetFirstParamValue(queryParams, "significant_terms_minimum", "33.34"), 64)
			criteria := map[string][]interface{}(criteria)
			startTime := time.Now()

			includeAllMatches := false

			if significantTermsField != "" {
				includeAllMatches = true
			}

			totalDocumentCount, documents, allDocuments := store.SearchDocuments(criteria, from, size, includeAllMatches)
			significantTerms := []map[string]interface{}{}

			// Optionally get significant terms
			if significantTermsField != "" {
				significantTerms = store.DiscoverSignificantTerms(&allDocuments, significantTermsField, significantTermsThreshold, significantTermsMinimumOccurrencePercentage)
			}

			timeTaken := (time.Since(startTime).Nanoseconds() / int64(time.Millisecond))
			info := map[string]interface{}{"total_matches": totalDocumentCount, "time_taken": timeTaken}
			searchResults := map[string]interface{}{"criteria": criteria, "information": info, "results": documents}

			// Optionally include significant terms
			if significantTermsField != "" {
				searchResults["significant_terms"] = significantTerms
			}

			output.WriteJSONResponse(response, searchResults, http.StatusOK)

		}

	})

	// Delete documents by criteria
	jsonserver.RegisterRoute("GET|POST", "/_delete", []jsonserver.Middleware{authMiddleware}, func(request *http.Request, response *http.ResponseWriter, body *[]byte, queryParams url.Values, routeParams jsonserver.RouteParams) {

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
				go messaging.RemoveDocument(documentID, true)
			}

			output.WriteJSONResponse(response, types.JSONDocument{"success": true, "message": strconv.Itoa(len(documentIds)) + " document(s) will be removed"}, http.StatusAccepted)

		}

	})

	// Get a document
	jsonserver.RegisterRoute("GET", "/{id}", []jsonserver.Middleware{authMiddleware}, func(request *http.Request, response *http.ResponseWriter, body *[]byte, queryParams url.Values, routeParams jsonserver.RouteParams) {

		id := routeParams["id"]
		document, err := store.GetDocument(id)

		if err != nil {

			output.WriteJSONErrorMessage(response, id, "Document does not exist", http.StatusNotFound)

		} else {

			output.WriteJSONResponse(response, document, http.StatusOK)

		}

	})

}

// GetFirstParamValue extracts the first value for a key from a route's URL parameters
func GetFirstParamValue(params url.Values, key string, fallback string) string {

	if _, ok := params[key]; ok {

		values := params[key]

		if len(values) >= 1 {
			return values[0]
		}

	}

	return fallback

}
