package main

import (
	"errors"
	"os"
	"runtime"
	"strings"
)

func WireguardSockPath() string {
	switch runtime.GOOS {
	case "linux":
		return "/var/run/wireguard/wg0.sock"
	case "darwin":
		return "/var/run/wireguard/wg0.sock"
	case "windows":
		return `\\.\pipe\wireguard\wg0.sock`
	default:
		return "/var/run/wireguard/wg0.sock"
	}
}

func GetHostname() (string, error) {
	name, err := os.Hostname()
	if err != nil {
		return "", err
	}

	// Check if hostname is empty
	if name == "" {
		return "", errors.New("failed to get hostname")
	}

	// Remove domain from hostname
	if strings.Contains(name, ".") {
		name = strings.Split(name, ".")[0]
	}

	return name, nil
}
