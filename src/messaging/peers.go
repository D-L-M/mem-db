package messaging

import (
	"bytes"
	"encoding/json"
	"net/http"

	"../crypt"
	"../types"
)

var hostname = ""
var peers = map[string]bool{}

// SetHostname sets a new hostname for the server
func SetHostname(newHostname string) {

	hostname = newHostname

}

// AddPeer adds/enables a peer host
func AddPeer(peerHostname string) {

	peers[peerHostname] = true

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

// ContactPeer sends a HMAC signed message to a peer server
func ContactPeer(message types.PeerMessage) bool {

	peerStatus := peers[message.To]

	if peerStatus == false {
		return false
	}

	message.From = hostname
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

	}

}
