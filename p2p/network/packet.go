package network

import (
	"encoding/binary"
	"errors"
	"io"

	crypto "github.com/wang900115/LCA/crypt"
	common "github.com/wang900115/LCA/p2p/com"
	// Removed import to avoid cycle: "github.com/wang900115/LCA/p2p"
)

// 1 bytes for command
// 4 bytes for length
// n bytes for payload
// 1+4 = 5 bytes header
// 4 bytes footer (checksum)

const (
	MaxPacketPayloadSize = 200
)

var (
	errPacketPayloadExceed = &decErr{"packet payload is exceed 200 bytes"}
	errPacketChecksumFail  = errors.New("packet checksum verification failed")
)

type Packet interface {
	GetCommand() common.Command
	Encode(w io.Writer) (int, error)
	Decode(r io.Reader) (int, error)
	Bytes() []byte
	Check() error
	Len() int
	Max() int
}

type PacketContent struct {
	Command    common.Command
	PayloadLen uint16
	Payload    [MaxPacketPayloadSize]byte
	CheckSum   uint64
}

func NewPacket(command common.Command, rpc RPC) (Packet, error) {
	if rpc.Len() > MaxPacketPayloadSize {
		return nil, errPacketPayloadExceed
	}
	var pkt PacketContent
	pkt.Command = command
	pkt.PayloadLen = uint16(rpc.Len())
	copy(pkt.Payload[:], rpc.Bytes())
	cmdBytes := []byte{byte(pkt.Command)}
	lenBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(lenBytes, pkt.PayloadLen)
	data := append(cmdBytes, lenBytes...)
	data = append(data, pkt.Payload[:pkt.PayloadLen]...)
	pkt.CheckSum = crypto.CRC64(data)
	return &pkt, nil
}

func (p *PacketContent) GetCommand() common.Command {
	return p.Command
}

func (p *PacketContent) Encode(w io.Writer) (int, error) {
	n := 0
	if err := write(w, p.Command); err != nil {
		return n, wrapFieldError(err, "command")
	}
	n += 1
	if err := write(w, p.PayloadLen); err != nil {
		return n, wrapFieldError(err, "payload length")
	}
	n += 2
	written, err := w.Write(p.Payload[:p.PayloadLen])
	if err != nil {
		return n, wrapFieldError(err, "payload")
	}
	n += written
	if err := write(w, p.CheckSum); err != nil {
		return n, wrapFieldError(err, "checksum")
	}
	n += 8
	return n, nil
}

func (p *PacketContent) Decode(r io.Reader) (int, error) {
	n := 0
	if err := read(r, &p.Command); err != nil {
		return n, wrapFieldError(err, "command")
	}
	n += 1
	if err := read(r, &p.PayloadLen); err != nil {
		return n, wrapFieldError(err, "payload length")
	}
	n += 2
	if p.PayloadLen > MaxPacketPayloadSize {
		return n, errPacketPayloadExceed
	}

	readBytes, err := io.ReadFull(r, p.Payload[:p.PayloadLen])
	if err != nil {
		return n, wrapFieldError(err, "payload")
	}
	n += readBytes
	if err := read(r, &p.CheckSum); err != nil {
		return n, wrapFieldError(err, "checksum")
	}
	n += 8
	return n, nil
}

func (p *PacketContent) Bytes() []byte {
	buf := make([]byte, 0, 1+2+int(p.PayloadLen)+8)
	buf = append(buf, byte(p.Command))
	tmp := make([]byte, 2)
	binary.BigEndian.PutUint16(tmp, uint16(p.PayloadLen))
	buf = append(buf, tmp...)
	buf = append(buf, p.Payload[:p.PayloadLen]...)
	var checksumBytes [8]byte
	binary.BigEndian.PutUint64(checksumBytes[:], p.CheckSum)
	buf = append(buf, checksumBytes[:]...)
	return buf
}

func (p *PacketContent) Check() error {
	cmdBytes := []byte{byte(p.Command)}
	lenBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(lenBytes, p.PayloadLen)
	data := append(cmdBytes, lenBytes...)
	data = append(data, p.Payload[:p.PayloadLen]...)
	if crypto.CRC64(data) == p.CheckSum {
		return nil
	}
	return errPacketChecksumFail
}

func (p *PacketContent) Len() int {
	return int(len(p.Bytes()))
}

func (p *PacketContent) Max() int {
	return MaxPacketPayloadSize
}
