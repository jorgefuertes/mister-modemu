package comm

import (
	"fmt"
	"io"
	"time"

	"github.com/jorgefuertes/mister-modemu/lib/console"
)

// SerialReader - Read loop
func SerialReader() {
	// log prefix
	prefix := "SER/RX"

	// buffers
	sBuf := make([]byte, 1, 2048) // cipsend buff
	cBuf := make([]byte, 1, 1024) // command buff
	cBuf = []byte{}               // command buff initialize
	rBuf := make([]byte, 1, 128)  // receiving buff

	// read loop
	console.Debug(prefix, "Read Loop Begin")
	for {
		var eof error
		var n int
		n, eof = m.port.Read(rBuf)
		if eof == io.EOF {
			console.Debug(prefix, `EOF`)
			if n < 1 {
				time.Sleep(200 * time.Millisecond)
				continue
			}
		} else if eof != nil {
			console.Error(prefix, eof.Error())
			continue
		}

		// if we are waiting for data to send
		if m.snd.on {
			for i := 0; i < n; i++ {
				m.snd.rec++
				console.Debug(prefix,
					"SNDBUF[", m.snd.rec, "/", m.snd.len, "]: ", rBuf[i])
				sBuf = append(sBuf, rBuf[i])
				if m.snd.rec >= m.snd.len {
					console.Debug(prefix, "Data set complete")
					// data transmission
					serialWrite(fmt.Sprintf("Rec %v bytes\r\n", m.snd.rec))
					if len(rBuf) > i {
						serialWrite("BUSY\r\n")
					}

					console.Debug(prefix, "Write BEGIN to link ", m.snd.ID)
					_, err := m.connections[m.snd.ID].conn.Write(sBuf)
					if err != nil {
						console.Error(prefix, err)
						serialWrite(er, cr, lf)
					}
					console.Debug(prefix, "Write END")
					clearSnd()
					serialWrite("SEND OK", cr, lf)
					break
				}
			}
			continue
		}

		// Delete
		if rBuf[0] == bs || rBuf[0] == del {
			if len(cBuf) > 0 {
				cBuf = cBuf[:len(cBuf)-1]
			}
			if m.echo {
				serialWrite(bs)
			}
			continue
		}

		// Echo
		if m.echo {
			serialWrite(rBuf)
			if rBuf[0] == cr {
				serialWrite(lf)
			}
		}

		// red buffer to command buffer
		for i := 0; i < n; i++ {
			cBuf = append(cBuf, rBuf[i])
			if len(cBuf) > 2048 {
				console.Error("CONN/RX", "Command buffer limit reached")
				serialWrite(er)
				cBuf = []byte{}
				break
			}
		}

		// process if finished
		if cBuf[len(cBuf)-1] == lf || cBuf[len(cBuf)-1] == cr {
			if string(cBuf[0:2]) == "AT" {
				cmd := bufToStr(&cBuf)
				res := parseCmd(cmd)
				if res != hush {
					console.Debug("CONN/REPLY", res)
					serialWriteLn(res)
				}
			}
			cBuf = []byte{}
		}
	}
}
