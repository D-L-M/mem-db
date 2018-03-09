package store

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"path/filepath"
	"strconv"

	"github.com/D-L-M/jsonserver"
	"github.com/D-L-M/mem-db/src/crypt"
	"github.com/D-L-M/mem-db/src/data"
	"github.com/D-L-M/mem-db/src/output"
)

// IndexFromFile reindexes a single document previously flushed to disk
func IndexFromFile(filename string) {

	// Read in and parse the JSON
	fileContents, err := ioutil.ReadFile(filename)

	if err == nil {

		var parsedDocument jsonserver.JSON

		err := json.Unmarshal(fileContents, &parsedDocument)

		if err == nil {

			// Check for required fields and index the document
			if id, ok := parsedDocument["id"].(string); ok {

				if document, ok := parsedDocument["document"].(string); ok {
					IndexDocument(id, []byte(document), false)
				}

			}

		}

	}

}

// GetDocumentFilePath gets a document's file path by its ID
func GetDocumentFilePath(documentID string) (string, error) {

	storageDirectory, err := data.GetStorageDirectory()

	if err != nil {
		return "", err
	}

	documentFilename := storageDirectory + "/" + crypt.Sha512([]byte(documentID)) + ".json"

	return documentFilename, nil

}

// IndexDocumentFromDisk reindexes a single document previously flushed to disk
// by its ID
func IndexDocumentFromDisk(documentID string) {

	documentFilename, err := GetDocumentFilePath(documentID)

	if err == nil {
		IndexFromFile(documentFilename)
	}

}

// IndexAllFromDisk reindexes all documents previously flushed to disk
func IndexAllFromDisk() {

	storageDirectory, err := data.GetStorageDirectory()

	if err != nil {
		log.Fatal(err)
	}

	// Iterate through all flushed JSON files
	files, err := filepath.Glob(storageDirectory + "/*.json")

	if err != nil {
		log.Fatal("Cannot read from storage directory")
	}

	data.SetState("recovering")

	for i, filename := range files {
		output.Log("Restoring index from disk: " + strconv.Itoa(i+1) + " / " + strconv.Itoa(len(files)))
		IndexFromFile(filename)
	}

	data.SetState("active")

}
