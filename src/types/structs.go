package types

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
	ID               string
	Document         []byte
	Action           string
	PropagateToPeers bool
}

// UserMessage structs inform a backround worker about changes to
// user accounts
type UserMessage struct {
	Username string
	Value    string
	Action   string
}

// PeerMessage structs contain instructional messages for peer servers
type PeerMessage struct {
	From       string
	To         string
	KnownPeers []string
	Action     string
	DocumentID string
}

// PeerList structs define additions and removals from the peer list
type PeerList struct {
	Hostname string
	Action   string
}
