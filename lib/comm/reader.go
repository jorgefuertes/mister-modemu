package comm

import (
	"fmt"
	"io"
	"strings"

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
	rBuf := make([]byte, 1, 1)   // receiving buff
	sBuf = []byte{}
	cBuf = []byte{}

	// read loop
	console.Debug(prefix, "Read Loop Begin")
	for {
		var eof error
		console.Debug(prefix, "Listening…")
		_, eof = m.port.Read(rBuf)
		if eof == io.EOF {
			console.Debug(prefix, `EOF`)
		} else if eof != nil {
			console.Error(prefix, eof.Error())
			continue
		}
		b := rBuf[0]
		if cfg.IsDev() {
			console.Debug(prefix, fmt.Sprintf("%03d %s", b, util.ByteToStr(b)))
		}

		// if we are waiting for data to send
		if m.snd.on {
			console.Debug(prefix, fmt.Sprintf("SendBuffer ADD [%v]: %03d", m.snd.len, b))
			sBuf = append(sBuf, b)
			if len(sBuf) == int(m.snd.len) {
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
				sBuf = []byte{}
				clearSnd()
				serialWriteLn("SEND OK")
			}
			continue
		} else {
			if len(cBuf) == 0 && b == 0 {
				continue
			}
		}

		// Delete
		if b == bs || b == del {
			if len(cBuf) > 0 {
				cBuf = cBuf[:len(cBuf)-1]
			}
			if m.echo {
				serialWrite(b)
			}
			continue
		}

		// Echo
		if m.echo {
			serialWrite(b)
		}

		// read buffer to command buffer
		cBuf = append(cBuf, b)

		// overflow
		if len(cBuf) > 255 {
			console.Error(prefix, "Command buffer limit reached")
			serialWriteLn(er)
			cBuf = []byte{}
			m.port.Flush()
			continue
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
			cBuf = []byte{}
		}
	}
}
