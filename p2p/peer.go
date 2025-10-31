package p2p

// import (
// 	"log"
// 	"sync"

// 	"github.com/wang900115/LCA/pkg/util/encode"
// )

// const (
// 	handshakeMsg = 0x00
// 	discMsg      = 0x01
// 	pingMsg      = 0x02
// 	pongMsg      = 0x03
// )

// type protoHandshake struct {
// 	Version    uint64
// 	Name       string
// 	Caps       []Cap
// 	ListenPort uint64
// 	ID         []byte
// }

// type PeerEventType string

// const (
// 	PeerEventTypeAdd     PeerEventType = "add"
// 	PeerEventTypeDrop    PeerEventType = "drop"
// 	PeerEventTypeMsgSend PeerEventType = "msgsend"
// 	PeerEventTypeMsgRecv PeerEventType = "msgrecv"
// )

// type PeerEvent struct {
// 	Type          PeerEventType
// 	Peer          encode.ID
// 	Error         string
// 	Protocol      string
// 	MsgCode       *uint64
// 	MsgSize       *uint32
// 	LocalAddress  string
// 	RemoteAddress string
// }

// type Peer struct {
// 	rw      *conn
// 	running map[string]*protoRW
// 	log     log.Logger

// 	wg       sync.WaitGroup
// 	protoErr chan error
// 	closed   chan struct{}
// 	pingRecv chan struct{}
// 	disc     chan DiscReason
// }

// type protoRW struct {
// 	Protocol
// 	in     chan Msg
// 	closed <-chan struct{}
// 	wstart <-chan struct{}
// 	werr   chan<- error
// 	offset uint64
// 	w      MsgWriter
// }
