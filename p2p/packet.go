package p2p

import "errors"

// 1 byte for type
// 1 bytes for command
// 4 bytes for length
// n bytes for payload
// 1 + 1 +4 = 6 bytes header
// 4 bytes footer (checksum)

const (
	INCOMMINGMESSAGE = 0x1
	INCOMMINGSTREAM  = 0x2
	OUTGOINGMESSAGE  = 0x3
	OUTGOINGSTREAM   = 0x4
)

const (
	HEARTBEAT      = 0x00
	PEERINFO       = 0x01
	PEERACK        = 0x02
	PEERERROR      = 0x03
	BETCREATE      = 0x04
	BETACK         = 0x05
	BETERROR       = 0x06
	RESETTLECREATE = 0x07
	RESETTLEACK    = 0x08
	RESETTLEERROR  = 0x09
	ROUNDSTART     = 0x0a
	ROUNDSTARTACK  = 0x0b
	ROUNDEND       = 0x0c
	ROUNDENDACK    = 0x0d
	ROUNDWAIT      = 0x0e
	ROUNDWAITACK   = 0x0f
)

const (
	MaxPacketSize = 4 * 1024 * 1024
)

var (
	ErrPacketSizeExceeds   = errors.New("packet size exceeds network limit")
	ErrPayloadSizeMisMatch = errors.New("packet length mismatch with payload")
)

type Packet struct {
	Command  byte
	Length   uint32
	Payload  []byte
	CheckSum uint32
}

// Check Packet length is less then network limit
func (p *Packet) Check() error {
	if p.Length > MaxPacketSize {
		return ErrPacketSizeExceeds
	}
	if uint32(len(p.Payload)) != p.Length {
		return ErrPayloadSizeMisMatch
	}
	return nil
}

// Encode Packet to binary packet
func (p *Packet) Encode() ([]byte, error) {
	return nil, nil
}

// Decode Packet from []byte
func Decode(data []byte) (*Packet, error) {
	return nil, nil
}
