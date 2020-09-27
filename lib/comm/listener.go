package comm

import (
	"fmt"

	"github.com/jorgefuertes/mister-modemu/lib/console"
)

func listener(id int) {
	var err error
	var n int
	var res string
	var netBuf []byte
	netBuf = make([]byte, 2048, 2048)
	prefix := fmt.Sprintf("NET/LISTEN/%v", id)
	console.Debug(prefix, "Listening…")
	for {
		n, err = m.connections[id].conn.Read(netBuf)
		if err != nil {
			console.Debug(prefix, err)
			console.Warn(prefix, "Finished")
			break
		}
		if n > 0 {
			cut := netBuf[0:n]
			console.Debug(prefix, "Received ", n, " bytes")
			console.Debug(prefix, cut[0], "…", cut[n-1])
			if m.cipmux == 0 {
				res = fmt.Sprintf("+IPD,%v:", n)
			} else {
				res = fmt.Sprintf("+IPD,%v,%v:", id, n)
			}
			serialWrite(res)
			serialWriteBytes(&cut)
			serialWrite(cr, lf)
			console.Debug(prefix, "EOT")
		}
	}
}
