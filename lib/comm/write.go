package comm

import (
	"github.com/jorgefuertes/mister-modemu/lib/console"
)

func write(data string) (int, error) {
	b := []byte(data)
	n, err := s.Write(b)
	if err != nil {
		console.Error("COMM/TX", err.Error())
	}

	return n, err
}

func writeByte(b *[]byte) (int, error) {
	n, err := s.Write(*b)
	if err != nil {
		console.Error("COMM/TX", err.Error())
	}

	return n, err
}
