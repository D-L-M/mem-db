package messaging

import (
	"github.com/D-L-M/mem-db/src/auth"
	"github.com/D-L-M/mem-db/src/types"
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

		go ContactAllPeers(types.PeerMessage{Action: "reload_users", DocumentID: ""})

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
