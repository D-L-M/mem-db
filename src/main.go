package main

import (
	"io/ioutil"
	"log"
	"net"
	"net/http"

	"./auth"
	"./output"
	"./routing"
	"./store"
	"./types"
)

// documentMessageQueue is a channel for document change messages
var documentMessageQueue = make(chan types.DocumentMessage)

// Entry point
func main() {

	// Register HTTP routes
	routing.RegisterRoutes(documentMessageQueue)

	// Set up a HTTP server
	requestHandler := &RequestHandler{}

	requestHandler.Start()

	// Reindex all documents previously flushed to disk
	store.IndexFromDisk()

	// Tell the disk indexer which channel to listen to for messages
	store.FlushToDisk(documentMessageQueue)

}

// RequestHandler defines the HTTP request handler
type RequestHandler struct{}

// Start is the TCP server initialiser
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

	body, error := ioutil.ReadAll(request.Body)

	if error != nil {

		output.WriteJSONErrorMessage(response, "", "Could not read request body", http.StatusBadRequest)

	} else {

		if auth.CheckBasic(request) == false {

			output.WriteJSONResponse(response, types.JSONDocument{"success": false, "message": "Not authorised"}, http.StatusUnauthorized)

		} else {

			method := request.Method
			path := request.URL.Path[:]
			id := request.URL.Path[1:]

			if routing.Dispatch(response, method, path, id, &body) == false {
				output.WriteJSONResponse(response, types.JSONDocument{"success": false, "message": "Unknown request"}, http.StatusBadRequest)
			}

		}

	}

}
