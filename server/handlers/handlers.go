package handlers

import (
	"fmt"
	"net/http"

	"github.com/Noah-Huppert/squad-up/server/models"
)

type Handler func(r *http.Request) models.HTTPResponse

func Register(mux *http.ServeMux, path string, handler Handler) {
	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		resp := handler(r)

		err := resp.Serve(w)
		if err != nil {
			fmt.Println("Error serving response: " + err.Error())
		}
	})
}

func LoadAll(mux *http.ServeMux) {
	mux.Handle("/lib/", http.StripPrefix("/lib/", http.FileServer(http.Dir("bower_components"))))
	mux.HandleFunc("/", ServeIndex)

	Register(mux, "/api/v1/auth/google/token", ExchangeTokenHandler)
}
