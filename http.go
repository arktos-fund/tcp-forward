package main

import (
	"net/http"
	"time"

	"github.com/arktos-venture/tcp-forward/config"
	"github.com/hashicorp/go-hclog"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func healthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func HTTP(logger hclog.Logger, c config.Health) error {
	var (
		log    hclog.Logger = logger.Named("http")
		server *http.Server = &http.Server{
			Addr:        c.Listen,
			ReadTimeout: c.Timeout,
			IdleTimeout: 3 * time.Second,
		}
		err error = nil
	)

	// metrics & healthz
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/healthz", healthz)

	log.With("address", c.Listen).Info("start webserver")
	if err = server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}
