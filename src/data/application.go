package data

import (
	"errors"
	"os"
	"sync"

	"../types"
)

// Application state
var state = "initialising"

// Closures to execute when the application becomes active
var executeWhenActive = []func(){}

// stateLock allows locking of the state variable during reads/writes
var stateLock = sync.RWMutex{}

// executeWhenActiveLock allows locking of the executeWhenActive slice during reads/writes
var executeWhenActiveLock = sync.RWMutex{}

// GetWelcomeMessage returns the welcome message object
func GetWelcomeMessage() types.JSONDocument {

	return types.JSONDocument{"engine": AppName, "version": AppVersion, "state": GetState()}

}

// ExecuteWhenActive stores a closure to execute when the application becomes
// active
func ExecuteWhenActive(callback func()) {

	executeWhenActiveLock.Lock()
	executeWhenActive = append(executeWhenActive, callback)
	executeWhenActiveLock.Unlock()

}

// SetState sets a new application state
func SetState(newState string) {

	stateLock.Lock()
	state = newState
	stateLock.Unlock()

	if newState == "active" {

		executeWhenActiveLock.RLock()

		for _, callback := range executeWhenActive {
			callback()
		}

		executeWhenActiveLock.RUnlock()

	}

}

// GetState gets the application state
func GetState() string {

	stateLock.RLock()
	defer stateLock.RUnlock()

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

	_, _, _, baseDirectory := GetOptions()

	if baseDirectory == "" {
		return "", errors.New("Base directory is not available")
	}

	subBaseDirectory := baseDirectory + "/.memdb"

	err := createDirectoryIfNotExists(subBaseDirectory)

	if err != nil {
		return "", err
	}

	return subBaseDirectory, nil

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
