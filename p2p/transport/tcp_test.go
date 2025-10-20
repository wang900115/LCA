package transport

import (
	"context"
	"errors"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/wang900115/LCA/p2p"
	"github.com/wang900115/LCA/p2p/network"
	"github.com/wang900115/LCA/p2p/node"
)

func TestNewTCPTransport(t *testing.T) {
	opts := TCPTransportOpts{
		ListenAddr: ":0",
		HandShake:  nil,
		InBoundLi:  5,
		OutBoundLi: 5,
	}

	transport := NewTCPTransport(opts)
	tcpTransport := transport.(*TCPTransport)

	if tcpTransport.ListenAddr != opts.ListenAddr {
		t.Errorf("Expected ListenAddr %s, got %s", opts.ListenAddr, tcpTransport.ListenAddr)
	}
	if tcpTransport.State == nil {
		t.Error("Expected State to be initialized")
	}
}

func TestTCPTransport_Addr(t *testing.T) {
	opts := TCPTransportOpts{
		ListenAddr: "localhost:8080",
		InBoundLi:  5,
		OutBoundLi: 5,
	}

	transport := NewTCPTransport(opts).(*TCPTransport)

	if transport.Addr() != "localhost:8080" {
		t.Errorf("Expected address localhost:8080, got %s", transport.Addr())
	}
}

func TestTCPTransport_ListenAndClose(t *testing.T) {
	opts := TCPTransportOpts{
		ListenAddr: ":0",
		InBoundLi:  5,
		OutBoundLi: 5,
	}

	transport := NewTCPTransport(opts).(*TCPTransport)
	ctx := context.Background()

	// Test Listen
	err := transport.Listen(ctx)
	if err != nil {
		t.Fatalf("Failed to start listening: %v", err)
	}

	if transport.listener == nil {
		t.Error("Expected listener to be set")
	}

	// Test Close
	err = transport.Close()
	if err != nil {
		t.Errorf("Failed to close transport: %v", err)
	}
}

func TestTCPTransport_Dial(t *testing.T) {
	// Setup server transport
	serverOpts := TCPTransportOpts{
		ListenAddr: ":0",
		InBoundLi:  5,
		OutBoundLi: 5,
	}
	serverTransport := NewTCPTransport(serverOpts).(*TCPTransport)
	ctx := context.Background()

	err := serverTransport.Listen(ctx)
	if err != nil {
		t.Fatalf("Failed to start server: %v", err)
	}
	defer serverTransport.Close()

	serverAddr := serverTransport.listener.Addr().String()

	// Setup client transport
	clientOpts := TCPTransportOpts{
		ListenAddr: ":0",
		InBoundLi:  5,
		OutBoundLi: 5,
	}
	clientTransport := NewTCPTransport(clientOpts).(*TCPTransport)

	// Test Dial
	err = clientTransport.Dial(ctx, serverAddr)
	if err != nil {
		t.Errorf("Failed to dial: %v", err)
	}

	// Give some time for connection handling
	time.Sleep(100 * time.Millisecond)
}

func TestTCPTransport_HandShake(t *testing.T) {
	handshakeCalled := false
	handshakeFunc := func(peer p2p.Peer) error {
		handshakeCalled = true
		return nil
	}

	// Setup server transport with handshake
	serverOpts := TCPTransportOpts{
		ListenAddr: ":0",
		HandShake:  handshakeFunc,
		InBoundLi:  5,
		OutBoundLi: 5,
	}
	serverTransport := NewTCPTransport(serverOpts).(*TCPTransport)
	ctx := context.Background()

	err := serverTransport.Listen(ctx)
	if err != nil {
		t.Fatalf("Failed to start server: %v", err)
	}
	defer serverTransport.Close()

	serverAddr := serverTransport.listener.Addr().String()

	// Setup client transport
	clientOpts := TCPTransportOpts{
		ListenAddr: ":0",
		InBoundLi:  5,
		OutBoundLi: 5,
	}
	clientTransport := NewTCPTransport(clientOpts).(*TCPTransport)

	// Dial to trigger handshake
	err = clientTransport.Dial(ctx, serverAddr)
	if err != nil {
		t.Errorf("Failed to dial: %v", err)
	}

	// Give some time for handshake
	time.Sleep(100 * time.Millisecond)

	if !handshakeCalled {
		t.Error("Expected handshake to be called")
	}
}

func TestTCPTransport_PeerManagement(t *testing.T) {
	opts := TCPTransportOpts{
		ListenAddr: ":0",
		InBoundLi:  5,
		OutBoundLi: 5,
	}
	transport := NewTCPTransport(opts).(*TCPTransport)

	// Create mock connections with different addresses
	// Start two listeners to get different addresses
	l1, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatalf("Failed to create listener 1: %v", err)
	}
	defer l1.Close()

	l2, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatalf("Failed to create listener 2: %v", err)
	}
	defer l2.Close()

	// Create connections to these listeners
	c1, err := net.Dial("tcp", l1.Addr().String())
	if err != nil {
		t.Fatalf("Failed to dial listener 1: %v", err)
	}
	defer c1.Close()

	c2, err := net.Dial("tcp", l2.Addr().String())
	if err != nil {
		t.Fatalf("Failed to dial listener 2: %v", err)
	}
	defer c2.Close()

	peer1 := node.NewPeer(c1, nil, network.TCPProtocol, 1, 1)
	peer2 := node.NewPeer(c2, nil, network.TCPProtocol, 1, 1)

	// Test AddOutPeer
	err = transport.AddOutPeer(peer1)
	if err != nil {
		t.Errorf("Failed to add outbound peer: %v", err)
	}

	// Test AddInPeer
	err = transport.AddInPeer(peer2)
	if err != nil {
		t.Errorf("Failed to add inbound peer: %v", err)
	}

	// Test Peers()
	peers := transport.Peers()
	if len(peers) != 2 {
		t.Errorf("Expected 2 peers, got %d", len(peers))
	}

	// Test hasPeer
	if !transport.hasPeer(peer1.Addr()) {
		t.Error("Expected to find peer1")
	}
	if !transport.hasPeer(peer2.Addr()) {
		t.Error("Expected to find peer2")
	}

	// Test RemoveOutPeer
	transport.RemoveOutPeer(peer1)
	if transport.hasPeer(peer1.Addr()) {
		t.Error("Expected peer1 to be removed")
	}

	// Test RemoveInPeer
	transport.RemoveInPeer(peer2)
	if transport.hasPeer(peer2.Addr()) {
		t.Error("Expected peer2 to be removed")
	}

	// Test Peers() after removal
	peers = transport.Peers()
	if len(peers) != 0 {
		t.Errorf("Expected 0 peers after removal, got %d", len(peers))
	}
}

func TestTCPTransport_HandleConn(t *testing.T) {
	opts := TCPTransportOpts{
		ListenAddr: ":0",
		InBoundLi:  5,
		OutBoundLi: 5,
	}
	transport := NewTCPTransport(opts).(*TCPTransport)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create mock connection
	c1, c2 := net.Pipe()
	defer c2.Close()

	// Handle outbound connection
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		transport.handleConn(ctx, c1, true)
	}()

	// Give some time for connection handling
	time.Sleep(100 * time.Millisecond)

	// Check if peer was added
	peers := transport.Peers()
	if len(peers) != 1 {
		t.Errorf("Expected 1 peer, got %d", len(peers))
	}

	// Close connection to trigger cleanup
	c1.Close()
	cancel()
	wg.Wait()

	// Give some time for cleanup
	time.Sleep(100 * time.Millisecond)

	// Check if peer was removed
	peers = transport.Peers()
	if len(peers) != 0 {
		t.Errorf("Expected 0 peers after cleanup, got %d", len(peers))
	}
}

func TestTCPTransport_DuplicatePeer(t *testing.T) {
	opts := TCPTransportOpts{
		ListenAddr: ":0",
		InBoundLi:  5,
		OutBoundLi: 5,
	}
	transport := NewTCPTransport(opts).(*TCPTransport)
	ctx := context.Background()

	// Create mock connections with same remote address
	c1, c2 := net.Pipe()
	defer c1.Close()
	defer c2.Close()

	// Add first peer
	peer1 := node.NewPeer(c1, nil, network.TCPProtocol, 1, 1)
	_ = transport.AddOutPeer(peer1)

	// Try to handle connection with same remote address
	// This should return early due to hasPeer check
	transport.handleConn(ctx, c2, false)

	// Should still have only 1 peer
	peers := transport.Peers()
	if len(peers) != 1 {
		t.Errorf("Expected 1 peer (duplicate should be ignored), got %d", len(peers))
	}
}

func TestTCPTransport_FailedHandshake(t *testing.T) {
	handshakeError := errors.New("handshake failed")
	handshakeFunc := func(peer p2p.Peer) error {
		return handshakeError
	}

	opts := TCPTransportOpts{
		ListenAddr: ":0",
		HandShake:  handshakeFunc,
		InBoundLi:  5,
		OutBoundLi: 5,
	}
	transport := NewTCPTransport(opts).(*TCPTransport)
	ctx := context.Background()

	// Create mock connection
	c1, c2 := net.Pipe()
	defer c2.Close()

	// Handle connection - should fail at handshake
	transport.handleConn(ctx, c1, true)

	// Give some time for processing
	time.Sleep(100 * time.Millisecond)

	// Should have no peers due to failed handshake
	peers := transport.Peers()
	if len(peers) != 0 {
		t.Errorf("Expected 0 peers due to failed handshake, got %d", len(peers))
	}
}
