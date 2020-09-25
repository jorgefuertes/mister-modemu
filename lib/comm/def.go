package comm

// returns
const ok = `OK`
const er = `ERROR`
const hush = `#NO#REPLY#`

// ascii
const cr = 0x0D
const lf = 0x0A
const sp = 0x20
const del = 0x7F
const bs = 0x08

var help = []string{
	`IMPLEMENTED AT COMMANDS:`,
	`------------------------`,
	`ATE`,
	`AT+HELP`,
	`AT+AUTHOR`,
	`AT+CIFSR`,
	`AT+CIPSTATUS`,
	`AT+CIPDOMAIN`,
	`AT+CIPMUX`,
	`AT+CIPSTART`,
	`AT+CIPSNED`,
	`AT+PING`,
	`AT+RST`,
	`AT+VERSION`,
}
