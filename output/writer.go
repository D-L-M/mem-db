package output


import (
	"net/http"
	"fmt"
	"encoding/json"
)


type JsonMessage map[string]interface {}


var outputWriter http.ResponseWriter


// Set the http.ResponseWriter that will be used to output messages
func SetWriter(writer http.ResponseWriter) {

	outputWriter = writer

}


// Write a message back to the user
func Write(message string) {

	fmt.Fprintf(outputWriter, message)

}


// Write a JSON 'success' message back to the user
func WriteJsonSuccessMessage(message string, success bool) {

	outputWriter.Header().Set("Content-Type", "application/json")

	jsonString, _ := json.Marshal(JsonMessage{"success": success, "message": message})

	fmt.Fprintf(outputWriter, string(jsonString[:]))

}