package models

// Config holds application configuration values
type Config struct {
    // Google API Client Id
    GAPIClientId string
    // URI used in JWTs to identify server
    JWTServerURI string
    // Key used to sign JWSs with HS512
    JWTHMACKey string
}
