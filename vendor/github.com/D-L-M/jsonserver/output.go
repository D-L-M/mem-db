package jsonserver

import (
	"encoding/json"
	"net/http"
)

// WriteResponse writes a JSON response back to the client
func WriteResponse(response *http.ResponseWriter, body JSON, statusCode int) {

	responseWriter := *response

	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(statusCode)

	jsonString, _ := json.Marshal(body)

	responseWriter.Write(jsonString)

}
