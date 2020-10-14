package main

import (
	"fmt"
	"os"

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
	kingpin.CommandLine.Help = "RetroWiki Mister ESP8266 AT Modem Emulator"
	v := kingpin.Flag("short_version", "Show short versión string and exit").Short('v').Default("false").Bool()
	kingpin.Parse()

	if *v {
		fmt.Println(build.VersionShort())
		os.Exit(0)
	}

	fmt.Println(cfg.Banner)
	fmt.Println(build.Version() + "\n")

	if cfg.IsDev() {
		console.Warn("CFG/ENV", "Development mode ON")
	} else {
		console.Info("CFG/ENV", "Production mode ON")
	}

	console.Info("CFG/PORT/BAUD", *cfg.Config.Port, " ", *cfg.Config.Baud)

	m := &modem.Modem{}
	if err := m.Open(cfg.Config.Port, cfg.Config.Baud); err != nil {
		console.Error("SER/OPEN", err.Error())
		console.Error("SER/OPEN", "Cannot open serial port!")
		os.Exit(1)
	}
	defer m.Close()
	m.Listen()
}
