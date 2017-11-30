package output

import (
	"encoding/json"
	"net/http"

	"../types"
)

// WriteJSONResponse writes a JSON response back to the client
func WriteJSONResponse(response *http.ResponseWriter, body map[string]interface{}, statusCode int) {

	responseWriter := *response

	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(statusCode)

	jsonString, _ := json.Marshal(body)

	responseWriter.Write(jsonString)

}

// writeJSONOutcomeMessage writes a JSON 'success' message back to the client
func writeJSONOutcomeMessage(response *http.ResponseWriter, id string, message string, success bool, statusCode int) {

	WriteJSONResponse(response, types.JSONDocument{"success": success, "id": id, "message": message}, statusCode)

}

// WriteJSONSuccessMessage writes a JSON 'success' message back to the client
func WriteJSONSuccessMessage(response *http.ResponseWriter, id string, message string, statusCode int) {

	writeJSONOutcomeMessage(response, id, message, true, statusCode)

}

// WriteJSONErrorMessage writes a JSON 'error' message back to the client
func WriteJSONErrorMessage(response *http.ResponseWriter, id string, message string, statusCode int) {

	writeJSONOutcomeMessage(response, id, message, false, statusCode)

}
