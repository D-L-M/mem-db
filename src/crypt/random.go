package crypt

import (
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/D-L-M/mem-db/src/data"
)

// secretKey is the secret key used for hashing
var secretKey = []byte("")

// GetRandomBytes gets a random byte array up to a specified length
func GetRandomBytes(length int) ([]byte, error) {

	bytes := make([]byte, length)
	length, err := io.ReadFull(rand.Reader, bytes)

	if err != nil {
		return []byte(""), err
	}

	if length != len(bytes) {
		return []byte(""), errors.New("Error generating random bytes")
	}

	return bytes, nil

}

// GenerateUUID generates a random UUID
func GenerateUUID() (string, error) {

	uuid, err := GetRandomBytes(16)

	if err != nil {
		return "", err
	}

	uuid[8] = uuid[8]&^0xc0 | 0x80
	uuid[6] = uuid[6]&^0xf0 | 0x40

	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil

}

// getSecretKeyFilePath gets the path to the secret key file
func getSecretKeyFilePath() (string, error) {

	baseDirectory, err := data.GetBaseDirectory()

	if err != nil {
		return "", err
	}

	secretKeyFilename := baseDirectory + "/.key"

	return secretKeyFilename, nil

}

// SecretKey gets the secret key (and creates one if necessary)
func SecretKey() []byte {

	// If the key does not already exist, attempt to load it from disk
	if len(secretKey) == 0 {

		secretKeyFilename, err := getSecretKeyFilePath()

		if err != nil {
			log.Fatal(err)
		}

		secretKey, err = ioutil.ReadFile(secretKeyFilename)

		// If it does not exist on disk, create and save a new one
		if err != nil || len(secretKey) != 32 {

			secretKey, err = GetRandomBytes(32)

			if err != nil {
				log.Fatal(err)
			}

			ioutil.WriteFile(secretKeyFilename, secretKey, os.FileMode(0600))

		}

	}

	return secretKey

}
