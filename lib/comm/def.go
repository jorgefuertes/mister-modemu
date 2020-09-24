package comm

import "github.com/tarm/serial"

// returns
const ok = `OK`
const er = `ERROR`

// ascii
const cr = 0x0D
const lf = 0x0A
const sp = 0x20
const del = 0x7F
const bs = 0x08

// vars
var s *serial.Port

var help = []string{
	`IMPLEMENTED AT COMMANDS:`,
	`------------------------`,
	`AT+HELP`,
	`AT+VERSION`,
	`AT+AUTHOR`,
	`AT+RST`,
	`ATE`,
	`AT+CIPSTATUS`,
	`AT+CIPDOMAIN`,
	`AT+CIPMUX`,
	`AT+CIPSTART`,
	`AT+CIFSR`,
}
