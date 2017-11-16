package output


import (
	"net/http"
	"encoding/json"
	"../types"
)


// Output writer which will be set here before is it used
var outputWriter http.ResponseWriter


// Set the http.ResponseWriter that will be used to output messages
func SetWriter(writer http.ResponseWriter) {

	outputWriter = writer

}


// Write a JSON response back to the user
func WriteJsonResponse(response map[string]interface {}, statusCode int) {
	
	outputWriter.Header().Set("Content-Type", "application/json")
	outputWriter.WriteHeader(statusCode)

	jsonString, _ := json.Marshal(response)

	outputWriter.Write(jsonString)

}


// Write a JSON 'success' message back to the user
func writeJsonOutcomeMessage(message string, success bool, statusCode int) {
	
	WriteJsonResponse(types.JsonDocument{"success": success, "message": message}, statusCode)

}


// Write a JSON 'success' message back to the user
func WriteJsonSuccessMessage(message string, statusCode int) {

	writeJsonOutcomeMessage(message, true, statusCode)

}


// Write a JSON 'error' message back to the user
func WriteJsonErrorMessage(message string, statusCode int) {
	
	writeJsonOutcomeMessage(message, false, statusCode)

}