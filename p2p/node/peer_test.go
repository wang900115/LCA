package node

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/wang900115/LCA/did"
	common "github.com/wang900115/LCA/p2p/com"
	"github.com/wang900115/LCA/p2p/network"
)

func TestPeer(t *testing.T) {
	// Create a new peer instance
	peer := NewPeer(nil, nil, did.VerifierConfig{}, network.TCPProtocol, 0, 0)
	// Test the peer's ID
	if peer.ID() == "" {
		t.Error("Expected peer ID to be non-empty")
	}
	// Test the peer's document
	if peer.Document() == nil {
		t.Error("Expected peer document to be non-nil")
	}
	// Test the peer's protocol information
	if peer.ProtocolInfo() == nil {
		t.Error("Expected peer protocol information to be non-nil")
	}
}

func TestConnPeer(t *testing.T) {
	c1, c2 := net.Pipe()
	defer c1.Close()
	defer c2.Close()

	peer := NewPeer(c1, nil, did.VerifierConfig{}, network.TCPProtocol, 1, 1)
	if peer.Addr() != c1.RemoteAddr().String() && peer.Addr() == c1.LocalAddr().String() {
		t.Errorf("Expected peer address to be %s, got %s", c1.RemoteAddr().String(), peer.Addr())
	}
	peer2 := NewPeer(c2, nil, did.VerifierConfig{}, network.TCPProtocol, 1, 1)
	if peer2.Addr() != c2.RemoteAddr().String() && peer2.Addr() == c2.LocalAddr().String() {
		t.Errorf("Expected peer address to be %s, got %s", c2.RemoteAddr().String(), peer2.Addr())
	}
}

func TestPeerPumps(t *testing.T) {
	c1, c2 := net.Pipe()
	defer c1.Close()
	defer c2.Close()

	p1 := NewPeer(c1, nil, did.VerifierConfig{}, network.TCPProtocol, 0, 0)
	p2 := NewPeer(c2, nil, did.VerifierConfig{}, network.TCPProtocol, 0, 0)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go p1.ReadPump(ctx)
	go p1.WritePump(ctx)
	go p2.ReadPump(ctx)
	go p2.WritePump(ctx)

	msg, err := network.NewMessageContent(common.PUBLIC, []byte("hello"), nil)
	if err != nil {
		t.Fatalf("failed to create message: %v", err)
	}
	testDID := did.NewDIDIdentifier(nil)
	rpc, err := network.NewRPCContent(msg, testDID)
	if err != nil {
		t.Fatalf("failed to create rpc: %v", err)
	}
	pkt, err := network.NewPacket(common.MESSAGESEND, rpc)
	if err != nil {
		t.Fatalf("failed to create packet: %v", err)
	}
	if err := p1.Send(pkt); err != nil {
		t.Fatalf("send failed: %v", err)
	}
	ch, _ := p2.Receive()
	select {
	case received := <-ch:
		if received.GetCommand() != common.MESSAGESEND {
			t.Fatalf("unexpected command: %v", received.GetCommand())
		}
	case <-time.After(10 * time.Second):
		t.Fatal("timeout waiting for packet")
	}
}
