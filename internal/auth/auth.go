package auth

import (
	"errors"
	"net/http"
	"strings"
)

// GetAPIKey extracts an API Key from the
// headers of an HTTP request
// Example:
// Authorization: ApiKey {insert apiKey here}
func GetAPIKey(headers http.Header) (string, error) {
	key := headers.Get("Authorization")
	if key == "" {
		return "", errors.New("no authentication info found")
	}

	vals := strings.Split(key, " ")
	if len(vals) != 2 {
		return "", errors.New("malformed auth header")
	}

	if vals[0] != "ApiKey" {
		return "", errors.New("malformed first part of auth header")
	}

	return vals[1], nil
}
