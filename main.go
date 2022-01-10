package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/arktos-venture/tcp-forward/config"
	"github.com/hashicorp/go-hclog"
	"golang.org/x/sync/errgroup"
)

var (
	logger  hclog.Logger   = hclog.Default()
	cfg     *config.Config = config.DefaultConfig
	Version string         = "unknown"
)

func init() {
	var err error

	flag.StringVar(&cfg.LogOptions.Level, "log.level", cfg.LogOptions.Level, "Log level values allowed [trace, debug, info, warn, error, fatal]")
	flag.StringVar(&cfg.LogOptions.Format, "log.fmt", cfg.LogOptions.Format, "Log format values allowed [logfmt, json]")
	flag.StringVar(&cfg.ConfigFile, "config", cfg.ConfigFile, "Config file name")

	flag.Parse()

	if cfg.ConfigFile != "" {
		if err = cfg.Reload(logger); err != nil {
			logger.With("error", err).Debug("failed to load config file")
		}
	}

	if (cfg.Socket.Listen.Address == "" || cfg.Socket.Listen.Protocol == "") ||
		(cfg.Socket.Destination.Address == "" || cfg.Socket.Destination.Protocol == "") {
		logger.Error("listen or destination is empty")
		os.Exit(1)
	}
}

func main() {
	var (
		ctx    context.Context = context.Background()
		g, _                   = errgroup.WithContext(ctx)
		stopCh                 = make(chan os.Signal)
		log    hclog.Logger    = cfg.LogOptions.LogFlagParse("main")
	)

	// Catch OS signal
	signal.Notify(stopCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-stopCh
		logger.Error("SIGTERM signal")
		os.Exit(1)
	}()

	log.With("version", Version).Info("starting")

	// Prometheus metrics & debug http server
	g.Go(func() error { return HTTP(log, cfg.Health) })
	g.Go(func() error { return Sock(log, cfg.Socket) })

	g.Wait()
}
