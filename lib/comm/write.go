package comm

import "github.com/jorgefuertes/mister-modemu/lib/console"

// Write - Write over port
func Write(snd string) (int, error) {
	b := []byte(snd)
	n, err := s.Write(b)
	if err != nil {
		console.Error("COMM/SND", err.Error())
	}

	return n, err
}
