package p2p

// import (
// 	"fmt"

// 	common "github.com/wang900115/LCA/p2p/com"
// 	"github.com/wang900115/LCA/p2p/network"
// )

// type HandShakeFunc func(Peer) error

// // NoopHandshakeFunc performs no handshake and immediately returns nil.
// func NoopHandshakeFunc(peer Peer) error {
// 	return nil
// }

// // BasicHandshakeFunc performs a simple handshake by sending a HEARTBEAT
// // command and expecting a PEERACK or HEARTBEAT response.
// func BasicHandshakeFunc(peer Peer) error {
// 	network.NewMessageContent(common.PRIVATE)
// 	network.NewRPCContent()
// 	pkt := network.NewPacket(common.HEARTBEAT)
// 	if err := peer.Send(pkt); err != nil {
// 		return fmt.Errorf("send heartbeat failed: %w", err)
// 	}
// 	pkCh, err := peer.Receive()
// 	if err != nil {
// 		return fmt.Errorf("receive failed: %w", err)
// 	}

// 	for pk := range pkCh {
// 		cmd := pk.GetCommand()
// 		if cmd != byte(common.PEERINFO) && cmd != byte(common.PEERACK) {
// 			return fmt.Errorf("invalid handshake response: 0x%x", cmd)
// 		}
// 		break
// 	}

// 	if err := peer.HandShake(); err != nil {
// 		return fmt.Errorf("handshake function failed: %w", err)
// 	}

// 	return nil
// }
// //
