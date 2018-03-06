package main

import (
	"github.com/D-L-M/jsonserver"
	"github.com/D-L-M/mem-db/src/auth"
	"github.com/D-L-M/mem-db/src/data"
	"github.com/D-L-M/mem-db/src/messaging"
	"github.com/D-L-M/mem-db/src/output"
	"github.com/D-L-M/mem-db/src/routing"
	"github.com/D-L-M/mem-db/src/store"
)

// Entry point
func main() {

	// Run goroutines that listen for messages on the various channels in use
	output.Log("Initialising channels")
	initialiseChannelListeners()

	// Reindex all documents previously flushed to disk
	output.Log("Restoring index from disk")
	store.IndexAllFromDisk()

	// Load authentication credentials into memory
	output.Log("Loading users")
	auth.Init()

	// Create a root user if one does not exist
	if auth.UserExists("root") == false {
		go messaging.AddUser("root", "password")
	}

	// Register HTTP routes
	output.Log("Registering routes")
	routing.RegisterRoutes()

	// Get start-up options
	port, hostname, peers, _, _ := data.GetOptions()

	// Set up a server
	output.Log("Starting server")
	jsonserver.Start(port)

	messaging.SetHostname(hostname)
	messaging.SetPeers(peers)

	// Block execution so the asynchronous code can handle requests
	output.Log("Listening for requests")
	select {}

}

// Launch goroutines for handling channel messages
func initialiseChannelListeners() {

	go messaging.ProcessUserMessages()
	go messaging.ProcessDocumentMessages()
	go messaging.ProcessPeerMessages()
	go messaging.ProcessPeerListMessages()

	data.ExecuteWhenActive(func() {
		messaging.ProcessPeerQueue()
	})

}
