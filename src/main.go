package main

import (
	"./auth"
	"./routing"
	"./server"
	"./store"
)

// Entry point
func main() {

	// Register HTTP routes
	routing.RegisterRoutes()

	// Set up a server
	server.InitTCP()

	// Reindex all documents previously flushed to disk
	store.IndexFromDisk()

	// Listen for user messages
	go auth.ProcessMessages()

	// Listen for document messages
	go store.ProcessMessages()

	// Block execution so the asynchronous code can handle requests
	select {}

}
