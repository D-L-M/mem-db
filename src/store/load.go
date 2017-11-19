package store


import (
	"path/filepath"
	"log"
	"io/ioutil"
	"encoding/json"
	"../data"
	"../types"
)


// Reindex all documents previously flushed to disk
func IndexFromDisk() {

	storageDirectory := data.GetStorageDirectory()

	// Iterate through all flushed JSON files
	files, error := filepath.Glob(storageDirectory + "/*.json")

	if error != nil {
		log.Fatal("Cannot read from storage directory")
	}

	for _, filename := range files {

		// Read in and parse the JSON
		fileContents, error := ioutil.ReadFile(filename)

		if error == nil {

			var parsedDocument types.JsonDocument

			error := json.Unmarshal(fileContents, &parsedDocument)

			if error == nil {

				// Check for required fields and index the document
				if id, ok := parsedDocument["id"].(string); ok {

					if document, ok := parsedDocument["document"].(string); ok {

						IndexDocument(id, []byte(document))

					}

				}

			}

		}

	}

}