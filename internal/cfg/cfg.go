package cfg

import (
	"sync"
)

type cfg struct {
	Env    *string
	Author *string
	Port   *string
	Baud   *int
}

// Config - Main configuration
var Config *cfg
var once sync.Once
var author = "Jorge Fuertes AKA Queru & Ramón Martinez AKA Rampa"
var port = "/dev/ttyS1"
var testPort = "/dev/ttyp5"
var baud = 115200

// TestPort2 - connect tester to here
var TestPort2 = "/dev/ptyp5"

func init() {
	once.Do(func() {
		Config = &cfg{
			Author: &author,
			Port:   &port,
			Baud:   &baud,
		}
	})
}

// TestInit - init test environment
func TestInit() {
	Config.Port = &testPort
}

// IsDev - Boolean
func IsDev() bool {
	return *Config.Env == "dev"
}

// IsTest - Boolean
func IsTest() bool {
	return *Config.Env == "test"
}

// IsProd - Boolean
func IsProd() bool {
	return *Config.Env == "prod"
}

// Banner - Application banner
const Banner = `
 ______         __               ________ __ __     __
|   __ \.-----.|  |_.----.-----.|  |  |  |__|  |--.|__|
|      <|  -__||   _|   _|  _  ||  |  |  |  |    < |  |
|___|__||_____||____|__| |_____||________|__|__|__||__|
_______________________________________________________

      ESP8266 AT Modem Emulator for ZX-Next core
_______________________________________________________
`
