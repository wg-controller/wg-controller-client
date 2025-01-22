package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/wg-controller/wg-controller-client/types"
)

// Creates a peer on the server is it doesn't exist
var initAttempts = 0

func InitPeerServer() {
	log.Println("Initiating peer on server...")
	if initAttempts > 10 {
		log.Fatal("Failed to init peer on server")
	}

	// Attempt to fetch peer with our UUID
	path := "https://" + State.Flags.Server_host + ":" + State.Flags.Server_port + "/api/v1/peers/" + State.UUID
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Set headers
	req.Header.Set("Authorization", State.Flags.Api_key)

	// Send request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	// Check response code
	if resp.StatusCode == 404 {
		log.Println("Peer not found")
		CreatePeer()
		return
	} else if resp.StatusCode != 200 {
		log.Println("Failed to init peer", resp.Body, resp.StatusCode)
		initAttempts++
		time.Sleep(3 * time.Second)
		InitPeerServer()
		return
	}

	log.Println("Peer found on server")
}

func GetInitPeer() types.PeerInit {
	// Get init values
	path := "https://" + State.Flags.Server_host + ":" + State.Flags.Server_port + "/api/v1/peers/init"
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Set headers
	req.Header.Set("Authorization", State.Flags.Api_key)

	// Send request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	// Check response code
	if resp.StatusCode != 200 {
		log.Fatal("Failed to init peer", resp.StatusCode)
	}

	// Get response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Parse response body
	peerInit := types.PeerInit{}
	err = json.Unmarshal(body, &peerInit)
	if err != nil {
		log.Fatal(err)
	}

	return peerInit
}

func CreatePeer() {
	log.Println("Creating peer on server...")

	// Get init values
	peerInit := GetInitPeer()

	// Create peer object
	peer := types.Peer{
		UUID:             State.UUID,
		Hostname:         Hostname,
		Enabled:          true,
		PrivateKey:       peerInit.PrivateKey,
		PublicKey:        peerInit.PublicKey,
		PreSharedKey:     peerInit.PreSharedKey,
		KeepAliveSeconds: 15,
		RemoteTunAddress: peerInit.RemoteTunAddress,
		AllowedSubnets:   []string{peerInit.ServerCIDR},
		Attributes:       []string{"wg-controller-client"},
	}

	// Create peer
	path := "https://" + State.Flags.Server_host + ":" + State.Flags.Server_port + "/api/v1/peers/" + peerInit.UUID
	jsonBody, err := json.Marshal(peer)
	if err != nil {
		log.Fatal(err)
	}

	// Create request
	req, err := http.NewRequest(http.MethodPut, path, bytes.NewReader(jsonBody))
	if err != nil {
		log.Fatal(err)
	}

	// Set headers
	req.Header.Set("Authorization", State.Flags.Api_key)

	// Send request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	// Check response code
	if resp.StatusCode != 200 {
		log.Fatal("Failed to create peer. resp:", resp.StatusCode)
	} else {
		log.Println("Peer created")
	}
}

func GetConfig() {
	path := "https://" + State.Flags.Server_host + ":" + State.Flags.Server_port + "/api/v1/peers/" + State.UUID
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Set headers
	req.Header.Set("Authorization", State.Flags.Api_key)

	// Send request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	// Check response code
	if resp.StatusCode != 200 {
		log.Fatal("Failed to get peer config")
	}

	// Get response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Parse response body
	peer := types.Peer{}
	err = json.Unmarshal(body, &peer)
	if err != nil {
		log.Fatal(err)
	}

	// Update peer config
	PeerConfigMU.Lock()
	PeerConfig = peer
	PeerConfigMU.Unlock()

	ApplyWireguardConfig()
}

func GetServerInfo() {
	path := "https://" + State.Flags.Server_host + ":" + State.Flags.Server_port + "/api/v1/serverinfo"
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Set headers
	req.Header.Set("Authorization", State.Flags.Api_key)

	// Send request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	// Check response code
	if resp.StatusCode != 200 {
		log.Fatal("Failed to get server info")
	}

	// Get response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Parse response body
	server := types.ServerInfo{}
	err = json.Unmarshal(body, &server)
	if err != nil {
		log.Fatal(err)
	}

	// Update server info
	ServerInfoMU.Lock()
	ServerInfo = server
	ServerInfoMU.Unlock()
}

func GetPeers() {
	path := "https://" + State.Flags.Server_host + ":" + State.Flags.Server_port + "/api/v1/peers"
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Set headers
	req.Header.Set("Authorization", State.Flags.Api_key)

	// Send request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	// Check response code
	if resp.StatusCode != 200 {
		log.Fatal("Failed to get hosts")
	}

	// Get response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Parse response body
	peers := []types.Peer{}
	err = json.Unmarshal(body, &peers)
	if err != nil {
		log.Fatal(err)
	}

	// Update peers
	PeersMU.Lock()
	Peers = peers
	PeersMU.Unlock()
}
