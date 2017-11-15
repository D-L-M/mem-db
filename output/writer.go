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
// TODO: Make a 'success' and 'error' version of this function
func WriteJsonSuccessMessage(message string, success bool) {

	WriteJsonResponse(types.JsonDocument{"success": success, "message": message})

}