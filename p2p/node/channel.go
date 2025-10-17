package node

import (
	"sync"

	"github.com/wang900115/LCA/p2p"
)

type channel struct {
	readCh  <-chan p2p.Packet
	writeCh chan<- p2p.Packet
	wg      *sync.WaitGroup
}

func NewChannel(readCh <-chan p2p.Packet, writeCh chan<- p2p.Packet) *channel {
	return &channel{
		readCh:  readCh,
		writeCh: writeCh,
		wg:      &sync.WaitGroup{},
	}
}

func (ch *channel) Consume() <-chan p2p.Packet {
	return ch.readCh
}

func (ch *channel) Produce() chan<- p2p.Packet {
	return ch.writeCh
}

func (ch *channel) OpenStream()  { ch.wg.Add(1) }
func (ch *channel) CloseStream() { ch.wg.Done() }
func (ch *channel) WaitStream()  { ch.wg.Wait() }
