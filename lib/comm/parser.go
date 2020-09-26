package comm

import (
	"fmt"
	"net"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/jorgefuertes/mister-modemu/lib/cfg"
	"github.com/jorgefuertes/mister-modemu/lib/console"
	"github.com/tatsushid/go-fastping"
)

// one arg line, even if it has colon sep args
func getArg(cmd *string) string {
	r := regexp.MustCompile(`^AT\+[A-Z]+\=\"*(?P<Arg>.*)\"*$`)
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

func bufToStr(buf *[]byte) string {
	var str string
	for _, b := range *buf {
		if b == cr || b == lf {
			break
		}
		if b == 0x00 {
			continue
		}
		str += string(b)
	}

	return str
}

// parser
func parseCmd(cmd string) string {
	// Log prefix
	prefix := `AT/PARSER`

	console.Debug(prefix, "'"+cmd+"'")

	// AT
	if cmd == "AT" {
		return ok
	}

	// AT+VERSION
	if cmd == "AT+VERSION" {
		serialWriteLn("+VERSION:", *cfg.Config.Version)
		return ok
	}

	// AT+AUTHOR
	if cmd == "AT+AUTHOR" {
		serialWriteLn("+AUTHOR:", *cfg.Config.Author)
		return ok
	}

	// AT+RST
	if cmd == "AT+RST" {
		resetStatus()
		return ok
	}

	// AT+HELP
	if cmd == "AT+HELP" {
		for _, line := range help {
			serialWriteLn(line)
		}
		return ok
	}

	// ATE
	if strings.HasPrefix(cmd, "ATE") {
		if strings.HasSuffix(cmd, "0") {
			m.echo = false
			return ok
		} else if strings.HasSuffix(cmd, "1") {
			m.echo = true
			return ok
		}
		return er
	}

	// AT+CIPSTATUS
	if cmd == "AT+CIPSTATUS" {
		serialWriteLn(fmt.Sprintf("STATUS:%v", m.status))
		for i, c := range m.connections {
			if c != nil {
				serialWriteLn(
					fmt.Sprintf("\r\n+CIPSTATUS:%v,%s,%s,%v,%v,%v", i, c.t, c.ip, c.port, 0, c.cs))
			}
		}

		return hush
	}

	// AT+CIPDOMAIN
	if strings.HasPrefix(cmd, "AT+CIPDOMAIN") {
		prefix = `CIPDOMAIN`
		name := getArg(&cmd)
		if name == "" {
			return er
		}
		if len(name) < 3 {
			return er
		}
		ips, err := net.LookupIP(name)
		if err != nil {
			console.Debug(prefix, err.Error())
			return "DNS Fail\r\nERROR"
		}
		serialWriteLn(fmt.Sprintf("+CIPDOMAIN:%s", ips[0]))

		return ok
	}

	// AT+CIPMUX
	if strings.HasPrefix(cmd, "AT+CIPMUX") {
		if cmd == "AT+CIPMUX?" {
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
	if strings.HasPrefix(cmd, "AT+CIFSR") {
		prefix = `CIFSR`

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

		serialWriteLn(fmt.Sprintf("+CIFSR:APIP,\"%s\"", ip.String()))
		serialWriteLn(fmt.Sprintf("+CIFSR:APMAC,\"%s\"", mac.String()))
		serialWriteLn(fmt.Sprintf("+CIFSR:STAIP,\"%s\"", ip.String()))
		serialWriteLn(fmt.Sprintf("+CIFSR:STAMAC,\"%s\"", mac.String()))

		return ok
	}

	// AT+CIPSTART
	if strings.HasPrefix(cmd, "AT+CIPSTART") {
		prefix = `CIPSTART`
		arg := getArg(&cmd)
		args := getArgs(&arg)

		var c *connection = &connection{}
		var id int
		var err error

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
		go listener(id)
		return ok
	}

	// AT+CIPSEND
	if strings.HasPrefix(cmd, "AT+CIPSEND") {
		prefix = `CIPSEND`
		var connNum int
		var sndLen int
		var err error

		arg := getArg(&cmd)
		args := getArgs(&arg)

		if m.cipmux == 0 {
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
			sndLen, err = strconv.Atoi(args[1])
			if err != nil {
				console.Warn(prefix, "Invalid param length")
				return er
			}
		}

		setSnd(uint8(connNum), uint(sndLen))

		console.Debug(prefix,
			fmt.Sprintf("SEND link %v waiting for %v bytes", m.snd.ID, m.snd.len))

		return ok
	}

	// AT+PING
	if strings.HasPrefix(cmd, "AT+PING") {
		prefix = `CIPSEND`
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
			serialWriteLn("+", rtt.String())
			serialWriteLn(ok)
			p.Stop()
		}

		err = p.Run()
		if err != nil {
			console.Warn("AT/PING", err.Error())
			return er
		}
	}

	// Fallback to OK
	return ok
}
