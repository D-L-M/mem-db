package crypt


import (
	"crypto/rand"
	"io"
	"fmt"
	"errors"
)


// Generate a random UUID
func GenerateUuid() (string, error) {

	uuid		  := make([]byte, 16)
	length, error := io.ReadFull(rand.Reader, uuid)

	if error != nil || length != len(uuid) {
		return "", errors.New("Error generating random bytes")
	}

	uuid[8] = uuid[8]&^0xc0 | 0x80
	uuid[6] = uuid[6]&^0xf0 | 0x40

	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil

}