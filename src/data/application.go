package data

import (
	"errors"
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

// GetBaseDirectory gets the directory in which to write any files
func GetBaseDirectory() (string, error) {

	user, error := user.Current()

	if error != nil {
		return "", errors.New("Could not determine base directory")
	}

	baseDirectory := user.HomeDir + "/.memdb"

	if _, error := os.Stat(baseDirectory); os.IsNotExist(error) {

		error := os.Mkdir(baseDirectory, os.FileMode(0700))

		if error != nil {
			return "", errors.New("Could not create base directory")
		}

	}

	return baseDirectory, nil

}

// GetStorageDirectory gets the directory in which to write any files
func GetStorageDirectory() (string, error) {

	baseDirctory, error := GetBaseDirectory()

	if error != nil {
		return "", error
	}

	storageDirectory := baseDirctory + "/documents"

	if _, error := os.Stat(storageDirectory); os.IsNotExist(error) {

		error := os.Mkdir(storageDirectory, os.FileMode(0700))

		if error != nil {
			return "", errors.New("Could not create storage directory")
		}

	}

	return storageDirectory, nil

}
