package messaging

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/D-L-M/jsonserver"
	"github.com/D-L-M/mem-db/src/store"
	"github.com/D-L-M/mem-db/src/types"
)

// ProcessDocumentMessages performs queued actions and flush document changes to disk
func ProcessDocumentMessages() {

	// Listen for messages to process
	for {

		message := <-DocumentMessageQueue
		documentFilename, err := store.GetDocumentFilePath(message.ID)

		if err != nil {
			continue
		}

		// Add a document to the index and write it to disk
		if message.Action == "add" {

			store.IndexDocument(message.ID, message.Document, true)

			documentFile, err := json.Marshal(jsonserver.JSON{"id": message.ID, "document": string(message.Document[:])})

			if err == nil {
				ioutil.WriteFile(documentFilename, documentFile, os.FileMode(0600))
			}

			if message.PropagateToPeers {
				go ContactAllPeers(types.PeerMessage{Action: "reindex_document", DocumentID: message.ID})
			}

		}

		// Reindex a document from disk
		if message.Action == "index_from_disk" {

			store.IndexDocumentFromDisk(message.ID)

			if message.PropagateToPeers {
				go ContactAllPeers(types.PeerMessage{Action: "reindex_document", DocumentID: message.ID})
			}

		}

		// Remove a document from the index and disk
		if message.Action == "remove" {

			// Remove all documents
			if message.ID == "_all" {

				store.RemoveAllDocuments(true)

				if message.PropagateToPeers {
					go ContactAllPeers(types.PeerMessage{Action: "remove_all_documents", DocumentID: ""})
				}

				// Remove a single document
			} else {

				store.RemoveDocument(message.ID, documentFilename, true)

				if message.PropagateToPeers {
					go ContactAllPeers(types.PeerMessage{Action: "remove_document", DocumentID: message.ID})
				}

			}

		}

		// Remove a document from the index
		if message.Action == "remove_from_memory" {

			// Remove all documents from memory
			if message.ID == "_all" {

				store.RemoveAllDocuments(false)

				// Remove a single document from memory
			} else {

				store.RemoveDocument(message.ID, documentFilename, false)

			}

		}

	}

}

// AddDocument adds a new document
func AddDocument(id string, body *[]byte, propagateToPeers bool) {

	DocumentMessageQueue <- types.DocumentMessage{ID: id, Document: *body, Action: "add", PropagateToPeers: propagateToPeers}

}

// IndexDocumentFromDisk reindexes a document from disk
func IndexDocumentFromDisk(id string, propagateToPeers bool) {

	DocumentMessageQueue <- types.DocumentMessage{ID: id, Document: []byte{}, Action: "index_from_disk", PropagateToPeers: propagateToPeers}

}

// RemoveDocument removes a document
func RemoveDocument(id string, propagateToPeers bool) {

	DocumentMessageQueue <- types.DocumentMessage{ID: id, Document: []byte{}, Action: "remove", PropagateToPeers: propagateToPeers}

}

// RemoveDocumentFromMemory removes a document from memory but not from disk
func RemoveDocumentFromMemory(id string, propagateToPeers bool) {

	DocumentMessageQueue <- types.DocumentMessage{ID: id, Document: []byte{}, Action: "remove_from_memory", PropagateToPeers: propagateToPeers}

}

// RemoveAllDocuments removes all documents
func RemoveAllDocuments(propagateToPeers bool) {

	DocumentMessageQueue <- types.DocumentMessage{ID: "_all", Document: []byte{}, Action: "remove", PropagateToPeers: propagateToPeers}

}
