package models

import "fmt"

type APIError struct {
	Id       string `json:"id"`
	Message  string `json:"message"`
	HTTPCode int    `json:"http_code"`
}

func (e APIError) Error() string {
	return fmt.Sprintf("%v (%v: %v)", e.Message, e.Id, e.HTTPCode)
}
