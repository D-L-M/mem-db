package store

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"path/filepath"

	"../data"
	"../types"
)

// IndexFromDisk reindexes all documents previously flushed to disk
func IndexFromDisk() {

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

	for _, filename := range files {

		// Read in and parse the JSON
		fileContents, err := ioutil.ReadFile(filename)

		if err == nil {

			var parsedDocument types.JSONDocument

			err := json.Unmarshal(fileContents, &parsedDocument)

			if err == nil {

				// Check for required fields and index the document
				if id, ok := parsedDocument["id"].(string); ok {

					if document, ok := parsedDocument["document"].(string); ok {
						IndexDocument(id, []byte(document))
					}

				}

			}

		}

	}

	data.SetState("active")

}
