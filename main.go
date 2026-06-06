package main

import (
	"log/slog"
	"os"

	"github.com/adarshkshitij/hivefs/p2p"
)

func main() {
	// Configure slog to use JSON handler for production-ready structured logging
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	listenAddr := os.Getenv("LISTEN_ADDR")
	if listenAddr == "" {
		listenAddr = ":3001"
	}

	opts := p2p.TCPTransportOpts{
		ListenAddr:    listenAddr,
		HandshakeFunc: p2p.NOPHandshakeFunc,
		Decoder:       p2p.DefaultDecoder{},
		OnPeer:        p2p.OnPeer,
	}

	tr := p2p.NewTCPTransport(opts)

	go func() {
		for {
			msg := <-tr.Consume()
			slog.Info("new message in channel", "payloadLen", len(msg.Payload))
		}
	}()

	slog.Info("starting node", "listenAddr", listenAddr)

	if err := tr.ListenAndAccept(); err != nil {
		slog.Error("server failed", "err", err)
		os.Exit(1)
	}

	select {}
}