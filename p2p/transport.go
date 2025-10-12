package p2p

import (
	"net"
)

// p2p.transport interface represents handles the communication between the nodes in the network
// ex: tcp, udp, websocket, rpc ...
type Transport interface {
	Addr() string
	Listen() error
	Dial(string) error
	Close() error
	Consume() chan Packet
	Peers() []Peer
}

// p2p.peer interface represents a peer in the network
type Peer interface {
	net.Conn
	SendPacket(Packet) error
	Receive() ([]byte, error)
	GetID() string
	GetMeta() map[string]string
	SetMeta(map[string]string)
	HandShake() error
	HandShakeWithData([]byte) error
	IsHandShake() bool
	OpenStream() (Peer, error)
	CloseStream()
}

// p2p.protocol interface represents the protocol used in the p2p network
type Protocol interface {
	IsVersionSupported() (bool, error)
	IsPortSupported() (bool, error)
	IsProtocolSupported() (bool, error)
	GetDefaultVersion() string
	GetDefaultPort() int
	GetDefaultProtocol() string
}
