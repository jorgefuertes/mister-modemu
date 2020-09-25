package comm

import "net"

type connection struct {
	t    string
	ip   string
	port int
	keep int
	conn *net.Conn
	cs   uint8
}

type st struct {
	st          uint8
	cipmux      uint8
	echo        bool
	connections [5]*connection
}

var status st

func resetStatus() {
	status = st{
		st:     5,
		cipmux: 0,
		echo:   true,
	}
	_, err := getLocalIP()
	if err != nil {
		status.st = 5
	} else {
		status.st = 2
	}
}
