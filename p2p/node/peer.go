package node

import (
	"net"
	"sync"
)

type TCPPeer struct {
	net.Conn
	outBound bool // if we dial and retrieve a conn => outbound is true else false
	wg       *sync.WaitGroup
}

func NewTCPPeer(conn net.Conn, outBound bool) *TCPPeer {
	return &TCPPeer{
		Conn:     conn,
		outBound: outBound,
		wg:       &sync.WaitGroup{},
	}
}

func (p *TCPPeer) AddWG() {
	p.wg.Add(1)
}

func (p *TCPPeer) WaitWG() {
	p.wg.Wait()
}

func (p *TCPPeer) CloseStream() {
	p.wg.Done()
}

func (p *TCPPeer) Send(b []byte) error {
	_, err := p.Conn.Write(b)
	return err
}
