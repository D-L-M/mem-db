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

	aesCipher, err := aes.NewCipher(data.SecretKey)

	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(aesCipher)

	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())

	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	return []byte(hex.EncodeToString(gcm.Seal(nonce, nonce, input, nil))), nil

}

// Decrypt decrypts data
func Decrypt(input []byte) ([]byte, error) {

	decodedInput, err := hex.DecodeString(string(input[:]))

	if err != nil {
		return nil, err
	}

	aesCipher, err := aes.NewCipher(data.SecretKey)

	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(aesCipher)

	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()

	if len(input) < nonceSize {
		return nil, errors.New("Malformed input ciphertext")
	}

	nonce, decodedInput := decodedInput[:nonceSize], decodedInput[nonceSize:]

	return gcm.Open(nil, nonce, decodedInput, nil)

}
