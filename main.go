package main

import (
	"github.com/jorgefuertes/mister-modemu/lib/cfg"
	"github.com/jorgefuertes/mister-modemu/lib/console"
	"gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	// command line flags and params
	cfg.Config.Env = kingpin.Flag(
		"environment",
		"prod or dev",
	).Short('e').Default("prod").String()
	cfg.Config.Serial.Port = kingpin.Flag(
		"port",
		"Serial port",
	).Short('p').Default("/dev/ttyS2").String()

	kingpin.UsageTemplate(kingpin.CompactUsageTemplate).Version(*cfg.Config.Version).Author(*cfg.Config.Author)
	kingpin.CommandLine.Help = "Mister Modem Emulator"
	kingpin.Parse()

	if cfg.IsDev() {
		console.Warn("CFG/ENV", "Development mode ON")
	} else {
		console.Info("CFG/ENV", "Production mode ON")
	}

	console.Info("CFG/PORT", *cfg.Config.Serial.Port)
}
