package modem

import "github.com/jorgefuertes/mister-modemu/internal/ascii"

// bufToStr - Return string with byte buffer trimmed
func (m *Status) bufToStr() string {
	var s string
	for i := 0; i < m.n; i++ {
		if m.b[i] > 31 && m.b[i] < 126 {
			s += string(m.b[i])
		}
	}

	return s
}

// bufToDebug - Return debyg string from buffer
func (m *Status) bufToDebug() string {
	var s string
	for i := 0; i < m.n; i++ {
		s += ascii.ByteToStr(m.b[i])
	}
	return s
}
