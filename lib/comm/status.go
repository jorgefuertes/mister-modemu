package comm

import (
	"net"
	"sync"

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
}

type modem struct {
	status      uint8
	cipmux      uint8
	echo        bool
	lock        *sync.Mutex
	connections [5]*connection
	snd         struct {
		ID  uint8
		on  bool
		len uint
	}
	port *serial.Port
	w    chan []interface{}
}

var m modem

func resetStatus() {
	m.lock.Lock()
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
	m.lock.Unlock()
}

func clearSnd() {
	m.lock.Lock()
	m.snd.on = false
	m.snd.ID = 0
	m.snd.len = 0
	m.lock.Unlock()
}

func setSnd(sndID uint8, sndLen uint) {
	m.lock.Lock()
	m.snd.on = true
	m.snd.ID = sndID
	m.snd.len = sndLen
	m.lock.Unlock()
}
