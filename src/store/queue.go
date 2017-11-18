package store


import (
	"../types"
	"os"
	"io/ioutil"
	"encoding/json"
	"../crypt"
	"../data"
)


// Perform queued actions and flush document changes to disk
func FlushToDisk(documentMessage chan types.DocumentMessage) {

	storageDirectory := data.GetStorageDirectory()
	
	// Listen for messages to process
	for {
		
		message          := <- documentMessage
		documentFilename := storageDirectory + "/" + crypt.Sha256([]byte(message.Id)) + ".json"

		// Add a document to the index and write it to disk
		if message.Action == "add" {

			IndexDocument(message.Id, message.Document)

			documentFile, error := json.Marshal(types.JsonDocument{"id": message.Id, "document": string(message.Document[:])})

			if (error == nil) {
				ioutil.WriteFile(documentFilename, documentFile, os.FileMode(0600))
			}
			
		}

		// Remove a document from the index and disk
		if message.Action == "remove" {

			RemoveDocument(message.Id)

			os.Remove(documentFilename)
			
		}

	}

}