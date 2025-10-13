package p2p

import "errors"

// 1 bytes for command
// 4 bytes for length
// n bytes for payload
// 1+4 = 5 bytes header
// 4 bytes footer (checksum)

type Command byte

const (
	HEARTBEAT       Command = 0x00
	PEERINFO        Command = 0x01
	PEERACK         Command = 0x02
	PEERERROR       Command = 0x03
	BETCREATE       Command = 0x10
	BETACK          Command = 0x11
	BETERROR        Command = 0x12
	RESETTLECREATE  Command = 0x20
	RESETTLEACK     Command = 0x21
	RESETTLEERROR   Command = 0x22
	ROUNDSTART      Command = 0x30
	ROUNDSTARTACK   Command = 0x31
	ROUNDSTARTERROR Command = 0x32
	ROUNDEND        Command = 0x33
	ROUNDENDACK     Command = 0x34
	ROUNDENDERROR   Command = 0x35
	ROUNDWAIT       Command = 0x36
	ROUNDWAITACK    Command = 0x37
	ROUNDWAITERROR  Command = 0x38
)

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
func NewPacket(command byte, playload []byte) Packet {
	return &PacketContent{
		Command:  command,
		Length:   uint32(len(playload)),
		Payload:  playload,
		CheckSum: 0,
	}
}

// Getters for Packet fields
func (p *PacketContent) GetCommand() Command {
	return Command(p.Command)
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
func Decode(data []byte) (Packet, error) {
	return nil, nil
}
