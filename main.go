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

    // Getting documents/data
    if request.Method == "GET" {

        // Index stats
        if (id == "/_stats") {

            output.WriteJsonResponse(store.GetStats(), http.StatusOK)

        // Single document
        } else {
            
            document, error := store.GetDocument(id)

            // Error getting the document
            if error != nil {

                output.WriteJsonErrorMessage("Document does not exist", http.StatusNotFound)

            // Document retrieved
            } else {

                output.WriteJsonResponse(document, http.StatusOK)

            }

        }

    }

    // Storing documents
    if request.Method == "PUT" {

        body, error := ioutil.ReadAll(request.Body)

        // Error reading the request body
        if error != nil {

            output.WriteJsonErrorMessage("Could not read request body", http.StatusBadRequest)

        // Request body received
        } else {

            if store.IndexDocument(id, body) {

                output.WriteJsonSuccessMessage("Document stored at " + id, http.StatusCreated)
                
            } else {

                output.WriteJsonErrorMessage("Document is not valid JSON", http.StatusBadRequest)

            }

        }

    }

    // Deleting documents
    if request.Method == "DELETE" {

        store.RemoveDocument(id)
        output.WriteJsonSuccessMessage("Document " + id + " removed", http.StatusOK)

    }

}