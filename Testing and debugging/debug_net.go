package main

import (
	"fmt"
	"net"
)

func main() {
	fmt.Println("--- SecureWrap Network Diagnostic ---")

	mac, err := getMACAddress()
	if err != nil {
		fmt.Printf("Error getting MAC: %v\n", err)
	} else {
		fmt.Printf("Detected MAC: %s\n", mac)
	}

	ip, err := getIPAddress()
	if err != nil {
		fmt.Printf("Error getting IP: %v\n", err)
	} else {
		fmt.Printf("Detected IP:  %s\n", ip)
	}

	fmt.Println("-------------------------------------")
	fmt.Println("Copy these values into your Creator App.")
}

// --- Exact copies of your project's helper functions ---

func getMACAddress() (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range interfaces {
		// We look for the first interface that is UP and NOT a loopback
		if iface.Flags&net.FlagUp != 0 && iface.Flags&net.FlagLoopback == 0 {
			if iface.HardwareAddr.String() != "" {
				return iface.HardwareAddr.String(), nil
			}
		}
	}
	return "", fmt.Errorf("no active MAC address found")
}

func getIPAddress() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}
	return "", fmt.Errorf("no non-loopback IPv4 address found")
}
