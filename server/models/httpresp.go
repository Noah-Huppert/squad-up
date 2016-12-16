package models

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type HTTPResponse struct {
	Status ResultEnum `json:"status"`
	Error  APIError   `json:"error"`
}

func (r HTTPResponse) Serve(w http.ResponseWriter) {
	// Encode response in jsson
	bytes, err := json.Marshal(r)
	if err != nil {
		// Handle encoding error
		fmt.Println("Error marshalling json: ", err)
		return
	}

	// Set headers
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	// Set status
	if r.Error == nil {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(r.Error.HTTPCode)
	}

	// Write body
	io.WriteString(w, bytes)
}
