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


type requestHandler struct{}


// Handle HTTP requests and route to the appropriate package
func (rh requestHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {

    output.SetWriter(response)

    // PUTTING a document
    if request.Method == "PUT" {

        body, err := ioutil.ReadAll(request.Body)

        // Error reading the request body
        if err != nil {

            output.Write("Could not read request body")

        // Request body received
        } else {

            if store.IndexDocument(request.URL.Path, body) {
                output.Write("PUT document to " + request.URL.Path)
            }

        }

    }

    

}