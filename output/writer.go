package output


import (
	"net/http"
	"fmt"
	"encoding/json"
	"../types"
)


// Output writer which will be set here before is it used
var outputWriter http.ResponseWriter


// Set the http.ResponseWriter that will be used to output messages
func SetWriter(writer http.ResponseWriter) {

	outputWriter = writer

}


// Write a message back to the user
func Write(message string) {

	fmt.Fprintf(outputWriter, message)

}


// Write a JSON response back to the user
func WriteJsonResponse(response map[string]interface {}) {
	
	outputWriter.Header().Set("Content-Type", "application/json")

	jsonString, _ := json.Marshal(response)

	fmt.Fprintf(outputWriter, string(jsonString[:]))

}


// Write a JSON 'success' message back to the user
func writeJsonOutcomeMessage(message string, success bool) {
	
	WriteJsonResponse(types.JsonDocument{"success": success, "message": message})

}


// Write a JSON 'success' message back to the user
func WriteJsonSuccessMessage(message string) {

	writeJsonOutcomeMessage(message, true)

}


// Write a JSON 'error' message back to the user
func WriteJsonErrorMessage(message string) {
	
	writeJsonOutcomeMessage(message, false)

}