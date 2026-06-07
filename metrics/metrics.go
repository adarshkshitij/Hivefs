package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// ActivePeers tracks the number of currently connected peers.
	ActivePeers = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: "hivefs",
		Name:      "active_peers",
		Help:      "Number of currently connected peer nodes.",
	})

	// BytesTransferredTotal tracks the total bytes received over the network.
	BytesTransferredTotal = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "hivefs",
		Name:      "bytes_transferred_total",
		Help:      "Total bytes received over TCP from peers.",
	})

	// MessagesReceivedTotal tracks the total number of RPC messages decoded.
	MessagesReceivedTotal = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "hivefs",
		Name:      "messages_received_total",
		Help:      "Total number of RPC messages received and decoded.",
	})

	// BytesStoredTotal tracks the total bytes written to local disk storage.
	BytesStoredTotal = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "hivefs",
		Name:      "bytes_stored_total",
		Help:      "Total bytes written to local disk storage.",
	})
)
