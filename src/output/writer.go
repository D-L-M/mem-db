package output

import (
	"../types"
	"encoding/json"
	"net/http"
)

// Write a JSON response back to the client
func WriteJsonResponse(response http.ResponseWriter, body map[string]interface{}, statusCode int) {

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(statusCode)

	jsonString, _ := json.Marshal(body)

	response.Write(jsonString)

}

// Write a JSON 'success' message back to the client
func writeJsonOutcomeMessage(response http.ResponseWriter, id string, message string, success bool, statusCode int) {

	WriteJsonResponse(response, types.JsonDocument{"success": success, "id": id, "message": message}, statusCode)

}

// Write a JSON 'success' message back to the client
func WriteJsonSuccessMessage(response http.ResponseWriter, id string, message string, statusCode int) {

	writeJsonOutcomeMessage(response, id, message, true, statusCode)

}

// Write a JSON 'error' message back to the client
func WriteJsonErrorMessage(response http.ResponseWriter, id string, message string, statusCode int) {

	writeJsonOutcomeMessage(response, id, message, false, statusCode)

}
