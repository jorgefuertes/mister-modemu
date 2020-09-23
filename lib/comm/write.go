package comm

import "github.com/jorgefuertes/mister-modemu/lib/console"

// Write - Write over port
func Write(data string) (int, error) {
	console.Debug("COMM/TX", data)
	b := []byte(data)
	n, err := s.Write(b)
	if err != nil {
		console.Error("COMM/TX", err.Error())
	}

	return n, err
}
