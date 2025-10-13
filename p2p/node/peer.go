package node

import (
	"net"
	"sync"

	"github.com/wang900115/LCA/p2p"
)

type Peer struct {
	net.Conn
	ID        string
	Protocol  p2p.Protocol
	Meta      map[string]string
	wg        *sync.WaitGroup
	outBound  bool
	handShake bool
	stream    bool

	rpcch <-chan p2p.RPC
}

func NewPeer(conn net.Conn, outBound bool) p2p.Peer {
	return &Peer{
		Conn:     conn,
		outBound: outBound,
		wg:       &sync.WaitGroup{},
		rpcch:    make(<-chan p2p.RPC),
	}
}

// GetID returns the peer ID
func (p *Peer) GetID() string {
	return p.ID
}

// GetMeta returns the peer metadata
func (p *Peer) GetMeta() map[string]string {
	return p.Meta
}

// OpenStream increments the waitgroup counter and returns a new peer
func (p *Peer) OpenStream() (p2p.Peer, error) {
	p.wg.Add(1)
	peer := *p
	peer.stream = true
	return &peer, nil
}

// WaitSream waits for all streams to be closed
func (p *Peer) WaitSream() {
	p.wg.Wait()
}

// CloseStream decrements the waitgroup counter
func (p *Peer) CloseStream() {
	p.wg.Done()
}

// IsStream returns true if the peer is a stream
func (p *Peer) IsStream() bool {
	return p.stream
}

// IsHandShake returns true if the handshake is done
func (p *Peer) IsHandShake() bool {
	return p.handShake
}

// HandShake marks the peer as handshaked
func (p *Peer) HandShake() error {
	p.handShake = true
	return nil
}

// HandShakeWithData marks the peer as handshaked and sets the metadata
func (p *Peer) HandShakeWithData(data []byte) error {
	p.handShake = true
	p.Meta = map[string]string{
		"handshake_data": string(data),
	}
	return nil
}

// SetMeta sets the peer metadata
func (p *Peer) SetMeta(meta map[string]string) {
	p.Meta = meta
}

// GetProtocol returns the peer protocol
func (p *Peer) GetProtocol() p2p.Protocol {
	return p.Protocol
}

// Close closes the peer connection
func (p *Peer) Close() error {
	return p.Conn.Close()
}

// Receive reads data from the peer connection
func (p *Peer) ReceivePacket() (p2p.Packet, error) {
	buf := make([]byte, 4096) // or any appropriate buffer size
	n, err := p.Conn.Read(buf)
	if err != nil {
		return nil, err
	}
	packet, err := p2p.Decode2Packet(buf[:n])
	if err != nil {
		return nil, err
	}
	return packet, nil
}

// SendPacket sends a packet to the peer
func (p *Peer) SendPacket(packet p2p.Packet) error {
	data, err := packet.Encode()
	if err != nil {
		return err
	}
	_, err = p.Conn.Write(data)
	if err != nil {
		return err
	}
	return nil
}

// IsOutBound returns true if the peer is an outbound connection
func (p *Peer) IsOutBound() bool {
	return p.outBound
}

// SetProtocol sets the peer protocol
func (p *Peer) SetProtocol(protocol p2p.Protocol) {
	p.Protocol = protocol
}

// Consume returns a channel to consume incoming RPCs
func (p *Peer) Consume() <-chan p2p.RPC {
	return p.rpcch
}
