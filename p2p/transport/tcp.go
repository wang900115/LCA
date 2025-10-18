package transport

import (
	"context"
	"errors"
	"fmt"
	"net"
	"sync"

	"github.com/wang900115/LCA/p2p"
	"github.com/wang900115/LCA/p2p/network"
	"github.com/wang900115/LCA/p2p/node"
)

// TCPTransportOpts holds configuration options for the TCPTransport.
type TCPTransportOpts struct {
	ListenAddr string
	HandShake  func(p2p.Peer) error
	InBoundLi  int
	OutBoundLi int
}

// TCPTransport implements a TCP-based transport layer for P2P communication.
type TCPTransport struct {
	TCPTransportOpts
	listener net.Listener
	State    *state
}

// NewTCPTransport creates a new TCPTransport with the given options.
func NewTCPTransport(opts TCPTransportOpts) p2p.Transport {
	return &TCPTransport{
		TCPTransportOpts: opts,
		State:            NewState(opts.OutBoundLi, opts.InBoundLi),
	}
}

// Addr is dial caller address
func (t *TCPTransport) Addr() string {
	return t.ListenAddr
}

// Close shuts down the TCP listener.
func (t *TCPTransport) Close() error {
	return t.listener.Close()
}

// Dial connects to a remote TCP address and starts handling the connection.
func (t *TCPTransport) Dial(ctx context.Context, addr string) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}
	go t.handleConn(ctx, conn, true)
	return nil
}

// Listen starts the TCP listener and begins accepting incoming connections.
func (t *TCPTransport) Listen(ctx context.Context) error {
	var err error
	t.listener, err = net.Listen("tcp", t.ListenAddr)
	if err != nil {
		return err
	}
	go t.startAcceptLoop(ctx)
	return nil
}

// startAcceptLoop continuously accepts incoming connections.
func (t *TCPTransport) startAcceptLoop(ctx context.Context) {
	for {
		conn, err := t.listener.Accept()
		if errors.Is(err, net.ErrClosed) {
			return
		}
		if err != nil {
			fmt.Printf("TCP accept error: %s\n", err)
		}
		go t.handleConn(ctx, conn, false)
	}
}

// handleConn performs the handshake and processes incoming packets for a connection.
func (t *TCPTransport) handleConn(ctx context.Context, conn net.Conn, outBound bool) {
	peer := node.NewPeer(conn, nil, network.TCPProtocol, 5, 5)
	if t.hasPeer(peer) {
		return
	}
	if t.HandShake != nil {
		if err := t.HandShake(peer); err != nil {
			peer.Close()
			return
		}
	}
	if outBound {
		if err := t.AddOutPeer(peer); err != nil {
			peer.Close()
			return
		}
	} else {
		if err := t.AddInPeer(peer); err != nil {
			peer.Close()
			return
		}
	}
	go func() {
		var wg sync.WaitGroup
		wg.Add(2)
		go func() { defer wg.Done(); peer.ReadPump(ctx) }()
		go func() { defer wg.Done(); peer.WritePump(ctx) }()
		wg.Wait()
		if outBound {
			t.RemoveOutPeer(peer)
		} else {
			t.RemoveInPeer(peer)
		}
		peer.Close()
	}()
}

// Add peer to the outbound peer map.
func (t *TCPTransport) AddOutPeer(peer p2p.Peer) error {
	err := t.State.AddOutPeer(peer)
	if err != nil {
		return err
	}
	t.State.IncOutBound()
	return nil
}

// Add peer to the inbound peer map.
func (t *TCPTransport) AddInPeer(peer p2p.Peer) error {
	err := t.State.AddInPeer(peer)
	if err != nil {
		return err
	}
	t.State.IncInBound()
	return nil
}

// Remove peer from the outbound peer map.
func (t *TCPTransport) RemoveOutPeer(peer p2p.Peer) {
	t.State.RemoveOutPeer(peer)
	t.State.DecOutBound()
}

// Remove peer from the inbound peer map.
func (t *TCPTransport) RemoveInPeer(peer p2p.Peer) {
	t.State.RemoveInPeer(peer)
	t.State.DecInBound()
}

// Including IN/OUT bounds peer
func (t *TCPTransport) Peers() map[string]p2p.Peer {
	combined := make(map[string]p2p.Peer)

	for id, p := range t.State.OutPeers() {
		combined[id] = p
	}
	for id, p := range t.State.InPeers() {
		combined[id] = p
	}
	return combined
}

// Existing peer in peerstable
func (t *TCPTransport) hasPeer(p p2p.Peer) bool {
	if _, exist := t.Peers()[p.ID()]; exist {
		return true
	}
	return false
}
