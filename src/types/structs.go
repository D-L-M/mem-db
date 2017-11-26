package types

import "net/http"

// DocumentIndex structs need to store both the document JSON byte array and an
// inverted index of the keys where its entries in the inverted search index
// can be found
type DocumentIndex struct {
	Document     []byte
	InvertedKeys []string
}

// DocumentMessage structs inform a backround worker about changes to
// individual documents so that the disk store can be kept up-to-date
type DocumentMessage struct {
	ID       string
	Document []byte
	Action   string
}

// Route structs define executable HTTP routes
type Route struct {
	Path  string
	Route func(response http.ResponseWriter, body *[]byte, id string)
}
