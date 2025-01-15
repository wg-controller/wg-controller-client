package main

import "runtime"

func AppendNameserver(nameServer string) error {
	switch runtime.GOOS {
	case "linux":
		return appendLinuxNameserver(nameServer)
	case "darwin":
		return appendDarwinNameserver(nameServer)
	default:
		return appendLinuxNameserver(nameServer)
	}
}

func appendLinuxNameserver(nameServer string) error {
	return nil
}

func appendDarwinNameserver(nameServer string) error {
	return nil
}

func CleanupNameservers() error {
	switch runtime.GOOS {
	case "linux":
		return cleanupLinuxNameservers()
	case "darwin":
		return cleanupDarwinNameservers()
	default:
		return cleanupLinuxNameservers()
	}
}

func cleanupLinuxNameservers() error {
	return nil
}

func cleanupDarwinNameservers() error {
	return nil
}
