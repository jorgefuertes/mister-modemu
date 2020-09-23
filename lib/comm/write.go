package comm

import (
	"fmt"

	"github.com/jorgefuertes/mister-modemu/lib/console"
)

func write(data string) (int, error) {
	console.Debug("COMM/TX", fmt.Sprintf("%q", data))
	b := []byte(data)
	n, err := s.Write(b)
	if err != nil {
		console.Error("COMM/TX", err.Error())
	}

	return n, err
}
