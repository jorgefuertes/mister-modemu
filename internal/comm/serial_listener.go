package comm

import (
	"fmt"

	"github.com/jorgefuertes/mister-modemu/lib/console"
	"github.com/jorgefuertes/mister-modemu/lib/util"
)

// SerialListener - Listener
func SerialListener() {
	prefix := "SER/LST"
	b := make([]byte, 1024, 1024) // receiving buffer
	for {
		n, err := m.port.Read(b)
		if err != nil {
			console.Warn(prefix, err.Error())
			continue
		}
		if n < 1 {
			continue
		}

		go console.Debug(prefix, n, " bytes: ", util.BufToDebug(b, n))
		if m.snd.on {
			go recData(b, n)
		} else {
			go serialEcho(b, n)
			go parse(b, n)
		}
	}
}

// serialEcho
func serialEcho(b []byte, n int) {
	if m.echo && !m.snd.on {
		serialWrite(b[0:n])
	}
}

func recData(b []byte, n int) {
	prefix := "SER/RX/LINK"
	// if we are waiting for data to send via remote
	if m.snd.on {
		console.Debug(prefix, "CIPSEND is on")
		for i := 0; i <= n; i++ {
			console.Debug(prefix, fmt.Sprintf("%04d: %02X %s", i, b[i], util.ByteToStr(b[i])))
		}
		if uint(n) > m.snd.len {
			serialWriteLn("BUSY")
		}

		// data complete
		if uint(n) >= m.snd.len {
			console.Debug("SER/RX/LINK", fmt.Sprintf("Data set complete with %v bytes", m.snd.len))
			// data transmission
			serialWriteLn(fmt.Sprintf("Rec %v bytes", m.snd.len))
			_, err := m.connections[m.snd.ID].conn.Write(b[0:m.snd.len])
			if err != nil {
				console.Error("LINK/TX", err)
				serialWriteLn(er)
			} else {
				console.Debug("LINK/TX", m.snd.len, " bytes sent to remote")
				serialWriteLn("SEND OK")
			}
			clearSnd()
			return
		}

		go console.Debug("SER/RX/LINK", fmt.Sprintf("Data set not complete with %v bytes", n))
		setSndLen(m.snd.len - uint(n))
		_, err := m.connections[m.snd.ID].conn.Write(b[0:n])
		if err != nil {
			console.Error("LINK/TX", err)
			serialWriteLn(er)
		} else {
			console.Debug("LINK/TX", n, " bytes sent to remote")
		}
	}
}
