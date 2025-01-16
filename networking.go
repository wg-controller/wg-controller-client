package main

import (
	"errors"
	"log"
	"net"
	"runtime"

	"github.com/vishvananda/netlink"
)

func ApplyNetworkConfiguration() error {
	log.Println("Applying network configuration...")
	// Apply IP address to WG interface
	err := OverwriteInterfaceIP(State.Flags.Wg_interface, PeerConfig.RemoteTunAddress+ServerInfo.Netmask)
	if err != nil {
		return err
	}

	// Set interface up
	err = SetInterfaceState(State.Flags.Wg_interface, true)
	if err != nil {
		return err
	}

	// Cleanup old routes
	err = CleanupRoutes()
	if err != nil {
		return err
	}

	// Add new routes
	err = AddRoutes(PeerConfig.AllowedSubnets, PeerConfig.RemoteTunAddress)
	if err != nil {
		return err
	}

	return nil
}

func OverwriteInterfaceIP(interfaceName string, ip string) error {
	switch runtime.GOOS {
	case "linux", "darwin":
		// Get interface
		link, err := netlink.LinkByName(interfaceName)
		if err != nil {
			return errors.New("error getting interface: " + err.Error())
		}

		// Remove existing IP addresses
		addrs, err := netlink.AddrList(link, 2)
		if err != nil {
			return errors.New("error getting IP: " + err.Error())
		}
		for _, a := range addrs {
			err = netlink.AddrDel(link, &a)
			if err != nil {
				return errors.New("error removing IP: " + err.Error())
			}
		}

		// Parse new IP address
		addr, err := netlink.ParseAddr(ip)
		if err != nil {
			return errors.New("error parsing IP: " + err.Error())
		}

		// Add new IP address
		err = netlink.AddrAdd(link, addr)
		if err != nil {
			return errors.New("error setting IP: " + err.Error())
		}
		return nil
	default:
		return errors.New("unsupported OS")
	}
}

func SetInterfaceState(interfaceName string, up bool) error {
	switch runtime.GOOS {
	case "linux", "darwin":
		link, err := netlink.LinkByName(interfaceName)
		if err != nil {
			return errors.New("error getting interface: " + err.Error())
		}
		if up {
			err = netlink.LinkSetUp(link)
		} else {
			err = netlink.LinkSetDown(link)
		}
		if err != nil {
			return errors.New("error setting interface state: " + err.Error())
		}
		return nil
	default:
		return errors.New("unsupported OS")
	}
}

func CleanupRoutes() error {
	cleanCount := 0
	switch runtime.GOOS {
	case "linux", "darwin":
		routes, _ := netlink.RouteList(nil, 2)
		for _, route := range routes {
			if route.Protocol == 171 {
				err := netlink.RouteDel(&route)
				if err == nil {
					cleanCount++
				}
			}
		}
		log.Println("Cleaned up", cleanCount, "routes")
		return nil
	default:
		return errors.New("unsupported OS")
	}
}

func AddRoutes(networks []string, gateway string) error {
	switch runtime.GOOS {
	case "linux", "darwin":
		for _, network := range networks {
			err := AddRoute(network, gateway)
			if err != nil {
				if err.Error() == "file exists" {
					log.Println("Route already exists for", network)
				} else {
					log.Println("Failed to add route:", err)
				}
			}
		}
	default:
		return errors.New("unsupported OS")
	}
	return nil
}

func AddRoute(destination string, gateway string) error {
	switch runtime.GOOS {
	case "linux", "darwin":
		_, dst, err := net.ParseCIDR(destination)
		if err != nil {
			return err
		}
		gw := net.ParseIP(gateway)
		if gw == nil {
			return errors.New("invalid gateway IP")
		}
		route := netlink.Route{
			Dst:      dst,
			Gw:       gw,
			Protocol: 171, // Identifies the route as a WireGuard route
		}
		return netlink.RouteAdd(&route)
	default:
		return errors.New("unsupported OS")
	}
}
