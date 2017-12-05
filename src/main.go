package main

import (
	"./auth"
	"./routing"
	"./server"
	"./store"
)

// Entry point
func main() {

	// Reindex all documents previously flushed to disk
	store.IndexFromDisk()

	// Listen for user messages
	go auth.ProcessMessages()

	// Listen for document messages
	go store.ProcessMessages()

	// Register HTTP routes
	routing.RegisterRoutes()

	// Set up a server
	server.InitTCP()

	// Block execution so the asynchronous code can handle requests
	select {}

}
