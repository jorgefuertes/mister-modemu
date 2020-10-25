package modem

import (
	"fmt"

	"github.com/jorgefuertes/mister-modemu/internal/ascii"
	"github.com/jorgefuertes/mister-modemu/internal/console"
)

// OK - write an OK line on serial
func (s *Status) OK() (int, error) {
	return s.WriteLn(ascii.OK)
}

// ERR - write an ERROR line on serial
func (s *Status) ERR() (int, error) {
	return s.WriteLn(ascii.ER)
}

// WriteLn - write a line on serial
func (s *Status) WriteLn(data ...interface{}) (int, error) {
	n, err := s.Write(data...)
	s.Write(ascii.CRLF)
	return n + 2, err
}

// WriteBytes - write a byte array on serial
func (s *Status) WriteBytes(b *[]byte) (int, error) {
	prefix := `SER/TX`

	console.Debug(prefix, fmt.Sprintf("Sending %v bytes to serial port", len(*b)))
	n, err := s.port.Write(*b)
	if err != nil {
		console.Error(prefix, err.Error())
		return n, err
	}
	console.Debug(prefix, n, " bytes written")
	return n, err
}

// Write - write interfaces on serial
func (s *Status) Write(data ...interface{}) (int, error) {
	prefix := `SER/TX/WRITE`

	var err error
	var bytes int

	for _, d := range data {
		var b int
		switch v := d.(type) {
		case nil:
			console.Debug(prefix, "Not sending nil")
		case int:
			console.Debug(prefix, ascii.ByteToStr(byte(v)))
			b, err = s.port.Write([]byte{byte(v)})
		case uint:
			console.Debug(prefix, ascii.ByteToStr(byte(v)))
			b, err = s.port.Write([]byte{byte(v)})
		case byte:
			console.Debug(prefix, ascii.ByteToStr(v))
			b, err = s.port.Write([]byte{v})
		case string:
			if v == ascii.CRLF {
				console.Debug(prefix, "[cr,lf]")
			} else {
				console.Debug(prefix, v)
			}
			b, err = s.port.Write([]byte(v))
		case []byte:
			s.WriteBytes(&v)
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
