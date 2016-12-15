// Main HTTP server package for Squad Up.
package main

// Import deps.
import (
	"fmt"
	"io/ioutil"
	"net/http"
	"encoding/json"
	"errors"
)

type ResultEnum int
const (
	SUCCESS ResultEnum = iota
	FAIL
)

func (this ResultEnum) MarshalJSON() ([]byte, error) {
	switch this {
	case SUCCESS:
		return []byte("success"), nil
	case FAIL:
		return []byte("fail"), nil
	}

	return nil, errors.New("Uknown value of ResultEnum: " + this)
}

func (this ResultEnum) UnmarshalJSON(data []byte) error {
	switch data {
	case []byte("success"):
		this = SUCCESS
		return nil
	case []byte("fail"):
		this = FAIL
		return nil
	}

	return errors.New("Uknown value for ResultEnum: " + data)
}

type HTTPResponse struct {
	Status ResultEnum `json:"status"`
	Errors []error `json:"errors"`
}

// serveIndex responds with index.html view
func serveIndex (w http.ResponseWriter, r *http.Request) {
	// File path (3rd arg) is relative to Go cli wkr dir (go/src/squad-up)
	http.ServeFile(w, r, "client/views/index.html")
}

// Main entry point of program.
func main() {
	// GAPI
	gapiClientId := "432144215744-2n6fha955i4f2en9jubvelfhmdsh1jcv.apps.googleusercontent.com"

	// New HTTP router.
	mux := http.NewServeMux()

	// Mux handlers
	mux.Handle("/lib/", http.StripPrefix("/lib/", http.FileServer(http.Dir("bower_components"))))
	mux.HandleFunc("/", serveIndex)

	mux.HandleFunc("/api/v1/auth/google/token", func (w http.ResponseWriter, r *http.Request) {
		// Get id_token passed in request
		idToken := r.PostFormValue("id_token")
		if len(idToken) == 0 {
			http.Error(w, "`id_token` must be provided as a post parameter", http.StatusUnprocessableEntity)
			return
		}

		// Make request to token info Gapi. This lets Google take care of
		// verifying the id token. If the token is valid it also provides
		// us with some basic profile info
		res, err := http.Get("https://www.googleapis.com/oauth2/v3/tokeninfo?id_token=" + idToken)
		if err != nil {
			fmt.Printf("Error sending HTTP request to verify id token: %s\n", err)
			http.Error(w, "Internal Error", http.StatusInternalServerError)
			return
		}

		// Read response body
		body, err := ioutil.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			fmt.Printf("Error reading body of response to verify id token %s\n", err)
			http.Error(w, "Internal Error", http.StatusInternalServerError)
			return
		}

		// Struct to unmarshal json resp into. Not all fields are
		// represented, only fields we care about.
		type GAPIIdTokenResp struct {
			// Id Token fields parsed by Gapi (Id Token is a JWT).
			// Audience field - Should be our apps Gapi client id.
			Aud string `json:"aud"`
			// Subject field - Authenticating user Gapi user id.
			Sub string `json:"sub"`

			// Google api fields
			// User email
			Email string `json:"email"`
			// If user has verified their email with Google
			EmailVerified bool `json:"email_verified,string"`
			// Url of profile picture
			Picture string `json:"picture"`
			// First name
			GivenName string `json:"given_name"`
			// Last name
			FamilyName string `json:"family_name"`
			// Locale string
			Locale string `json:"locale"`
		}

		// Decode response into json
		var resp GAPIIdTokenResp

		err = json.Unmarshal(body, &resp)
		if err != nil {
			fmt.Printf("Error decoding json response: %s\n", err)
			http.Error(w, "Internal Error", http.StatusInternalServerError)
			return
		}

		// Check
		// Check that aud is our client id
		if resp.Aud != gapiClientId {
			http.Error(w, "Invalid id token", http.StatusUnauthorized)
			return
		}

		// Check that email is verified
		if resp.EmailVerified == false {
			http.Error(w, "Email not verified", http.StatusUnauthorized)
			return
		}

		fmt.Fprintf(w, )
	})

	// Start listening on any host, port 5000.
	err := http.ListenAndServe(":5000", mux)
	if err != nil { // If err print to console.
		fmt.Println("Error starting HTTP server on :5000 -> {}", err)
	} else { // Else print normal OK status message.
		fmt.Println("Listening on :5000")
	}
}
