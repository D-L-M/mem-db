package main

import (
	"./auth"
	"./data"
	"./messaging"
	"./routing"
	"./server"
	"./store"
)

// Entry point
func main() {

	// Run goroutines that listen for messages on the various channels in use
	initialiseChannelListeners()

	// Reindex all documents previously flushed to disk
	store.IndexAllFromDisk()

	// Load authentication credentials into memory
	auth.Init()

	// Create a root user if one does not exist
	if auth.UserExists("root") == false {
		go messaging.AddUser("root", "password")
	}

	// Register HTTP routes
	routing.RegisterRoutes()

	// Get start-up options
	port, hostname, peers, _ := data.GetOptions()

	messaging.SetHostname(hostname)
	messaging.SetPeers(peers)

	// Set up a server
	server.InitTCP(port)

	// Block execution so the asynchronous code can handle requests
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
