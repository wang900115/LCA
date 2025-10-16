package p2p

import (
	"net"

	"github.com/wang900115/LCA/crypt/did"
)

// p2p.Transport interface represents handles the communication between the nodes in the network
// ex: tcp, udp, websocket, rpc ...
type Transport interface {
	Addr() string
	Listen() error
	Dial(string) error
	Close() error
	Consume() chan Packet
	Peers() []Peer
}

// p2p.Peer interface represents a peer in the network
type Peer interface {
	net.Conn
	SendPacket(Packet) error
	ReceivePacket() (Packet, error)
	GetID() string
	GetDocument() *did.DIDDocument
	GetMeta() map[string]string
	SetMeta(map[string]string)
	HandShake() error
	HandShakeWithData([]byte) error
	IsHandShake() bool
	OpenStream() (Peer, error)
	CloseStream()
	IsStream() bool
	WaitSream()
	IsOutBound() bool
	SetProtocol(Protocol)
	GetProtocol() Protocol
	Consume() <-chan RPC
}

// p2p.Protocol interface represents the protocol used in the p2p network
type Protocol interface {
	IsVersionSupported() (bool, error)
	IsPortSupported() (bool, error)
	IsProtocolSupported() (bool, error)
	GetDefaultVersion() string
	GetDefaultPort() int
	GetDefaultProtocol() string
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
