package main

import (
	"bufio"
	"errors"
	"log"
	"os"
	"runtime"
	"strings"
)

func PopulateHostsFile() error {
	switch runtime.GOOS {
	case "linux":
		return PopulateLinuxHostsFile()
	case "darwin":
		return errors.New("unsupported OS")
	case "windows":
		return errors.New("unsupported OS")
	default:
		return errors.New("unsupported OS")
	}
}

func CleanupHostsFile() error {
	switch runtime.GOOS {
	case "linux":
		return CleanupLinuxHostsFile()
	case "darwin":
		return errors.New("unsupported OS")
	case "windows":
		return errors.New("unsupported OS")
	default:
		return errors.New("unsupported OS")
	}
}

func PopulateLinuxHostsFile() error {
	// Open the hosts file for reading and writing
	file, err := os.OpenFile("/etc/hosts", os.O_RDWR, 0644)
	if err != nil {
		return errors.New("error opening hosts file: " + err.Error())
	}
	defer file.Close()

	// Read the entire file
	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return errors.New("error reading hosts file: " + err.Error())
	}

	// Append known hosts
	PeersMU.Lock()
	for _, peer := range Peers {
		// Skip self
		if peer.Hostname == Hostname {
			continue
		}

		// Append entry
		newEntry := peer.RemoteTunAddress + " " + peer.Hostname + " # wg-controller"
		lines = append(lines, newEntry)
	}
	PeersMU.Unlock()

	// Append server host
	newEntry := ServerInfo.ServerInternalIP + " " + ServerInfo.ServerInternalName + " # wg-controller"
	lines = append(lines, newEntry)

	// Rewrite the file
	if err := file.Truncate(0); err != nil {
		return errors.New("error truncating hosts file: " + err.Error())
	}
	if _, err := file.Seek(0, 0); err != nil {
		return errors.New("error seeking hosts file: " + err.Error())
	}
	writer := bufio.NewWriter(file)
	for _, line := range lines {
		if _, err := writer.WriteString(line + "\n"); err != nil {
			return errors.New("error writing to hosts file: " + err.Error())
		}
	}
	return writer.Flush()
}

func CleanupLinuxHostsFile() error {
	// Open the hosts file for reading and writing
	file, err := os.OpenFile("/etc/hosts", os.O_RDWR, 0644)
	if err != nil {
		return errors.New("error opening hosts file: " + err.Error())
	}
	defer file.Close()

	// Read the entire file
	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return errors.New("error reading hosts file: " + err.Error())
	}

	// Remove wg-controller entries
	var newLines []string
	count := 0
	for _, line := range lines {
		if !strings.Contains(line, "# wg-controller") {
			newLines = append(newLines, line)
		} else {
			count++
		}
	}

	// Rewrite the file
	if err := file.Truncate(0); err != nil {
		return errors.New("error truncating hosts file: " + err.Error())
	}
	if _, err := file.Seek(0, 0); err != nil {
		return errors.New("error seeking hosts file: " + err.Error())
	}
	writer := bufio.NewWriter(file)
	for _, line := range newLines {
		if _, err := writer.WriteString(line + "\n"); err != nil {
			return errors.New("error writing to hosts file: " + err.Error())
		}
	}
	if err := writer.Flush(); err != nil {
		return errors.New("error flushing hosts file: " + err.Error())
	}

	log.Println("Cleaned up", count, "hosts file entries")
	return nil
}
