package comm

import (
	"fmt"
	"strings"
	"time"

	"github.com/jorgefuertes/mister-modemu/lib/cfg"
	"github.com/jorgefuertes/mister-modemu/lib/console"
)

// ReadLoop - Read loop
func ReadLoop() {
	console.Info("CONN/READ", "Listening…")
	buf := make([]byte, 128)

	// read loop
	for {
		n, err := s.Read(buf)
		if err != nil {
			if err.Error() == "EOF" {
				time.Sleep(50 * time.Millisecond)
			} else {
				console.Error("CONN/RECV", err.Error())
			}
		} else {
			recv := fmt.Sprintf("%q", buf[:n])
			if cfg.IsDev() {
				console.Debug("CONN/RECV", recv)
			}
			if strings.HasPrefix(recv, "AT") {
				console.Debug("COMMAND", recv)
				Write("OK")
			} else {
				console.Debug("WHAT?", recv)
			}
		}
	}
}
