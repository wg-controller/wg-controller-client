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
	err := LoadState()
	if err != nil {
		log.Println("Failed to load state, creating new uuid")
		State.UUID = uuid.New().String()
	}

	loadFlags()

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
	wg_interface := ""
	flag.StringVar(&wg_interface, "wg-interface", "", "Wireguard interface name (optional)")
	server_host := ""
	flag.StringVar(&server_host, "server-host", "", "Hostname or IP of the server API")
	server_port := ""
	flag.StringVar(&server_port, "server-port", "", "Port of the server API (optional)")
	api_key := ""
	flag.StringVar(&api_key, "api-key", "", "API key for the server")

	flag.Parse()

	if wg_interface != "" {
		State.Flags.Wg_interface = wg_interface
	}
	if server_host != "" {
		State.Flags.Server_host = server_host
	}
	if server_port != "" {
		State.Flags.Server_port = server_port
	}
	if api_key != "" {
		State.Flags.Api_key = api_key
	}
}

func checkRequiredFlags() {
	if State.Flags.Wg_interface == "" {
		State.Flags.Wg_interface = "wg0"
	}
	if State.Flags.Server_host == "" {
		log.Fatal("--server-host is required")
	}
	if State.Flags.Server_port == "" {
		State.Flags.Server_port = "443"
	}
	if State.Flags.Api_key == "" {
		log.Fatal("--api-key is required")
	}
}
