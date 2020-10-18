package console

import (
	"fmt"
	"log"

	"github.com/jorgefuertes/mister-modemu/internal/cfg"
)

// Info - Log info line
func Info(prefix string, data ...interface{}) {
	Log("info", prefix, data...)
}

// Debug - Log debug line
func Debug(prefix string, data ...interface{}) {
	Log("debug", prefix, data...)
}

// Warn - Debug warn line
func Warn(prefix string, data ...interface{}) {
	Log("warn", prefix, data...)
}

// Error - Debug error line
func Error(prefix string, data ...interface{}) {
	Log("error", prefix, data...)
}

// Log - Log line to console only, no DB
//
// Arguments:
//  - level: info, warn, error, debug
//  - prefix: Any string as log prefix
//  - data: One or more types convertibles by fmt.Sprint
func Log(level string, prefix string, data ...interface{}) {
	switch level {
	case "warn":
		prefix = fmt.Sprintf("★ [%s]", prefix)
	case "debug":
		if cfg.IsProd() {
			return
		}
		prefix = fmt.Sprintf("● [%s]", prefix)
	case "error":
		prefix = fmt.Sprintf("⚠ [%s]", prefix)
	default:
		prefix = fmt.Sprintf("ℹ [%s]", prefix)
	}
	log.SetFlags(0)
	log.Println(prefix, fmt.Sprint(data...))
}
