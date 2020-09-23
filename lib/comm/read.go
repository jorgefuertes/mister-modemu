package comm

import (
	"fmt"
	"time"

	"github.com/jorgefuertes/mister-modemu/lib/cfg"
	"github.com/jorgefuertes/mister-modemu/lib/console"
)

// ReadLoop - Read loop
func ReadLoop() {
	console.Info("CONN/RX", "Listening…")
	buf := make([]byte, 128)

	// read loop
	for {
		n, err := s.Read(buf)
		if err != nil {
			if err.Error() == "EOF" {
				time.Sleep(100 * time.Millisecond)
			} else {
				console.Error("CONN/RX", err.Error())
			}
		} else {
			if buf[0] == 13 || buf[0] == 10 {
				continue
			}
			if cfg.IsDev() {
				console.Debug("CONN/RX", fmt.Sprintf("%q", buf[:n]))
			}
			if string(buf[0:2]) == "AT" {
				write(parseCmd(buf[:n]) + "\r\n")
			}
		}
	}
}
