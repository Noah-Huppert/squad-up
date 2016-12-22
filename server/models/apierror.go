package models

import "fmt"

// APIError provides detail about an error that occurred while handling an endpoint
//
// IMPORTANT:
// Update the manual field below (API_ERROR_MANUAL_MARSHALLED_ERROR_MARSHALLING_HTTP_RESPONSE) if any changes are made to the field
// or their tags.
type APIError struct {
	Id       string `json:"id"`
	Message  string `json:"message"`
	HTTPCode int    `json:"http_code"`
}

// String representing the following APIError in JSON form:
//     Id: error_marshalling_http_response
//     Message: An internal error occurred while generating the response
//     HTTPCode: 500
//
// This field is here for use when an error occurs marshalling an HTTPResponse object to send to a client.
// It was put in this file (Although used elsewhere) so that is is kept up to date with any field changes of APIError.
// TODO: Rename to CamelCase
var API_ERROR_MANUAL_MARSHALLED_ERROR_MARSHALLING_HTTP_RESPONSE = "{" +
                                                            "\"error\":" +
                                                                "{" +
                                                                    "\"id\":\"error_marshalling_http_response\"," +
                                                                    "\"message\":\"An internal error occurred while generating the response\"," +
                                                                    "\"http_code\": 500" +
                                                                "}" +
                                                        "}"

func (e APIError) Error() string {
	return fmt.Sprintf("%v (%v: %v)", e.Message, e.Id, e.HTTPCode)
}
