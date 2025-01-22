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

func PopulateLinuxHostsFile() error {
	// Open hosts file
	file, err := os.OpenFile("/etc/hosts", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return errors.New("error opening hosts file: " + err.Error())
	}
	defer file.Close()

	// Split the file into lines
	scanner := bufio.NewScanner(file)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	// Append known hosts to the file
	PeersMU.Lock()
	for _, peer := range Peers {
		newEntry := peer.RemoteTunAddress + " " + peer.Hostname + " # wg-controller"
		lines = append(lines, newEntry)
	}
	PeersMU.Unlock()

	// Write the lines back to the file
	file.Truncate(0)
	file.Seek(0, 0)
	for _, line := range lines {
		_, err = file.WriteString(line + "\n")
		if err != nil {
			return err
		}
	}
	return err
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

func CleanupLinuxHostsFile() error {
	// Open hosts file
	file, err := os.OpenFile("/etc/hosts", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return errors.New("error opening hosts file: " + err.Error())
	}
	defer file.Close()

	// Split the file into lines
	scanner := bufio.NewScanner(file)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
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

	// Write the lines back to the file
	file.Truncate(0)
	file.Seek(0, 0)
	for _, line := range newLines {
		_, err = file.WriteString(line + "\n")
		if err != nil {
			return err
		}
	}
	log.Println("Cleaned up", count, "hosts file entries")
	return err
}
