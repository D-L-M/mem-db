package crypt

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
)

// Sha512 generates a SHA512 hex representation of some input data
func Sha512(input []byte) string {

	hasher := sha512.New()

	hasher.Write(input)

	return hex.EncodeToString(hasher.Sum(nil))

}

// Sha512HMAC generates a SHA512 HMAC hash using the application's secret key
func Sha512HMAC(input []byte) string {

	hasher := hmac.New(sha512.New, SecretKey())

	hasher.Write(input)

	return hex.EncodeToString(hasher.Sum(nil))

}
