package messaging

import (
	"../auth"
	"../types"
)

// ProcessUserMessages performs queued actions
func ProcessUserMessages() {

	// Listen for messages to process
	for {

		message := <-UserMessageQueue

		if message.Action == "create" {

			auth.AddUser(message.Username, message.Value)

		}

		if message.Action == "delete" {

			auth.DeleteUser(message.Username)

		}

	}

}

// AddUser adds a new user
func AddUser(username string, password string) {

	UserMessageQueue <- types.UserMessage{Username: username, Value: password, Action: "create"}

}

// DeleteUser removes a user
func DeleteUser(username string) {

	UserMessageQueue <- types.UserMessage{Username: username, Value: "", Action: "delete"}

}
