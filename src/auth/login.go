package auth

import (
	"net/http"
)

// CheckCredentials checks whether any kind of authentication has been
// successful
func CheckCredentials(request *http.Request, body *[]byte) bool {

	if CheckBasic(request) || CheckHMAC(request, body) {
		return true
	}

	return false

}
