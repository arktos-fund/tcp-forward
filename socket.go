package main

import (
	"io"
	"net"
	"os"

	"github.com/arktos-venture/tcp-forward/config"
	"github.com/hashicorp/go-hclog"
)

func Sock(logger hclog.Logger, c config.Socket) error {
	var (
		log hclog.Logger = logger.Named("sock")
		l   net.Listener
		err error
	)

	log.With("sock", c.Listen).Info("start sock")
	if l, err = net.Listen(c.Listen.Protocol, c.Listen.Address); err != nil {
		log.With("error", err, "sock", c.Listen).Error("failed to open socket")
		os.Exit(1)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			log.With("error", err, "protocol", c.Listen.Protocol, "address", c.Listen.Address).Warn("accept failed")
			continue
		}

		log.With("client", conn.RemoteAddr().String()).Info("connected")
		go fwd(log, c.Destination, conn)
	}
}

func fwd(logger hclog.Logger, cfg config.SockAddress, src net.Conn) error {
	var (
		log hclog.Logger = logger.Named("fwd")
		dst net.Conn
		err error
	)

	if dst, err = net.Dial(cfg.Protocol, cfg.Address); err != nil {
		log.With("error", err, "protocol", cfg.Protocol, "address", cfg.Address).Error("dial failed")
		os.Exit(1)
	}

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
