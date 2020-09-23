package comm

import (
	"fmt"
	"time"

	"github.com/jorgefuertes/mister-modemu/lib/cfg"
	"github.com/jorgefuertes/mister-modemu/lib/console"
)

const cr = 0x0D
const lf = 0x0A

// ReadLoop - Read loop
func ReadLoop() {
	console.Info("CONN/RX", "Listening…")
	cBuf := make([]byte, 1, 128)
	rBuf := make([]byte, 128)

	// read loop
	for {
		n, err := s.Read(rBuf)
		if err != nil {
			if err.Error() != "EOF" {
				console.Error("CONN/RX", err.Error())
			}
			time.Sleep(100 * time.Millisecond)
			continue
		}
		if status.echo {
			writeByte(&rBuf)
		}
		for i := 0; i < n; i++ {
			cBuf = append(cBuf, rBuf[i])
			if len(cBuf) > 2048 {
				console.Error("CONN/RX", "Command buffer limit reached")
				write(er)
				cBuf = []byte{}
				break
			}
		}
		if cfg.IsDev() {
			console.Debug("CONN/RX/CBUF", len(cBuf), ": ", string(cBuf))
		}
		console.Debug("LAST", fmt.Sprintf("%q", cBuf))
		if cBuf[len(cBuf)-1] == lf {
			_ = cBuf[:len(cBuf)-2]
			if string(cBuf[0:2]) == "AT" {
				write(parseCmd(cBuf))
			}
			cBuf = []byte{}
		}
	}
}
