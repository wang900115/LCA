package node

import (
	"sync"

	"github.com/wang900115/LCA/p2p/network"
)

type channel struct {
	readCh  <-chan network.Packet
	writeCh chan<- network.Packet
	wg      *sync.WaitGroup
}

func NewChannel(readCh <-chan network.Packet, writeCh chan<- network.Packet) *channel {
	return &channel{
		readCh:  readCh,
		writeCh: writeCh,
		wg:      &sync.WaitGroup{},
	}
}

func (ch *channel) Consume() <-chan network.Packet {
	return ch.readCh
}

func (ch *channel) Produce() chan<- network.Packet {
	return ch.writeCh
}

func (ch *channel) OpenStream()  { ch.wg.Add(1) }
func (ch *channel) CloseStream() { ch.wg.Done() }
func (ch *channel) WaitStream()  { ch.wg.Wait() }
