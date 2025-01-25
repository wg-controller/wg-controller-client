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

func GetSystemString() string {
	switch runtime.GOOS {
	case "linux":
		file, err := os.Open("/etc/os-release")
		if err != nil {
			return "unknown"
		}
		defer file.Close()

		buf := make([]byte, 1024)
		n, err := file.Read(buf)
		if err != nil {
			return "unknown"
		}

		if strings.Contains(string(buf[:n]), "ubuntu") {
			return "ubuntu"
		} else if strings.Contains(string(buf[:n]), "debian") {
			return "debian"
		} else if strings.Contains(string(buf[:n]), "centos") {
			return "centos"
		} else if strings.Contains(string(buf[:n]), "fedora") {
			return "fedora"
		} else if strings.Contains(string(buf[:n]), "arch") {
			return "arch"
		} else if strings.Contains(string(buf[:n]), "alpine") {
			return "alpine"
		} else if strings.Contains(string(buf[:n]), "rhel") {
			return "rhel"
		} else if strings.Contains(string(buf[:n]), "gentoo") {
			return "gentoo"
		} else if strings.Contains(string(buf[:n]), "kali") {
			return "kali"
		} else {
			return "unknown"
		}

	case "darwin":
		return "darwin"
	case "windows":
		return "windows"
	default:
		return "unknown"
	}
}
