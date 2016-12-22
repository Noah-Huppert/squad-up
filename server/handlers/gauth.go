package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
    "fmt"

	"github.com/Noah-Huppert/squad-up/server/models"
    "github.com/Noah-Huppert/squad-up/server/models/db"
)

type ExchangeTokenHandler struct {}

type exchangeResponse struct {
    User db.User
}

// Exchange users Google Id Token for a Squad Up API token, essentially the "login" endpoint.
func (h ExchangeTokenHandler) Serve (ctx *models.AppContext, r *http.Request) (interface{}, *models.APIError) {
	httpResp := exchangeResponse{}

	// Get id_token passed in request
	idToken := r.PostFormValue("id_token")
	if len(idToken) == 0 {
		err := &models.APIError{"missing_param", "`id_token` must be provided as a post parameter", http.StatusUnprocessableEntity}
		return nil, err
	}

	// Make request to token info Gapi. This lets Google take care of
	// verifying the id token. If the token is valid it also provides
	// us with some basic profile info
	res, err := http.Get("https://www.googleapis.com/oauth2/v3/tokeninfo?id_token=" + idToken)
	if err != nil {
		fmt.Printf("Error sending HTTP request to verify id token: %s\n", err)

		err := &models.APIError{"http_err_verifying_id_token", "An error occured while contacting Google servers to verify your identity", http.StatusInternalServerError}
		return nil, err
	}

	// Read response body
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		fmt.Printf("Error reading body of response to verify id token %s\n", err)

		err := &models.APIError{"body_read_err_verifying_id_token", "We couldn't read the Google server's response while verifying your identity", http.StatusInternalServerError}
		return nil, err
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

		err := &models.APIError{"json_parse_err_verifying_id_token", "We couldn't understand the response the Google server's gave us while verifying your identity", http.StatusInternalServerError}
		return nil, err
	}

	// Check
	// Check that aud is our client id
	if resp.Aud != models.GapiConf.ClientId {
		err := &models.APIError{"invalid_id_token", "Google login not valid", http.StatusUnauthorized}
		return nil, err
	}

	// Check that email is verified
	if resp.EmailVerified == false {
		err := &models.APIError{"email_not_verified", "Your email is not verifeid with Google", http.StatusUnauthorized}
		return nil, err
	}

    // Try and find Google user in our system
    var user db.User
    ctx.Db.First(&user, "email = ?", resp.Email)

    httpResp.User = user

	return httpResp, nil
}
