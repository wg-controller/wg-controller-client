package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"time"

	"golang.zx2c4.com/wireguard/conn"
	"golang.zx2c4.com/wireguard/device"
	"golang.zx2c4.com/wireguard/ipc"
	"golang.zx2c4.com/wireguard/tun"
	"golang.zx2c4.com/wireguard/wgctrl"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

// Starts wireguard-go
func StartWireguard() {
	tun, err := tun.CreateTUN(State.Flags.Wg_interface, 0)
	if err == nil {
		realInterfaceName, err2 := tun.Name()
		if err2 == nil {
			State.Flags.Wg_interface = realInterfaceName
		}
	} else {
		log.Printf("Failed to create TUN device: %v", err)
		os.Exit(1)
	}

	logger := device.NewLogger(
		device.LogLevelVerbose,
		fmt.Sprintf("(%s) ", State.Flags.Wg_interface),
	)
	device := device.NewDevice(tun, conn.NewDefaultBind(), logger)
	err = device.Up()
	if err != nil {
		log.Printf("Failed to bring up device: %v", err)
		os.Exit(1)
	}
	logger.Verbosef("Device started")

	sockFile, err := ipc.UAPIOpen(State.Flags.Wg_interface)
	if err != nil {
		log.Printf("Failed to open uapi socket: %v", err)
		os.Exit(1)
	}

	uapi, err := ipc.UAPIListen(State.Flags.Wg_interface, sockFile)
	if err != nil {
		log.Printf("Failed to listen on uapi socket: %v", err)
		os.Exit(1)
	}

	errs := make(chan error)
	term := make(chan os.Signal, 1)

	go func() {
		for {
			conn, err := uapi.Accept()
			if err != nil {
				errs <- err
				return
			}
			go device.IpcHandle(conn)
		}
	}()
	log.Println("UAPI listener started")

	// wait for program to terminate
	signal.Notify(term, os.Interrupt)
	signal.Notify(term, os.Kill)

	select {
	case <-term:
	case <-errs:
	case <-device.Wait():
	}

	// clean up
	uapi.Close()
	device.Close()
	os.Exit(0)
}

func ApplyWireguardConfig() {
	PeerConfigMU.Lock()
	defer PeerConfigMU.Unlock()

	ServerInfoMU.Lock()
	defer ServerInfoMU.Unlock()

	client, err := wgctrl.New()
	if err != nil {
		log.Fatal("Unable to connect to wireguard-go")
	}

	// Parse PublicKey
	publicKey, err := wgtypes.ParseKey(ServerInfo.PublicKey)
	if err != nil {
		log.Fatal("Failed to parse public key")
	}

	// Parse PreSharedKey
	preSharedKey, err := wgtypes.ParseKey(PeerConfig.PreSharedKey)
	if err != nil {
		log.Fatal("Failed to parse pre-shared key")
	}

	// Convert KeepAliveSeconds to time.Duration
	keepAliveDuration := time.Duration(PeerConfig.KeepAliveSeconds) * time.Second

	// Parse allowed subnets
	allowedIPs := []net.IPNet{}
	for _, subnet := range PeerConfig.RemoteSubnets {
		_, ipNet, err := net.ParseCIDR(subnet)
		if err != nil {
			break
		}
		allowedIPs = append(allowedIPs, *ipNet)
	}
	// Append peer's own subnet
	_, ipNet, err := net.ParseCIDR(PeerConfig.RemoteTunAddress + "/32")
	if err != nil {
		log.Println("Error parsing peer's own subnet:", err)
	} else {
		allowedIPs = append(allowedIPs, *ipNet)
	}

	// Parse endpoint
	host, port, err := net.SplitHostPort(ServerInfo.PublicEndpoint)
	if err != nil {
		log.Fatal("Failed to parse endpoint")
	}
	// Resolve host
	ips, err := net.LookupIP(host)
	if err != nil || len(ips) == 0 {
		log.Fatal("Failed to resolve host")
	}
	portNum, err := strconv.Atoi(port)
	if err != nil {
		log.Fatal("Failed to parse port")
	}
	endpoint := &net.UDPAddr{
		IP:   ips[0],
		Port: portNum,
	}

	wgPeer := wgtypes.PeerConfig{
		PublicKey:                   publicKey,
		PresharedKey:                &preSharedKey,
		Endpoint:                    endpoint,
		PersistentKeepaliveInterval: &keepAliveDuration,
		AllowedIPs:                  allowedIPs,
		ReplaceAllowedIPs:           true,
	}

	// Parse private key
	privateKey, err := wgtypes.ParseKey(PeerConfig.PrivateKey)
	if err != nil {
		log.Fatal("Failed to parse private key")
	}

	err = client.ConfigureDevice(State.Flags.Wg_interface, wgtypes.Config{
		ReplacePeers: true,
		Peers:        []wgtypes.PeerConfig{wgPeer},
		PrivateKey:   &privateKey,
	})
	if err != nil {
		log.Fatal("Failed to apply configuration")
	}
}
