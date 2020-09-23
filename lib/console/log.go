package console

import (
	"fmt"
	"log"

	"github.com/fatih/color"
	"github.com/jorgefuertes/mister-modemu/cfg"
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
	var red = color.New(color.FgRed).SprintFunc()
	var cyan = color.New(color.FgCyan).SprintFunc()
	var magenta = color.New(color.FgMagenta).SprintFunc()
	var yellow = color.New(color.FgYellow).SprintFunc()

	switch level {
	case "warn":
		log.Println(yellow("★ ", "[", prefix, "]"), fmt.Sprint(data...))
	case "debug":
		if cfg.IsDev() {
			log.Println(magenta("● ", "[", prefix, "]"), fmt.Sprint(data...))
		}
	case "error":
		log.Println(red("⚠ ", "[", prefix, "]"), fmt.Sprint(data...))
	default:
		log.Println(cyan("ℹ ", "[", prefix, "]"), fmt.Sprint(data...))
	}
}
