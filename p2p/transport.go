package p2p

import "net"

// p2p.Peer  interface represents the remote node
type Peer interface {
	net.Conn
	Send([]byte) error
	CloseStream()
}

// p2p.transport interface represents handles the communication between the nodes in the network
// ex: tcp, udp, websocket, rpc ...
type Transport interface {
	Addr() string
	Dial(string) error
	ListenAndAccep() error
	Consume()
	Close() error
}
