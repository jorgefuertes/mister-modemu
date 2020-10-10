package cfg

import (
	"sync"
)

type cfg struct {
	Env     *string
	Version *string
	Author  *string
	Port    *string
	Baud    *int
}

// Config - Main configuration
var Config *cfg
var once sync.Once
var version = "v0.1.0b"
var author = "Jorge Fuertes AKA Queru & Ramón Martinez AKA Rampa"
var port = "/dev/ttyS1"
var baud = 115200

func init() {
	once.Do(func() {
		Config = &cfg{
			Version: &version,
			Author:  &author,
			Port:    &port,
			Baud:    &baud,
		}
	})
}

// IsDev - Boolean
func IsDev() bool {
	return *Config.Env == "dev"
}

// IsProd - Boolean
func IsProd() bool {
	return *Config.Env == "prod"
}
