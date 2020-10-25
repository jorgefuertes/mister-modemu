package modem

import "github.com/jorgefuertes/mister-modemu/internal/console"

// Close - Closes the port
func (m *Status) Close() {
	console.Debug(`SER/PORT`, "Closing")
	m.port.Flush()
	m.port.Close()
}
