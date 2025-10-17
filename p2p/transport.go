package p2p

import (
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
	Consume() chan Packet
}

// p2p.Peer interface represents a peer in the network
type Peer interface {
	net.Conn
	ID() string
	Document() *did.DIDDocument
	ProtocolInfo() *network.ProtocolInfo
	Send(Packet) error
	Receive() (<-chan Packet, error)
	Peers() map[string]Peer
	// HandShake() error
}

// p2p.Packet interface represents a packet in the network
type Packet interface {
	GetCommand() byte
	GetLength() uint32
	GetPayload() []byte
	GetCheckSum() uint32
	Encode() ([]byte, error)
}

// p2p.RPC interface represents the rpc mechanism used in the p2p network
type RPC interface {
	Call(method string, args interface{}, reply interface{}) error
	Encode(interface{}) (Packet, error)
}
