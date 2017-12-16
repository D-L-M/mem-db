package server

import (
	"io/ioutil"
	"log"
	"net"
	"net/http"

	"../auth"
	"../output"
	"../routing"
	"../types"
)

// tcpRequestHandler defines the HTTP request handler
type tcpRequestHandler struct{}

// Start is the TCP server initialiser
func (requestHandler *tcpRequestHandler) Start(port int) {

	http.HandleFunc("/", requestHandler.dispatcher)

	server := &http.Server{}
	listener, err := net.ListenTCP("tcp", &net.TCPAddr{IP: net.IPv4(0, 0, 0, 0), Port: port})

	if err != nil {
		log.Fatal("Error creating TCP listener")
	}

	go server.Serve(listener)

}

// Handle incoming requests and route to the appropriate package
func (requestHandler *tcpRequestHandler) dispatcher(response http.ResponseWriter, request *http.Request) {

	body, err := ioutil.ReadAll(request.Body)

	if err != nil {

		output.WriteJSONErrorMessage(&response, "", "Could not read request body", http.StatusBadRequest)

	} else {

		if auth.CheckCredentials(request, &body) == false {

			output.WriteJSONResponse(&response, types.JSONDocument{"success": false, "message": "Not authorised"}, http.StatusUnauthorized)

		} else {

			method := request.Method
			path := request.URL.Path[:]
			id := request.URL.Path[1:]
			success, err := routing.Dispatch(request, &response, method, path, id, &body)

			// Root user only route, but user is not root
			if err != nil {

				output.WriteJSONResponse(&response, types.JSONDocument{"success": false, "message": "Not authorised"}, http.StatusUnauthorized)

				// No matching routes found
			} else if success == false {

				output.WriteJSONResponse(&response, types.JSONDocument{"success": false, "message": "Unknown request"}, http.StatusBadRequest)

			}

		}

	}

}

// InitTCP initialises the TCP server
func InitTCP(port int) {

	requestHandler := &tcpRequestHandler{}

	requestHandler.Start(port)

}
