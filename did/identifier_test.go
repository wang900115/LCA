package did

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDocument(t *testing.T) {
	did := NewDID([]ServiceEndpoint{
		{
			ID:   "service1",
			Type: "Messaging",
			URL:  "https://example.com/msg",
		},
	})
	doc := did.ToDocument()
	assert.Equal(t, doc.ID, did.DID().ID)
	assert.Len(t, doc.VerificationMethod, 1)
	assert.Len(t, doc.KeyAgreement, 1)
	assert.Len(t, doc.Service, 1)
	t.Logf("Document: %+v\n", doc)
}

func TestSignAndVerifyDocument(t *testing.T) {
	did := NewDID([]ServiceEndpoint{
		{
			ID:   "service1",
			Type: "Messaging",
			URL:  "https://example.com/msg",
		},
	})
	signature, err := did.SignDocument()
	assert.NoError(t, err)
	valid, err := did.VerifyDocument(signature)
	assert.NoError(t, err)
	assert.True(t, valid)
}
