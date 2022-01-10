package config

import (
	"errors"
	"io/ioutil"
	"os"
	"regexp"

	"github.com/hashicorp/go-hclog"
	"github.com/imdario/mergo"
	"gopkg.in/yaml.v2"
)

var (
	levels *regexp.Regexp = regexp.MustCompile("^(trace|debug|info|warn|error|fatal)$")
	format *regexp.Regexp = regexp.MustCompile("^(json)$")

	// DefaultConfig values for CLI
	DefaultConfig = &Config{
		ConfigFile: "/config.yaml",
		LogOptions: LogOptions{
			Level:  "info",
			Format: "logfmt",
		},
		Socket: Socket{},
	}
)

// Reload configuration
func (c *Config) Reload(logger hclog.Logger) error {
	var (
		v   *Config = DefaultConfig
		err error
	)

	// Parse config file if needed
	if c, err = loadFile(logger, c.ConfigFile); err != nil {
		return err
	}

	// Merge overwritting
	if err = mergo.Merge(v, c, mergo.WithOverride); err != nil {
		return err
	}

	return nil
}

// LogFlagParse level logs
func (c *LogOptions) LogFlagParse(name string) hclog.Logger {
	var level string

	if levels.MatchString(c.Level) {
		level = c.Level
	} else {
		level = "INFO"
	}

	return hclog.New(&hclog.LoggerOptions{
		Name:       name,
		Level:      hclog.LevelFromString(level),
		JSONFormat: format.MatchString(c.Format),
		Output:     os.Stdout,
	})
}

// LoadFile parses the given YAML file into a Config.
func loadFile(logger hclog.Logger, filename string) (*Config, error) {
	var (
		cfg     *Config = DefaultConfig
		content []byte
		err     error
	)

	if filename != "" {
		if content, err = ioutil.ReadFile(filename); err != nil {
			return nil, err
		}

		if err = yaml.Unmarshal(content, cfg); err != nil {
			return nil, err
		}

		return cfg, nil
	}

	return nil, errors.New("filename is empty")
}
