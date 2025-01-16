package main

import (
	"errors"
	"log"
	"os"
	"os/exec"
	"runtime"
)

func InstallService() error {
	switch runtime.GOOS {
	case "linux":
		return installLinuxService()
	case "darwin":
		return errors.New("unsupported OS")
	case "windows":
		return errors.New("unsupported OS")
	default:
		return errors.New("unsupported OS")
	}
}

func UninstallService() error {
	switch runtime.GOOS {
	case "linux":
		return uninstallLinuxService()
	case "darwin":
		return errors.New("unsupported OS")
	case "windows":
		return errors.New("unsupported OS")
	default:
		return errors.New("unsupported OS")
	}
}

func installLinuxService() error {
	// Check if systemd is installed
	_, err := exec.LookPath("systemctl")
	if err != nil {
		return err
	}

	// Check if service file exists
	_, err = os.Stat("/etc/systemd/system/wg-controller.service")
	if err == nil {
		return nil
	} else {
		log.Println("Creating service file")
	}

	// Create service file
	serviceFile, err := os.OpenFile("/etc/systemd/system/wg-controller.service",
		os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer serviceFile.Close()

	// Get current directory
	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	// Copy binary to /usr/local/bin
	err = exec.Command("cp", dir+"/wg-controller", "/usr/local/bin/wg-controller").Run()
	if err != nil {
		return err
	}

	// Create service file
	lines := []string{
		"[Unit]",
		"Description=Wireguard Controller",
		"After=network.target",
		"",
		"[Service]",
		"Type=simple",
		"ExecStart=/usr/local/bin/wg-controller",
		"Restart=always",
		"RestartSec=5",
		"",
		"[Install]",
		"WantedBy=multi-user.target",
	}

	// Write lines to service file
	for _, line := range lines {
		_, err := serviceFile.WriteString(line + "\n")
		if err != nil {
			return err
		}
	}

	// Reload systemd
	if err := exec.Command("systemctl", "daemon-reload").Run(); err != nil {
		return err
	}
	return nil
}

func uninstallLinuxService() error {
	// Stop service
	exec.Command("systemctl", "stop", "wg-controller").Run()

	// Disable service
	exec.Command("systemctl", "disable", "wg-controller").Run()

	// Remove service file
	os.Remove("/etc/systemd/system/wg-controller.service")

	// Remove binary
	os.Remove("/usr/local/bin/wg-controller")

	// Reload systemd
	exec.Command("systemctl", "daemon-reload").Run()

	return nil
}
