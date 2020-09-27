package util

// ascii
const cr = 0x0D
const lf = 0x0A
const del = 0x7F
const bs = 0x08

// ByteToStr - Transform byte into readable string
func ByteToStr(b byte) string {
	switch b {
	case lf:
		return `LF`
	case cr:
		return `CR`
	case bs:
		return `BS`
	case del:
		return `DEL`
	default:
		if b >= 32 && b < 127 {
			return string(b)
		}
		return "•"
	}
}
