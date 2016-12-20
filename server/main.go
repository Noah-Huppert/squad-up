// Main HTTP server package for Squad Up.
package main

// Import deps.
import (
	"fmt"
	"net/http"

	"github.com/Noah-Huppert/squad-up/server/handlers"
)

// serveIndex responds with index.html view
func serveIndex(w http.ResponseWriter, r *http.Request) {
	// File path (3rd arg) is relative to Go cli wkr dir (go/src/squad-up)
	http.ServeFile(w, r, "client/views/index.html")
}

// Main entry point of program.
func main() {
	// New HTTP router.
	mux := http.NewServeMux()

	// Mux handlers
	mux.Handle("/lib/", http.StripPrefix("/lib/", http.FileServer(http.Dir("bower_components"))))
	mux.HandleFunc("/", serveIndex)

	mux.HandleFunc("/api/v1/auth/google/token", handlers.ExchangeToken)

	// Start listening on any host, port 5000.
	fmt.Println("Listening on :5000")
	err := http.ListenAndServe(":5000", mux)
	if err != nil { // If err print to console.
		fmt.Println("Error starting HTTP server on :5000: " + err.Error())
	}
}
