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

// createDirectoryIfNotExists creates a directory if it does not exist
func createDirectoryIfNotExists(path string) error {

	if _, error := os.Stat(path); os.IsNotExist(error) {

		error := os.Mkdir(path, os.FileMode(0700))

		if error != nil {
			return error
		}

	}

	return nil

}

// GetBaseDirectory gets the directory in which to write any files
func GetBaseDirectory() (string, error) {

	user, error := user.Current()

	if error != nil {
		return "", errors.New("Could not determine base directory")
	}

	baseDirectory := user.HomeDir + "/.memdb"

	error = createDirectoryIfNotExists(baseDirectory)

	if error != nil {
		return "", error
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

	error = createDirectoryIfNotExists(storageDirectory)

	if error != nil {
		return "", error
	}

	return storageDirectory, nil

}
