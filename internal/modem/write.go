package modem

import (
	"fmt"

	"github.com/jorgefuertes/mister-modemu/internal/console"
)

// WriteLn - write a line on serial
func (m *Modem) writeLn(data ...interface{}) (int, error) {
	n, err := m.write(data...)
	m.write(cr, lf)
	return n + 2, err
}

// WriteBytes - write a byte array on serial
func (m *Modem) writeBytes(b *[]byte) (int, error) {
	prefix := `SER/TX`

	console.Debug(prefix, fmt.Sprintf("Sending %v bytes to serial port", len(*b)))
	n, err := m.port.Write(*b)
	if err != nil {
		console.Error(prefix, err.Error())
		return n, err
	}
	console.Debug(prefix, n, " bytes written")
	return n, err
}

// Write - write interfaces on serial
func (m *Modem) write(data ...interface{}) (int, error) {
	prefix := `SER/TX`

	var err error
	var bytes int

	for _, d := range data {
		var b int
		switch v := d.(type) {
		case nil:
			console.Debug(prefix, "Not sending nil")
		case int:
			console.Debug(prefix, byteToStr(byte(v)))
			b, err = m.port.Write([]byte{byte(v)})
		case uint:
			console.Debug(prefix, byteToStr(byte(v)))
			b, err = m.port.Write([]byte{byte(v)})
		case byte:
			console.Debug(prefix, byteToStr(v))
			b, err = m.port.Write([]byte{v})
		case string:
			console.Debug(prefix, v)
			b, err = m.port.Write([]byte(v))
		case []byte:
			m.writeBytes(&v)
		default:
			err = fmt.Errorf("I don't know how to write this: %q(%t)", d, d)
		}

		if err != nil {
			console.Error(prefix, err.Error())
			return bytes, err
		}
		bytes += b
	}

	return bytes, nil
}
