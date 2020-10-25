package ascii

import "fmt"

// ByteToStr - Transform byte into readable string
func ByteToStr(b byte) string {
	switch b {
	case LF:
		return `[lf]`
	case CR:
		return `[cr]`
	case BS:
		return `[bs]`
	case DEL:
		return `[dlt]`
	default:
		if b >= 32 && b < 127 {
			return string(b)
		}
		return fmt.Sprintf("[%x]", b)
	}
}
