package modem

import (
	"fmt"
	"net"
	"time"

	"github.com/jorgefuertes/mister-modemu/internal/cfg"
	"github.com/jorgefuertes/mister-modemu/internal/console"
)

func (m *Modem) prefix(id uint8) string {
	return fmt.Sprintf("NET/LISTEN/%v/%v:%v", id, m.connections[id].ip, m.connections[id].port)
}

func (m *Modem) listenLink(id uint8) {
	var err error
	var res string
	var b []byte = make([]byte, 2048)
	var n int

	console.Debug(m.prefix(id), "Listening")

	for {
		// Check if close on
		if m.connections[id].close {
			m.close(id)
			break
		}
		// Set timeout
		err = m.connections[id].conn.SetReadDeadline(time.Now().Add(30 * time.Second))
		if err != nil {
			console.Warn(m.prefix(id), "SetReadDeadline failed: ", err)
			console.Warn(m.prefix(id), "Closed")
			return
		}
		// Read
		n, err = m.connections[id].conn.Read(b)
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				console.Debug(m.prefix(id), "TimeOut - resuming")
				continue
			}
			console.Debug(m.prefix(id), err.Error())
			m.close(id)
			break
		}
		// Something received?
		// (+CIPMUX=0)+IPD,<len>[,<remote IP>,<remote port>]:<data>
		// (+CIPMUX=1)+IPD,<link ID>,<len>[,<remote IP>,<remote port>]:<data>
		console.Debug(m.prefix(id), "Received ", n, " bytes")
		if n > 0 {
			cut := b[0:n]
			if cfg.IsDev() {
				// Debug received data
				var count int
				var hex string
				var str string
				for i := 0; i < n; i++ {
					count++
					hex += fmt.Sprintf("%02X", cut[i])
					str += byteToStr(cut[i])
					if count == 20 || i == n-1 {
						console.Debug(m.prefix(id), hex, "| ", str)
						count = 0
						hex = ""
						str = ""
					}
				}
			}
			if m.cipmux == 0 {
				res = fmt.Sprintf("+IPD,%v", n)
			} else {
				res = fmt.Sprintf("+IPD,%v,%v", id, n)
			}
			if m.cipinfo {
				res += fmt.Sprintf(",%s,%v", m.connections[id].ip, m.connections[id].port)
			}
			m.write(res + ":")
			m.writeBytes(&cut)
			m.write(cr, lf)
			console.Debug(m.prefix(id), "Internal EOD")
		}
	}
}

// close and remove connection
func (m *Modem) close(id uint8) {
	console.Debug(m.prefix(id), "Link closing")
	m.writeLn(fmt.Sprintf("%v,CLOSED", id))
	console.Debug(m.prefix(id), "Removing connection")
	m.connections[id] = nil
}
