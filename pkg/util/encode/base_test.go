package encode

import "testing"

func TestBase58Encode(t *testing.T) {
	input := []byte("hello world")
	encoded := Base58Encode(input)
	expected := "StV1DL6CwTryKyV"
	if encoded != expected {
		t.Errorf("Base58Encode(%s) = %s; want %s", input, encoded, expected)
	}
}
func TestBase58Decode(t *testing.T) {
	input := "StV1DL6CwTryKyV"
	decoded := Base58Decode(input)
	expected := []byte("hello world")
	if string(decoded) != string(expected) {
		t.Errorf("Base58Decode(%s) = %s; want %s", input, decoded, expected)
	}
}
