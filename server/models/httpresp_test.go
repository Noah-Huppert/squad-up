package models

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"testing"
)

// Mock http.ResponseWriter for testing
type MockResponseWriter struct {
	mock.Mock

	Body       string
	RespHeader http.Header
}

// Mock write function to make struct a Writer
func (m *MockResponseWriter) Write(p []byte) (n int, err error) {
	m.Body += string(p[:])
	return len(p), nil
}

// Getter for header field, used in HTTPResponse.Serve
func (m *MockResponseWriter) Header() http.Header {
	return m.RespHeader
}

// Mock function which simulates sending an HTTP response header with the specified status code
func (m *MockResponseWriter) WriteHeader(code int) {
	m.Called(code)
}

func TestHTTPResponse_Serve(t *testing.T) {
	// Real version of Pseudo object "r"
	type MatrixItem struct {
		// APIError value to construct test HTTPResponse with
		Error *APIError
		// HTTP status code to expect MockResponseWriter.WriteHeader to be called with
		// Set to -1 if MockResponseWriter.WriteHeader isn't expected to be called
		code int
		// Expected value of error returned by HTTPResponse.Serve
		err error
	}

	// Make a test error to use
	testErr := APIError{"testerr", "msg", http.StatusInternalServerError}

	// Make test matrix
	matrix := []MatrixItem{
		MatrixItem{nil, http.StatusOK, nil},
		MatrixItem{&testErr, testErr.HTTPCode, nil},
	}

	// Test each case
	for _, item := range matrix {
		// Make resp
		resp := HTTPResponse{item.Error}

		// Setup expect for HTTP status code
		writer := new(MockResponseWriter)
		writer.RespHeader = make(map[string][]string)

		// If code wasn't set to -1 (Which is considered the nil value of the field)
		if item.code != -1 {
			writer.On("WriteHeader", item.code)
		}

		// Call Serve function
		err := resp.Serve(writer)

		// Assert
		a := assert.New(t)

		// Check err
		if item.err != nil {
			a.Equal(item.err, err, "Expected: " + item.err.Error() + ", Got: " + err.Error())
		}

		// Check code
		writer.AssertExpectations(t)

		// Check writer.RespHeader for correct Content Type value
		a.Equal([]string{"application/json; charset=utf-8"}, writer.RespHeader["Content-Type"])

		// Check body was marshalled properly, don't check body if this case expects an error (As it would never
		// sent so the body doesn't matter)
		if item.err == nil {
			bytes, err := json.Marshal(resp)
			if err != nil {
				// Handle encoding error
				t.Error("Error marshalling json: " + err.Error())
			}
			a.Equal(string(bytes[:]), writer.Body, "writer.Body should be equal to marshalled resp object")
		}
	}
}
