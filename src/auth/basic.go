package auth

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/D-L-M/mem-db/src/crypt"
	"github.com/D-L-M/mem-db/src/data"
	"golang.org/x/crypto/bcrypt"
)

// userPasswords will hold the user credentials
var userPasswords = map[string]string{}

// userPasswordsLock allows locking of the passwords map during reads/writes
var userPasswordsLock = sync.RWMutex{}

// GetCredentials gets the current username and password
func GetCredentials(request *http.Request) (string, string, error) {

	authHeader := strings.SplitN(request.Header.Get("Authorization"), " ", 2)

	if len(authHeader) != 2 || authHeader[0] != "Basic" {
		return "", "", errors.New("Invalid Basic authorisation header")
	}

	decodedAuth, err := base64.StdEncoding.DecodeString(authHeader[1])

	if err != nil {
		return "", "", err
	}

	authParts := strings.SplitN(string(decodedAuth), ":", 2)

	if len(authParts) != 2 {
		return "", "", errors.New("Malformed Basic authorisation header")
	}

	return authParts[0], authParts[1], nil

}

// CheckBasic checks whether basic authentication has been successful
func CheckBasic(request *http.Request) bool {

	username, password, err := GetCredentials(request)

	if err != nil {
		return false
	}

	return isUsernameAndPasswordValid(username, password)

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

	userPasswordsLock.Lock()
	userPasswords[username] = string(hashedPassword)
	userPasswordsLock.Unlock()

	savePasswordFile()

	return nil

}

// DeleteUser removes a user
func DeleteUser(username string) {

	userPasswordsLock.Lock()
	delete(userPasswords, username)
	userPasswordsLock.Unlock()
	savePasswordFile()

}

// savePasswordFile updates the password file based on the user passwords in
// memory
func savePasswordFile() {

	passwordFilename, err := getPasswordFilePath()

	if err != nil {
		log.Fatal(err)
	}

	userPasswordsLock.RLock()
	passwordFile, err := json.Marshal(userPasswords)
	userPasswordsLock.RUnlock()

	if err != nil {
		log.Fatal(err)
	}

	ioutil.WriteFile(passwordFilename, passwordFile, os.FileMode(0600))

}

// Init loads the password file into a map in memory if it has not already been
// loaded
func Init() {

	// Load Basic authentication credentials
	userPasswordsLock.Lock()

	userPasswords = map[string]string{}
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

	userPasswordsLock.Unlock()

	// Load HMAC authentication secret key
	crypt.SecretKey()

}

// isUsernameAndPasswordValid checks whether a username and password pair
// exists
func isUsernameAndPasswordValid(username string, password string) bool {

	// Look up the user's password and see if the hash is valid
	userPasswordsLock.RLock()

	if hashedPassword, ok := userPasswords[username]; ok {

		userPasswordsLock.RUnlock()

		err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))

		if err == nil {
			return true
		}

	} else {
		userPasswordsLock.RUnlock()
	}

	return false

}

// UserExists checks whether a user exists
func UserExists(username string) bool {

	userPasswordsLock.RLock()
	defer userPasswordsLock.RUnlock()

	if _, ok := userPasswords[username]; ok {
		return true
	}

	return false

}
