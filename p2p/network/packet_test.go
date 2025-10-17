package network

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wang900115/LCA/crypt/did"
	common "github.com/wang900115/LCA/p2p/com"
)

func TestPacketEncodeDecode(t *testing.T) {
	msg, err := NewMessageContent(common.PUBLIC, []byte("Hello Packet"), []byte("SHARED"))
	assert.NoError(t, err)
	d := did.NewDID([]did.ServiceEndpoint{})
	rpc, err := NewRPCContent(msg, d)
	assert.NoError(t, err)
	packet, err := NewPacket(common.HEARTBEAT, rpc)
	assert.NoError(t, err)
	assert.NotNil(t, packet)
	var buf bytes.Buffer
	n, err := packet.Encode(&buf)
	assert.NoError(t, err)
	assert.Equal(t, n, packet.Len())
	var decoded PacketContent
	n2, err := decoded.Decode(&buf)
	assert.NoError(t, err)
	assert.Equal(t, n2, packet.Len())
	assert.NoError(t, decoded.Check())
	assert.Equal(t, packet.(*PacketContent).Payload[:packet.(*PacketContent).PayloadLen],
		decoded.Payload[:decoded.PayloadLen])
}

func TestPacketChecksumFail(t *testing.T) {
	msg, _ := NewMessageContent(common.PUBLIC, []byte("Test"), []byte("SHARED"))
	d := did.NewDID([]did.ServiceEndpoint{})
	rpc, _ := NewRPCContent(msg, d)
	packet, _ := NewPacket(common.HEARTBEAT, rpc)

	packetContent := packet.(*PacketContent)
	packetContent.Payload[0] ^= 0xFF

	err := packetContent.Check()
	assert.Error(t, err)
	assert.Equal(t, errPacketChecksumFail, err)
}
