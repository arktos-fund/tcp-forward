package main

import (
	"io"
	"net"
	"os"

	"github.com/arktos-venture/tcp-forward/config"
	"github.com/hashicorp/go-hclog"
)

func Sock(logger hclog.Logger, cfg config.Socket) error {
	var (
		log hclog.Logger = logger.Named("sock")
		l   net.Listener
		err error
	)

	log.With("sock", cfg.Listen).Info("start sock")
	if l, err = net.Listen(cfg.Listen.Protocol, cfg.Listen.Address); err != nil {
		log.With("error", err, "sock", cfg.Listen).Error("failed to open socket")
		os.Exit(1)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			log.With("error", err, "protocol", cfg.Listen.Protocol, "address", cfg.Listen.Address).Warn("accept failed")
			continue
		}

		if err = fwd(log, cfg.Destination, conn); err != nil {
			log.With("error", err, "protocol", cfg.Destination.Protocol, "address", cfg.Destination.Address).Warn("dial failed")
		}
	}
}

func fwd(logger hclog.Logger, cfg config.SockAddress, src net.Conn) error {
	var (
		log hclog.Logger = logger.Named("fwd")
		dst net.Conn
		err error
	)

	if dst, err = net.Dial(cfg.Protocol, cfg.Address); err != nil {
		return err
	}

	log.With("client", src.RemoteAddr().String()).Info("connected")
	done := make(chan struct{})

	go func() {
		defer src.Close()
		defer dst.Close()
		io.Copy(dst, src)
		done <- struct{}{}
	}()

	go func() {
		defer src.Close()
		defer dst.Close()
		io.Copy(src, dst)
		done <- struct{}{}
	}()

	<-done
	<-done

	return nil
}
