package modem

// Init - Init or reset modem status
func (m *Modem) init() {
	m.status = 5
	m.cipmux = 0
	m.cw = 1
	m.ate = false
	_, err := getLocalIP()
	if err != nil {
		m.status = 5
	} else {
		m.status = 2
	}
	m.b = make([]byte, 2048, 2048)
	if m.port != nil {
		m.port.Flush()
	}
	m.clearSnd()
	// Closing connections
	for _, c := range m.connections {
		if c != nil {
			c.close = true
		}
	}
}

// clearSnd - clear the send status
func (m *Modem) clearSnd() {
	m.snd.on = false
	m.snd.ts = false
	m.snd.id = 0
	m.snd.len = 0
}

// SetSnd - sets the send len metadata
func (m *Modem) setSnd(sndID uint8, sndLen uint) {
	m.snd.on = true
	m.snd.id = sndID
	m.snd.len = sndLen
}

// SetSndPacket - sets the transparent packet mode
func (m *Modem) setSndPacket() {
	m.snd.on = true
	m.snd.ts = true
	m.snd.id = 0
	m.snd.len = 0
}
