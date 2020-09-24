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
		return fmt.Sprintf("STATUS:%v", status.st)
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
		status.cipmux = int8(mode)
		return ok
	}

	// AT+CIFSR - Gets the local IP address
	if strings.HasPrefix(cmd, "AT+CIFSR") {
		ip, err := getOutboundIP()
		if err != nil {
			return er
		}
		return "+CIFSR:APIP," + ip.String()
	}

	// AT+CIPSTART
	if strings.HasPrefix(cmd, "AT+CIPSTART") {
		arg := getArg(&cmd)
		args := getArgs(&arg)

		if status.cipmux == 0 {
			// single conn
			c := &connection{}
			c.t, c.ip = args[0], args[1]

			port, err := strconv.Atoi(args[2])
			if err != nil {
				console.Debug("CONN/START", "Invalid port")
				return er
			}
			c.port = int16(port)

			if len(args) > 3 {
				keep, err := strconv.Atoi(args[3])
				if err != nil {
					console.Debug("CONN/START", "Invalid keep alive")
					return er
				}
				c.keep = int16(keep)
			}

			status.connections[0] = c
		} else {
			// multiple conn
			id, err := strconv.Atoi(args[0])
			if err != nil {
				return er
			}
			c := &connection{}
			c.t, c.ip = args[1], args[2]

			port, err := strconv.Atoi(args[3])
			if err != nil {
				return er
			}
			c.port = int16(port)

			if len(args) > 4 {
				keep, err := strconv.Atoi(args[4])
				if err != nil {
					return er
				}
				c.keep = int16(keep)
			}

			status.connections[id] = c
		}

		return ok

	}

	// Fallback to OK
	return ok
}

// AT+CIPDOMAIN
// AT+CIPSTART
// AT+CIPSSLSIZE
// Description
// AT+CIPSSLCONF
// AT+CIPSEND
// AT+CIPSENDEX
// AT+CIPSENDBUF
// AT+CIPBUFRESET
// AT+CIPBUFSTATUS
// AT+CIPCHECKSEQ
// AT+CIPCLOSE
// AT+CIFSR
// AT+CIPMUX
// AT+CIPSERVER
// AT+CIPSERVERMAXCONN
// AT+CIPMODE
// AT+SAVETRANSLINK
// AT+CIPSTO
// AT+PING
// AT+CIUPDATE
// AT+CIPDINFO
// AT+IPD
// AT+CIPRECVMODE
// AT+CIPRECVDATA
// AT+CIPRECVLEN
// AT+CIPSNTPCFG
// AT+CIPSNTPTIME
// AT+CIPDNS_CUR
// AT+CIPDNS_DEF
