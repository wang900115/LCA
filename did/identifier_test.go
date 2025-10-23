package did

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDocument(t *testing.T) {
	did := NewDIDIdentifier([]ServiceEndpoint{
		{
			ID:              "service1",
			Type:            "Messaging",
			ServiceEndpoint: "https://example.com/msg",
		},
	})
	doc := did.Document()
	assert.Len(t, doc.VerificationMethod, 2)
	assert.Len(t, doc.KeyAgreement, 1)
	assert.Len(t, doc.Service, 1)
	t.Logf("Document: %+v\n", doc)
}

func TestSign(t *testing.T) {
	did := NewDIDIdentifier([]ServiceEndpoint{
		{
			ID:              "service1",
			Type:            "Messaging",
			ServiceEndpoint: "https://example.com/msg",
		},
	})
	signature, err := did.SignDocument()
	assert.NoError(t, err)
	t.Logf("Signature: %x\n", signature)
}
