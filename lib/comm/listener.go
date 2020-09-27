package comm

import (
	"fmt"

	"github.com/jorgefuertes/mister-modemu/lib/console"
)

func listener(id int) {
	var err error
	var n int
	var res string
	var rBuf []byte
	rBuf = make([]byte, 1024, 2048)
	prefix := fmt.Sprintf("NET/LISTEN/%v", id)
	console.Debug(prefix, "Listening…")
	for {
		n, err = m.connections[id].conn.Read(rBuf)
		if err != nil {
			console.Debug(prefix, err)
			console.Warn(prefix, "Finished")
			break
		}
		if n > 0 {
			console.Debug(prefix, "Received ", n, " bytes")
			if m.cipmux == 0 {
				res = fmt.Sprintf("+IPD,%v:", n)
			} else {
				res = fmt.Sprintf("+IPD,%v,%v:", id, n)
			}
			serialWrite(res)
			serialWrite(rBuf)
			rBuf = nil
			serialWrite(cr, lf)
		}
	}
}
