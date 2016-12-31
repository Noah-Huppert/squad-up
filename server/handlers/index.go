package handlers

import "net/http"

// ServeIndex responds with index.html view
func ServeIndex(w http.ResponseWriter, r *http.Request) {
    // This handler is registed with ServeMux at "/". As a result ServeMux will use this handler in the event that "/"
    // is requested or in the event that a request does not have a registered handler. If the later is the case we will
    // respond with 404.
    if r.URL.Path != "/" {
        http.NotFound(w, r)
        return
    }

    // File path (3rd arg) is relative to Go cli wkr dir (go/src/squad-up)
    http.ServeFile(w, r, "client/views/index.html")
}
