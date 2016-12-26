package models

import "fmt"

// APIError provides detail about an error that occurred while handling an endpoint
//
// IMPORTANT:
// Update the manual fields below (APIErrorErrorMarshallingHTTPResponse and APIErrorManualMarshalledErrorMarshallingHTTPResponse)
// if any changes are made to the field
// or their tags.
type APIError struct {
	Id       string `json:"id"`
	Message  string `json:"message"`
	HTTPCode int    `json:"http_code"`
}

// Error served when there is an error encoding the provided data into json for a response.
var APIErrorErrorMarshallingHTTPResponse = &APIError{"error_marshalling_http_response", "An internal error occured while generating the response", 500}

// String representing APIErrorErrorMarshallingHTTPResponse in JSON form
//
// This field is here for use when an error occurs marshalling an HTTPResponse object to send to a client.
// It was put in this file (Although used elsewhere) so that is is kept up to date with any field changes of APIError.
var APIErrorManualMarshalledErrorMarshallingHTTPResponse = "" +
                "{" +
				    "\"error\":" +
					"{" +
					    "\"id\":\"" + APIErrorErrorMarshallingHTTPResponse.Id + "\"," +
					    "\"message\":\"" + APIErrorErrorMarshallingHTTPResponse.Message + "\"," +
					    "\"http_code\": " + string(APIErrorErrorMarshallingHTTPResponse.HTTPCode) +
					"}" +
				"}"

func (e APIError) Error() string {
	return fmt.Sprintf("%v (%v: %v)", e.Message, e.Id, e.HTTPCode)
}
