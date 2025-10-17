package network

import (
	"encoding/binary"
	"errors"
	"io"
	"time"
)

type Message interface {
	Encode(w io.Writer) (int, error)
	Decode(r io.Reader) (int, error)
	Bytes() []byte
}

type MessageContent struct {
	Type      byte   // Type of message (e.g., "text", "file", etc.)
	Payload   []byte // Actual message content or data
	CreatedAt int64  // Timestamp of message creation
}

func NewMessageContent(msgType byte, payload []byte) *MessageContent {
	return &MessageContent{
		Type:      msgType,
		Payload:   payload,
		CreatedAt: time.Now().UTC().UnixNano(),
	}
}

func (m *MessageContent) Encode(w io.Writer) (int, error) {
	n := 0
	if err := write(w, m.Type); err != nil {
		return n, err
	}
	n += 1
	if written, err := w.Write(m.Payload); err != nil {
		return n, err
	} else {
		n += written
	}
	if err := write(w, m.CreatedAt); err != nil {
		return n, err
	}
	n += 8
	return n, nil
}

func (m *MessageContent) Decode(r io.Reader) (int, error) {
	n := 0
	if err := read(r, m.Type); err != nil {
		return n, err
	}
	n += 1
	payload, err := io.ReadAll(r)
	if err != nil {
		return n, err
	}
	if len(payload) < 8 {
		return n, errors.New("invalid message: missing CreatedAt")
	}
	payloadLen := len(payload) - 8
	m.Payload = payload[:payloadLen]
	m.CreatedAt = int64(binary.BigEndian.Uint64(payload[payloadLen:]))
	n += len(payload)
	return n, nil
}

func (m *MessageContent) Bytes() []byte {
	buf := make([]byte, 0, 1+len(m.Payload)+8)
	buf = append(buf, m.Type)
	buf = append(buf, m.Payload...)

	ts := make([]byte, 8)
	binary.BigEndian.PutUint64(ts, uint64(m.CreatedAt))
	buf = append(buf, ts...)
	return buf
}

func write(w io.Writer, data interface{}) error {
	if err := binary.Write(w, binary.BigEndian, data); err != nil {
		return err
	}
	return nil
}

func read(r io.Reader, dst interface{}) error {
	if err := binary.Read(r, binary.BigEndian, &dst); err != nil {
		return err
	}
	return nil
}
