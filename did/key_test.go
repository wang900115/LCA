package did

import (
	"crypto/rand"
	"testing"
)

func TestKey(t *testing.T) {
	_, err := NewPeerKeyPair(rand.Reader)
	if err != nil {
		t.Fatalf("Failed to generate PeerKeyPair: %v", err)
	}
}

func TestGenerateID(t *testing.T) {
	keyPair, err := NewPeerKeyPair(rand.Reader)
	if err != nil {
		t.Fatalf("Failed to generate PeerKeyPair: %v", err)
	}
	did := keyPair.GenerateDID()
	if len(did) == 0 {
		t.Fatalf("Generated DID is empty")
	}
	t.Logf("Generated DID: %s", did)
}

func TestGenerateAddr(t *testing.T) {
	keyPair, err := NewPeerKeyPair(rand.Reader)
	if err != nil {
		t.Fatalf("Failed to generate PeerKeyPair: %v", err)
	}
	addr := keyPair.GenerateAddr()
	if len(addr) == 0 {
		t.Fatalf("Generated address is empty")
	}
	t.Logf("Generated address: %s", addr)
}
