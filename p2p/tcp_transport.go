package p2p

import (
	"errors"
	"io"
	"log/slog"
	"net"

	"github.com/adarshkshitij/hivefs/metrics"
)

type TCPPeer struct {
	conn net.Conn
	outbound bool
}

type TCPTransport struct {
	TCPTransportOpts
	Listener net.Listener
	rpcch chan RPC

	// peerLock sync.RWMutex
	// peer     map[net.Addr]Peer
	OnPeer func(Peer) error
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		conn,
		outbound,
	}
}

//Auto implements close from Peer
func (p *TCPPeer) Close() error {
	return p.conn.Close()
}

type TCPTransportOpts struct {
	ListenAddr    string
	HandshakeFunc HandshakeFunc
	Decoder       Decoder
	OnPeer func(Peer) error
}

func NewTCPTransport(opts TCPTransportOpts) *TCPTransport {
	return &TCPTransport{
		TCPTransportOpts: opts,
		rpcch: make(chan RPC),
	}
}

//Consume implements the transport interface, which will return read only channel
func (t *TCPTransport) Consume() <-chan RPC  {
	return t.rpcch
}

func (t *TCPTransport) ListenAndAccept() error {
	ln, err := net.Listen("tcp", t.ListenAddr)
	if err != nil {
		return err
	}
	t.Listener = ln

	slog.Info("TCP server started", "addr", t.ListenAddr)

	go t.startAcceptLoop()
	return nil
}

func (t *TCPTransport) startAcceptLoop() {
	for {
		conn, err := t.Listener.Accept()
		if err != nil {
			slog.Error("accept loop error", "err", err)
			continue
		}

		slog.Info("new peer connected", "remoteAddr", conn.RemoteAddr())

		go t.handleConn(conn)
	}
}

func (t *TCPTransport) handleConn(conn net.Conn) {
	var err error

	// Track active peer connection
	metrics.ActivePeers.Inc()
	defer func() {
		metrics.ActivePeers.Dec()
		if err != nil && !errors.Is(err, io.EOF) {
			slog.Error("dropping peer connection", "err", err, "remoteAddr", conn.RemoteAddr())
		}
		conn.Close()
	}()

	peer := NewTCPPeer(conn, true)
	if err := t.HandshakeFunc(peer); err != nil {
		slog.Error("TCP handshake error", "err", err, "remoteAddr", conn.RemoteAddr())
		return
	}

	if t.OnPeer != nil {
		if err := t.OnPeer(peer); err != nil {
			slog.Warn("onPeer error, dropping connection", "err", err, "remoteAddr", conn.RemoteAddr())
			return
		}
	}

	rpc := &RPC{}
	for {
		if err := t.Decoder.Decode(conn, rpc); err != nil {
			if errors.Is(err, io.EOF) {
				slog.Info("connection closed by peer", "remoteAddr", conn.RemoteAddr())
				break
			}
			slog.Error("TCP message decode error", "err", err, "remoteAddr", conn.RemoteAddr())
			continue
		}
		rpc.From = conn.RemoteAddr()
		t.rpcch <- *rpc

		// Track message and bytes transferred
		metrics.MessagesReceivedTotal.Inc()
		metrics.BytesTransferredTotal.Add(float64(len(rpc.Payload)))

		slog.Info("message received", "from", rpc.From, "payloadLen", len(rpc.Payload))
	}
}
