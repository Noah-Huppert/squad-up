package models

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type HTTPResponse struct {
	Status ResultEnum `json:"status"`
	Error  *APIError  `json:"error"`
}

// Serve HTTPResponse struct as JSON with appropriate HTTP status code.
func (r *HTTPResponse) Serve(w http.ResponseWriter) error {
	/*
		Sanity check:
			If Status == SUCCESS then Error == nil
			If Status == FAIL then Error != nil

		Allows us to make assumptions when setting the response status coded

		Also avoids useless API responses where Status is fail but no error is given or
		confusing API responses where Status is success but an error is given
	*/
	// Error if Status == SUCCESS but Error isn't nil
	if r.Status == SUCCESS && r.Error != nil {
		return errors.New("If response \"Status\" is \"SUCCESS\" then an \"Error\" can not be provided")
	}

	// Error if Status == FAIL but Error is nil
	if r.Status == FAIL && r.Error == nil {
		return errors.New("If response \"Status\" is \"FAIL\" then an \"Error\" must be provided")
	}

	// Encode response in json
	bytes, err := json.Marshal(r)
	if err != nil {
		// Handle encoding error
		return errors.New("Error marshalling json: " + err.Error())
	}

	// Set headers
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	// Set body
	io.WriteString(w, string(bytes[:]))

	// Send response with custom status code
	if r.Error == nil {
		// Ok status code
		w.WriteHeader(http.StatusOK)
	} else {
		// Custom status code depending on error
		w.WriteHeader(r.Error.HTTPCode)
	}

	return nil
}

func (r *HTTPResponse) WithError(id, message string, code int) *HTTPResponse {
	r.Error = APIError{id, message, code}
	return r
}
