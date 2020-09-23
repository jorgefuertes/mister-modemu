package comm

import (
	"fmt"
	"time"

	"github.com/jorgefuertes/mister-modemu/lib/cfg"
	"github.com/jorgefuertes/mister-modemu/lib/console"
	"github.com/tarm/serial"
)

// Read - Read loop
func Read() {
	c := &serial.Config{Name: *cfg.Config.Port, Baud: *cfg.Config.Baud}
	s, err := serial.OpenPort(c)
	if err != nil {
		console.Error("CONN/OPEN", err.Error())
	}
	console.Info("CONN/READ", "Listening…")

	buf := make([]byte, 128)
	for {
		n, err := s.Read(buf)
		if err != nil {
			if err.Error() == "EOF" {
				time.Sleep(50 * time.Millisecond)
			} else {
				console.Error("CONN/RECV", err.Error())
			}
		} else {
			if cfg.IsDev() {
				recv := fmt.Sprintf("%q", buf[:n])
				console.Debug("CONN/RECV", recv)
			}
		}
	}
}
