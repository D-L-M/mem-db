package crypt

import (
	"crypto/sha512"
	"encoding/hex"

	"../data"
)

// Sha512 generates a SHA512 hex representation of some input data
func Sha512(input []byte) string {

	hasher := sha512.New()

	hasher.Write(input)

	return hex.EncodeToString(hasher.Sum(nil))

}

// SaltedSha512 generates a SHA512 hex representation of some input data with a
// salt taken from the application configuration
func SaltedSha512(input []byte) string {

	return Sha512([]byte(string(data.SecretKey[:]) + string(input[:])))

}
