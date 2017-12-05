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

// Entry point
func main() {

	// Register HTTP routes
	routing.RegisterRoutes()

	// Set up a HTTP server
	requestHandler := &RequestHandler{}

	requestHandler.Start()

	// Reindex all documents previously flushed to disk
	store.IndexFromDisk()

	// Listen for user messages
	go auth.ProcessMessages()

	// Listen for document messages
	go store.ProcessMessages()

	// Block execution so the asynchronous code can handle requests
	select {}

}

// RequestHandler defines the HTTP request handler
type RequestHandler struct{}

// Start is the TCP server initialiser
func (requestHandler *RequestHandler) Start() {

	http.HandleFunc("/", requestHandler.dispatcher)

	server := &http.Server{}
	listener, err := net.ListenTCP("tcp", &net.TCPAddr{IP: net.IPv4(0, 0, 0, 0), Port: 9999})

	if err != nil {
		log.Fatal("Error creating TCP listener")
	}

	go server.Serve(listener)

}

// Handle incoming requests and route to the appropriate package
func (requestHandler *RequestHandler) dispatcher(response http.ResponseWriter, request *http.Request) {

	body, err := ioutil.ReadAll(request.Body)

	if err != nil {

		output.WriteJSONErrorMessage(&response, "", "Could not read request body", http.StatusBadRequest)

	} else {

		if auth.CheckBasic(request) == false {

			output.WriteJSONResponse(&response, types.JSONDocument{"success": false, "message": "Not authorised"}, http.StatusUnauthorized)

		} else {

			method := request.Method
			path := request.URL.Path[:]
			id := request.URL.Path[1:]

			if routing.Dispatch(&response, method, path, id, &body) == false {
				output.WriteJSONResponse(&response, types.JSONDocument{"success": false, "message": "Unknown request"}, http.StatusBadRequest)
			}

		}

	}

}
