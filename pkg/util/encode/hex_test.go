// Copyright 2016 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package encode

import (
	"fmt"
	"math/big"
	"reflect"
	"testing"
)

// TestEncode tests the Encode function
func TestEncode(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected string
	}{
		{"empty", []byte{}, "0x"},
		{"single byte", []byte{0x42}, "0x42"},
		{"multiple bytes", []byte{0x01, 0x23, 0x45}, "0x012345"},
		{"all zeros", []byte{0x00, 0x00}, "0x0000"},
		{"all max", []byte{0xFF, 0xFF}, "0xffff"},
		{"hello world", []byte("hello"), "0x68656c6c6f"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Encode(tt.input)
			if result != tt.expected {
				t.Errorf("Encode(%x) = %s, want %s", tt.input, result, tt.expected)
			}
		})
	}
}

// TestDecode tests the Decode function
func TestDecode(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		expected  []byte
		shouldErr bool
		errType   error
	}{
		{"valid hex", "0x42", []byte{0x42}, false, nil},
		{"valid multi-byte", "0x012345", []byte{0x01, 0x23, 0x45}, false, nil},
		{"valid empty", "0x", []byte{}, false, nil},
		{"valid lowercase", "0xabcdef", []byte{0xab, 0xcd, 0xef}, false, nil},
		{"valid uppercase", "0xABCDEF", []byte{0xab, 0xcd, 0xef}, false, nil},
		{"missing prefix", "42", nil, true, ErrMissingPrefix},
		{"empty string", "", nil, true, ErrEmptyString},
		{"odd length", "0x123", nil, true, ErrOddLength},
		{"invalid hex", "0xZZ", nil, true, ErrSyntax},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Decode(tt.input)
			if tt.shouldErr {
				if err == nil {
					t.Errorf("Decode(%s) expected error, got nil", tt.input)
				}
				if tt.errType != nil && err != tt.errType {
					t.Errorf("Decode(%s) expected error %v, got %v", tt.input, tt.errType, err)
				}
			} else {
				if err != nil {
					t.Errorf("Decode(%s) unexpected error: %v", tt.input, err)
				}
				if !reflect.DeepEqual(result, tt.expected) {
					t.Errorf("Decode(%s) = %x, want %x", tt.input, result, tt.expected)
				}
			}
		})
	}
}

// TestMustDecode tests the MustDecode function
func TestMustDecode(t *testing.T) {
	// Test valid case
	result := MustDecode("0x42")
	expected := []byte{0x42}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("MustDecode(\"0x42\") = %x, want %x", result, expected)
	}

	// Test panic case
	defer func() {
		if r := recover(); r == nil {
			t.Error("MustDecode with invalid input should panic")
		}
	}()
	MustDecode("invalid")
}

// TestEncodeUint64 tests the EncodeUint64 function
func TestEncodeUint64(t *testing.T) {
	tests := []struct {
		name     string
		input    uint64
		expected string
	}{
		{"zero", 0, "0x0"},
		{"small number", 42, "0x2a"},
		{"large number", 0xDEADBEEF, "0xdeadbeef"},
		{"max uint64", ^uint64(0), "0xffffffffffffffff"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := EncodeUint64(tt.input)
			if result != tt.expected {
				t.Errorf("EncodeUint64(%d) = %s, want %s", tt.input, result, tt.expected)
			}
		})
	}
}

// TestDecodeUint64 tests the DecodeUint64 function
func TestDecodeUint64(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		expected  uint64
		shouldErr bool
		errType   error
	}{
		{"zero", "0x0", 0, false, nil},
		{"small number", "0x2a", 42, false, nil},
		{"large number", "0xdeadbeef", 0xDEADBEEF, false, nil},
		{"max uint64", "0xffffffffffffffff", ^uint64(0), false, nil},
		{"missing prefix", "42", 0, true, ErrMissingPrefix},
		{"empty", "", 0, true, ErrEmptyString},
		{"empty number", "0x", 0, true, ErrEmptyNumber},
		{"leading zero", "0x01", 0, true, ErrLeadingZero},
		{"too large", "0x10000000000000000", 0, true, ErrUint64Range},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := DecodeUint64(tt.input)
			if tt.shouldErr {
				if err == nil {
					t.Errorf("DecodeUint64(%s) expected error, got nil", tt.input)
				}
				if tt.errType != nil && err != tt.errType {
					t.Errorf("DecodeUint64(%s) expected error %v, got %v", tt.input, tt.errType, err)
				}
			} else {
				if err != nil {
					t.Errorf("DecodeUint64(%s) unexpected error: %v", tt.input, err)
				}
				if result != tt.expected {
					t.Errorf("DecodeUint64(%s) = %d, want %d", tt.input, result, tt.expected)
				}
			}
		})
	}
}

// TestMustDecodeUint64 tests the MustDecodeUint64 function
func TestMustDecodeUint64(t *testing.T) {
	// Test valid case
	result := MustDecodeUint64("0x2a")
	expected := uint64(42)
	if result != expected {
		t.Errorf("MustDecodeUint64(\"0x2a\") = %d, want %d", result, expected)
	}

	// Test panic case
	defer func() {
		if r := recover(); r == nil {
			t.Error("MustDecodeUint64 with invalid input should panic")
		}
	}()
	MustDecodeUint64("invalid")
}

// TestEncodeBig tests the EncodeBig function
func TestEncodeBig(t *testing.T) {
	tests := []struct {
		name     string
		input    *big.Int
		expected string
	}{
		{"zero", big.NewInt(0), "0x0"},
		{"positive small", big.NewInt(42), "0x2a"},
		{"positive large", big.NewInt(0).SetBytes([]byte{0xDE, 0xAD, 0xBE, 0xEF}), "0xdeadbeef"},
		{"negative small", big.NewInt(-42), "-0x2a"},
		{"negative large", big.NewInt(0).Neg(big.NewInt(0).SetBytes([]byte{0xDE, 0xAD, 0xBE, 0xEF})), "-0xdeadbeef"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := EncodeBig(tt.input)
			if result != tt.expected {
				t.Errorf("EncodeBig(%s) = %s, want %s", tt.input.String(), result, tt.expected)
			}
		})
	}
}

// TestDecodeBig tests the DecodeBig function
func TestDecodeBig(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		expected  *big.Int
		shouldErr bool
		errType   error
	}{
		{"zero", "0x0", big.NewInt(0), false, nil},
		{"positive small", "0x2a", big.NewInt(42), false, nil},
		{"positive large", "0xdeadbeef", big.NewInt(0).SetBytes([]byte{0xde, 0xad, 0xbe, 0xef}), false, nil},
		{"missing prefix", "42", nil, true, ErrMissingPrefix},
		{"empty", "", nil, true, ErrEmptyString},
		{"empty number", "0x", nil, true, ErrEmptyNumber},
		{"leading zero", "0x01", nil, true, ErrLeadingZero},
		{"too large", "0x" + string(make([]byte, 65)), nil, true, ErrBig256Range},
		{"invalid hex", "0xZZ", nil, true, ErrSyntax},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := DecodeBig(tt.input)
			if tt.shouldErr {
				if err == nil {
					t.Errorf("DecodeBig(%s) expected error, got nil", tt.input)
				}
				if tt.errType != nil && err != tt.errType {
					t.Errorf("DecodeBig(%s) expected error %v, got %v", tt.input, tt.errType, err)
				}
			} else {
				if err != nil {
					t.Errorf("DecodeBig(%s) unexpected error: %v", tt.input, err)
				}
				if result.Cmp(tt.expected) != 0 {
					t.Errorf("DecodeBig(%s) = %s, want %s", tt.input, result.String(), tt.expected.String())
				}
			}
		})
	}
}

// TestMustDecodeBig tests the MustDecodeBig function
func TestMustDecodeBig(t *testing.T) {
	// Test valid case
	result := MustDecodeBig("0x2a")
	expected := big.NewInt(42)
	if result.Cmp(expected) != 0 {
		t.Errorf("MustDecodeBig(\"0x2a\") = %s, want %s", result.String(), expected.String())
	}

	// Test panic case
	defer func() {
		if r := recover(); r == nil {
			t.Error("MustDecodeBig with invalid input should panic")
		}
	}()
	MustDecodeBig("invalid")
}

// TestHas0xPrefix tests the has0xPrefix function
func TestHas0xPrefix(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"with 0x", "0x123", true},
		{"with 0X", "0X123", true},
		{"without prefix", "123", false},
		{"empty", "", false},
		{"only 0", "0", false},
		{"x only", "x123", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := has0xPrefix(tt.input)
			if result != tt.expected {
				t.Errorf("has0xPrefix(%s) = %t, want %t", tt.input, result, tt.expected)
			}
		})
	}
}

// TestDecodeNibble tests the decodeNibble function
func TestDecodeNibble(t *testing.T) {
	tests := []struct {
		name     string
		input    byte
		expected uint64
	}{
		{"digit 0", '0', 0},
		{"digit 9", '9', 9},
		{"upper A", 'A', 10},
		{"upper F", 'F', 15},
		{"lower a", 'a', 10},
		{"lower f", 'f', 15},
		{"invalid G", 'G', badNibble},
		{"invalid z", 'z', badNibble},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := decodeNibble(tt.input)
			if result != tt.expected {
				t.Errorf("decodeNibble(%c) = %d, want %d", tt.input, result, tt.expected)
			}
		})
	}
}

// TestRoundTrip tests encode/decode round trips
func TestRoundTrip(t *testing.T) {
	tests := [][]byte{
		{},
		{0x00},
		{0xFF},
		{0x01, 0x23, 0x45, 0x67, 0x89, 0xAB, 0xCD, 0xEF},
		[]byte("hello world"),
	}

	for i, data := range tests {
		t.Run(fmt.Sprintf("test_%d", i), func(t *testing.T) {
			encoded := Encode(data)
			decoded, err := Decode(encoded)
			if err != nil {
				t.Errorf("Round trip failed: %v", err)
			}
			if !reflect.DeepEqual(data, decoded) {
				t.Errorf("Round trip failed: %x != %x", data, decoded)
			}
		})
	}
}

// TestBigWordNibbles tests that the init function sets bigWordNibbles correctly
func TestBigWordNibbles(t *testing.T) {
	if bigWordNibbles != 8 && bigWordNibbles != 16 {
		t.Errorf("bigWordNibbles = %d, expected 8 or 16", bigWordNibbles)
	}
}

// Benchmark tests
func BenchmarkEncode(b *testing.B) {
	data := []byte("hello world this is a test string for benchmarking")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Encode(data)
	}
}

func BenchmarkDecode(b *testing.B) {
	hex := "0x68656c6c6f20776f726c642074686973206973206120746573742073747269"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Decode(hex)
	}
}
