package comm

import (
	"fmt"

	"github.com/jorgefuertes/mister-modemu/lib/console"
)

func serialWriteLn(data ...interface{}) (int, error) {
	n, err := serialWrite(data...)
	serialWrite(cr, lf)
	return n, err
}

func serialWrite(data ...interface{}) (int, error) {
	prefix := `SER/TX`
	// console.Debug(prefix, data)

	var err error
	var bytes int

	for _, d := range data {
		var b int
		switch v := d.(type) {
		case nil:
			console.Debug(prefix, "Not sending nil")
		case int:
			b, err = m.port.Write([]byte{byte(v)})
		case uint:
			b, err = m.port.Write([]byte{byte(v)})
		case byte:
			b, err = m.port.Write([]byte{v})
		case string:
			b, err = m.port.Write([]byte(v))
		case []byte:
			b, err = m.port.Write(v)
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
