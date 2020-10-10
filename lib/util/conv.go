package util

import (
	"fmt"
)

// ascii
const cr = 0x0D
const lf = 0x0A
const del = 0x7F
const bs = 0x08

// ByteToStr - Transform byte into readable string
func ByteToStr(b byte) string {
	switch b {
	case lf:
		return `[lf]`
	case cr:
		return `[cr]`
	case bs:
		return `[bs]`
	case del:
		return `[dlt]`
	default:
		if b >= 32 && b < 127 {
			return string(b)
		}
		return fmt.Sprintf("[%x]", b)
	}
}

// BufToStr - Transform byte buffer to trimmed string
func BufToStr(b *[]byte, n int) string {
	var s string
	for i, v := range *b {
		if v > 31 && v < 126 {
			s += string(v)
		}
		if i == n-1 {
			break
		}
	}
	return s
}

// BufToDebug - Buffer to debug string
func BufToDebug(b []byte, n int) string {
	var s string
	for i := 0; i < n; i++ {
		s += ByteToStr(b[i])
	}
	return s
}
