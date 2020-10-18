package modem

import (
	"fmt"

	"github.com/jorgefuertes/mister-modemu/internal/console"
)

// Listen - listen neverending loop
func (m *Modem) Listen() {
	prefix := "SER/LST"
	var err error
	for {
		if m.snd.on {
			m.write(">")
		}
		m.n, err = m.port.Read(m.b)
		if err != nil {
			console.Warn(prefix, err.Error())
			continue
		}
		if m.n < 1 {
			continue
		}

		console.Debug(prefix, m.n, " bytes: ", m.bufToDebug())
		if m.snd.on {
			if m.snd.ts {
				m.recPacket()
			} else {
				m.recData()
			}
		} else {
			m.echo()
			m.parse()
		}
	}
}

func (m *Modem) echo() {
	if m.ate && !m.snd.on {
		m.write(m.b[0:m.n])
	}
}

func (m *Modem) recData() {
	prefix := "SER/RX/LINK"
	// len bytes mode
	console.Debug(prefix, "CIPSEND ON (len bytes mode)")
	// cheking for ATE0, that will be a lost connection or a reseted computer
	if string(m.b[0:4]) == `ATE0` && m.snd.len != 5 {
		// let's guess its a reset
		console.Debug(prefix, "Unexpected ATE0: Guessing a computer's RESET")
		m.init()
		// simulate ATE0
		m.ate = false
		m.writeLn(ok)
		return
	}

	for i := 0; i <= m.n; i++ {
		if uint(i) == m.snd.len {
			m.writeLn("BUSY")
		}
		console.Debug(prefix, fmt.Sprintf("%04d: %02X %s", i, m.b[i], byteToStr(m.b[i])))
	}

	// data complete
	if uint(m.n) >= m.snd.len {
		console.Debug("SER/RX/LINK", fmt.Sprintf("Data set complete with %v bytes", m.snd.len))
		// data transmission
		m.writeLn(fmt.Sprintf("Rec %v bytes", m.snd.len))
		_, err := m.connections[m.snd.id].conn.Write(m.b[0:m.snd.len])
		if err != nil {
			console.Error("LINK/TX", err)
			m.writeLn(er)
		} else {
			console.Debug("LINK/TX", m.snd.len, " bytes sent to remote")
			m.writeLn("SEND OK")
		}
		m.clearSnd()
		return
	}

	console.Debug("SER/RX/LINK", fmt.Sprintf("Data set not complete with %v bytes", m.n))
	m.snd.len -= uint(m.n)
	_, err := m.connections[m.snd.id].conn.Write(m.b[0:m.n])
	if err != nil {
		console.Error("LINK/TX", err)
		m.writeLn(er)
	} else {
		console.Debug("LINK/TX", m.n, " bytes sent to remote")
	}
}

func (m *Modem) recPacket() {
	prefix := "SER/RX/LINK"
	// packet mode
	console.Debug(prefix, "CIPSEND ON (packet mode)")
	for i := 0; i < m.n; i++ {
		console.Debug(prefix, fmt.Sprintf("%04d: %02X %s", i, m.b[i], byteToStr(m.b[i])))
	}

	if m.bufToStr() == "+++" {
		// back to command mode
		console.Debug(prefix, "Return to command mode")
		m.clearSnd()
		m.writeLn(ok)
		return
	}

	_, err := m.connections[m.snd.id].conn.Write(m.b[0 : m.n-1])
	if err != nil {
		console.Error("LINK/TX", err)
		m.writeLn(er)
	} else {
		console.Debug("LINK/TX", m.n, " bytes sent to remote")
	}
}
