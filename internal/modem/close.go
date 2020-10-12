package modem

// Close - Closes the port
func (m *Modem) Close() {
	m.port.Flush()
	m.port.Close()
}
