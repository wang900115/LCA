package transport

import (
	"errors"
	"fmt"
	"log"
	"net"

	"github.com/wang900115/LCA/p2p"
	"github.com/wang900115/LCA/p2p/encode"
	"github.com/wang900115/LCA/p2p/node"
	"github.com/wang900115/LCA/p2p/packet"
)

type TCPTransportOpts struct {
	ListenAddr    string
	HandShakeFunc p2p.HandShakeFunc
	OnPeer        func(p2p.Peer) error
}

type TCPTransport struct {
	TCPTransportOpts
	listener net.Listener
	Decoder  encode.Decoder
	rpcch    chan packet.RPC
}

func NewTCPTransport(opts TCPTransportOpts) *TCPTransport {
	return &TCPTransport{
		TCPTransportOpts: opts,
		rpcch:            make(chan packet.RPC, 1024),
	}
}

func (t *TCPTransport) Addr() string {
	return t.ListenAddr
}

func (t *TCPTransport) Consume() <-chan packet.RPC {
	return t.rpcch
}

func (t *TCPTransport) Close() error {
	return t.listener.Close()
}

func (t *TCPTransport) Dial(addr string) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}
	go t.handleConn(conn, true)
	return nil
}

func (t *TCPTransport) ListenAndAccept() error {
	var err error
	t.listener, err = net.Listen("tcp", t.ListenAddr)
	if err != nil {
		return err
	}

	go t.startAcceptLoop()
	return nil
}

func (t *TCPTransport) startAcceptLoop() {
	for {
		conn, err := t.listener.Accept()
		if errors.Is(err, net.ErrClosed) {
			return
		}
		if err != nil {
			fmt.Printf("TCP accept error: %s\n", err)
		}
		go t.handleConn(conn, false)
	}
}

func (t *TCPTransport) handleConn(conn net.Conn, outBound bool) {
	var err error
	defer func() {
		conn.Close()
	}()

	peer := node.NewPeer(conn, outBound)

	if err = t.HandShakeFunc(peer); err != nil {
		return
	}

	if t.OnPeer(peer) != nil {
		if err = t.OnPeer(peer); err != nil {
			return
		}
	}

	for {
		rpc := packet.RPC{}
		err := t.Decoder.Decode(conn, &rpc)
		if err != nil {
			return
		}

		rpc.From = conn.RemoteAddr().String()

		if rpc.Stream {
			if p, err := peer.OpenStream(); err != nil {
				log.Printf("open stream error with peer: %+v\n", p)
				return
			}
			fmt.Printf("[%s] incoming stream, waiting...\n", conn.RemoteAddr())
			peer.WaitSream()
			fmt.Printf("[%s] stream closed, resuming read loop\n", conn.RemoteAddr())
			continue
		}

		t.rpcch <- rpc
	}
}
