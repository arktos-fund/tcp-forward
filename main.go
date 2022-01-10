package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/hashicorp/go-hclog"
	flag "github.com/spf13/pflag"

	conf "github.com/arktos-venture/grpc-apis/config"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	cfg *conf.Bootstrap = conf.DefaultConfig
	// Name is the name of the compiled software.
	Name string
	// Version is the version of the compiled software.
	Version string

	id, _ = os.Hostname()
)

func init() {
	flag.StringVar(&cfg.LogOptions.Level, "log.level", cfg.LogOptions.Level, "Log level values allowed [trace, debug, info, warn, error, fatal]")
	flag.StringVar(&cfg.LogOptions.Format, "log.fmt", cfg.LogOptions.Format, "Log format values allowed [logfmt, json]")
	flag.StringVar(&cfg.ConfigFile, "conf", cfg.ConfigFile, "config path, eg: -conf config.yaml")

	flag.Parse()
}

func main() {
	var (
		log hclog.Logger = conf.LogFlagParse("Example", cfg.LogOptions).With(
			"service.id", id,
			"service.version", Version,
		)
		c config.Config = config.New(
			config.WithSource(
				file.NewSource(cfg.GetConfigFile()),
			),
		)
		stopCh chan os.Signal = make(chan os.Signal)
		bc     conf.Bootstrap
		err    error
	)

	if err = c.Load(); err != nil {
		log.With("error", err).Error("failed to config file")
		os.Exit(1)
	}

	if err = c.Scan(&bc); err != nil {
		log.With("error", err).Error("failed to scan file")
		os.Exit(1)
	}

	// Catch OS signal
	signal.Notify(stopCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-stopCh
		log.Error("SIGTERM signal")
		os.Exit(1)
	}()

	log.With("version", Version).Info("starting")
}
