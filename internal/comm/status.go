package comm

import (
	"fmt"
	"net"
	"sync"

	"github.com/jorgefuertes/mister-modemu/internal/console"
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
	params      *sync.Mutex
	connections [5]*connection
	snd         struct {
		ID  uint8
		on  bool
		len uint
		lst bool
	}
	port *serial.Port
}

var m modem

func resetStatus() {
	m.params.Lock()
	defer m.params.Unlock()
	console.Debug("MODEM/STATUS", "Reseting status")
	m.status = 5
	m.cipmux = 0
	m.echo = false
	_, err := getLocalIP()
	if err != nil {
		m.status = 5
	} else {
		m.status = 2
	}
	m.port.Flush()
}

func clearSnd() {
	m.params.Lock()
	defer m.params.Unlock()
	m.snd.on = false
	m.snd.ID = 0
	m.snd.len = 0
	console.Debug("CLEARSND", fmt.Sprintf("ON:%v ID:%v LEN:%v", m.snd.on, m.snd.ID, m.snd.len))
}

func setSnd(sndID uint8, sndLen uint) {
	m.params.Lock()
	defer m.params.Unlock()
	m.snd.on = true
	m.snd.ID = sndID
	m.snd.len = sndLen
	console.Debug("SETSND", fmt.Sprintf("ON:%v ID:%v LEN:%v", m.snd.on, m.snd.ID, m.snd.len))
}

func setSndLen(n uint) {
	m.params.Lock()
	defer m.params.Unlock()
	m.snd.len = n
}
