package messaging

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"../crypt"
	"../data"
	"../store"
	"../types"
)

// ProcessDocumentMessages performs queued actions and flush document changes to disk
func ProcessDocumentMessages() {

	storageDirectory, err := data.GetStorageDirectory()

	if err != nil {
		log.Fatal(err)
	}

	// Listen for messages to process
	for {

		message := <-DocumentMessageQueue
		documentFilename := storageDirectory + "/" + crypt.Sha512([]byte(message.ID)) + ".json"

		// Add a document to the index and write it to disk
		if message.Action == "add" {

			store.IndexDocument(message.ID, message.Document)

			documentFile, err := json.Marshal(types.JSONDocument{"id": message.ID, "document": string(message.Document[:])})

			if err == nil {
				ioutil.WriteFile(documentFilename, documentFile, os.FileMode(0600))
			}

		}

		// Remove a document from the index and disk
		if message.Action == "remove" {

			// Remove all documents
			if message.ID == "_all" {

				store.RemoveAllDocuments()

				// Remove a single document
			} else {

				store.RemoveDocument(message.ID, documentFilename)

			}

		}

	}

}

// AddDocument adds a new user
func AddDocument(id string, body *[]byte) {

	DocumentMessageQueue <- types.DocumentMessage{ID: id, Document: *body, Action: "add"}

}

// RemoveDocument removes a document
func RemoveDocument(id string) {

	DocumentMessageQueue <- types.DocumentMessage{ID: id, Document: []byte{}, Action: "remove"}

}

// RemoveAllDocuments removes all documents
func RemoveAllDocuments() {

	DocumentMessageQueue <- types.DocumentMessage{ID: "_all", Document: []byte{}, Action: "remove"}

}
