package auth

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"../data"
	"golang.org/x/crypto/bcrypt"
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
func AddUser(username string, password string) error {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	userPasswords[username] = string(hashedPassword)

	savePasswordFile()

	return nil

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

// Init loads the password file into a map in memory if it has not already been
// loaded
func Init() {

	passwordFilename, err := getPasswordFilePath()

	if err != nil {
		log.Fatal(err)
	}

	passwordFile, err := ioutil.ReadFile(passwordFilename)

	if err == nil {

		err = json.Unmarshal(passwordFile, &userPasswords)

		if err != nil {
			log.Fatal(err)
		}

	}

}

// isUsernameAndPasswordValid checks whether a username and password pair
// exists
func isUsernameAndPasswordValid(username string, password string) bool {

	// Look up the user's password and see if the hash is valid
	if hashedPassword, ok := userPasswords[username]; ok {

		err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))

		if err == nil {
			return true
		}

	}

	return false

}

// UserExists checks whether a user exists
func UserExists(username string) bool {

	if _, ok := userPasswords[username]; ok {
		return true
	}

	return false

}
