package handlers

import (
    "net/http"

    "github.com/Noah-Huppert/squad-up/server/models"
    "github.com/Noah-Huppert/squad-up/server/models/utils"

    "fmt"
    "encoding/json"
    "io"
    "github.com/fatih/structs"
)

// Custom handler type for our web app. Provides a models.AppContext struct which allows the handler to interact with
// the state of the application.
type EndpointHandler interface {
    // Serve is called by the `handler` struct defined in this file (handlers.go)
    // Given application context and an http request.
    // Returns an interface to serve back to the client (Cannot contain the "error" key) and a point to an models.APIError
    // both can be nil
    Serve (ctx *models.AppContext, r *http.Request) (interface{}, *models.APIError)
}

// handler is the custom http.Handler used to serve our web app endpoints. Calls the EndpointHandler.Serve method to get
// the response to serve back to the client.
type handler struct {
    // Embedded EndpointHandler to handle client requests
    EndpointHandler
    // Embedded AppContextProvider used to get the context for the EndpointHandler.Serve method.
    models.AppContextProvider
}

// ServeHTTP calls the custom EndpointHandler to handle the request and serves the result.
func (h handler) ServeHTTP (w http.ResponseWriter, r *http.Request) {
    // Call EndpointHandler.Serve
    hdlrRes, hdlrErr := h.Serve(h.Ctx(), r)

    // convert endpoint handler result into a map
    var resMap map[string]interface{}

    // Set to value if error occurs in the conversion process, if not nil will replace
    // hdlrErr in response
    var convertErr *models.APIError

    // Check that result is indeed a struct
    if structs.IsStruct(hdlrRes) == false {// If hdlrRes is not a struct print and set error
        fmt.Println("Endpoint handler returned invalid type (Non struct) as result")
        convertErr = &models.APIError{"endpoint_handler_invalid_result_type", "The handler for this endpoint returned a result with an invalid type", http.StatusInternalServerError}
    } else {// If hdlrRes is a struct convert to map[string]interface{}
        resStruct := structs.New(hdlrRes)

        m, err := utils.ToMap(resStruct)
        if err != nil {
            fmt.Println("Error converting resStruct into map: " + err.Error())
            convertErr = &models.APIError{"err_converting_endpoint_handler_result", "An internal error occured in an intermedierary json conversion step", http.StatusInternalServerError}
        }

        resMap = m
    }

    // Replace hdlrErr with any conversion error that may have occurred.
    if convertErr != nil {
        // Set result map to nothing
        resMap = make(map[string]interface{}, 0)

        // Override hdlrErr
        hdlrErr = convertErr
    }

    // Assign error to "error" key
    resMap["error"] = hdlrErr

    // Encode response in json
	bytes, err := json.Marshal(resMap)
	if err != nil {// Handle encoding error
		fmt.Println("Error marshalling response into JSON: " + err.Error())

        // Manually serve encoding error
        bytes = []byte(models.APIErrorManualMarshalledErrorMarshallingHTTPResponse)

        // Set HTTP response code
        hdlrErr = models.APIErrorErrorMarshallingHTTPResponse
        return
	}

	// Set headers
	w.Header().Set("Content-Type", "application/json")


	// Send response with custom status code
	if hdlrErr == nil {
		// Ok status code
		w.WriteHeader(http.StatusOK)
	} else {
		// Custom status code depending on error
		w.WriteHeader(hdlrErr.HTTPCode)
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
func (l Loader) registerEndpoint(path string, eHdlr EndpointHandler) {
    hdlr := handler{eHdlr, l}

    l.mux.Handle(path, hdlr)
}

func (l Loader) registerFile (path, file string) {
    l.mux.HandleFunc(path, func (w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, file)
    })
}

func (l Loader) registerDir (path, dir string) {
    l.mux.Handle(path, http.StripPrefix(path, http.FileServer(http.Dir(dir))))
}

func (l Loader) Load() {
    // Resources
    l.registerDir("/lib/", "bower_components")
    l.registerDir("/js/", "dist/js")
    l.registerDir("/css/", "client/css")
    l.registerDir("/components/", "dist/components")

    // Pages
    l.mux.HandleFunc("/", ServeIndex)

    // API
    l.registerEndpoint("/api/v1/auth/token/google", ExchangeTokenHandler{})
}
