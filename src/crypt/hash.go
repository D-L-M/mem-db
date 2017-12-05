package crypt

import (
	"crypto/sha512"
	"encoding/hex"
)

// Sha512 generates a SHA512 hex representation of some input data
func Sha512(input []byte) string {

	hasher := sha512.New()

	hasher.Write(input)

	return hex.EncodeToString(hasher.Sum(nil))

}
