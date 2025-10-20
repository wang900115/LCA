package node

import (
	"context"
	"errors"
	"io"
	"net"

	"github.com/wang900115/LCA/did"
	"github.com/wang900115/LCA/p2p"
	"github.com/wang900115/LCA/p2p/network"
)

// Peer represents a peer in the P2P network.
type Peer struct {
	net.Conn
	DID       did.PeerDID
	Protocol  network.Protocol
	Transport network.Packet
	Channel   *channel
	Meta      map[string]string
}

// NewPeer creates a new peer instance.
func NewPeer(conn net.Conn, services []did.ServiceEndpoint, transport network.TransportProtocol, inBoundLi, outBoundLi int) p2p.Peer {
	did := did.NewDID(services)
	protocol := network.NewProtocolInfo(transport)
	inCh := make(chan network.Packet, 1024)
	outCh := make(chan network.Packet, 1024)
	channel := NewChannel(inCh, outCh)

	return &Peer{
		Conn:     conn,
		DID:      did,
		Channel:  channel,
		Protocol: protocol,
		Meta:     map[string]string{},
	}
}

// Addr return the peer remote address
func (p *Peer) Addr() string {
	return p.Conn.RemoteAddr().String()
}

// ID returns the unique identifier of the peer.
func (p *Peer) ID() string {
	return p.DID.DID().ID
}

// Document returns the DID document of the peer.
func (p *Peer) Document() *did.DIDDocument {
	return p.DID.ToDocument()
}

// Protocol returns the protocol information of the peer.
func (p *Peer) ProtocolInfo() *network.ProtocolInfo {
	return p.Protocol.ProtocolInfo()
}

// SendPacket sends a packet to the peer.
func (p *Peer) Send(packet network.Packet) error {
	p.Channel.Produce() <- packet
	return nil
}

// ReceivePacket returns a channel to receive packets from the peer.
func (p *Peer) Receive() (<-chan network.Packet, error) {
	return p.Channel.Consume(), nil
}

// ReadPump pumps packets from the peer connection to the channel.
func (p *Peer) ReadPump(ctx context.Context) {
	defer func() {
		err := p.Conn.Close()
		if err != nil {
			panic(err)
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			var pkt network.PacketContent
			_, err := pkt.Decode(p.Conn)
			if err != nil {
				if errors.Is(err, io.EOF) {
					return
				}
				return
			}
			select {
			case p.Channel.In() <- &pkt:
			case <-ctx.Done():
				return
			}
		}
	}
}

// WritePump pumps packets from the channel to the peer connection.
func (p *Peer) WritePump(ctx context.Context) {
	defer func() {
		err := p.Conn.Close()
		if err != nil {
			panic(err)
		}
	}()
	for {
		select {
		case packet, ok := <-p.Channel.Out():
			if !ok {
				return
			}
			_, err := packet.Encode(p.Conn)
			if err != nil {
				return
			}
		case <-ctx.Done():
			return
		}
	}
}
