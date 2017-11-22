package crypt

import (
	"crypto/sha256"
	"encoding/hex"
)

// Sha256 generates a SHA256 hex representation of some input data
func Sha256(input []byte) string {

	hasher := sha256.New()

	hasher.Write(input)

	return hex.EncodeToString(hasher.Sum(nil))

}
