package modem

import (
	"errors"
	"net"
	"strings"
)

// getLocalIP - Returns local IP
func getLocalIP() (net.IP, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP, nil
			}
		}
	}

	return nil, errors.New("Cannot get IP address")
}

// getLocalMac - Return local MAC
func getLocalMac(ip net.IP) (net.HardwareAddr, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	for _, ifc := range interfaces {
		if addrs, err := ifc.Addrs(); err == nil {
			for _, addr := range addrs {
				if strings.Contains(addr.String(), ip.String()) {
					return ifc.HardwareAddr, nil
				}
			}
		}
	}

	return nil, errors.New("Cannot find hardware address")
}
