package modem

import (
	"errors"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/jorgefuertes/mister-modemu/internal/ascii"
	"github.com/jorgefuertes/mister-modemu/internal/cfg"
	"github.com/jorgefuertes/mister-modemu/internal/console"
)

// conn log prefix
func (c *connection) prefix() string {
	return fmt.Sprintf("NET/LISTEN/%v:%v", c.IP, c.Port)
}

// run link listener
func (c *connection) listen(s *Status) {
	var err error
	var res string
	c.b = make([]byte, 2048)
	c.n = 0

	console.Debug(c.prefix(), "Listening")

	for {
		if c.Closed {
			s.WriteLn(fmt.Sprintf("%d,CLOSED", c.ID))
			break
		}

		// Set timeout
		c.conn.SetReadDeadline(time.Now().Add(30 * time.Second))

		// Read
		c.n, err = c.conn.Read(c.b)

		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				console.Debug(c.prefix(), "TimeOut - resuming")
				continue
			}
			console.Debug(c.prefix(), err.Error())
			c.close()
			continue
		}

		// Something received?
		// (+CIPMUX=0)+IPD,<len>[,<remote IP>,<remote port>]:<data>
		// (+CIPMUX=1)+IPD,<link ID>,<len>[,<remote IP>,<remote port>]:<data>
		console.Debug(c.prefix(), "Received ", c.n, " bytes")
		if c.n > 0 {
			cut := c.b[0:c.n]
			if cfg.IsDev() {
				// Debug received data
				var count int
				var hex string
				var str string
				for i := 0; i < c.n; i++ {
					count++
					hex += fmt.Sprintf("%02X", cut[i])
					str += ascii.ByteToStr(cut[i])
					if count == 20 || i == c.n-1 {
						console.Debug(c.prefix(), hex, "| ", str)
						count = 0
						hex = ""
						str = ""
					}
				}
			}
			if s.CipMux == false {
				res = fmt.Sprintf("+IPD,%v", c.n)
			} else {
				res = fmt.Sprintf("+IPD,%v,%v", c.ID, c.n)
			}
			if s.CipInfo {
				res += fmt.Sprintf(",%s,%v", c.IP, c.Port)
			}

			s.Write(res + ":")
			s.WriteBytes(&cut)
			s.Write(ascii.CRLF)
			console.Debug(c.prefix(), "Internal EOD")
		}
	}
}

// GetConn - connection by ID
func (m *Status) GetConn(id uint8) (*connection, error) {
	for _, c := range m.Connections {
		if c.ID == id && c.Closed == false {
			return c, nil
		}
	}

	return nil, errors.New("Connection not found")
}

// NewConn - define new connection
func (m *Status) NewConn(id uint8, t string, ip string, port int, keep int) error {
	var err error
	if _, err = m.GetConn(id); err == nil {
		return errors.New("Connection already in use")
	}

	if id > 4 {
		return errors.New("Connection ID should be >= 0 and <= 4")
	}

	m.Connections[id] = &connection{
		ID:     id,
		T:      strings.ToLower(t),
		IP:     ip,
		Port:   port,
		Keep:   keep,
		Closed: false,
	}

	console.Debug(`INET/CONN`, "New connection ", m.Connections[id].String())

	m.Connections[id].conn, err = net.Dial(t, fmt.Sprintf("%s:%d", ip, port))

	if err != nil {
		console.Warn(`INET/CONN`, fmt.Sprintf("Cannot dial %s %s:%d", t, ip, port))
		return err
	}

	return nil
}

// String - Connection to string
func (c *connection) String() string {
	return fmt.Sprintf("ID %d %s %s:%d", c.ID, c.T, c.IP, c.Port)
}

// Close - Link close
func (c *connection) close() {
	console.Debug(`INET/CONN`, "Closing link ", c.String())
	c.conn.Close()
	c.Closed = true
}

// Connect - Link connect
func (c *connection) Connect() error {
	var err error

	// connect
	c.conn, err = net.Dial(strings.ToLower(c.T), fmt.Sprintf("%s:%v", c.IP, c.Port))
	if err != nil {
		console.Warn(c.prefix(), "Cannot dial")
		return err
	}
	console.Debug(c.prefix(), "Connected")
	return nil
}

// Connect - Link connect from modem
func (m *Status) Connect(id uint8) error {
	c, err := m.GetConn(id)
	if err != nil {
		return err
	}
	if err = c.Connect(); err != nil {
		return err
	}
	return nil
}
