package p2p

import "net"

//RPC holds the arbitary data that is being sent over the each
//tcp between 2 nodes in network
type RPC struct {
	Payload []byte
	From  net.Addr
	
}