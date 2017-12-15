package auth

import (
	"net/http"
)

// CheckCredentials checks whether any kind of authentication has been
// successful
func CheckCredentials(request *http.Request) bool {

	if CheckBasic(request) || CheckHMAC(request) {
		return true
	}

	return false

}
