package crypt

import (
	"crypto/rand"
	"errors"
	"fmt"
	"io"
)

// GenerateUUID generates a random UUID
func GenerateUUID() (string, error) {

	uuid := make([]byte, 16)
	length, err := io.ReadFull(rand.Reader, uuid)

	if err != nil || length != len(uuid) {
		return "", errors.New("Error generating random bytes")
	}

	uuid[8] = uuid[8]&^0xc0 | 0x80
	uuid[6] = uuid[6]&^0xf0 | 0x40

	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil

}
