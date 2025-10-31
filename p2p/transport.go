package p2p

import (
	"context"
)

// p2p.Transport interface represents handles the communication between the nodes in the network
// ex: tcp, udp, websocket, rpc ...
type Transport interface {
	Addr() string
	Listen(context.Context) error
	Dial(context.Context, string) error
	Close() error
	// AddOutPeer(Peer) error
	// AddInPeer(Peer) error
	// RemoveOutPeer(Peer)
	// RemoveInPeer(Peer)
	// Peers() map[string]Peer
}

// p2p.Peer interface represents a peer in the network
// type Peer interface {
// 	net.Conn
// 	Addr() string
// 	ID() string
// 	Document() *did.Document
// 	ProtocolInfo() *network.ProtocolInfo
// 	Send(network.Packet) error
// 	Receive() (<-chan network.Packet, error)
// 	ReadPump(context.Context)
// 	WritePump(context.Context)
// }
