package p2p

import (
	"errors"
	// Removed import to avoid cycle: "github.com/wang900115/LCA/p2p"
)

// 1 bytes for command
// 4 bytes for length
// n bytes for payload
// 1+4 = 5 bytes header
// 4 bytes footer (checksum)

const (
	MaxPacketSize = 4 * 1024 * 1024
)

var (
	ErrPacketSizeExceeds   = errors.New("packet size exceeds network limit")
	ErrPayloadSizeMisMatch = errors.New("packet length mismatch with payload")
)

type PacketContent struct {
	Command  byte
	Length   uint32
	Payload  []byte
	CheckSum uint32
}
