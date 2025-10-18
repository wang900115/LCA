package network

import (
	"encoding/binary"
	"fmt"
	"io"
	"time"

	crypto "github.com/wang900115/LCA/crypt"
	common "github.com/wang900115/LCA/p2p/com"
	"golang.org/x/crypto/sha3"
)

const (
	MaxMessagePayloadSize = 32
)

type decErr struct{ msg string }

func (d *decErr) Error() string { return d.msg }

var (
	errMessagePayloadExceed = &decErr{"message payload is exceed 32 bytes"}
)

type Message interface {
	Encode(w io.Writer) (int, error)
	Decode(r io.Reader) (int, error)
	Bytes() []byte
	Verify(sharedKey []byte) bool
	Len() int
	Max() int
}

type MessageContent struct {
	Type       common.Message              // Type of message
	PayloadLen uint8                       // Real payload length
	Payload    [MaxMessagePayloadSize]byte // Actual message content or data
	Hash       [32]byte                    // Payload hash
	CreatedAt  int64                       // Timestamp of message creation
}

func NewMessageContent(msgType common.Message, payload []byte, sharedKey []byte) (Message, error) {
	if len(payload) > MaxMessagePayloadSize {
		return nil, errMessagePayloadExceed
	}
	var fixedPayload [MaxMessagePayloadSize]byte
	copy(fixedPayload[:], payload)

	var msg MessageContent
	msg.Type = msgType
	msg.PayloadLen = uint8(len(payload))
	copy(msg.Payload[:], payload)
	msg.CreatedAt = time.Now().UTC().UnixNano()
	msg.Hash = msg.computeHash(sharedKey)
	return &msg, nil
}

func (m *MessageContent) Encode(w io.Writer) (int, error) {
	n := 0
	if err := write(w, m.Type); err != nil {
		return n, wrapFieldError(err, "type")
	}
	n += 1
	if err := write(w, m.PayloadLen); err != nil {
		return n, wrapFieldError(err, "payload length")
	}
	n += 1
	written, err := w.Write(m.Payload[:m.PayloadLen])
	if err != nil {
		return n, wrapFieldError(err, "payload")
	}
	n += written
	if err := write(w, m.CreatedAt); err != nil {
		return n, wrapFieldError(err, "createdAt")
	}
	n += 8
	written, err = w.Write(m.Hash[:])
	if err != nil {
		return n, wrapFieldError(err, "hash")
	}
	n += written
	return n, nil
}

func (m *MessageContent) Decode(r io.Reader) (int, error) {
	n := 0
	if err := read(r, &m.Type); err != nil {
		return n, wrapFieldError(err, "type")
	}
	n += 1

	if err := read(r, &m.PayloadLen); err != nil {
		return n, wrapFieldError(err, "payload length")
	}
	n += 1
	if m.PayloadLen > MaxMessagePayloadSize {
		return n, errMessagePayloadExceed
	}
	if _, err := io.ReadFull(r, m.Payload[:m.PayloadLen]); err != nil {
		return n, wrapFieldError(err, "payload")
	}
	n += int(m.PayloadLen)
	if err := read(r, &m.CreatedAt); err != nil {
		return n, wrapFieldError(err, "createdAt")
	}
	n += 8
	readBytes, err := io.ReadFull(r, m.Hash[:])
	if err != nil {
		return n, wrapFieldError(err, "hash")
	}
	n += readBytes
	return n, nil
}

func (m *MessageContent) Bytes() []byte {
	buf := make([]byte, 0, 1+1+int(m.PayloadLen)+8)
	buf = append(buf, byte(m.Type))
	buf = append(buf, m.PayloadLen)
	buf = append(buf, m.Payload[:m.PayloadLen]...)
	ts := make([]byte, 8)
	binary.BigEndian.PutUint64(ts, uint64(m.CreatedAt))
	buf = append(buf, ts...)
	return buf
}

func (m *MessageContent) Verify(sharedKey []byte) bool {
	expected := m.computeHash(sharedKey)
	return crypto.HMACVerify(sha3.New256, sharedKey, m.Hash[:], expected[:])
}

func (m *MessageContent) computeHash(sharedKey []byte) [32]byte {
	sum := crypto.HMACSign(sha3.New256, sharedKey, m.Bytes())
	var out [32]byte
	copy(out[:], sum)
	return out
}

func (m *MessageContent) Len() int {
	return len(m.Bytes())
}

func (m *MessageContent) Max() int {
	return MaxMessagePayloadSize
}

func write(w io.Writer, data interface{}) error {
	return binary.Write(w, binary.BigEndian, data)
}

func read(r io.Reader, dst interface{}) error {
	return binary.Read(r, binary.BigEndian, dst)
}

func wrapFieldError(err error, field string) error {
	if err == nil {
		return nil
	}
	return &decErr{msg: fmt.Sprintf("field %s: %s", field, err.Error())}
}
