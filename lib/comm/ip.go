package comm

import "net"

func getOutboundIP() (net.IP, error) {
	conn, err := net.Dial("udp", "1.1.1.1:80")
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP, nil
}
