package comm

import (
	"fmt"
	"net"
	"regexp"
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

func parseCmd(buf []byte) string {
	n := len(buf) - 2
	cmd := string(buf[:n])
	console.Debug("COMM/PARSE", cmd)

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
