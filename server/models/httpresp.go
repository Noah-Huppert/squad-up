package models

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type HTTPResponse struct {
    Error  *APIError  `json:"error"`
}

// Serve HTTPResponse struct as JSON with appropriate HTTP status code.
func (r *HTTPResponse) Serve(w http.ResponseWriter) error {
	// Encode response in json
	bytes, err := json.Marshal(r)
	if err != nil {
		// Handle encoding error
		return errors.New("Error marshalling json: " + err.Error())
	}

	// Set headers
	w.Header().Set("Content-Type", "application/json; charset=utf-8")


	// Send response with custom status code
	if r.Error == nil {
		// Ok status code
		w.WriteHeader(http.StatusOK)
	} else {
		// Custom status code depending on error
		w.WriteHeader(r.Error.HTTPCode)
	}

	// Set body
	io.WriteString(w, string(bytes[:]))

	return nil
}

func (r *HTTPResponse) WithError(id, message string, code int) *HTTPResponse {
	r.Error = &APIError{id, message, code}
	return r
}
