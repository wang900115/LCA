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

// NewPacket creates a new packet with the given command and payload
func NewPacket(command byte, playload []byte) *PacketContent {
	return &PacketContent{
		Command:  command,
		Length:   uint32(len(playload)),
		Payload:  playload,
		CheckSum: 0,
	}
}

// Getters for Packet fields
func (p *PacketContent) GetCommand() byte {
	return p.Command
}

// Getters for Packet fields
func (p *PacketContent) GetLength() uint32 {
	return p.Length
}

// Getters for Packet fields
func (p *PacketContent) GetPayload() []byte {
	return p.Payload
}

// Getters for Packet fields
func (p *PacketContent) GetCheckSum() uint32 {
	return p.CheckSum
}

// Check Packet length is less then network limit
func (p *PacketContent) Check() error {
	if p.Length > MaxPacketSize {
		return ErrPacketSizeExceeds
	}
	if uint32(len(p.Payload)) != p.Length {
		return ErrPayloadSizeMisMatch
	}
	return nil
}

// Encode Packet to binary packet
func (p *PacketContent) Encode() ([]byte, error) {
	return nil, nil
}

// Decode binary packet to Packet
func Decode2Packet(data []byte) (*PacketContent, error) {
	return nil, nil
}
