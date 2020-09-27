package comm

import (
	"os"
	"sync"

	"github.com/jorgefuertes/mister-modemu/lib/console"
	"github.com/tarm/serial"
)

// Open - Open the serial port
func Open(port *string, baud *int) {
	prefix := "SER/OPEN"
	if _, err := os.Stat(*port); os.IsNotExist(err) {
		console.Error(prefix, "Cannot find port ", *port)
		os.Exit(4)
	}
	console.Debug(prefix, "Opening serial port")
	var err error
	m.port, err = serial.OpenPort(&serial.Config{Name: *port, Baud: *baud})
	if err != nil {
		console.Error(prefix, err.Error())
		os.Exit(1)
	}
	console.Debug(prefix, "Serial port open")
	m.lock = new(sync.Mutex)
	resetStatus()
}

// Close - Closes the port
func Close() {
	m.lock.Lock()
	m.port.Flush()
	m.port.Close()
	m.lock.Unlock()
}
