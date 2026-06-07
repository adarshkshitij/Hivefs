package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus/promhttp"

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

	metricsAddr := os.Getenv("METRICS_ADDR")
	if metricsAddr == "" {
		metricsAddr = ":9090"
	}

	// Start Prometheus metrics HTTP server in a goroutine
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		slog.Info("metrics server started", "addr", metricsAddr)
		if err := http.ListenAndServe(metricsAddr, nil); err != nil {
			slog.Error("metrics server failed", "err", err)
		}
	}()

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

	slog.Info("starting node", "listenAddr", listenAddr, "metricsAddr", metricsAddr)

	if err := tr.ListenAndAccept(); err != nil {
		slog.Error("server failed", "err", err)
		os.Exit(1)
	}

	select {}
}