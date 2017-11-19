package data

import (
	"../types"
	"log"
	"os"
	"os/user"
)

// Welcome message
var WelcomeMessage = types.JsonDocument{"engine": AppName, "version": AppVersion}

// Get the directory in which to flush documents
func GetStorageDirectory() string {

	user, error := user.Current()

	if error != nil {
		log.Fatal("Could not determine storage directory")
	}

	storageDirectory := user.HomeDir + "/.memdb"

	if _, error := os.Stat(storageDirectory); os.IsNotExist(error) {

		error := os.Mkdir(storageDirectory, os.FileMode(0700))

		if error != nil {
			log.Fatal("Could not create storage directory")
		}

	}

	return storageDirectory

}
