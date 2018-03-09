package jsonserver

import (
	"encoding/json"
	"net/http"
)

// WriteResponse writes a JSON response back to the client
func WriteResponse(response http.ResponseWriter, body *JSON, statusCode int) {

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(statusCode)

	jsonString, _ := json.Marshal(*body)

	response.Write(jsonString)

}
