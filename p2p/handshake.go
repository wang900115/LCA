package p2p

import (
	"fmt"

	common "github.com/wang900115/LCA/p2p/com"
	"github.com/wang900115/LCA/p2p/network"
)

type HandShakeFunc func(Peer) error

// NoopHandshakeFunc performs no handshake and immediately returns nil.
func NoopHandshakeFunc(peer Peer) error {
	return nil
}

// BasicHandshakeFunc performs a simple handshake by sending a HEARTBEAT
// command and expecting a PEERACK or HEARTBEAT response.
func BasicHandshakeFunc(peer Peer) error {

	pkt, err := network.NewPacket(common.HEARTBEAT, nil)
	if err != nil {
		return err
	}
	if err := peer.Send(pkt); err != nil {
		return err
	}
	pkCh, err := peer.Receive()
	if err != nil {
		return err
	}
	for pk := range pkCh {
		cmd := pk.GetCommand()
		if cmd == common.PEERINFO || cmd == common.PEERACK {
			return nil
		}
	}
	return fmt.Errorf("handshake failed")
}
