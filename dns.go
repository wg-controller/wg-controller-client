package main

import (
	"errors"
	"os/exec"
	"runtime"
	"strings"
)

func AppendNameserver(nameServer string) error {
	switch runtime.GOOS {
	case "linux":
		return appendLinuxNameserver(nameServer)
	case "darwin":
		return appendDarwinNameserver(nameServer)
	case "windows":
		return errors.New("unsupported OS")
	default:
		return errors.New("unsupported OS")
	}
}

func appendLinuxNameserver(nameServer string) error {
	if ResolvconfUtilityInstalled() {
		in := "nameserver " + nameServer
		cmd := exec.Command("resolvconf", "-a", State.Flags.Wg_interface, "-m", "0", "-x")
		cmd.Stdin = strings.NewReader(in)
		return cmd.Run()
	} else {
		return errors.New("unsupported OS")
	}
}

func appendDarwinNameserver(nameServer string) error {
	return errors.New("unsupported OS")
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
	if ResolvconfUtilityInstalled() {
		cmd := exec.Command("resolvconf", "-d", State.Flags.Wg_interface)
		return cmd.Run()
	} else {
		return errors.New("unsupported OS")
	}
}

func cleanupDarwinNameservers() error {
	return errors.New("unsupported OS")
}

func ResolvconfUtilityInstalled() bool {
	_, err := exec.LookPath("resolvconf")
	return err == nil
}
