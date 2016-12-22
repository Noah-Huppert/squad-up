package handlers

import (
    "net/http"

    "github.com/Noah-Huppert/squad-up/server/models"
    "fmt"
    "encoding/json"
    "io"
    "github.com/fatih/structs"
)

type EndpointHandler interface {
    Serve (ctx *models.AppContext, r *http.Request) (interface{}, *models.APIError)
}

type handler struct {
    EndpointHandler
    models.AppContextProvider
}

func (h handler) ServeHTTP (w http.ResponseWriter, r *http.Request) {
    // Call logic for endpoint
    respData, respErr := h.Serve(h.Ctx(), r)

    // Convert respData to a map
    var respDataMap map[string]interface{}
    if respData != nil && structs.IsStruct(respData) {
        respDataMap = structs.Map(respData)
    } else {
        respDataMap = make(map[string]interface{}, 0)

        if structs.IsStruct(respData) == false {
            fmt.Println("Response returned by endpoint handler was not a struct, setting to empty map")
            respErr = &models.APIError{"data_not_struct", "The response returned by the endpoint handler was an illegal format", http.StatusInternalServerError}
        }
    }

    // Check that the interface to serve as response returned by h.Serve doesn't contain the key "error"
    if _, ok := respDataMap["error"]; ok == true {
        fmt.Println("EndpointHandler.Serve cannot return data which contains the \"error\" field. Replacing with empty interface")

        if respErr != nil {
            fmt.Println("Replacing existing error (" + respErr.Error() + ") with custom to notify of issue")
        }

        respDataMap = make(map[string]interface{}, 0)
        respErr = &models.APIError{"data_contains_error_key", "The response returned by the endpoint handler contained an illegal value", http.StatusInternalServerError}
    }

    // Add error to response
    respDataMap["error"] = respErr

    // Encode response in json
	bytes, err := json.Marshal(respDataMap)
	if err != nil {
		// Handle encoding error
		fmt.Println("Error marshalling response into JSON: " + err.Error())
        http.Error(w, models.API_ERROR_MANUAL_MARSHALLED_ERROR_MARSHALLING_HTTP_RESPONSE, http.StatusInternalServerError)
        return
	}

	// Set headers
	w.Header().Set("Content-Type", "application/json; charset=utf-8")


	// Send response with custom status code
	if respErr == nil {
		// Ok status code
		w.WriteHeader(http.StatusOK)
	} else {
		// Custom status code depending on error
		w.WriteHeader(respErr.HTTPCode)
	}

	// Set body
	io.WriteString(w, string(bytes[:]))
}

// HandlerLoader interface provides methods useful for loading large numbers of http handlers
type HandlerLoader interface {
    Load ()
}

// Loader is a data structure for this file's HandlerLoader implementation.
type Loader struct {
    mux *http.ServeMux
    ctx *models.AppContext
}

func NewLoader (mux *http.ServeMux, ctx *models.AppContext) Loader {
    l := Loader{}
    l.mux = mux
    l.ctx = ctx

    return l
}

func (l Loader) Ctx () *models.AppContext {
    return l.ctx
}

// register's the provided handler for the provided path with the http.ServeMux
func (l Loader) register (path string, eHdlr EndpointHandler) {
    hdlr := handler{eHdlr, l}

    l.mux.Handle(path, hdlr)
}

func (l Loader) Load() {
    // Static paths
    l.mux.Handle("/lib/", http.StripPrefix("/lib/", http.FileServer(http.Dir("bower_components"))))
    l.mux.HandleFunc("/", ServeIndex)

    // API
    l.register("/api/v1/auth/google/token", ExchangeTokenHandler{})
}
