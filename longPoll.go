package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/wg-controller/wg-controller-client/types"
)

type LP_Message struct {
	Topic      string            `json:"topic"`
	Data       string            `json:"data"`
	Attributes map[string]string `json:"attributes"`
	Config     types.Peer        `json:"config,omitempty"`
	Peers      []types.Peer      `json:"peers,omitempty"`
}

const pollTimeout = 15 * time.Second
const pollPause = 2 * time.Second
const pollErrorPause = 10 * time.Second

func InitLongPoll() {
	log.Println("Initiating long poll...")
	poll()
}

func poll() {
	path := "https://" + State.Flags.Server_host + ":" + State.Flags.Server_port + "/api/v1/poll" + "?uuid=" + State.UUID
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		log.Println(err)
		time.Sleep(pollErrorPause)
		poll()
		return
	}

	// Set headers
	req.Header.Set("Authorization", State.Flags.Api_key)

	// Create client
	client := http.Client{
		Timeout: pollTimeout,
	}

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		time.Sleep(pollErrorPause)
		poll()
		return
	}

	// Check response code
	if resp.StatusCode == 200 {
		// Get response body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
			time.Sleep(pollErrorPause)
			poll()
			return
		}

		// Parse response
		var message LP_Message
		err = json.Unmarshal(body, &message)
		if err != nil {
			log.Println(err)
			time.Sleep(pollErrorPause)
			poll()
			return
		}

		handleIncomingPollMsg(message)
		time.Sleep(pollPause)
		poll()
		return
	} else if resp.StatusCode == 204 {
		// No new messages
		time.Sleep(pollPause)
		poll()
		return
	} else {
		log.Println("Unexpected poll resp code:", resp.StatusCode)
		time.Sleep(pollErrorPause)
		poll()
		return
	}
}

func handleIncomingPollMsg(message LP_Message) {
	switch message.Topic {
	case "peerConfig":
		log.Println("Received LP peerConfig message")
		// Apply config
		PeerConfigMU.Lock()
		PeerConfig = message.Config
		PeerConfigMU.Unlock()
		ApplyWireguardConfig()
	case "peers":
		log.Println("Received LP peers message")
		// Apply hosts
		PeersMU.Lock()
		Peers = message.Peers
		PeersMU.Unlock()
		CleanupHostsFile()
		PopulateHostsFile()
	default:
		log.Println("Unknown LP message type:", message.Topic)
	}
}
