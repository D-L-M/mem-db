package main


import (
    "net/http"
    "log"
    "io/ioutil"
    "./store"
    "./output"
    "./types"
    "net"
)


// Channel for document change messages
var documentMessage = make(chan types.DocumentMessage)


// Entry point
func main() {

    // Set up a HTTP server
    requestHandler := &RequestHandler{}
    
    requestHandler.Start()

    // Tell the disk indexer which channel to listen to for messages
    store.IndexOnDisk(documentMessage)

    select{}

}


// Define HTTP request handler type
type RequestHandler struct{}


// TCP server initialiser
func (requestHandler *RequestHandler) Start() {

	http.HandleFunc("/", requestHandler.dispatcher)

	server          := &http.Server{}
	listener, error := net.ListenTCP("tcp", &net.TCPAddr{IP: net.IPv4(0, 0, 0, 0), Port: 9999})

	if error != nil {
		log.Fatal("Error creating TCP listener")
	}

    go server.Serve(listener)
    
}


// Handle incoming requests and route to the appropriate package
func (requestHandler *RequestHandler) dispatcher(response http.ResponseWriter, request *http.Request) {
    
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

                documentMessage <- types.DocumentMessage{id, "add"}

                output.WriteJsonSuccessMessage("Document stored at " + id, http.StatusCreated)
                
            } else {

                output.WriteJsonErrorMessage("Document is not valid JSON", http.StatusBadRequest)

            }

        }

    }

    // Deleting documents
    if request.Method == "DELETE" {

        store.RemoveDocument(id)

        documentMessage <- types.DocumentMessage{id, "remove"}

        output.WriteJsonSuccessMessage("Document " + id + " removed", http.StatusOK)

    }

}