package main

import (
	"./auth"
	"./messaging"
	"./routing"
	"./server"
	"./store"
)

// Entry point
func main() {

	// Reindex all documents previously flushed to disk
	store.IndexFromDisk()

	// Listen for user messages
	go messaging.ProcessUserMessages()

	// Listen for document messages
	go messaging.ProcessDocumentMessages()

	// Load authentication credentials into memory
	auth.Init()

	// Create a root user if one does not exist
	if auth.UserExists("root") == false {
		go messaging.AddUser("root", "password")
	}

	// Register HTTP routes
	routing.RegisterRoutes()

	// Set up a server
	server.InitTCP()

	// Block execution so the asynchronous code can handle requests
	select {}

}
