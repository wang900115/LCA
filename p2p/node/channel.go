package node

import (
	"sync"

	"github.com/wang900115/LCA/p2p/network"
)

// channel implements the Channel interface for peer communication.
type channel struct {
	// inbound channel:  packets from connection -> consumed by app
	readCh chan network.Packet
	// outbound channel: packets from app -> written to connection
	writeCh chan network.Packet
	// wait group to track open streams
	wg *sync.WaitGroup
}

// Ensure channel implements network.Channel interface
func NewChannel(readCh chan network.Packet, writeCh chan network.Packet) *channel {
	return &channel{
		readCh:  readCh,
		writeCh: writeCh,
		wg:      &sync.WaitGroup{},
	}
}

// Consume returns the inbound channel for consuming decoded packets.
func (ch *channel) Consume() <-chan network.Packet {
	return ch.readCh
}

// Produce returns the outbound channel for producing packets to be encoded and sent.
func (ch *channel) Produce() chan<- network.Packet {
	return ch.writeCh
}

// In returns the inbound channel as a send-capable channel so internal pumps can
// place decoded packets into it.
func (ch *channel) In() chan<- network.Packet {
	return ch.readCh
}

// Out returns the outbound channel as a receive-capable channel so internal pumps can
// read packets to write to the connection.
func (ch *channel) Out() <-chan network.Packet {
	return ch.writeCh
}

// OpenStream indicates that a new stream has been opened on this channel.
func (ch *channel) OpenStream() { ch.wg.Add(1) }

// CloseStream indicates that a stream has been closed on this channel.
func (ch *channel) CloseStream() { ch.wg.Done() }

// WaitStream blocks until all opened streams have been closed.
func (ch *channel) WaitStream() { ch.wg.Wait() }
