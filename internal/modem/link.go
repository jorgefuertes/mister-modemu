package modem

import (
	"fmt"
	"net"
	"time"

	"github.com/jorgefuertes/mister-modemu/internal/console"
)

func (m *Modem) listenLink(id uint8) {
	var err error
	var res string
	var b []byte = make([]byte, 2048)
	var n int

	prefix := fmt.Sprintf("NET/LISTEN/%v/%v:%v", id, m.connections[id].ip, m.connections[id].port)

	console.Debug(prefix, "Listening")

	for {
		// Check if close on
		if m.connections[id].close {
			console.Debug(prefix, "Closing")
			break
		}
		// Set timeout
		err = m.connections[id].conn.SetReadDeadline(time.Now().Add(30 * time.Second))
		if err != nil {
			console.Warn(prefix, "SetReadDeadline failed: ", err)
			console.Warn(prefix, "Closed")
			return
		}
		// Read
		n, err = m.connections[id].conn.Read(b)
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				console.Debug(prefix, "TimeOut - resuming")
				continue
			} else {
				console.Debug(prefix, err.Error())
				break
			}
		}
		// Something received?
		console.Debug(prefix, "Received ", n, " bytes")
		if n > 0 {
			cut := b[0:n]
			if m.cipmux == 0 {
				res = fmt.Sprintf("+IPD,%v:", n)
			} else {
				res = fmt.Sprintf("+IPD,%v,%v:", id, n)
			}
			m.write(res)
			m.writeBytes(&cut)
			m.write(cr, lf)
			console.Debug(prefix, "EOT")
		}
	}
}
