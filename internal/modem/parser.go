package modem

import (
	"fmt"
	"net"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/jorgefuertes/mister-modemu/internal/build"
	"github.com/jorgefuertes/mister-modemu/internal/cfg"
	"github.com/jorgefuertes/mister-modemu/internal/console"
	"github.com/tatsushid/go-fastping"
)

// parse entry func
// does nothing if buffer doesn't start with AT
func (m *Modem) parse() {
	cmd := m.bufToStr()
	if strings.HasPrefix(cmd, "AT") {
		res := parseAT(m, cmd)
		if res != hush {
			m.writeLn(res)
		}
	}
}

// one arg line, even if it has colon sep args
func getArg(cmd *string) string {
	r := regexp.MustCompile(`^[A-Z]+\=\"*(?P<Arg>.*)\"*$`)
	m := r.FindStringSubmatch(*cmd)
	if len(m) < 2 {
		return ""
	}
	return m[1]
}

// slice of args from colon sep args
func getArgs(argLine *string) []string {
	prefix := `AT/CMD/ARG`
	args := strings.Split(*argLine, ",")
	for i, a := range args {
		args[i] = strings.Trim(a, `"`)
		console.Debug(prefix, args[i])
	}

	return args
}

func removeAT(cmd string) string {
	return strings.TrimPrefix(cmd, "AT+")
}

// at parser
func parseAT(m *Modem, cmd string) string {
	// Log prefix
	prefix := `AT/PARSER`

	cmd = removeAT(cmd)
	console.Debug(prefix, "'"+cmd+"'")

	// AT
	if cmd == "AT" {
		return ok
	}

	// AT+VERSION
	if cmd == "VERSION" {
		return "+VERSION:" + build.Version() + crlf + ok
	}

	// AT+AUTHOR
	if cmd == "AUTHOR" {
		return "+AUTHOR:" + *cfg.Config.Author + crlf + ok
	}

	// AT+RST
	if cmd == "RST" {
		m.init()
		return ok
	}

	// AT+HELP
	if cmd == "HELP" {
		var s string
		for _, line := range help {
			s += line + crlf
		}
		return s + ok
	}

	// AT+HELP GAMES
	if cmd == "HELP GAMES" {
		var s string
		for _, line := range helpGames {
			s += line + crlf
		}
		return s + ok
	}

	// AT+HELP GAMES
	if cmd == "LIST GAMES" {
		var s string
		for _, line := range listGames {
			s += line + crlf
		}
		return s + ok
	}

	// ATE
	if strings.HasPrefix(cmd, "ATE") {
		if strings.HasSuffix(cmd, "0") {
			m.ate = false
			return ok
		} else if strings.HasSuffix(cmd, "1") {
			m.ate = true
			return ok
		}
		return er
	}

	// AT+CIPSTATUS
	if cmd == "CIPSTATUS" {
		s := fmt.Sprintf("STATUS:%v", m.status)
		for i, c := range m.connections {
			if c != nil {
				s += fmt.Sprintf("\r\n+CIPSTATUS:%v,%s,%s,%v,%v,%v",
					i, c.t, c.ip, c.port, 0, c.cs) + crlf
			}
		}

		return s
	}

	// AT+CIPDOMAIN
	if strings.HasPrefix(cmd, "CIPDOMAIN") {
		name := getArg(&cmd)
		if name == "" {
			return er
		}
		if len(name) < 3 {
			return er
		}
		ips, err := net.LookupIP(name)
		if err != nil {
			return "DNS Fail\r\nERROR"
		}
		return fmt.Sprintf("+CIPDOMAIN:%s", ips[0]) + crlf + ok
	}

	// AT+CIPMUX
	if strings.HasPrefix(cmd, "AT+CIPMUX") {
		if cmd == "CIPMUX?" {
			return fmt.Sprintf("+CIPMUX:%v\r\nOK", m.cipmux)
		}
		mode, err := strconv.Atoi(getArg(&cmd))
		if err != nil {
			return er
		}
		if mode > 1 || mode < 0 {
			return er
		}
		m.cipmux = uint8(mode)
		return ok
	}

	// AT+CIFSR - Gets the local IP address
	if strings.HasPrefix(cmd, "CIFSR") {
		prefix := "PARSER/CIFSR"
		ip, err := getLocalIP()
		if err != nil {
			console.Error(prefix, err.Error())
			return er
		}

		mac, err := getLocalMac(ip)
		if err != nil {
			console.Error(prefix, err.Error())
			return er
		}

		return fmt.Sprintf("+CIFSR:APIP,\"%s\"\r\n", ip.String()) +
			fmt.Sprintf("+CIFSR:APMAC,\"%s\"\r\n", mac.String()) +
			fmt.Sprintf("+CIFSR:STAIP,\"%s\"\r\n", ip.String()) +
			fmt.Sprintf("+CIFSR:STAMAC,\"%s\"\r\n", mac.String()) +
			ok
	}

	// AT+CIPSTART
	if strings.HasPrefix(cmd, "CIPSTART") {
		prefix = `PARSER/CIPSTART`
		arg := getArg(&cmd)
		args := getArgs(&arg)

		var c *connection = &connection{}
		var err error
		var id int

		if m.cipmux == 0 {
			// single conn
			id = 0

			// type
			c.t = strings.ToUpper(args[0])
			if c.t != "TCP" && c.t != "UDP" {
				console.Error(prefix, "Unimplemented conn type")
				return er
			}

			// remote IP
			c.ip = args[1]

			// port
			c.port, err = strconv.Atoi(args[2])
			if err != nil {
				console.Warn(prefix, "Invalid port")
				return er
			}

			// keep alive
			if len(args) > 3 {
				keep, err := strconv.Atoi(args[3])
				if err != nil {
					console.Warn(prefix, "Invalid keep alive")
					return er
				}
				c.keep = keep
			}
		} else {
			// multiple conn
			id, err = strconv.Atoi(args[0])
			if err != nil {
				console.Warn(prefix, "Invalid port")
				return er
			}

			// type
			c.t = args[1]

			// remote IP
			c.ip = args[2]

			// port
			c.port, err = strconv.Atoi(args[3])
			if err != nil {
				console.Warn(prefix, "Invalid port")
				return er
			}

			// keep alive
			if len(args) > 4 {
				keep, err := strconv.Atoi(args[4])
				if err != nil {
					console.Warn(prefix, "Invalid keep alive")
					return er
				}
				c.keep = keep
			}
		}

		// connect
		c.conn, err = net.Dial(strings.ToLower(c.t), c.ip+":"+strconv.Itoa(c.port))
		if err != nil {
			console.Warn(prefix,
				"Cannot dial '"+strings.ToLower(c.t), c.ip+":"+strconv.Itoa(c.port)+"'")
			return er
		}

		m.connections[id] = c
		m.status = 3
		// start listener
		go m.listenLink(uint8(id))
		return ok
	}

	// AT+CIPSEND
	if cmd == "CIPSEND" {
		prefix = "CIPSEND"
		m.snd.id = uint8(0)
		m.snd.len = uint(0)
		m.snd.on = true
		m.snd.ts = false
		console.Debug(prefix,
			fmt.Sprintf("SEND link %v transparent packet mode ON", m.snd.id))

		return ok
	}

	// AT+CIPSEND=<params>
	if strings.HasPrefix(cmd, "CIPSEND") {
		prefix = `CIPSEND2`
		var connNum int
		var sndLen int
		var err error

		arg := getArg(&cmd)
		args := getArgs(&arg)

		if m.cipmux == 0 {
			if m.connections[0] == nil {
				return er
			}
			connNum = 0
			sndLen, err = strconv.Atoi(args[0])
			if err != nil {
				console.Warn(prefix, "Invalid param length")
				return er
			}
		} else {
			if len(args) < 2 {
				console.Warn(prefix, "Missing link_id")
				return er
			}
			connNum, err := strconv.Atoi(args[0])
			if err != nil || connNum > 4 {
				console.Warn(prefix, "Invalid link_id")
				return er
			}
			if m.connections[connNum] == nil {
				return er
			}
			sndLen, err = strconv.Atoi(args[1])
			if err != nil {
				console.Warn(prefix, "Invalid param length")
				return er
			}
		}

		m.snd.id = uint8(connNum)
		m.snd.len = uint(sndLen)
		m.snd.on = true
		m.snd.ts = false

		console.Debug(prefix,
			fmt.Sprintf("SEND link %v waiting for %v bytes", m.snd.id, m.snd.len))

		return ok
	}

	// AT+PING
	if strings.HasPrefix(cmd, "PING") {
		prefix = `PING`
		host := getArg(&cmd)
		if host == "" {
			return er
		}

		p := fastping.NewPinger()
		ra, err := net.ResolveIPAddr("ip4:icmp", host)
		if err != nil {
			return er
		}

		p.AddIPAddr(ra)
		p.Network("udp")
		p.OnRecv = func(addr *net.IPAddr, rtt time.Duration) {
			m.writeLn("+", rtt.String())
			m.writeLn(ok)
			p.Stop()
		}

		err = p.Run()
		if err != nil {
			console.Warn("AT/PING", err.Error())
			return er
		}
	}

	// AT+CIPCLOSE
	if strings.HasPrefix(cmd, "CIPCLOSE") {
		if m.cipmux == 0 {
			if m.connections[0] != nil {
				m.connections[0].conn.Close()
				m.connections[0] = nil
				m.init()
			}

			return ok
		}

		n, err := strconv.Atoi(getArg(&cmd))
		if err != nil {
			return er
		}

		if n == 5 {
			for i := 0; i < 5; i++ {
				if m.connections[i] != nil {
					m.connections[i].conn.Close()
					m.connections[i] = nil
				}
			}
			m.init()

			return ok
		}

		if n < 5 {
			if m.connections[n] != nil {
				m.connections[n].conn.Close()
				m.connections[n] = nil
				m.init()
			}

			return ok
		}

		return er
	}

	// Fallback to OK
	return ok
}
