package store

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"../crypt"
	"../data"
	"../types"
)

// FlushToDisk performs queued actions and flush document changes to disk
func FlushToDisk(documentMessageQueue chan types.DocumentMessage) {

	storageDirectory := data.GetStorageDirectory()

	// Listen for messages to process
	for {

		message := <-documentMessageQueue
		documentFilename := storageDirectory + "/" + crypt.Sha256([]byte(message.ID)) + ".json"

		// Add a document to the index and write it to disk
		if message.Action == "add" {

			IndexDocument(message.ID, message.Document)

			documentFile, error := json.Marshal(types.JSONDocument{"id": message.ID, "document": string(message.Document[:])})

			if error == nil {
				ioutil.WriteFile(documentFilename, documentFile, os.FileMode(0600))
			}

		}

		// Remove a document from the index and disk
		if message.Action == "remove" {

			// Remove all documents
			if message.ID == "_all" {

				RemoveAllDocuments()

				// Remove a single document
			} else {

				RemoveDocument(message.ID, documentFilename)

			}

		}

	}

}
