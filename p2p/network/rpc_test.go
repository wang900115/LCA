package network

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wang900115/LCA/crypt/did"
	common "github.com/wang900115/LCA/p2p/com"
)

func TestRPCEncode(t *testing.T) {
	message, err := NewMessageContent(common.PUBLIC, []byte("Ping"), []byte("SHARED"))
	assert.NoError(t, err)
	assert.NotNil(t, message)

	d := did.NewDID([]did.ServiceEndpoint{})
	assert.NotNil(t, d)

	rpc, err := NewRPCContent(message, d)
	assert.NoError(t, err)
	assert.NotNil(t, rpc)

	var buf bytes.Buffer
	n, err := rpc.Encode(&buf)
	assert.NoError(t, err)
	assert.Greater(t, n, 0)
	assert.Equal(t, n, rpc.Len(), "encoded length should match rpc.Len()")
	assert.LessOrEqual(t, n, MaxPacketPayloadSize, "payload must not exceed limit")

	t.Logf("Encoded RPC (hex): %x\n", rpc.Bytes())
}

func TestRPCDecode(t *testing.T) {
	message, err := NewMessageContent(common.PUBLIC, []byte("Ping"), []byte("SHARED"))
	assert.NoError(t, err)
	d := did.NewDID([]did.ServiceEndpoint{})
	assert.NotNil(t, d)

	original, err := NewRPCContent(message, d)
	assert.NoError(t, err)

	var buf bytes.Buffer
	_, err = original.Encode(&buf)
	assert.NoError(t, err)

	var decoded RPCContent
	_, err = decoded.Decode(&buf)
	assert.NoError(t, err)

	var expected [50]byte
	copy(expected[:], []byte(d.DIDInfo().Address))

	assert.Equal(t, expected, decoded.From, "source address mismatch")
	assert.Equal(t, uint8(message.Len()), decoded.PayloadLen, "payload length mismatch")
	assert.Equal(t, message.Bytes(), decoded.Payload[:decoded.PayloadLen], "payload mismatch")

	t.Logf("Decoded RPC: %+v", decoded)
}

func TestRPCVerify(t *testing.T) {
	message, _ := NewMessageContent(common.PUBLIC, []byte("Ping"), []byte("SHARED"))
	d := did.NewDID([]did.ServiceEndpoint{})
	rpc, _ := NewRPCContent(message, d)
	err := rpc.Verify(d.DIDInfo().KeyPair.EdPublic)
	assert.NoError(t, err)
	otherD := did.NewDID([]did.ServiceEndpoint{})
	err = rpc.Verify(otherD.DIDInfo().KeyPair.EdPublic)
	assert.Error(t, err)
}
