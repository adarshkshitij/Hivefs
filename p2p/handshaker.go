package p2p


import "log/slog"

//Handshake Func
type HandshakeFunc func (any) error

func NOPHandshakeFunc(any) error {
	return nil
}

func OnPeer(p Peer) error {
	slog.Info("peer accepted", "peer", p)
	return nil
}