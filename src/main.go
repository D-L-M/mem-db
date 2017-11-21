package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net"
	"net/http"

	"./crypt"
	"./data"
	"./output"
	"./store"
	"./types"
)

// Channel for document change messages
var documentMessage = make(chan types.DocumentMessage)

// Entry point
func main() {

	// Set up a HTTP server
	requestHandler := &RequestHandler{}

	requestHandler.Start()

	// Reindex all documents previously flushed to disk
	store.IndexFromDisk()

	// Tell the disk indexer which channel to listen to for messages
	store.FlushToDisk(documentMessage)

}

// Define HTTP request handler type
type RequestHandler struct{}

// TCP server initialiser
func (requestHandler *RequestHandler) Start() {

	http.HandleFunc("/", requestHandler.dispatcher)

	server := &http.Server{}
	listener, error := net.ListenTCP("tcp", &net.TCPAddr{IP: net.IPv4(0, 0, 0, 0), Port: 9999})

	if error != nil {
		log.Fatal("Error creating TCP listener")
	}

	go server.Serve(listener)

}

// Handle incoming requests and route to the appropriate package
func (requestHandler *RequestHandler) dispatcher(response http.ResponseWriter, request *http.Request) {

	// The document ID is the path
	id := request.URL.Path[1:]

	// Getting documents/data
	if request.Method == "GET" {

		// Welcome message
		if id == "" {

			output.WriteJsonResponse(response, data.GetWelcomeMessage(), http.StatusOK)

			// Index stats
		} else if id == "_stats" {

			output.WriteJsonResponse(response, store.GetStats(), http.StatusOK)

			// Single document
		} else {

			document, error := store.GetDocument(id)

			// Error getting the document
			if error != nil {

				output.WriteJsonErrorMessage(response, id, "Document does not exist", http.StatusNotFound)

				// Document retrieved
			} else {

				output.WriteJsonResponse(response, document, http.StatusOK)

			}

		}

	}

	if request.Method == "GET" || request.Method == "POST" {

		if id == "_search" {

			body, error := ioutil.ReadAll(request.Body)

			// Error reading the request body
			if error != nil {

				output.WriteJsonErrorMessage(response, "", "Could not read request body", http.StatusBadRequest)

				// Request body received
			} else {

				// If no body sent, assume an empty criteria
				if string(body[:]) == "" {
					body = []byte("{}")
				}

				// Get the actual JSON criteria
				var criteria map[string][]interface{}

				error := json.Unmarshal(body, &criteria)

				if error != nil {

					output.WriteJsonErrorMessage(response, "", "Search criteria is not valid JSON", http.StatusBadRequest)

					// Retrieve documents matching the search criteria
				} else {

					criteria := map[string][]interface{}(criteria)
					documents := store.SearchDocuments(criteria)
					searchResults := map[string]interface{}{}

					searchResults["total_count"] = len(documents)
					searchResults["documents"] = documents

					output.WriteJsonResponse(response, searchResults, http.StatusOK)

				}

			}

		}

	}

	// Storing documents
	if request.Method == "PUT" {

		// If an ID was not provided, create one
		if id == "" {

			id, _ = crypt.GenerateUuid()

			if id == "" {
				output.WriteJsonErrorMessage(response, "", "An error occurred whilst generating a document ID", http.StatusInternalServerError)
			}

		}

		// If an ID was provided or has been generated, attempt to store the
		// document under it
		if id != "" {

			body, error := ioutil.ReadAll(request.Body)

			// Error reading the request body
			if error != nil {

				output.WriteJsonErrorMessage(response, id, "Could not read request body", http.StatusBadRequest)

				// Request body received
			} else {

				_, error = store.ParseDocument(body)

				// Malformed document
				if error != nil {

					output.WriteJsonErrorMessage(response, id, "Document is not valid JSON", http.StatusBadRequest)

					// Everything is okay, so store the document
				} else {

					documentMessage <- types.DocumentMessage{Id: id, Document: body, Action: "add"}

					output.WriteJsonSuccessMessage(response, id, "Document will be stored", http.StatusAccepted)

				}

			}

		}

	}

	// Deleting documents
	if request.Method == "DELETE" {

		document, error := store.GetRawDocument(id)

		if error != nil {

			output.WriteJsonErrorMessage(response, id, "Document does not exist", http.StatusNotFound)

		} else {

			documentMessage <- types.DocumentMessage{Id: id, Document: document, Action: "remove"}

			output.WriteJsonSuccessMessage(response, id, "Document will be removed", http.StatusAccepted)

		}

	}

}
