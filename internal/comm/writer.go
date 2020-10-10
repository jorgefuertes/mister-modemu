package comm

import (
	"fmt"

	"github.com/jorgefuertes/mister-modemu/lib/console"
	"github.com/jorgefuertes/mister-modemu/lib/util"
)

func serialWriteLn(data ...interface{}) (int, error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	n, err := serialWrite(data...)
	serialWrite(cr, lf)
	return n, err
}

func serialWriteBytes(b *[]byte) (int, error) {
	prefix := `SER/TX`

	m.lock.Lock()
	defer m.lock.Unlock()

	console.Debug(prefix, fmt.Sprintf("Sending %v bytes to serial port", len(*b)))
	n, err := m.port.Write(*b)
	if err != nil {
		console.Error(prefix, err.Error())
		return n, err
	}
	console.Debug(prefix, n, " bytes written")
	return n, err
}

func serialWrite(data ...interface{}) (int, error) {
	prefix := `SER/TX`

	m.lock.Lock()
	defer m.lock.Unlock()

	var err error
	var bytes int

	for _, d := range data {
		var b int
		switch v := d.(type) {
		case nil:
			console.Debug(prefix, "Not sending nil")
		case int:
			console.Debug(prefix, util.ByteToStr(byte(v)))
			b, err = m.port.Write([]byte{byte(v)})
		case uint:
			console.Debug(prefix, util.ByteToStr(byte(v)))
			b, err = m.port.Write([]byte{byte(v)})
		case byte:
			console.Debug(prefix, util.ByteToStr(v))
			b, err = m.port.Write([]byte{v})
		case string:
			console.Debug(prefix, v)
			b, err = m.port.Write([]byte(v))
		case []byte:
			serialWriteBytes(&v)
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
