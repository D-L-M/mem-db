package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"

	"../data"
)

// Encrypt encrypts data
func Encrypt(input []byte) ([]byte, error) {

	aesCipher, error := aes.NewCipher(data.SecretKey)

	if error != nil {
		return nil, error
	}

	gcm, error := cipher.NewGCM(aesCipher)

	if error != nil {
		return nil, error
	}

	nonce := make([]byte, gcm.NonceSize())

	if _, error = io.ReadFull(rand.Reader, nonce); error != nil {
		return nil, error
	}

	return []byte(hex.EncodeToString(gcm.Seal(nonce, nonce, input, nil))), nil

}

// Decrypt decrypts data
func Decrypt(input []byte) ([]byte, error) {

	decodedInput, error := hex.DecodeString(string(input[:]))

	if error != nil {
		return nil, error
	}

	aesCipher, error := aes.NewCipher(data.SecretKey)

	if error != nil {
		return nil, error
	}

	gcm, error := cipher.NewGCM(aesCipher)

	if error != nil {
		return nil, error
	}

	nonceSize := gcm.NonceSize()

	if len(input) < nonceSize {
		return nil, errors.New("Malformed input ciphertext")
	}

	nonce, decodedInput := decodedInput[:nonceSize], decodedInput[nonceSize:]

	return gcm.Open(nil, nonce, decodedInput, nil)

}
