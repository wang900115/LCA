package network

import (
	"crypto/ed25519"
	"errors"
	"io"

	crypto "github.com/wang900115/LCA/crypt"
	"github.com/wang900115/LCA/crypt/did"
)

const (
	MaxRPCPayloadSize = 50
)

var (
	errRPCPayloadExceed = &decErr{"rpc payload is exceed 50 bytes"}
	errRPCPayloadVerify = errors.New("rpc payload signature verification failed")
)

type RPC interface {
	Encode(w io.Writer) (int, error)
	Decode(r io.Reader) (int, error)
	Verify(pubKey ed25519.PublicKey) error
	Bytes() []byte
	Len() int
	Max() int
}

type RPCContent struct {
	From       [50]byte
	PayloadLen uint8
	Payload    [MaxRPCPayloadSize]byte
	Sig        [64]byte
}

func NewRPCContent(msg Message, d did.PeerDID) (RPC, error) {
	if msg.Len() > MaxRPCPayloadSize {
		return nil, errRPCPayloadExceed
	}
	var rpc RPCContent
	copy(rpc.From[:], []byte(d.DIDInfo().Address))
	copy(rpc.Payload[:], msg.Bytes())
	rpc.PayloadLen = uint8(msg.Len())
	signature, err := crypto.ED25519Sign(d.DIDInfo().KeyPair.EdPrivate, rpc.dataToSign())
	if err != nil {
		return nil, err
	}
	copy(rpc.Sig[:], signature)
	return &rpc, nil
}

func (rpc *RPCContent) Encode(w io.Writer) (int, error) {
	n := 0
	written, err := w.Write(rpc.From[:])
	if err != nil {
		return n, wrapFieldError(err, "from")
	}
	n += written
	if err := write(w, rpc.PayloadLen); err != nil {
		return n, wrapFieldError(err, "payload length")
	}
	n += 1
	written, err = w.Write(rpc.Payload[:rpc.PayloadLen])
	if err != nil {
		return n, wrapFieldError(err, "payload")
	}
	n += written
	written, err = w.Write(rpc.Sig[:])
	if err != nil {
		return n, wrapFieldError(err, "signature")
	}
	n += written
	return n, nil
}

func (rpc *RPCContent) Decode(r io.Reader) (int, error) {
	n := 0
	readBytes, err := io.ReadFull(r, rpc.From[:])
	if err != nil {
		return n, wrapFieldError(err, "from")
	}
	n += readBytes
	var lengthBuf [1]byte
	if _, err := io.ReadFull(r, lengthBuf[:]); err != nil {
		return n, wrapFieldError(err, "payload length")
	}
	rpc.PayloadLen = lengthBuf[0]
	n += 1
	if rpc.PayloadLen > MaxRPCPayloadSize {
		return n, errRPCPayloadExceed
	}
	readBytes, err = io.ReadFull(r, rpc.Payload[:rpc.PayloadLen])
	if err != nil {
		return n, wrapFieldError(err, "payload")
	}
	n += readBytes
	readBytes, err = io.ReadFull(r, rpc.Sig[:])
	if err != nil {
		return n, wrapFieldError(err, "signature")
	}
	n += readBytes
	return n, nil
}

func (rpc *RPCContent) dataToSign() []byte {
	return append(rpc.From[:], rpc.Payload[:rpc.PayloadLen]...)
}

func (rpc *RPCContent) Verify(pub ed25519.PublicKey) error {
	ok, err := crypto.ED25519Verify(pub, rpc.dataToSign(), rpc.Sig[:])
	if err != nil {
		return err
	}
	if ok {
		return nil
	}
	return errRPCPayloadVerify
}

func (rpc *RPCContent) Bytes() []byte {
	b := make([]byte, 0, 50+1+int(rpc.PayloadLen)+64)
	b = append(b, rpc.From[:]...)
	b = append(b, rpc.PayloadLen)
	b = append(b, rpc.Payload[:rpc.PayloadLen]...)
	b = append(b, rpc.Sig[:]...)
	return b
}

func (rpc *RPCContent) Len() int {
	return int(len(rpc.Bytes()))
}

func (rpc *RPCContent) Max() int {
	return MaxRPCPayloadSize
}
