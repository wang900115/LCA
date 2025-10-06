package encode

import (
	"encoding/gob"
	"io"

	"github.com/wang900115/LCA/p2p/packet"
)

type Decoder interface {
	Decode(io.Reader, *packet.RPC) error
}

type GOBDecoder struct{}

func (dec *GOBDecoder) Decode(r io.Reader, msg *packet.RPC) error {
	return gob.NewDecoder(r).Decode(msg)
}

type DefaultDecoder struct{}

func (dec DefaultDecoder) Decode(r io.Reader, msg *packet.RPC) error {
	peekBuf := make([]byte, 1)
	if _, err := r.Read(peekBuf); err != nil {
		return nil
	}

	stream := peekBuf[0] == packet.INCOMMINGSTREAM
	if stream {
		msg.Stream = true
		return nil
	}

	buf := make([]byte, 1028)
	n, err := r.Read(buf)
	if err != nil {
		return err
	}
	msg.Payload = buf[:n]
	return nil
}
