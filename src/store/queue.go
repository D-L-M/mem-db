package store


import (
	"../types"
	"os"
	"os/user"
	"io/ioutil"
	"log"
	"../crypt"
)


// Perform queued actions and flush document changes to disk
func IndexOnDisk(documentMessage chan types.DocumentMessage) {

	// Get the user's home directory and set up a storage directory if one does
	// not already exist
	user, error := user.Current()

    if error != nil {

		log.Fatal(error)
		
	}

	storageDirectory := user.HomeDir + "/.memdb"
	
	if _, error := os.Stat(storageDirectory); os.IsNotExist(error) {

		os.Mkdir(storageDirectory, os.FileMode(0700))
		
	}
	
	// Listen for messages to process
	for {
		
		message          := <- documentMessage
		documentFilename := storageDirectory + "/" + crypt.Sha256([]byte(message.Id)) + ".json"

		// Add a document to the index and write it to disk
		if message.Action == "add" {

			IndexDocument(message.Id, message.Document)

			_ = ioutil.WriteFile(documentFilename, message.Document, os.FileMode(0600))
			
		}

		// Remove a document from the index and disk
		if message.Action == "remove" {

			RemoveDocument(message.Id)

			_ = os.Remove(documentFilename)
			
		}

	}

}