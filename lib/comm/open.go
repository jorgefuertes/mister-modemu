package comm

import (
	"os"

	"github.com/jorgefuertes/mister-modemu/lib/console"
	"github.com/tarm/serial"
)

var s *serial.Port

// Open - Open the serial port
func Open(port *string, baud *int) {
	if _, err := os.Stat(*port); os.IsNotExist(err) {
		console.Error("COMM/OPEN", "Cannot find port ", *port)
		os.Exit(4)
	}
	c := &serial.Config{Name: *port, Baud: *baud}
	var err error
	s, err = serial.OpenPort(c)
	if err != nil {
		console.Error("COMM/OPEN", err.Error())
		os.Exit(1)
	}
}
