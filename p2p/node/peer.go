package node

import (
	"context"
	"errors"
	"io"
	"net"

	"github.com/wang900115/LCA/crypt/did"
	"github.com/wang900115/LCA/p2p"
	"github.com/wang900115/LCA/p2p/network"
)

type Peer struct {
	net.Conn
	DID       did.PeerDID
	Protocol  network.Protocol
	Transport network.Packet
	State     *state
	Channel   *channel
	Meta      map[string]string
}

func NewPeer(conn net.Conn, services []did.ServiceEndpoint, transport network.TransportProtocol, inBoundLi, outBoundLi int) p2p.Peer {
	did := did.NewDID(services)
	protocol := network.NewProtocolInfo(transport)
	state := NewState(outBoundLi, inBoundLi)
	channel := NewChannel(make(chan network.Packet, 1024), make(chan network.Packet, 1024))

	return &Peer{
		Conn:     conn,
		DID:      did,
		State:    state,
		Channel:  channel,
		Protocol: protocol,
	}
}

// ID returns the unique identifier of the peer.
func (p *Peer) ID() string {
	return p.DID.DIDInfo().ID
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

// Add peer to the outbound peer map.
func (p *Peer) AddOutPeer(peer p2p.Peer) error {
	err := p.State.AddOutPeer(peer)
	if err != nil {
		return err
	}
	p.State.IncOutBound()
	return nil
}

// Add peer to the inbound peer map.
func (p *Peer) AddInPeer(peer p2p.Peer) error {
	err := p.State.AddInPeer(peer)
	if err != nil {
		return err
	}
	p.State.IncInBound()
	return nil
}

// Remove peer from the outbound peer map.
func (p *Peer) RemoveOutPeer(peer p2p.Peer) {
	p.State.RemoveOutPeer(peer)
	p.State.DecOutBound()
}

// Remove peer from the inbound peer map.
func (p *Peer) RemoveInPeer(peer p2p.Peer) {
	p.State.RemoveInPeer(peer)
	p.State.DecInBound()
}

func (p *Peer) Peers() map[string]p2p.Peer {
	comBined := make(map[string]p2p.Peer)
	for k, v := range p.State.OutPeers() {
		comBined[k] = v
	}
	for k, v := range p.State.InPeers() {
		comBined[k] = v
	}
	return comBined
}

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
			case p.Channel.Produce() <- &pkt:
			case <-ctx.Done():
				return
			}
		}
	}
}

func (p *Peer) WritePump(ctx context.Context) {

	defer func() {
		err := p.Conn.Close()
		if err != nil {
			panic(err)
		}
	}()

	for {
		select {
		case packet, ok := <-p.Channel.Consume():
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
