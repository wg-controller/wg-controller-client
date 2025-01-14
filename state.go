package main

import (
	"encoding/json"
	"os"
	"runtime"
)

func LoadState() error {
	// Check if file exists
	if _, err := os.Stat(stateFilePath()); os.IsNotExist(err) {
		return err
	} else {
		// Load state
		file, err := os.Open(stateFilePath())
		if err != nil {
			return err
		}

		// Decode JSON
		decoder := json.NewDecoder(file)
		err = decoder.Decode(&State)
		if err != nil {
			return err
		}

		return nil
	}
}

func SaveState() {
	// Create state file and directorys
	file, err := os.Create(stateFilePath())
	if err != nil {
		panic(err)
	}

	// Encode JSON
	encoder := json.NewEncoder(file)
	err = encoder.Encode(State)
	if err != nil {
		panic(err)
	}
}

func stateFilePath() string {
	switch runtime.GOOS {
	case "linux":
		return "/var/run/wireguard/wg-state.json"
	case "darwin":
		return "/var/run/wireguard/wg-state.json"
	case "windows":
		return `\\.\pipe\wireguard\wg-state.json`
	default:
		return "/var/run/wireguard/wg-state.json"
	}
}
