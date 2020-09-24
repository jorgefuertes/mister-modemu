package comm

import (
	"fmt"
	"net"
	"regexp"
	"strconv"
	"strings"

	"github.com/jorgefuertes/mister-modemu/lib/console"
)

const ok = `OK`
const er = `ERROR`

func getArg(cmd *string) string {
	r := regexp.MustCompile(`^AT\+[A-Z]+\=(?P<Arg>.*$)`)
	m := r.FindStringSubmatch(*cmd)
	if len(m) < 2 {
		return ""
	}
	return m[1]
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

func parseCmd(cmd string) string {
	console.Debug("COMM/PARSE", "'"+cmd+"'")

	// AT
	if cmd == "AT" {
		return ok
	}

	// AT+RST
	if cmd == "AT+RST" {
		resetStatus()
		return ok
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

	// AT+CIPSTART
	if strings.HasPrefix(cmd, "AT+CIPSTART") {
		arg := getArg(&cmd)
		args := strings.Split(arg, ",")
		if status.cipmux == 0 {
			matched, err := regexp.MatchString("^(TCP|UDP|SSL),([0-9.]{7,15}),([0-9]{1,5})([,0-9]*)$", arg)
			if !matched || err != nil {
				return er
			}

			c := &connection{}
			c.t, c.ip = args[0], args[1]

			port, err := strconv.Atoi(args[2])
			if err != nil {
				return er
			}
			c.port = int16(port)

			if len(args) > 3 {
				keep, err := strconv.Atoi(args[3])
				if err != nil {
					return er
				}
				c.keep = int16(keep)
			}

			status.connections[0] = c
		} else {
			matched, err := regexp.MatchString("^[0-4]{1},(TCP|UDP|SSL),([0-9.]{7,15}),([0-9]{1,5})([,0-9]*)$", arg)
			if !matched || err != nil {
				return er
			}

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
