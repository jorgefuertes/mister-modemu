package comm

import (
	"net"

	"github.com/jorgefuertes/mister-modemu/lib/console"
	"github.com/tarm/serial"
)

type connection struct {
	t    string
	ip   string
	port int
	keep int
	conn net.Conn
	cs   uint8
	ch   *chan byte
}

type modem struct {
	status      uint8
	cipmux      uint8
	echo        bool
	connections [5]*connection
	snd         struct {
		ID  uint8
		on  bool
		len uint
		rec uint
	}
	port *serial.Port
	w    chan []interface{}
}

var m modem

func resetStatus() {
	console.Debug("MODEM/STATUS", "Reseting status")
	m.status = 5
	m.cipmux = 0
	m.echo = true
	_, err := getLocalIP()
	if err != nil {
		m.status = 5
	} else {
		m.status = 2
	}
	m.port.Flush()
}

func clearSnd() {
	m.snd.on = false
	m.snd.ID = 0
	m.snd.len = 0
	m.snd.rec = 0
}

func setSnd(sndID uint8, sndLen uint) {
	m.snd.on = true
	m.snd.ID = sndID
	m.snd.len = sndLen
	m.snd.rec = 0
}
