package auth

import (
	"net/http"

	"github.com/D-L-M/mem-db/src/crypt"
)

// CheckHMAC checks whether HMAC authentication has been successful
func CheckHMAC(request *http.Request, body *[]byte) bool {

	hmacAuth := request.Header.Get("x-hmac-auth")
	hmacNonce := request.Header.Get("x-hmac-nonce")

	if hmacAuth == "" || hmacNonce == "" {
		return false
	}

	inputToHash := []byte(string((*body)[:]) + hmacNonce)
	calculatedHash := crypt.Sha512HMAC(inputToHash)

	if calculatedHash == hmacAuth {
		return true
	}

	return false

}
