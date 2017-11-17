package store


import (
	"../types"
)


// Flush document changes to disk
func IndexOnDisk(documentMessage chan types.DocumentMessage) {
	
	for {
		_ = <- documentMessage
	}

}