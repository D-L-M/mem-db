package messaging

import (
	"bytes"
	"encoding/json"
	"net/http"
	"sync"

	"github.com/D-L-M/mem-db/src/auth"
	"github.com/D-L-M/mem-db/src/crypt"
	"github.com/D-L-M/mem-db/src/data"
	"github.com/D-L-M/mem-db/src/output"
	"github.com/D-L-M/mem-db/src/types"
)

// Hostname of the running application
var hostname = ""

// Hostnames of running peer servers
var peers = map[string]bool{}

// Messages queued for processing during application start-up
var queuedMessages = []types.PeerMessage{}

// peersLock allows locking of the peers map during reads/writes
var peersLock = sync.RWMutex{}

// SetHostname sets a new hostname for the server
func SetHostname(newHostname string) {

	hostname = newHostname

}

// AddPeer adds/enables a peer host
func AddPeer(peerHostname string) {

	if peerHostname != "" && peerHostname != hostname {
		PeerListQueue <- types.PeerList{Hostname: peerHostname, Action: "add"}
	}

}

// RemovePeer removes/disables a peer host
func RemovePeer(peerHostname string) {

	PeerListQueue <- types.PeerList{Hostname: peerHostname, Action: "remove"}

}

// SetPeers overwrites the peer list with a new one
func SetPeers(peerHostnames []string) {

	peersLock.Lock()
	peers = map[string]bool{}
	peersLock.Unlock()

	for _, peerHostname := range peerHostnames {
		AddPeer(peerHostname)
	}

}

// GetPeers gets a list of known active peers
func GetPeers() []string {

	peersLock.RLock()

	knownPeers := []string{}

	for peerHostname, activeState := range peers {

		if activeState {
			knownPeers = append(knownPeers, peerHostname)
		}

	}

	peersLock.RUnlock()

	return knownPeers

}

// ContactAllPeers sends a HMAC signed message to all peer servers
func ContactAllPeers(message types.PeerMessage) {

	activePeers := GetPeers()

	for _, peer := range activePeers {

		peerMessage := message
		peerMessage.To = peer

		go ContactPeer(peerMessage)

	}

}

// ContactPeer sends a HMAC signed message to a peer server
func ContactPeer(message types.PeerMessage) bool {

	peersLock.RLock()
	peerStatus := peers[message.To]
	peersLock.RUnlock()

	if peerStatus == false {
		return false
	}

	message.From = hostname
	message.KnownPeers = GetPeers()
	payload, err := json.Marshal(message)

	if err != nil {
		return false
	}

	nonce, err := crypt.GenerateUUID()

	if err != nil {
		return false
	}

	signature := crypt.Sha512HMAC([]byte(string(payload[:]) + nonce))

	go sendPeerMessage(message.To, payload, signature, nonce)

	return true

}

// sendPeerMessage contacts a peer with a signed message
func sendPeerMessage(peerHostname string, message []byte, signature string, nonce string) {

	request, err := http.NewRequest("POST", peerHostname+"/_peer-message", bytes.NewBuffer(message))

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("x-hmac-auth", signature)
	request.Header.Set("x-hmac-nonce", nonce)

	client := &http.Client{}
	response, err := client.Do(request)

	if err == nil {

		defer response.Body.Close()

		if response.StatusCode != 202 {
			RemovePeer(peerHostname)
		}

	} else {
		RemovePeer(peerHostname)
	}

}

// ProcessPeerQueue redrives the processing queue to the channel
func ProcessPeerQueue() {

	for _, message := range queuedMessages {
		PeerMessageQueue <- message
	}

	queuedMessages = []types.PeerMessage{}

}

// ProcessPeerListMessages handles addition and removal instructions for the
// peer list
func ProcessPeerListMessages() {

	// Listen for messages to process
	for {

		message := <-PeerListQueue

		if message.Action == "add" {

			peersLock.Lock()

			peerAlreadyKnown := peers[message.Hostname]
			peers[message.Hostname] = true

			output.Log(message.Hostname + " added as a peer")
			peersLock.Unlock()

			// Forward peers list to peers if peer is not already known
			if peerAlreadyKnown == false {
				ContactPeer(types.PeerMessage{To: message.Hostname, Action: "update_peers", DocumentID: ""})
			}

		}

		if message.Action == "remove" {

			peersLock.Lock()

			peers[message.Hostname] = false

			output.Log(message.Hostname + " removed as a peer")
			peersLock.Unlock()
		}

	}

}

// ProcessPeerMessages performs queued peer instructions
func ProcessPeerMessages() {

	// Listen for messages to process
	for {

		message := <-PeerMessageQueue

		// If the application is not active, queue any peer messages for now
		if data.GetState() != "active" {

			queuedMessages = append(queuedMessages, message)

		} else {

			// Update the peers list
			if message.Action == "update_peers" {

				AddPeer(message.From)

				for _, peerHostname := range message.KnownPeers {
					AddPeer(peerHostname)
				}

			}

			// Reindex a document from disk
			if message.Action == "reindex_document" {
				output.Log(message.From + " instructed to reindex document '" + message.DocumentID + "' from disk")
				IndexDocumentFromDisk(message.DocumentID, false)
			}

			// Remove a document from memory
			if message.Action == "remove_document" {
				output.Log(message.From + " instructed to remove document '" + message.DocumentID + "' from memory")
				RemoveDocumentFromMemory(message.DocumentID, false)
			}

			// Remove all documents from memory
			if message.Action == "remove_all_documents" {
				output.Log(message.From + " instructed to remove all documents from memory")
				RemoveDocumentFromMemory("_all", false)
			}

			// Reload the user's list
			if message.Action == "reload_users" {
				output.Log(message.From + " instructed to reload users list")
				auth.Init()
			}

		}

	}

}
