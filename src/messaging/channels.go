package messaging

import (
	"../types"
)

// UserMessageQueue is a channel for user change messages
var UserMessageQueue = make(chan types.UserMessage)
