package auth

import (
	"../crypt"
	"../data"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

// CheckBasic checks whether basic authentication has been successful
func CheckBasic(request *http.Request) bool {

	authHeader := strings.SplitN(request.Header.Get("Authorization"), " ", 2)

	if len(authHeader) != 2 || authHeader[0] != "Basic" {
		return false
	}

	decodedAuth, err := base64.StdEncoding.DecodeString(authHeader[1])

	if err != nil {
		return false
	}

	authParts := strings.SplitN(string(decodedAuth), ":", 2)

	if len(authParts) != 2 {
		return false
	}

	return isUsernameAndPasswordValid(authParts[0], authParts[1])

}

// isUsernameAndPasswordValid checks whether a username and password pair exists
func isUsernameAndPasswordValid(username string, password string) bool {

	// Get the username/password file
	baseDirectory, err := data.GetBaseDirectory()

	if err != nil {
		return false
	}

	passwordFilename := baseDirectory + "/passwd"
	passwordFile, err := ioutil.ReadFile(passwordFilename)

	// If no username/password file exists, create one with the root user
	if err != nil {

		hashedRootPassword := crypt.SaltedSha256([]byte("password"))
		passwordFile = []byte("{\"root\": \"" + hashedRootPassword + "\"}")

		ioutil.WriteFile(passwordFilename, passwordFile, os.FileMode(0600))

	}

	var passwords map[string]string

	err = json.Unmarshal(passwordFile, &passwords)

	if err != nil {
		return false
	}

	// Look up the user's password and see if the hashes match
	if userPassword, ok := passwords[username]; ok {

		if crypt.SaltedSha256([]byte(password)) == userPassword {
			return true
		}

	}

	return false

}
