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
				time.Sleep(50 * time.Millisecond)
			} else {
				console.Error("CONN/RX", err.Error())
			}
		} else {
			if cfg.IsDev() {
				console.Debug("CONN/RX", fmt.Sprintf("%q", buf[:n]))
			}
			if string(buf[0:1]) == "AT" {
				console.Debug("COMMAND", fmt.Sprintf("%q", buf[:n]))
				Write("OK")
			} else {
				Write("WHAT?")
			}
		}
	}
}
