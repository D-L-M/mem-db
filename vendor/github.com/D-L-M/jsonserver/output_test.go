package jsonserver

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestWriteResponse tests writing a JSON response back to the client
func TestWriteResponse(t *testing.T) {

	responseWriter := httptest.NewRecorder()

	body := JSON{"foo": "bar"}
	statusCode := http.StatusAccepted

	WriteResponse(responseWriter, &body, statusCode)

	if responseWriter.Body.String() != `{"foo":"bar"}` {
		t.Errorf("Incorrect body (expected: %v, actual: %v)", `{"foo":"bar"}`, responseWriter.Body.String())
	}

	if responseWriter.Header().Get("Content-Type") != "application/json" {
		t.Errorf("Incorrect content-type header (expected: %v, actual: %v)", "application/json", responseWriter.Header().Get("Content-Type"))
	}

	if responseWriter.Code != http.StatusAccepted {
		t.Errorf("Incorrect status code (expected: %v, actual: %v)", http.StatusAccepted, responseWriter.Code)
	}

}
