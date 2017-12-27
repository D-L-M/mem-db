package messaging

import (
	"bytes"
	"encoding/json"
	"net/http"

	"../crypt"
	"../data"
	"../types"
)

// Hostname of the running application
var hostname = ""

// Hostnames of running peer servers
var peers = map[string]bool{}

// Messages queued for processing during application start-up
var queuedMessages = []types.PeerMessage{}

// SetHostname sets a new hostname for the server
func SetHostname(newHostname string) {

	hostname = newHostname

}

// AddPeer adds/enables a peer host
func AddPeer(peerHostname string) {

	if peerHostname != "" && peerHostname != hostname {

		peerAlreadyKnown := peers[peerHostname]
		peers[peerHostname] = true

		// Forward peers list to peers if peer is not already known
		if peerAlreadyKnown == false {
			ContactPeer(types.PeerMessage{To: peerHostname, Action: "update_peers", DocumentID: ""})
		}

	}

}

// RemovePeer removes/disables a peer host
func RemovePeer(peerHostname string) {

	peers[peerHostname] = false

}

// SetPeers overwrites the peer list with a new one
func SetPeers(peerHostnames []string) {

	peers = map[string]bool{}

	for _, peerHostname := range peerHostnames {
		AddPeer(peerHostname)
	}

}

// GetPeers gets a list of known peers
func GetPeers() []string {

	knownPeers := []string{}

	for peerHostname, activeState := range peers {

		if activeState {
			knownPeers = append(knownPeers, peerHostname)
		}

	}

	return knownPeers

}

// ContactPeer sends a HMAC signed message to a peer server
func ContactPeer(message types.PeerMessage) bool {

	peerStatus := peers[message.To]

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

		}

	}

}
