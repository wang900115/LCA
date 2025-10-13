package transport

import (
	"errors"
	"fmt"
	"log"
	"net"

	"github.com/wang900115/LCA/p2p"
	"github.com/wang900115/LCA/p2p/node"
)

// TCPTransportOpts holds configuration options for the TCPTransport.
type TCPTransportOpts struct {
	ListenAddr    string
	HandShakeFunc p2p.HandShakeFunc
	OnPeer        func(p2p.Peer) error
}

// TCPTransport implements a TCP-based transport layer for P2P communication.
type TCPTransport struct {
	TCPTransportOpts
	listener net.Listener
	ch       chan p2p.Packet
}

// NewTCPTransport creates a new TCPTransport with the given options.
func NewTCPTransport(opts TCPTransportOpts) *TCPTransport {
	return &TCPTransport{
		TCPTransportOpts: opts,
		ch:               make(chan p2p.Packet, 1024),
	}
}

// Addr returns the listening address of the TCPTransport.
func (t *TCPTransport) Addr() string {
	return t.ListenAddr
}

// Consume returns a channel to receive incoming packets.
func (t *TCPTransport) Consume() <-chan p2p.Packet {
	return t.ch
}

// Close shuts down the TCP listener.
func (t *TCPTransport) Close() error {
	return t.listener.Close()
}

// Dial connects to a remote TCP address and starts handling the connection.
func (t *TCPTransport) Dial(addr string) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}
	go t.handleConn(conn, true)
	return nil
}

// Listen starts the TCP listener and begins accepting incoming connections.
func (t *TCPTransport) Listen() error {
	var err error
	t.listener, err = net.Listen("tcp", t.ListenAddr)
	if err != nil {
		return err
	}

	go t.startAcceptLoop()
	return nil
}

// startAcceptLoop continuously accepts incoming connections.
func (t *TCPTransport) startAcceptLoop() {
	for {
		conn, err := t.listener.Accept()
		if errors.Is(err, net.ErrClosed) {
			return
		}
		if err != nil {
			fmt.Printf("TCP accept error: %s\n", err)
		}
		go t.handleConn(conn, false)
	}
}

// handleConn performs the handshake and processes incoming packets for a connection.
func (t *TCPTransport) handleConn(conn net.Conn, outBound bool) {
	var err error
	defer func() {
		conn.Close()
	}()

	peer := node.NewPeer(conn, outBound)
	if err = t.HandShakeFunc(peer); err != nil {
		return
	}
	if err = peer.HandShake(); err != nil {
		return
	}
	if t.OnPeer(peer) != nil {
		if err = t.OnPeer(peer); err != nil {
			return
		}
	}
	for {
		buf := make([]byte, 4096)
		n, err := conn.Read(buf)
		if err != nil {
			return
		}
		pk, err := p2p.Decode2Packet(buf[:n])
		if err != nil {
			return
		}
		if peer.IsStream() {
			if p, err := peer.OpenStream(); err != nil {
				log.Printf("open stream error with peer: %+v\n", p)
				return
			}
			fmt.Printf("[%s] incoming stream, waiting...\n", conn.RemoteAddr())
			peer.WaitSream()
			fmt.Printf("[%s] stream closed, resuming read loop\n", conn.RemoteAddr())
			continue
		}
		t.ch <- pk
	}
}
