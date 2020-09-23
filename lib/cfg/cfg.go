package cfg

import (
	"sync"
)

type cfg struct {
	Env     *string
	Version *string
	Author  *string
	Serial  *struct {
		Port *string
	}
}

// Config - Main configuration
var Config cfg
var once sync.Once

func init() {
	once.Do(func() {
		*Config.Version = "v0.1.0b"
		*Config.Author = "Jorge Fuertes AKA Queru & Ramón Martinez AKA Rampa"
		*Config.Serial.Port = "/dev/ttyS2"
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
