package store

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"../crypt"
	"../data"
	"../messaging"
	"../types"
)

// ProcessMessages performs queued actions and flush document changes to disk
func ProcessMessages() {

	storageDirectory, err := data.GetStorageDirectory()

	if err != nil {
		log.Fatal(err)
	}

	// Listen for messages to process
	for {

		message := <-messaging.DocumentMessageQueue
		documentFilename := storageDirectory + "/" + crypt.Sha512([]byte(message.ID)) + ".json"

		// Add a document to the index and write it to disk
		if message.Action == "add" {

			IndexDocument(message.ID, message.Document)

			documentFile, err := json.Marshal(types.JSONDocument{"id": message.ID, "document": string(message.Document[:])})

			if err == nil {
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
