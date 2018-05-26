package jsonserver

import (
	"io/ioutil"
	"net/http"
	"testing"
)

// TestServerCanListen tests starting and making a request to the server
func TestServerCanListen(t *testing.T) {

	testRouteSetUp()

	server, err := Start(9999)

	if err != nil {
		t.Errorf("Unable to make start server")
	}

	defer server.Close()

	response, err := http.Get("http://127.0.0.1:9999/")

	if err != nil {
		t.Errorf("Unable to make request")
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		t.Errorf("Unexpected error thrown when attempting to read response")
	}

	bodyString := string(body)

	if bodyString != "GET /" {
		t.Errorf("Could not reach route")
	}

	testRouteTearDown()

}

// TestServerReturnsNotFound tests receiving a 404 response from the server for a bad route
func TestServerReturnsNotFound(t *testing.T) {

	testRouteSetUp()

	server, err := Start(9999)

	if err != nil {
		t.Errorf("Unable to make start server")
	}

	defer server.Close()

	response, err := http.Get("http://127.0.0.1:9999/404")

	if err != nil {
		t.Errorf("Unable to make request")
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	code := response.StatusCode

	if err != nil {
		t.Errorf("Unexpected error thrown when attempting to read response")
	}

	bodyString := string(body)

	if bodyString != `{"message":"Could not find /404","success":false}` {
		t.Errorf("Route did not return 'not found' message")
	}

	if code != 404 {
		t.Errorf("Route did not return 404 HTTP code")
	}

	testRouteTearDown()

}

// TestServerReturnsOutputFromDenyingMiddleware tests receiving a middleware denial response
func TestServerReturnsOutputFromDenyingMiddleware(t *testing.T) {

	testRouteSetUp()

	server, err := Start(9999)

	if err != nil {
		t.Errorf("Unable to make start server")
	}

	defer server.Close()

	response, err := http.Get("http://127.0.0.1:9999/middleware_deny")

	if err != nil {
		t.Errorf("Unable to make request")
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	code := response.StatusCode

	if err != nil {
		t.Errorf("Unexpected error thrown when attempting to read response")
	}

	bodyString := string(body)

	if bodyString != `{"message":"Access denied","success":false}` {
		t.Errorf("Route did not return middleware denial message")
	}

	if code != 401 {
		t.Errorf("Route did not return middleware denial HTTP code")
	}

	testRouteTearDown()

}

// TestServerFailsOnOutOfBoundsPort tests that creating a server fails with an out-of-bounds port number
func TestServerFailsOnOutOfBoundsPort(t *testing.T) {

	server, err := Start(99999)

	if err == nil {
		defer server.Close()
		t.Errorf("Server with out-of-bounds port number unexpectedly started")
	}

	testRouteTearDown()

}
