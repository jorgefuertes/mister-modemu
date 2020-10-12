package main

import (
	"github.com/jorgefuertes/mister-modemu/internal/build"
	"github.com/jorgefuertes/mister-modemu/internal/cfg"
	"github.com/jorgefuertes/mister-modemu/internal/console"
	"github.com/jorgefuertes/mister-modemu/internal/modem"

	"gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	// command line flags and params
	cfg.Config.Env = kingpin.Flag(
		"environment",
		"prod or dev",
	).Short('e').Default("prod").String()
	cfg.Config.Port = kingpin.Flag(
		"port",
		"Serial port",
	).Short('p').Default("/dev/ttyS1").String()
	cfg.Config.Baud = kingpin.Flag(
		"baud",
		"Serial Speed",
	).Short('b').Default("115200").Int()

	kingpin.UsageTemplate(kingpin.CompactUsageTemplate).Version(build.Version()).Author(*cfg.Config.Author)
	kingpin.CommandLine.Help = "Mister Modem Emulator"
	kingpin.Parse()

	if cfg.IsDev() {
		console.Warn("CFG/ENV", "Development mode ON")
	} else {
		console.Info("CFG/ENV", "Production mode ON")
	}

	console.Info("CFG/PORT/BAUD", *cfg.Config.Port, " ", *cfg.Config.Baud)

	m := &modem.Modem{}
	m.Open(cfg.Config.Port, cfg.Config.Baud)
	defer m.Close()
	m.Listen()
}
