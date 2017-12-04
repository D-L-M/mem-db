package crypt

import (
	"../data"
	"crypto/sha256"
	"encoding/hex"
)

// Sha256 generates a SHA256 hex representation of some input data
func Sha256(input []byte) string {

	hasher := sha256.New()

	hasher.Write(input)

	return hex.EncodeToString(hasher.Sum(nil))

}

// SaltedSha256 generates a SHA256 hex representation of some input data with a
// salt taken from the application configuration
func SaltedSha256(input []byte) string {

	return Sha256([]byte(string(data.SecretKey[:]) + string(input[:])))

}
