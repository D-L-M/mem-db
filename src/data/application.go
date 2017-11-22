package data

import (
	"log"
	"os"
	"os/user"

	"../types"
)

// Application state
var state = "initialising"

// GetWelcomeMessage returns the welcome message object
func GetWelcomeMessage() types.JSONDocument {

	return types.JSONDocument{"engine": AppName, "version": AppVersion, "state": GetState()}

}

// SetState sets a new application state
func SetState(newState string) {

	state = newState

}

// GetState gets the application state
func GetState() string {

	return state

}

// GetStorageDirectory gets the directory in which to flush documents
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
