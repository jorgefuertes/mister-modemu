package at

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/jorgefuertes/mister-modemu/internal/build"
	"github.com/jorgefuertes/mister-modemu/internal/cfg"
	"github.com/jorgefuertes/mister-modemu/internal/inet"
	"github.com/jorgefuertes/mister-modemu/internal/modem"
	"github.com/tatsushid/go-fastping"
)

// Esp8266 - original AT commands
func Esp8266(m *modem.Status) {
	// AT
	m.Parser.AT("AT", func(m *modem.Status) {
		m.OK()
	})

	// AT+VERSION
	m.Parser.AT("AT+VERSION", func(m *modem.Status) {
		m.WriteLn("+VERSION:" + build.Version())
		m.OK()
	})

	// AT+AUTHOR
	m.Parser.AT("AT+AUTHOR", func(m *modem.Status) {
		m.WriteLn("+AUTHOR:" + *cfg.Config.Author)
		m.OK()
	})

	// AT+RST
	m.Parser.AT("AT+RST", func(m *modem.Status) {
		m.Reset()
		m.OK()
	})

	// AT+HELP
	m.Parser.AT("AT+HELP", func(m *modem.Status) {
		for _, line := range help {
			m.WriteLn(line)
		}
		m.OK()
	})

	// AT+HELP GAMES
	m.Parser.AT("AT+HELP GAMES", func(m *modem.Status) {
		for _, line := range helpGames {
			m.WriteLn(line)
		}
		m.OK()
	})

	// AT+LIST GAMES
	m.Parser.AT("AT+LIST GAMES", func(m *modem.Status) {
		for _, line := range listGames {
			m.WriteLn(line)
		}
		m.OK()
	})

	// AT+CIPSTATUS
	m.Parser.AT("AT+CIPSTATUS", func(m *modem.Status) {
		m.WriteLn(fmt.Sprintf("STATUS:%v", m.Sta))
		for _, c := range m.Connections {
			if c != nil {
				m.WriteLn(fmt.Sprintf("\r\n+CIPSTATUS:%v,%s,%s,%v,%v,%v",
					c.ID, c.T, c.IP, c.Port, 0, c.Cs))
			}
		}
	})

	// AT+CIPDOMAIN
	m.Parser.AT("AT+DOMAIN=*", func(m *modem.Status) {
		domain := m.Parser.GetArg()
		if domain == "" || len(domain) < 3 {
			m.Parser.Error("Domain too short")
			m.ERR()
		}
		ips, err := net.LookupIP(domain)
		if err != nil {
			m.WriteLn("DNS Fail")
			m.Parser.Error(err)
			m.ERR()
		}
		m.WriteLn(fmt.Sprintf("+CIPDOMAIN:%s", ips[0]))
		m.OK()
	})

	// ATE - echo on/off
	m.Parser.AT("ATE*", func(m *modem.Status) {
		m.Ate = (m.Parser.GetArg() == "1")
		m.OK()
	})

	// AT+GMR
	m.Parser.AT("AT+GMR", func(m *modem.Status) {
		m.WriteLn(build.ATVersion())
		m.WriteLn(build.SDKVersion())
		m.WriteLn(build.CompileTime())
		m.WriteLn(build.BinVersion())
		m.OK()
	})

	// AT+CWJAP?, AT+CWJAP_CUR?, AT+CWJAP_DEF?
	m.Parser.AT("AT+CWJAP*?", func(m *modem.Status) {
		// +CWJAP_CUR:<ssid>,<bssid>,<channel>,<rssi> OK
		m.WriteLn(fmt.Sprintf("+CWJAP_CUR:%s,%v,%d,%d",
			inet.GetLocalInterface().Name, inet.GetLocalMac(), 11, -34))
		m.OK()
	})

	// AT+CIPSTA?, AT+CIPSTA_CUR?, AT+CIPSTA_DEF?
	m.Parser.AT("AT+CIPSTA*?", func(m *modem.Status) {
		m.WriteLn("+CIPSTA_CUR:ip:", inet.GetLocalIP().String())
		m.WriteLn("+CIPSTA_CUR:gateway:", inet.GetGateway().String())
		m.WriteLn("+CIPSTA_CUR:netmask:", inet.GetLocalMask())
		m.OK()
	})

	// AT+CIPMUX
	m.Parser.AT("AT+CIPMUX=*", func(m *modem.Status) {
		m.CipMux = (m.Parser.GetArg() == "1")
		m.OK()
	})

	// AT+CIFSR - Gets the local IP address
	m.Parser.AT("AT+CIFSR", func(m *modem.Status) {
		m.WriteLn("+CIFSR:APIP,", inet.GetLocalIP().String())
		m.WriteLn("+CIFSR:APIP,", inet.GetLocalIP().String())
		m.WriteLn("+CIFSR:APMAC,", inet.GetLocalMac().String())
		m.WriteLn("+CIFSR:STAIP,", inet.GetLocalIP().String())
		m.WriteLn("+CIFSR:STAMAC,", inet.GetLocalMac().String())
		m.OK()
	})

	// AT+CIPDINFO=<0|1>
	m.Parser.AT("AT+CIPDINFO=*", func(m *modem.Status) {
		m.CipInfo = (m.Parser.GetArg() == "1")
		m.OK()
	})

	// AT+CIPDNS*?
	m.Parser.AT("AT+CIPDNS*?", func(m *modem.Status) {
		m.WriteLn("+CIPDNS_CUR:", inet.GetNameServer().String())
		m.WriteLn("+CIPDNS_CUR:", inet.GetNameServer().String())
		m.OK()
	})

	// AT+PING=<IP>
	m.Parser.AT("AT+PING=*", func(m *modem.Status) {
		host := m.Parser.GetArg()

		p := fastping.NewPinger()
		ra, err := net.ResolveIPAddr("ip4:icmp", host)
		if err != nil {
			m.Parser.Error(err)
			m.ERR()
			return
		}

		p.AddIPAddr(ra)
		p.Network("udp")
		p.OnRecv = func(addr *net.IPAddr, rtt time.Duration) {
			m.WriteLn("+", rtt.String())
			m.OK()
			p.Stop()
		}

		err = p.Run()
		if err != nil {
			m.Parser.Error(err)
			m.ERR()
		}
	})

	// AT+CWMODE_ (Query)
	m.Parser.AT("AT+CWMODE_*?", func(m *modem.Status) {
		m.WriteLn("+CWMODE_CUR:", m.Cw)
		m.OK()
	})

	// AT+CWMODE_SET, AT+CWMODE_CUR
	m.Parser.AT("AT+CWMODE_*=*", func(m *modem.Status) {
		// Set
		arg := m.Parser.GetArg()
		mode, err := strconv.Atoi(arg)
		if err != nil {
			m.ERR()
			m.Parser.Error("err")
			return
		}
		m.Cw = uint8(mode)
		m.OK()
	})

	// AT+CIPSTART
	m.Parser.AT("AT+CIPSTART=*", func(m *modem.Status) {
		var id int
		var t string
		var ip string
		var port int
		var keep int
		var err error

		next := 0
		if !m.CipMux {
			// One conn, and one less argument
			id = 0
		} else {
			// multiple conn
			id, err = strconv.Atoi(m.Parser.GetArgs()[0])
			if err != nil {
				m.Parser.Error("Invalid connection ID")
				m.ERR()
				return
			}
			next++
		}

		// type
		t = strings.ToUpper(m.Parser.GetArgs()[next])
		if t != "TCP" && t != "UDP" {
			m.Parser.Error("Unimplemented conn type")
			m.ERR()
			return
		}
		next++

		// remote IP
		ip = m.Parser.GetArgs()[next]
		next++

		// port
		port, err = strconv.Atoi(m.Parser.GetArgs()[next])
		if err != nil {
			m.Parser.Error("Invalid port")
			m.ERR()
			return
		}
		next++

		// keep alive
		if len(m.Parser.GetArgs()) > next {
			keep, err = strconv.Atoi(m.Parser.GetArgs()[next])
			if err != nil {
				m.Parser.Error("Invalid keep alive")
				m.ERR()
				return
			}
		}

		if err = m.NewConn(uint8(id), t, ip, port, keep); err != nil {
			m.Parser.Error(err)
			if err.Error() == "Connection already in use" {
				m.WriteLn("ALREADY CONNECTED")
				return
			}
			m.ERR()
			return
		}

		// Connect
		if err = m.Connect(uint8(id)); err != nil {
			m.Parser.Error(err)
			m.ERR()
			return
		}

		m.OK()
	})

	// AT+CIPSEND
	m.Parser.AT("AT+CIPSEND", func(m *modem.Status) {
		m.CipPacketOn()
		m.OK()
	})

	// AT+CIPSEND=<params>
	m.Parser.AT("AT+CIPSEND=*", func(m *modem.Status) {
		var id int
		var err error
		var clen int
		args := m.Parser.GetArgs()

		if !m.CipMux {
			id = 0
			clen, err = strconv.Atoi(args[0])
			if err != nil {
				m.Parser.Error("Invalid param length")
				m.ERR()
				return
			}
		} else {
			if len(args) < 2 {
				m.Parser.Error("Missing link_id")
				m.ERR()
				return
			}

			id, err = strconv.Atoi(args[0])
			if err != nil {
				m.Parser.Error("Invalid link_id")
				m.ERR()
				return
			}

			clen, err = strconv.Atoi(args[1])
			if err != nil {
				m.Parser.Error("Invalid param length")
				m.ERR()
				return
			}
		}

		m.CipSet(uint8(id), uint(clen))
		m.OK()
	})

	// AT+CIPCLOSE
	m.Parser.AT("AT+CIPCLOSE=*", func(m *modem.Status) {
		if !m.CipMux {
			m.CipClose(0)
			m.CipClearAll()
		} else {
			id, err := strconv.Atoi(m.Parser.GetArg())
			if err != nil {
				m.Parser.Error(err)
				m.ERR()
				return
			}
			m.CipClose(uint8(id))
			m.CipClear(uint8(id))
		}

		m.OK()
	})
}
