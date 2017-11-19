package data

import (
	"../types"
	"log"
	"os"
	"os/user"
)

// Application state
var state = "initialising"

// Get the welcome message
func GetWelcomeMessage() types.JsonDocument {

	return types.JsonDocument{"engine": AppName, "version": AppVersion, "state": GetState()}

}

/// Set a new application state
func SetState(newState string) {

	state = newState

}

// Get the application state
func GetState() string {

	return state

}

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
