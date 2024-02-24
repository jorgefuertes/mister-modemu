package modem

import (
	"fmt"
	"net"
	"regexp"

	"github.com/jorgefuertes/mister-modemu/internal/console"
	"github.com/tarm/serial"
)

// inet links
type connection struct {
	ID     uint8    // Connection ID 0-4
	T      string   // type
	IP     string   // remote IP
	Port   int      // remote port
	Keep   int      // keep alive
	conn   net.Conn // connection
	Cs     int      // cipstatus
	Closed bool     // mark to close
	b      []byte   // link buffer
	n      int      // n bytes recv
}

type route struct {
	path string         // Command path
	e    *regexp.Regexp // Regexp
	cb   func(s *Status)
}

type parser struct {
	routes []route
	Cmd    string
	Err    error
}

// Status - general modem status
type Status struct {
	Sta     uint8 // modem status
	CipMux  bool  // status (false: one, true: multiple)
	CipInfo bool  // true to show the Remote IP and Port with +IPD
	Ate     bool  // echo on/off
	Cw      uint8 // cwmode (1: Station, 2: SoftAP, 3: SoftAP+Station)
	cipsend struct {
		on  bool  // cipsend on/off
		id  uint8 // connection ID
		ts  bool  // transparent mode on/off
		len uint  // expected len
	}
	Connections [5]*connection // Internet connections
	port        *serial.Port   // serial port
	sconf       struct {
		port *string // serial port string
		baud *int    // speed
	}
	b      []byte  // serial port buffer
	n      int     // n bytes received at serial port
	Parser *parser // AT Parser
}

func (s *Status) init() {
	console.Debug(`STATUS/GENERAL`, "Initializing")
	s.Reset()
}

// Reset - modem reset
func (s *Status) Reset() {
	console.Debug(`STATUS/GENERAL`, "Reset")
	// modem state
	s.Ate = false
	s.Sta = 2
	s.CipMux = false
	s.Cw = 1
	// buffer
	s.b = make([]byte, 2048)
	s.n = 0
	s.CipClearAll()
	// port
	if s.port != nil {
		s.port.Flush()
	}
	// links
	for _, c := range s.Connections {
		if c != nil && !c.Closed {
			c.conn.Close()
			c.Closed = true
		}
	}
	// parser
	if s.Parser == nil {
		s.Parser = new(parser)
	}
}

// CipClear - clear the send status for ID
func (s *Status) CipClear(id uint8) {
	if s.cipsend.id != id {
		return
	}
	console.Debug(fmt.Sprintf("LINK/CLEAR/%v", id), "Clear CIPSEND status")
	s.CipClearAll()
}

// CipClearAll - clear the send status for any its ID
func (s *Status) CipClearAll() {
	s.cipsend.id = 0
	s.cipsend.on = false
	s.cipsend.ts = false
	s.cipsend.len = 0
	console.Debug("LINK/CLEAR/0", "CIPSEND status cleared")
}

// CipSet - sets the send len metadata
func (s *Status) CipSet(id uint8, len uint) error {
	_, err := s.GetConn(id)
	if err != nil {
		return err
	}
	s.cipsend.on = true
	s.cipsend.id = id
	s.cipsend.len = len
	console.Debug(fmt.Sprintf("LINK/SET/%d", id), fmt.Sprintf("ID: %d, LEN: %v, ON", id, len))
	return nil
}

// CipPacketOn - sets the transparent packet mode on
func (s *Status) CipPacketOn() error {
	_, err := s.GetConn(0)
	if err != nil {
		return err
	}
	s.cipsend.on = true
	s.cipsend.id = 0
	s.cipsend.ts = true
	s.cipsend.len = 0
	console.Debug("LINK/PACKET/0", "Packet mode: ON")
	return nil
}

// CipClose - close a link
func (s *Status) CipClose(id uint8) error {
	c, err := s.GetConn(id)
	if err != nil {
		return err
	}

	c.close()
	return nil
}
