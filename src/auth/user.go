package auth

import (
	"../messaging"
)

// ProcessMessages performs queued actions
func ProcessMessages() {

	// Listen for messages to process
	for {

		message := <-messaging.UserMessageQueue

		if message.Action == "create" {

			AddUser(message.Username, message.Value)

		}

	}

}
