package comm

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/jorgefuertes/mister-modemu/lib/cfg"
	"github.com/jorgefuertes/mister-modemu/lib/console"
	"github.com/jorgefuertes/mister-modemu/lib/util"
)

// SerialReader - Read loop
func SerialReader() {
	// log prefix
	prefix := "SER/RX"

	// buffers
	sBuf := make([]byte, 1, 255) // cipsend buff
	cBuf := make([]byte, 1, 255) // command buff
	rBuf := make([]byte, 1, 128) // receiving buff
	sBuf = []byte{}
	cBuf = []byte{}

	// read loop
	console.Debug(prefix, "Read Loop Begin")
	for {
		var eof error
		var n int
		n, eof = m.port.Read(rBuf)
		if eof == io.EOF {
			console.Debug(prefix, `EOF`)
		} else if eof != nil {
			console.Error(prefix, eof.Error())
			continue
		}
		if n == 0 {
			time.Sleep(200 * time.Millisecond)
			continue
		}
		if cfg.IsDev() {
			for _, d := range rBuf {
				console.Debug(prefix, fmt.Sprintf("%03d %s", d, util.ByteToStr(d)))
			}

		}

		// if we are waiting for data to send
		if m.snd.on {
			for i, d := range rBuf {
				console.Debug(prefix,
					fmt.Sprintf("SendBuffer [%v/%v]: %03d", i, m.snd.len, d))
				sBuf = append(sBuf, d)
				if len(sBuf) == int(m.snd.len) {
					console.Debug(prefix,
						fmt.Sprintf("Data set complete with %v bytes", len(sBuf)))
					// data transmission
					serialWriteLn(fmt.Sprintf("Rec %v bytes", len(sBuf)))
					if len(rBuf) > int(m.snd.len) {
						serialWriteLn("BUSY")
					}
					console.Debug(fmt.Sprintf("NET/TX/%v", m.snd.ID), "BEGIN")
					for i1, d1 := range sBuf {
						console.Debug(
							fmt.Sprintf("NET/TX/%v", m.snd.ID),
							fmt.Sprintf(" [%v/%v] %s %03d", i1, len(sBuf), util.ByteToStr(d1), d1))
					}
					_, err := m.connections[m.snd.ID].conn.Write(sBuf)
					if err != nil {
						console.Error(prefix, err)
						serialWriteLn(er)
					}
					console.Debug(fmt.Sprintf("NET/TX/%v", m.snd.ID), "END")
					sBuf, cBuf = []byte{}, []byte{}
					clearSnd()
					serialWriteLn("SEND OK")
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
			if rBuf[len(rBuf)-1] == cr {
				serialWrite(lf)
			}
		}

		// read buffer to command buffer
		for _, d := range rBuf {
			cBuf = append(cBuf, d)

			if len(cBuf) > 255 {
				console.Error(prefix, "Command buffer limit reached")
				serialWriteLn(er)
				cBuf = []byte{}
				m.port.Flush()
				break
			}

			// process if finished
			if d == lf || d == cr {
				m.port.Flush()
				console.Debug("BUF/CMD", string(cBuf))
				if strings.HasPrefix(string(cBuf), `AT`) {
					cmd := bufToStr(&cBuf)
					res := parseCmd(cmd)
					if res != hush {
						console.Debug("SER/REPLY", res)
						serialWriteLn(res)
					}
				}
				cBuf = []byte{}
			}
		}
	}
}
