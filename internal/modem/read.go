package modem

import (
	"fmt"
	"time"

	"github.com/jorgefuertes/mister-modemu/internal/ascii"
	"github.com/jorgefuertes/mister-modemu/internal/console"
)

// Listen - listen neverending loop
func (s *Status) Listen() {
	prefix := `SER/LST`
	var err error
	for {
		if s.cipsend.on {
			s.Write(">")
		}
		s.n, err = s.port.Read(s.b)
		if err != nil {
			console.Warn(prefix, err.Error())
			if err.Error() == "EOF" {
				s.Close()
				time.Sleep(250 * time.Millisecond)
				s.Open(s.sconf.port, s.sconf.baud)
			}
			continue
		}
		if s.n < 1 {
			continue
		}

		console.Debug(prefix, s.n, " bytes: ", s.bufToDebug())
		if s.cipsend.on {
			if s.cipsend.ts {
				s.recPacket()
			} else {
				s.recData()
			}
		} else {
			s.echo()
			s.Parser.Parse(s, s.bufToStr())
		}
	}
}

func (s *Status) echo() {
	if s.Ate && !s.cipsend.on {
		s.Write(s.b[0:s.n])
	}
}

func (s *Status) recData() {
	prefix := `SER/RDATA`
	// len bytes mode
	console.Debug(prefix, "CIPSEND ON (len bytes mode)")
	// cheking for ATE0, that will be a lost connection or a reseted computer
	if string(s.b[0:4]) == `ATE0` && s.cipsend.len != 5 {
		// let's guess its a reset
		console.Debug(prefix, "Unexpected ATE0: Guessing a computer's RESET")
		s.Reset()
		// simulate ATE0
		s.Ate = false
		s.WriteLn(ascii.OK)
		return
	}

	for i := 0; i <= s.n; i++ {
		if uint(i) == s.cipsend.len {
			s.WriteLn("BUSY")
		}
		console.Debug(prefix, fmt.Sprintf("%04d: %02X %s", i, s.b[i], ascii.ByteToStr(s.b[i])))
	}

	// get connection
	c, err := s.GetConn(s.cipsend.id)
	if err != nil {
		console.Error(prefix, err)
		s.WriteLn(ascii.ER)
	}

	// data complete
	if uint(s.n) >= s.cipsend.len {
		console.Debug(prefix, fmt.Sprintf("Data set complete with %v bytes", s.cipsend.len))
		// data transmission
		s.WriteLn(fmt.Sprintf("Rec %v bytes", s.cipsend.len))
		if _, err = c.conn.Write(s.b[0:s.cipsend.len]); err != nil {
			console.Error(prefix, err)
			s.WriteLn(ascii.ER)
		} else {
			console.Debug(prefix, s.cipsend.len, " bytes sent to remote")
			s.WriteLn("SEND OK")
		}

		s.CipClearAll()
		return
	}

	console.Debug(prefix, fmt.Sprintf("Data set not complete with %v bytes", s.n))
	s.cipsend.len -= uint(s.n)
	if _, err = c.conn.Write(s.b[0:s.n]); err != nil {
		console.Error(prefix, err)
		s.WriteLn(ascii.ER)
	} else {
		console.Debug(prefix, s.n, " bytes sent to remote")
	}
}

func (s *Status) recPacket() {
	prefix := `SER/RPACKET`
	// packet mode
	console.Debug(prefix, "CIPSEND ON (packet mode)")
	for i := 0; i < s.n; i++ {
		console.Debug(prefix, fmt.Sprintf("%04d: %02X %s", i, s.b[i], ascii.ByteToStr(s.b[i])))
	}

	if s.bufToStr() == "+++" {
		// back to command mode
		console.Debug(prefix, "Return to command mode")
		s.CipClearAll()
		s.WriteLn(ascii.OK)
		return
	}

	// get connection
	c, err := s.GetConn(s.cipsend.id)
	if err != nil {
		console.Error(prefix, err)
		s.WriteLn(ascii.ER)
	}

	if _, err = c.conn.Write(s.b[0 : s.n-1]); err != nil {
		console.Error(prefix, err)
		s.WriteLn(ascii.ER)
	} else {
		console.Debug(prefix, s.n, " bytes sent to remote")
	}
}
