package modem

import (
	"fmt"
)

// byteToStr - Transform byte into readable string
func byteToStr(b byte) string {
	switch b {
	case lf:
		return `[lf]`
	case cr:
		return `[cr]`
	case bs:
		return `[bs]`
	case del:
		return `[dlt]`
	default:
		if b >= 32 && b < 127 {
			return string(b)
		}
		return fmt.Sprintf("[%x]", b)
	}
}

// bufToStr - Transform byte buffer to trimmed string
func (m *Modem) bufToStr() string {
	var s string
	for i := 0; i < m.n; i++ {
		if m.b[i] > 31 && m.b[i] < 126 {
			s += string(m.b[i])
		}
	}

	return s
}

// bufToDebug - Buffer to debug string
func (m *Modem) bufToDebug() string {
	var s string
	for i := 0; i < m.n; i++ {
		s += byteToStr(m.b[i])
	}
	return s
}
