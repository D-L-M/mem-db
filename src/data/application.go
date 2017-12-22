package data

import (
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

	if _, err := os.Stat(path); os.IsNotExist(err) {

		err := os.Mkdir(path, os.FileMode(0700))

		if err != nil {
			return err
		}

	}

	return nil

}

// GetBaseDirectory gets the directory in which to write any files
func GetBaseDirectory() (string, error) {

	user, err := user.Current()

	if err != nil {
		return "", err
	}

	baseDirectory := user.HomeDir + "/.memdb"

	err = createDirectoryIfNotExists(baseDirectory)

	if err != nil {
		return "", err
	}

	return baseDirectory, nil

}

// GetStorageDirectory gets the directory in which to write any files
func GetStorageDirectory() (string, error) {

	baseDirectory, err := GetBaseDirectory()

	if err != nil {
		return "", err
	}

	storageDirectory := baseDirectory + "/documents"

	err = createDirectoryIfNotExists(storageDirectory)

	if err != nil {
		return "", err
	}

	return storageDirectory, nil

}
