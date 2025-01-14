package main

import (
	"flag"
	"log"
	"sync"

	"github.com/google/uuid"
	"github.com/wg-controller/wg-controller-client/types"
)

// Version
var IMAGE_TAG string = "0.0.0"

// Flags
type SysFlags struct {
	Wg_interface string `json:"wg_interface"`
	Server_host  string `json:"server_host"`
	Server_port  string `json:"server_port"`
	Api_key      string `json:"api_key"`
}

type StateStruct struct {
	UUID  string   `json:"uuid"`
	Flags SysFlags `json:"flags"`
}

var Hostname string
var PeerConfig types.Peer
var PeerConfigMU sync.Mutex
var ServerInfo types.ServerInfo
var ServerInfoMU sync.Mutex
var State StateStruct

func main() {
	loadFlags()

	err := LoadState()
	if err != nil {
		log.Println("Failed to load state, creating new uuid")
		State.UUID = uuid.New().String()
	}

	SaveState()

	checkRequiredFlags()

	go StartWireguard()

	Hostname, err = GetHostname()
	if err != nil {
		log.Fatal(err)
	}

	InitPeerServer()

	GetServerInfo()

	GetConfig()

	select {}
}

func loadFlags() {
	log.Println("Wireguard client ver", IMAGE_TAG)
	flag.StringVar(&State.Flags.Wg_interface, "wg-interface", "wg0", "Wireguard interface name (optional)")
	flag.StringVar(&State.Flags.Server_host, "server-host", "", "Hostname or IP of the server API")
	flag.StringVar(&State.Flags.Server_port, "server-port", "443", "Port of the server API (optional)")
	flag.StringVar(&State.Flags.Api_key, "api-key", "", "API key for the server")
	flag.Parse()
}

func checkRequiredFlags() {
	if State.Flags.Server_host == "" {
		log.Fatal("Server host is required")
	}

	if State.Flags.Api_key == "" {
		log.Fatal("API key is required")
	}
}
