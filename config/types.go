package config

import "time"

// LogOptions struct parameters for logging level and format
type LogOptions struct {
	Level  string `yaml:"level,omitempty"`
	Format string `yaml:"format,omitempty"`
}

type Socket struct {
	Listen      string        `yaml:"listen,omitempty"`
	Destination string        `yaml:"destination,omitempty"`
	Timeout     time.Duration `yaml:"timeout,omitempty"`
}

type Health struct {
	Listen  string        `yaml:"listen,omitempty"`
	Timeout time.Duration `yaml:"timeout,omitempty"`
}

// Config CLI structs
type Config struct {
	ConfigFile string
	LogOptions LogOptions `yaml:"logOptions"`
	Health     Health
	Socket     Socket `yaml:"socket"`

	// Catches all undefined fields and must be empty after parsing.
	XXX map[string]interface{} `yaml:",inline"`
}
