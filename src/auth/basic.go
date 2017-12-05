package auth

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"../crypt"
	"../data"
)

// userPasswords will hold the user credentials
var userPasswords = map[string]string{}

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

// getPasswordFilePath gets the path to the user password file
func getPasswordFilePath() (string, error) {

	baseDirectory, err := data.GetBaseDirectory()

	if err != nil {
		return "", err
	}

	passwordFilename := baseDirectory + "/.passwd"

	return passwordFilename, nil

}

// AddUser adds a new user
func AddUser(username string, password string) {

	userPasswords[username] = crypt.SaltedSha512([]byte(password))

	savePasswordFile()

}

// savePasswordFile updates the password file based on the user passwords in
// memory
func savePasswordFile() {

	passwordFilename, err := getPasswordFilePath()

	if err != nil {
		log.Fatal(err)
	}

	passwordFile, err := json.Marshal(userPasswords)

	if err != nil {
		log.Fatal(err)
	}

	ioutil.WriteFile(passwordFilename, passwordFile, os.FileMode(0600))

}

// loadPasswordFileIfRequired loads the password file into a map in memory if
// it has not already been loaded
func loadPasswordFileIfRequired() {

	if len(userPasswords) == 0 {

		passwordFilename, err := getPasswordFilePath()

		if err != nil {
			log.Fatal(err)
		}

		passwordFile, err := ioutil.ReadFile(passwordFilename)

		// If no username/password file exists, add a fallback root user
		if err != nil {

			AddUser("root", "password")

			// Otherwise load the actual file into memory
		} else {

			err = json.Unmarshal(passwordFile, &userPasswords)

			if err != nil {
				log.Fatal(err)
			}

		}

	}

}

// isUsernameAndPasswordValid checks whether a username and password pair
// exists
func isUsernameAndPasswordValid(username string, password string) bool {

	// Ensure that a user password file has been loaded into memory
	loadPasswordFileIfRequired()

	// Look up the user's password and see if the hashes match
	if userPassword, ok := userPasswords[username]; ok {

		if crypt.SaltedSha512([]byte(password)) == userPassword {
			return true
		}

	}

	return false

}
