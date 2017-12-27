package messaging

import (
	"../types"
)

// UserMessageQueue is a channel for user change messages
var UserMessageQueue = make(chan types.UserMessage)

// DocumentMessageQueue is a channel for document change messages
var DocumentMessageQueue = make(chan types.DocumentMessage)

// PeerMessageQueue is a channel for instructional peer server messages
var PeerMessageQueue = make(chan types.PeerMessage)
