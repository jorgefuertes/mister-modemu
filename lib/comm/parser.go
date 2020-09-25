package comm

import (
	"fmt"
	"net"
	"regexp"
	"strconv"
	"strings"

	"github.com/jorgefuertes/mister-modemu/lib/cfg"
	"github.com/jorgefuertes/mister-modemu/lib/console"
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
	args := strings.Split(*argLine, ",")
	for i, a := range args {
		args[i] = strings.Trim(a, `"`)
		console.Debug("ARG", args[i])
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
	console.Debug("COMM/PARSE", "'"+cmd+"'")

	// AT
	if cmd == "AT" {
		return ok
	}

	// AT+VERSION
	if cmd == "AT+VERSION" {
		return *cfg.Config.Version
	}

	// AT+AUTHOR
	if cmd == "AT+AUTHOR" {
		return *cfg.Config.Author
	}

	// AT+RST
	if cmd == "AT+RST" {
		resetStatus()
		return ok
	}

	// AT+HELP
	if cmd == "AT+HELP" {
		var output string
		for _, line := range help {
			output += line + "\r\n"
		}
		return output
	}

	// ATE
	if strings.HasPrefix(cmd, "ATE") {
		if strings.HasSuffix(cmd, "0") {
			status.echo = false
			return ok
		} else if strings.HasSuffix(cmd, "1") {
			status.echo = true
			return ok
		}
		return er
	}

	// AT+CIPSTATUS
	if cmd == "AT+CIPSTATUS" {
		var output string
		output += fmt.Sprintf("STATUS:%v", status.st)
		for i, c := range status.connections {
			output += fmt.Sprintf("\r\n+CIPSTATUS:%v,%s,%s,%v,%v,%v",
				i, c.t, c.ip, c.port, 0, c.cs)
		}

		return output
	}

	// AT+CIPDOMAIN
	if strings.HasPrefix(cmd, "AT+CIPDOMAIN") {
		name := getArg(&cmd)
		if name == "" {
			return er
		}
		if len(name) < 3 {
			return er
		}
		ips, err := net.LookupIP(name)
		if err != nil {
			console.Debug("DNS/FAIL", err.Error())
			return "DNS Fail\r\nERROR"
		}
		return fmt.Sprintf("+CIPDOMAIN:%s", ips[0])
	}

	// AT+CIPMUX
	if strings.HasPrefix(cmd, "AT+CIPMUX") {
		if cmd == "AT+CIPMUX?" {
			return fmt.Sprintf("+CIPMUX:%v\r\nOK", status.cipmux)
		}
		mode, err := strconv.Atoi(getArg(&cmd))
		if err != nil {
			return er
		}
		if mode > 1 || mode < 0 {
			return er
		}
		status.cipmux = uint8(mode)
		return ok
	}

	// AT+CIFSR - Gets the local IP address
	if strings.HasPrefix(cmd, "AT+CIFSR") {
		ip, err := getLocalIP()
		if err != nil {
			console.Error("AT+CIFSR", err.Error())
			return er
		}
		mac, err := getLocalMac(ip)
		if err != nil {
			console.Error("AT+CIFSR", err.Error())
			return er
		}
		return fmt.Sprintf(
			"+CIFSR:APIP,\"%s\"\r\n+CIFSR:APMAC,\"%s\"\r\n+CIFSR:STAIP,\"%s\"\r\n"+
				"+CIFSR:STAMAC,\"%s\"\r\nOK\r\n",
			ip.String(), mac.String(), ip.String(), mac.String(),
		)
	}

	// AT+CIPSTART
	if strings.HasPrefix(cmd, "AT+CIPSTART") {
		arg := getArg(&cmd)
		args := getArgs(&arg)

		if status.cipmux == 0 {
			// single conn
			c := &connection{}

			// type
			c.t = args[0]
			// remote IP
			c.ip = args[1]
			// port
			port, err := strconv.Atoi(args[2])
			if err != nil {
				console.Warn("CONN/START", "Invalid port")
				return er
			}
			c.port = port
			// keep alive
			if len(args) > 3 {
				keep, err := strconv.Atoi(args[3])
				if err != nil {
					console.Warn("CONN/START", "Invalid keep alive")
					return er
				}
				c.keep = keep
			}

			// connect
			d, err := net.Dial(strings.ToLower(c.t), c.ip+":"+strconv.Itoa(c.port))
			if err != nil {
				console.Warn("CONN/START",
					"Cannot dial '"+strings.ToLower(c.t), c.ip+":"+strconv.Itoa(c.port)+"'")
				return er
			}
			c.conn = &d

			status.connections[0] = c
		} else {
			// multiple conn
			c := &connection{}

			// id
			id, err := strconv.Atoi(args[0])
			if err != nil {
				console.Warn("CONN/START", "Invalid port")
				return er
			}
			// type
			c.t = args[1]
			// remote IP
			c.ip = args[2]
			// port
			port, err := strconv.Atoi(args[3])
			if err != nil {
				console.Warn("CONN/START", "Invalid port")
				return er
			}
			c.port = port
			// keep alive
			if len(args) > 4 {
				keep, err := strconv.Atoi(args[4])
				if err != nil {
					console.Warn("CONN/START", "Invalid keep alive")
					return er
				}
				c.keep = keep
			}

			// connect
			d, err := net.Dial(strings.ToLower(c.t), c.ip+":"+strconv.Itoa(c.port))
			if err != nil {
				console.Warn("CONN/START",
					"Cannot dial '"+strings.ToLower(c.t), c.ip+":"+strconv.Itoa(c.port)+"'")
				return er
			}
			c.conn = &d

			status.connections[id] = c
		}

		status.st = 3
		return ok
	}

	// Fallback to OK
	return ok
}
