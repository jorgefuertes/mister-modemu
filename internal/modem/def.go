package modem

import (
	"net"

	"github.com/tarm/serial"
)

// returns
const ok = `OK`
const er = `ERROR`
const hush = `#NO#REPLY#`

// ascii
const cr = 0x0D  // carriage return
const lf = 0x0A  // line feed
const sp = 0x20  // space
const del = 0x7F // delete
const bs = 0x08  // backspace
const crlf string = string(cr) + string(lf)

type connection struct {
	t     string   // type
	ip    string   // remote IP
	port  int      // remote port
	keep  int      // keep alive
	conn  net.Conn // connection
	cs    int      // cipstatus
	close bool     // mark to close
}

// Modem object
type Modem struct {
	status      uint8 // modem status
	cipmux      uint8 // cipmux status
	ate         bool  // echo on/off
	connections [5]*connection
	// cipsend structure
	snd struct {
		id  uint8 // connection id
		on  bool  // on/off
		ts  bool  // Transparent mode on/off
		len uint  // expected len
	}
	port *serial.Port // serial port
	b    []byte       // serial port buffer
	n    int          // n bytes received at serial port
}

// wargames help
var help = []string{`HELP NOT AVAILABLE`}
var helpGames = []string{
	`'GAMES' REFERS TO MODELS, SIMULATIONS, AND GAMES WICH HAVE ` +
		`TACTICAL AND STRATEGIC APPLICATIONS.`,
}
var listGames = []string{
	`FALKEN'S MAZE`,
	`BLACK JACK`,
	`GIN RUMMY`,
	`HEARTS`,
	`BRIDGE`,
	`CHECKERS`,
	`CHESS`,
	`POKER`,
	`FIGHTER COMBAT`,
	`GUERRILLA ENGAGEMENT`,
	`DESERT WARFARE`,
	`AIR-TO-GROUND ACTIONS`,
	`THEATERWIDE TACTICAL WARFARE`,
	`THEATERWIDE BIOTOXIC AND CHEMICAL WARFARE`,
	``,
	`GLOBAL THERMONUCLEAR WAR`,
}
