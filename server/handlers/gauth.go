package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/Noah-Huppert/squad-up/server/models"
	"fmt"
)

// Exchange users Google Id Token for a Squad Up API token, essentially the "login" endpoint.
func ExchangeTokenHandler (r *http.Request) models.HTTPResponse {
	httpResp := models.HTTPResponse{}

	// Get id_token passed in request
	idToken := r.PostFormValue("id_token")
	if len(idToken) == 0 {
		httpResp.WithError("missing_param", "`id_token` must be provided as a post parameter", http.StatusUnprocessableEntity)
		return httpResp
	}

	// Make request to token info Gapi. This lets Google take care of
	// verifying the id token. If the token is valid it also provides
	// us with some basic profile info
	res, err := http.Get("https://www.googleapis.com/oauth2/v3/tokeninfo?id_token=" + idToken)
	if err != nil {
		fmt.Printf("Error sending HTTP request to verify id token: %s\n", err)
		httpResp.WithError("http_err_verifying_id_token", "An error occured while contacting Google servers to verify your identity", http.StatusInternalServerError)
		return httpResp
	}

	// Read response body
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		fmt.Printf("Error reading body of response to verify id token %s\n", err)
		httpResp.WithError("body_read_err_verifying_id_token", "We couldn't read the Google server's response while verifying your identity", http.StatusInternalServerError)
		return httpResp
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
		httpResp.WithError("json_parse_err_verifying_id_token", "We couldn't understand the response the Google server's gave us while verifying your identity", http.StatusInternalServerError)
		return httpResp
	}

	// Check
	// Check that aud is our client id
	if resp.Aud != models.GapiConf.ClientId {
		httpResp.WithError("invalid_id_token", "Google login not valid", http.StatusUnauthorized)
		return httpResp
	}

	// Check that email is verified
	if resp.EmailVerified == false {
		httpResp.WithError("email_not_verified", "Your email is not verifeid with Google", http.StatusUnauthorized)
		return httpResp
	}

	return httpResp
}
