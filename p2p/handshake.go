package p2p

import (
	"errors"

	common "github.com/wang900115/LCA/p2p/com"
)

type HandShakeFunc func(Peer) error

// NoopHandshakeFunc performs no handshake and immediately returns nil.
func NoopHandshakeFunc(peer Peer) error {
	return nil
}

// BasicHandshakeFunc performs a simple handshake by sending a HEARTBEAT
// command and expecting a PEERACK or HEARTBEAT response.
func BasicHandshakeFunc(peer Peer) error {
	pkt := &PacketContent{
		Command:  byte(common.HEARTBEAT),
		Length:   0,
		Payload:  nil,
		CheckSum: 0,
	}
	if err := peer.SendPacket(pkt); err != nil {
		return err
	}
	resp, err := peer.ReceivePacket()
	if err != nil {
		return err
	}
	if resp.GetCommand() != 0x02 && resp.GetCommand() != 0x00 {
		return errors.New("invalid handshake response")
	}
	if err := peer.HandShake(); err != nil {
		return err
	}
	return nil
}
