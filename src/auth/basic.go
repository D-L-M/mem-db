package auth

import (
	"encoding/base64"
	"net/http"
	"strings"
)

// CheckBasic checks whether basic authentication has been successful
func CheckBasic(request *http.Request) bool {

	authHeader := strings.SplitN(request.Header.Get("Authorization"), " ", 2)

	if len(authHeader) != 2 || authHeader[0] != "Basic" {
		return false
	}

	decodedAuth, _ := base64.StdEncoding.DecodeString(authHeader[1])
	authParts := strings.SplitN(string(decodedAuth), ":", 2)

	if len(authParts) != 2 {
		return false
	}

	if authParts[0] != "root" || authParts[1] != "password" {
		return false
	}

	return true

}
