package jsonserver

import (
	"io/ioutil"
	"log"
	"net"
	"net/http"
)

// TCP server
type server struct{}

// Start is the TCP server initialiser
func (requestHandler *server) Start(port int) {

	http.HandleFunc("/", requestHandler.dispatcher)

	server := &http.Server{}
	listener, err := net.ListenTCP("tcp", &net.TCPAddr{IP: net.IPv4(0, 0, 0, 0), Port: port})

	if err != nil {
		log.Fatal("Error creating TCP listener")
	}

	go server.Serve(listener)

}

// Handle incoming requests and route to the appropriate package
func (requestHandler *server) dispatcher(response http.ResponseWriter, request *http.Request) {

	body, err := ioutil.ReadAll(request.Body)

	if err != nil {

		WriteResponse(&response, JSON{"success": false, "message": "Could not read request body"}, http.StatusBadRequest)

	} else {

		method := request.Method
		path := request.URL.Path[:]
		params := request.URL.RawQuery
		success, err := dispatch(request, &response, method, path, params, &body)

		// Access denied by middleware
		if err != nil {

			WriteResponse(&response, JSON{"success": false, "message": "Access denied"}, http.StatusForbidden)

			// No matching routes found
		} else if success == false {

			WriteResponse(&response, JSON{"success": false, "message": "Could not find " + path}, http.StatusNotFound)

		}

	}

}

// Start initialises the TCP server
func Start(port int) {

	requestHandler := &server{}

	requestHandler.Start(port)

}
