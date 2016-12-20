package handlers

import "net/http"

// ServeIndex responds with index.html view
func ServeIndex(w http.ResponseWriter, r *http.Request) {
	// File path (3rd arg) is relative to Go cli wkr dir (go/src/squad-up)
	http.ServeFile(w, r, "client/views/index.html")
}
