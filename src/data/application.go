package data

import (
	"os"
	"os/user"

	"../types"
)

// Application state
var state = "initialising"

// Closures to execute when the application becomes active
var executeWhenActive = []func(){}

// GetWelcomeMessage returns the welcome message object
func GetWelcomeMessage() types.JSONDocument {

	return types.JSONDocument{"engine": AppName, "version": AppVersion, "state": GetState()}

}

// ExecuteWhenActive stores a closure to execute when the application becomes
// active
func ExecuteWhenActive(callback func()) {

	executeWhenActive = append(executeWhenActive, callback)

}

// SetState sets a new application state
func SetState(newState string) {

	state = newState

	if state == "active" {

		for _, callback := range executeWhenActive {
			callback()
		}

	}

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
