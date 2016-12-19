package models

import (
	"testing"
	"github.com/stretchr/testify/mock"
	"net/http"
	"github.com/stretchr/testify/assert"
	"errors"
	"encoding/json"
)

// Mock http header with kv
type MockHeader map[string][]string

// Mock function used in Serve method
func (h MockHeader) Set(key, value string) {
	h.Set(key, value)
}



// Mock http.ResponseWriter for testing
type MockResponseWriter struct {
	mock.Mock

	Body string
	RespHeader MockHeader
}

// Mock write function to make struct a Writer
func (m *MockResponseWriter) Write(p []byte) (n int, err error) {
	m.Body = append(m.Body, p...)
	return len(p), nil
}

// Getter for header field, used in HTTPResponse.Serve
func (m *MockResponseWriter) Header() MockHeader {
	return m.RespHeader
}

// Mock function which simulates sending an HTTP response header with the specified status code
func (m *MockResponseWriter) WriteHeader(code int) {
	m.Called(code)
}

/*
HTTPResponse.Serve test matrix:
	Pseudo objects:
		r {
			Status: Either "SUCCESS" or "FAIL"
			Error: Either "nil" (unset) or "testerr" which is the models.APIError struct:
				models.APIError { Id: "testerr", Message: "msg", HTTPCode: 500 }
		}
		code: HTTP status code that Serve() method tried to respond with (nil if method never responds)
		err: Error returned by Serve() method

	Matrix:
	- r.Status: SUCCESS
		- r.Error: nil ~> assert(code == 200), assert(err == nil)
		- r.Error: testerr ~> assert(code == nil), assert(err == errors.New("If response \"Status\" is \"SUCCESS\" then an \"Error\" can not be provided"))
	- r.Status: FAIL
		- r.Error: nil ~> assert(code == nil), assert(err == errors.New("If response \"Status\" is \"FAIL\" then an \"Error\" must be provided")
		- r.Error: testerr ~> assert(code == 500), assert(err == nil)
 */
func TestHTTPResponse_Serve(t *testing.T) {
	// Real version of Pseudo object "r"
	type MatrixItem struct {
		// ResultEnum value to construct test HTTPResponse with
		Status ResultEnum
		// APIError value to construct test HTTPResponse with
		Error APIError
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
		MatrixItem{SUCCESS, nil, http.StatusOK, nil},
		MatrixItem{SUCCESS, testErr, testErr.HTTPCode, errors.New("If response \"Status\" is \"SUCCESS\" then an \"Error\" can not be provided")},
		MatrixItem{FAIL, nil, -1, errors.New("If response \"Status\" is \"FAIL\" then an \"Error\" must be provided")},
		MatrixItem{FAIL, testErr, testErr.HTTPCode, nil},
	}

	// Test each case
	for _, item := range matrix {
		// Make resp
		resp := HTTPResponse{item.Status, item.Error}

		// Setup expect for HTTP status code
		writer := new(MockResponseWriter)

		// If code wasn't set to -1 (Which is considered the nil value of the field)
		if item.code != -1 {
			writer.On("WriteHeader", item.code)
		}

		// Call Serve function
		err := resp.Serve(writer)

		// Assert
		assert := assert.New(t)

		// Check err
		assert.Equal(item.err, err, "[Status: " + item.Status + ", Error: " + item.Error + "] Expected: " +
			item.err + ", Got: " + err)

		// Check code
		writer.AssertExpectations(t)

		// Check writer.RespHeader for correct Content Type value
		assert.Equal(writer.RespHeader["Content-Type"], "application/json; charset=utf-8")

		// Check body was marshalled properly, don't check body if this case expects an error (As it would never
		// sent so the body doesn't matter)
		if item.err == nil {
			bytes, err := json.Marshal(resp)
			if err != nil {
				// Handle encoding error
				t.Error("[Status: " + item.Status + ", Error: " + item.Error + "] Error marshalling json: " + err)
			}
			assert.Equal(string(bytes[:]), writer.Body, "writer.Body should be equal to marshalled resp object")
		}
	}
}
