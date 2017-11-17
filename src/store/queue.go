package store


import (
	"../types"
)


// Perform queued actions and flush document changes to disk
func IndexOnDisk(documentMessage chan types.DocumentMessage) {
	
	for {
		
		message := <- documentMessage

		if message.Action == "add" {

			IndexDocument(message.Id, message.Document)
			
		}

		if message.Action == "remove" {

			RemoveDocument(message.Id)
			
		}

	}

}