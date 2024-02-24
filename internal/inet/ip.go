package inet

import (
	"net"
	"strconv"

	"github.com/jackpal/gateway"
)

// GetNameServer - Namerserver IP
func GetNameServer() net.IP {
	return net.ParseIP("127.0.0.1")
}

// GetGateway - Local gateway IP
func GetGateway() net.IP {
	gw, err := gateway.DiscoverGateway()
	if err != nil {
		return net.ParseIP("127.0.0.1")
	}

	return gw
}

// GetLocalInterface - Returns the local network interface
func GetLocalInterface() net.Interface {
	ifs, err := net.Interfaces()
	if err != nil {
		return net.Interface{}
	}

	for _, i := range ifs {
		addrs, err := i.Addrs()
		if err != nil {
			return net.Interface{}
		}

		for _, a := range addrs {
			if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					return i
				}
			}
		}
	}
	return net.Interface{}
}

// GetLocalIP - Returns local IP
func GetLocalIP() net.IP {
	defIP := net.ParseIP("127.0.0.1")
	i := GetLocalInterface()

	addrs, err := i.Addrs()
	if err != nil {
		return defIP
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.To4()
			}
		}
	}

	return defIP
}

// GetLocalMac - Return local MAC
func GetLocalMac() net.HardwareAddr {
	return GetLocalInterface().HardwareAddr
}

// GetLocalMask - Return local mask
func GetLocalMask() string {
	var mask string
	m := GetLocalIP().DefaultMask()
	for i, v := range m {
		if i > 0 {
			mask += `.`
		}
		mask += strconv.Itoa(int(v))
	}
	return mask
}
