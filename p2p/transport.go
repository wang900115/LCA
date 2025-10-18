package p2p

import (
	"context"
	"net"

	"github.com/wang900115/LCA/crypt/did"
	"github.com/wang900115/LCA/p2p/network"
)

// p2p.Transport interface represents handles the communication between the nodes in the network
// ex: tcp, udp, websocket, rpc ...
type Transport interface {
	Addr() string
	Listen() error
	Dial(string) error
	Close() error
	Consume() chan network.Packet
}

// p2p.Peer interface represents a peer in the network
type Peer interface {
	net.Conn
	ID() string
	Document() *did.DIDDocument
	ProtocolInfo() *network.ProtocolInfo
	Send(network.Packet) error
	Receive() (<-chan network.Packet, error)
	Peers() map[string]Peer
	// HandShake() error
	ReadPump(context.Context)
	WritePump(context.Context)
}
