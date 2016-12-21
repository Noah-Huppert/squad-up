package handlers

import (
    "net/http"

    "github.com/Noah-Huppert/squad-up/server/models"
    "fmt"
)

// Custom HandlerFunc function type. Includes the app context and automatically serves the returned HTTPResponse.
type HandlerFunc func (ctx *models.AppContext, r *http.Request) models.HTTPResponse

// HandlerLoader is a helper struct with methods which help load our custom handler functions
type HandlerLoader struct {
    // The http.Mux to register the routes on
    Mux *http.ServeMux
    // Application context
    Ctx *models.AppContext
}

// Register is a helper function which registers one of our custom HandlerFunc functions to a path.
func (h *HandlerLoader) register (path string, f HandlerFunc) {
    // Register a function which calls our special HandlerFunc typed function and deals with the result
    h.Mux.HandleFunc(path, func (w http.ResponseWriter, r *http.Request) {
        // Call HandlerFunc with context and request
        httpResp := f(h.Ctx, r)

        // Serve HTTPResponse
        err := httpResp.Serve(w)
        // Deal with error in serving http response
		if err != nil {
            fmt.Println("Error serving HTTPResponse: " + err.Error())
            http.Error(w, models.API_ERROR_MANUAL_MARSHALLED_ERROR_MARSHALLING_HTTP_RESPONSE, http.StatusInternalServerError)
		}
    })
}

func (h *HandlerLoader) Load() {
    // Static paths
    h.Mux.Handle("/lib/", http.StripPrefix("/lib/", http.FileServer(http.Dir("bower_components"))))
    h.Mux.HandleFunc("/", ServeIndex)

    // API
    h.register("/api/v1/auth/google/token", ExchangeTokenHandler)
}
