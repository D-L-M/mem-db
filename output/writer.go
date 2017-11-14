package output


import (
	"net/http"
	"fmt"
)


var outputWriter http.ResponseWriter


// Set the http.ResponseWriter that will be used to output messages
func SetWriter(writer http.ResponseWriter) {

	outputWriter = writer

}


// Write a message back to the user
func Write(message string) {

	fmt.Fprintf(outputWriter, message)

}