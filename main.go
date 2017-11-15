package main


import (
    "net/http"
    "log"
    "io/ioutil"
    "./store"
    "./output"
)


// Set up a HTTP server and log any errors
func main() {

    err := http.ListenAndServe(":9999", requestHandler{})
    
    log.Fatal(err)
    
}


// Define the request handler
type requestHandler struct{}


// Handle HTTP requests and route to the appropriate package
func (rh requestHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {

    output.SetWriter(response)

    // The document ID is the path
    id := request.URL.Path;

    // GETTING a document
    if request.Method == "GET" {
        
        document, error := store.GetDocument(id)

        // Error getting the document
        if error != nil {

            output.WriteJsonSuccessMessage("Document does not exist", false)

        // Document retrieved
        } else {

            output.WriteJsonResponse(document)

        }

    }

    // PUTTING a document
    if request.Method == "PUT" {

        body, error := ioutil.ReadAll(request.Body)

        // Error reading the request body
        if error != nil {

            output.WriteJsonSuccessMessage("Could not read request body", true)

        // Request body received
        } else {

            if store.IndexDocument(id, body) {

                output.WriteJsonSuccessMessage("PUT document to " + id, true)
                
            } else {

                output.WriteJsonSuccessMessage("Document is not valid JSON", false)

            }

        }

    }

}