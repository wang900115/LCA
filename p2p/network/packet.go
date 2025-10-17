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
	Encode(w io.Writer) (int, error)
	Decode(r io.Reader) (int, error)
	Bytes() []byte
	Check() error
	Len() int
	Max() int
}

type PacketContent struct {
	Command    common.Command
	PayloadLen uint8
	Payload    [MaxPacketPayloadSize]byte
	CheckSum   uint64
}

func NewPacket(command common.Command, rpc RPC) (Packet, error) {
	if rpc.Len() > MaxPacketPayloadSize {
		return nil, errPacketPayloadExceed
	}
	checkSum := crypto.CRC64(rpc.Bytes())
	var pkt PacketContent
	pkt.Command = command
	pkt.PayloadLen = uint8(rpc.Len())
	copy(pkt.Payload[:], rpc.Bytes())
	pkt.CheckSum = checkSum
	return &pkt, nil
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
	n += 1
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
	n += 1
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
	buf := make([]byte, 0, 1+1+int(p.PayloadLen)+8)
	buf = append(buf, byte(p.Command))
	buf = append(buf, p.PayloadLen)
	buf = append(buf, p.Payload[:p.PayloadLen]...)
	checksumBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(checksumBytes, p.CheckSum)
	buf = append(buf, checksumBytes...)
	return buf
}

func (p *PacketContent) Check() error {
	ok := crypto.VerifyCRC64(p.Payload[:p.PayloadLen], p.CheckSum)
	if ok {
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
