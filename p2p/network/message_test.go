package network

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	common "github.com/wang900115/LCA/p2p/com"
)

func TestEncodeMessage(t *testing.T) {
	message, err := NewMessageContent(common.PUBLIC, []byte("Ping"), []byte("SHARED"))
	assert.Nil(t, err)
	assert.NotNil(t, message)
	var buf bytes.Buffer
	n, err := message.Encode(&buf)
	assert.Nil(t, err)
	assert.Greater(t, n, 0)
	assert.Less(t, message.Len(), message.Max())
	t.Logf("Encoded Message (hex): %x\n", message.Bytes())
}

func TestDecodeMessage(t *testing.T) {
	original, err := NewMessageContent(common.PUBLIC, []byte("Ping"), []byte("SHARED"))
	assert.Nil(t, err)
	var buf bytes.Buffer
	_, err = original.Encode(&buf)
	assert.Nil(t, err)
	var decoded MessageContent
	_, err = decoded.Decode(&buf)
	assert.Nil(t, err)
	assert.Equal(t, common.PUBLIC, decoded.Type)
	assert.Equal(t, uint8(len([]byte("Ping"))), decoded.PayloadLen)
	assert.Equal(t, []byte("Ping"), decoded.Payload[:decoded.PayloadLen])
	t.Logf("Decoded Message: %+v", decoded)
}
