package modem

import (
	"os"

	"github.com/jorgefuertes/mister-modemu/internal/console"
	"github.com/tarm/serial"
)

// Open - Open the serial port
func (s *Status) Open(port *string, baud *int) error {
	prefix := "SER/OPEN"
	if _, err := os.Stat(*port); os.IsNotExist(err) {
		return err
	}
	console.Debug(prefix, "Opening serial port")
	s.Reset()
	var err error
	s.port, err = serial.OpenPort(&serial.Config{Name: *port, Baud: *baud})
	if err != nil {
		return err
	}
	console.Debug(prefix, "Serial port open")

	return nil
}
