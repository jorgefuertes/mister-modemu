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
	cBuf := make([]byte, 0, 255)     // command buff
	rBuf := make([]byte, 1024, 1024) // receiving buff
	sBuf := make([]byte, 0, 1024)    // receiving buff

	// read loop
	console.Debug(prefix, "Read Loop Begin")
	for {
		console.Debug(prefix, "Listening…")
		n, eof := m.port.Read(rBuf)
		if eof == io.EOF {
			console.Debug(prefix, `EOF`)
			time.Sleep(250 * time.Millisecond)
			continue
		} else if eof != nil {
			console.Error(prefix, eof.Error())
			continue
		}
		if cfg.IsDev() {
			console.Debug(prefix, n, " bytes - ", bufToStr(&rBuf))
		}

		// if we are waiting for data to send
		if m.snd.on {
			for i := 0; i < n; i++ {
				b := rBuf[i]
				if len(sBuf) > int(m.snd.len) {
					serialWriteLn("BUSY")
					break
				}
				console.Debug("SNDBUFF/ADD",
					fmt.Sprintf("[%v/%v] %03d %s", len(sBuf), m.snd.len, b, util.ByteToStr(b)))
				sBuf = append(sBuf, b)
			}
			if len(sBuf) >= int(m.snd.len) {
				console.Debug(prefix,
					fmt.Sprintf("Data set complete with %v bytes", len(sBuf)))
				// data transmission
				serialWriteLn(fmt.Sprintf("Rec %v bytes", len(sBuf)))
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
				sBuf = nil
				clearSnd()
				serialWriteLn("SEND OK")
			}
			continue
		}

		// Echo
		if m.echo {
			for i := 0; i < n; i++ {
				serialWrite(rBuf[i])
			}
		}

		// read buffer to command buffer
		for i := 0; i < n; i++ {
			b := rBuf[i]
			cBuf = append(cBuf, b)
			console.Debug("CBUF/ADD",
				fmt.Sprintf("[%v/%v] %03d - %s", len(cBuf), cap(cBuf), b, util.ByteToStr(b)))
			// overflow
			if len(cBuf) == 1024 {
				console.Error(prefix, "Command buffer limit reached")
				serialWriteLn(er)
				cBuf = nil
				m.port.Flush()
				break
			}
			// process if finished
			if b == lf {
				cmd := bufToStr(&cBuf)
				console.Debug("BUF/CMD", fmt.Sprintf("'%s'", cmd))
				if cmd == "TE0" {
					cmd = "ATE0"
				}
				if strings.HasPrefix(cmd, "AT") {
					res := parseCmd(cmd)
					if res != hush {
						console.Debug("SER/REPLY", res)
						serialWriteLn(res)
					}
				} else {
					console.Debug("BUF/CMD", cmd, ": Not an AT command")
					for i, chr := range cmd {
						console.Debug("BUF/CMD", i, "/", len(cmd), ": ", chr)
					}
				}
				cBuf = nil
			}
		}
	}
}
