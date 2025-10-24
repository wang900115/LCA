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
	"encoding/json"
	"math/big"
	"reflect"
	"testing"

	"github.com/holiman/uint256"
)

// TestBytesMarshaling tests Bytes marshaling and unmarshaling
func TestBytesMarshaling(t *testing.T) {
	tests := []struct {
		name     string
		input    Bytes
		expected string
	}{
		{"empty", Bytes{}, `"0x"`},
		{"single byte", Bytes{0x42}, `"0x42"`},
		{"multiple bytes", Bytes{0x01, 0x23, 0x45}, `"0x012345"`},
		{"hello", Bytes("hello"), `"0x68656c6c6f"`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := json.Marshal(tt.input)
			if err != nil {
				t.Errorf("Marshal(%x) unexpected error: %v", tt.input, err)
			}
			if string(result) != tt.expected {
				t.Errorf("Marshal(%x) = %s, want %s", tt.input, result, tt.expected)
			}
		})
	}
}

func TestBytesUnmarshaling(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		expected  Bytes
		shouldErr bool
	}{
		{"valid empty", `"0x"`, Bytes{}, false},
		{"valid single", `"0x42"`, Bytes{0x42}, false},
		{"valid multiple", `"0x012345"`, Bytes{0x01, 0x23, 0x45}, false},
		{"valid lowercase", `"0xabcdef"`, Bytes{0xab, 0xcd, 0xef}, false},
		{"valid uppercase", `"0xABCDEF"`, Bytes{0xab, 0xcd, 0xef}, false},
		{"invalid non-string", `42`, nil, true},
		{"invalid missing prefix", `"42"`, nil, true},
		{"invalid odd length", `"0x123"`, nil, true},
		{"invalid hex chars", `"0xZZ"`, nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result Bytes
			err := json.Unmarshal([]byte(tt.input), &result)
			if tt.shouldErr {
				if err == nil {
					t.Errorf("Unmarshal(%s) expected error, got nil", tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("Unmarshal(%s) unexpected error: %v", tt.input, err)
				}
				if !reflect.DeepEqual(result, tt.expected) {
					t.Errorf("Unmarshal(%s) = %x, want %x", tt.input, result, tt.expected)
				}
			}
		})
	}
}

func TestBytesString(t *testing.T) {
	b := Bytes{0x01, 0x23, 0x45}
	expected := "0x012345"
	result := b.String()
	if result != expected {
		t.Errorf("String() = %s, want %s", result, expected)
	}
}

func TestBytesGraphQL(t *testing.T) {
	var b Bytes

	// Test ImplementsGraphQLType
	if !b.ImplementsGraphQLType("Bytes") {
		t.Error("ImplementsGraphQLType(\"Bytes\") should return true")
	}
	if b.ImplementsGraphQLType("String") {
		t.Error("ImplementsGraphQLType(\"String\") should return false")
	}

	// Test UnmarshalGraphQL with string
	err := b.UnmarshalGraphQL("0x42")
	if err != nil {
		t.Errorf("UnmarshalGraphQL(\"0x42\") unexpected error: %v", err)
	}
	expected := Bytes{0x42}
	if !reflect.DeepEqual(b, expected) {
		t.Errorf("UnmarshalGraphQL(\"0x42\") = %x, want %x", b, expected)
	}

	// Test UnmarshalGraphQL with invalid type
	err = b.UnmarshalGraphQL(42)
	if err == nil {
		t.Error("UnmarshalGraphQL(42) expected error, got nil")
	}
}

// TestBigMarshaling tests Big marshaling and unmarshaling
func TestBigMarshaling(t *testing.T) {
	tests := []struct {
		name     string
		input    *Big
		expected string
	}{
		{"zero", (*Big)(big.NewInt(0)), `"0x0"`},
		{"small positive", (*Big)(big.NewInt(42)), `"0x2a"`},
		{"large positive", (*Big)(big.NewInt(0xDEADBEEF)), `"0xdeadbeef"`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := json.Marshal(tt.input)
			if err != nil {
				t.Errorf("Marshal(%s) unexpected error: %v", tt.input.String(), err)
			}
			if string(result) != tt.expected {
				t.Errorf("Marshal(%s) = %s, want %s", tt.input.String(), result, tt.expected)
			}
		})
	}
}

func TestBigUnmarshaling(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		expected  *big.Int
		shouldErr bool
	}{
		{"zero", `"0x0"`, big.NewInt(0), false},
		{"small positive", `"0x2a"`, big.NewInt(42), false},
		{"large positive", `"0xdeadbeef"`, big.NewInt(0xDEADBEEF), false},
		{"invalid non-string", `42`, nil, true},
		{"invalid missing prefix", `"42"`, nil, true},
		{"invalid leading zero", `"0x01"`, nil, true},
		{"too large", `"0x` + string(make([]byte, 65)) + `"`, nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result Big
			err := json.Unmarshal([]byte(tt.input), &result)
			if tt.shouldErr {
				if err == nil {
					t.Errorf("Unmarshal(%s) expected error, got nil", tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("Unmarshal(%s) unexpected error: %v", tt.input, err)
				}
				if result.ToInt().Cmp(tt.expected) != 0 {
					t.Errorf("Unmarshal(%s) = %s, want %s", tt.input, result.String(), tt.expected.String())
				}
			}
		})
	}
}

func TestBigGraphQL(t *testing.T) {
	var b Big

	// Test ImplementsGraphQLType
	if !b.ImplementsGraphQLType("BigInt") {
		t.Error("ImplementsGraphQLType(\"BigInt\") should return true")
	}

	// Test UnmarshalGraphQL with string
	err := b.UnmarshalGraphQL("0x2a")
	if err != nil {
		t.Errorf("UnmarshalGraphQL(\"0x2a\") unexpected error: %v", err)
	}
	expected := big.NewInt(42)
	if b.ToInt().Cmp(expected) != 0 {
		t.Errorf("UnmarshalGraphQL(\"0x2a\") = %s, want %s", b.String(), expected.String())
	}

	// Test UnmarshalGraphQL with int32
	err = b.UnmarshalGraphQL(int32(100))
	if err != nil {
		t.Errorf("UnmarshalGraphQL(100) unexpected error: %v", err)
	}
	expected = big.NewInt(100)
	if b.ToInt().Cmp(expected) != 0 {
		t.Errorf("UnmarshalGraphQL(100) = %s, want %s", b.String(), expected.String())
	}
}

// TestUint64Marshaling tests Uint64 marshaling and unmarshaling
func TestUint64Marshaling(t *testing.T) {
	tests := []struct {
		name     string
		input    Uint64
		expected string
	}{
		{"zero", Uint64(0), `"0x0"`},
		{"small", Uint64(42), `"0x2a"`},
		{"large", Uint64(0xDEADBEEF), `"0xdeadbeef"`},
		{"max", Uint64(^uint64(0)), `"0xffffffffffffffff"`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := json.Marshal(tt.input)
			if err != nil {
				t.Errorf("Marshal(%d) unexpected error: %v", tt.input, err)
			}
			if string(result) != tt.expected {
				t.Errorf("Marshal(%d) = %s, want %s", tt.input, result, tt.expected)
			}
		})
	}
}

func TestUint64Unmarshaling(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		expected  Uint64
		shouldErr bool
	}{
		{"zero", `"0x0"`, Uint64(0), false},
		{"small", `"0x2a"`, Uint64(42), false},
		{"large", `"0xdeadbeef"`, Uint64(0xDEADBEEF), false},
		{"max", `"0xffffffffffffffff"`, Uint64(^uint64(0)), false},
		{"invalid non-string", `42`, Uint64(0), true},
		{"invalid missing prefix", `"42"`, Uint64(0), true},
		{"invalid leading zero", `"0x01"`, Uint64(0), true},
		{"too large", `"0x10000000000000000"`, Uint64(0), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result Uint64
			err := json.Unmarshal([]byte(tt.input), &result)
			if tt.shouldErr {
				if err == nil {
					t.Errorf("Unmarshal(%s) expected error, got nil", tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("Unmarshal(%s) unexpected error: %v", tt.input, err)
				}
				if result != tt.expected {
					t.Errorf("Unmarshal(%s) = %d, want %d", tt.input, result, tt.expected)
				}
			}
		})
	}
}

func TestUint64GraphQL(t *testing.T) {
	var u Uint64

	// Test ImplementsGraphQLType
	if !u.ImplementsGraphQLType("Long") {
		t.Error("ImplementsGraphQLType(\"Long\") should return true")
	}

	// Test UnmarshalGraphQL with string
	err := u.UnmarshalGraphQL("0x2a")
	if err != nil {
		t.Errorf("UnmarshalGraphQL(\"0x2a\") unexpected error: %v", err)
	}
	if u != Uint64(42) {
		t.Errorf("UnmarshalGraphQL(\"0x2a\") = %d, want %d", u, 42)
	}

	// Test UnmarshalGraphQL with int32
	err = u.UnmarshalGraphQL(int32(100))
	if err != nil {
		t.Errorf("UnmarshalGraphQL(100) unexpected error: %v", err)
	}
	if u != Uint64(100) {
		t.Errorf("UnmarshalGraphQL(100) = %d, want %d", u, 100)
	}
}

// TestUintMarshaling tests Uint marshaling and unmarshaling
func TestUintMarshaling(t *testing.T) {
	tests := []struct {
		name     string
		input    Uint
		expected string
	}{
		{"zero", Uint(0), `"0x0"`},
		{"small", Uint(42), `"0x2a"`},
		{"large", Uint(0xDEADBEEF), `"0xdeadbeef"`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := json.Marshal(tt.input)
			if err != nil {
				t.Errorf("Marshal(%d) unexpected error: %v", tt.input, err)
			}
			if string(result) != tt.expected {
				t.Errorf("Marshal(%d) = %s, want %s", tt.input, result, tt.expected)
			}
		})
	}
}

func TestUintUnmarshaling(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		expected  Uint
		shouldErr bool
	}{
		{"zero", `"0x0"`, Uint(0), false},
		{"small", `"0x2a"`, Uint(42), false},
		{"large", `"0xdeadbeef"`, Uint(0xDEADBEEF), false},
		{"invalid non-string", `42`, Uint(0), true},
		{"invalid missing prefix", `"42"`, Uint(0), true},
		{"invalid leading zero", `"0x01"`, Uint(0), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result Uint
			err := json.Unmarshal([]byte(tt.input), &result)
			if tt.shouldErr {
				if err == nil {
					t.Errorf("Unmarshal(%s) expected error, got nil", tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("Unmarshal(%s) unexpected error: %v", tt.input, err)
				}
				if result != tt.expected {
					t.Errorf("Unmarshal(%s) = %d, want %d", tt.input, result, tt.expected)
				}
			}
		})
	}
}

// TestU256Marshaling tests U256 marshaling and unmarshaling
func TestU256Marshaling(t *testing.T) {
	tests := []struct {
		name     string
		input    *U256
		expected string
	}{
		{"zero", (*U256)(uint256.NewInt(0)), `"0x0"`},
		{"small", (*U256)(uint256.NewInt(42)), `"0x2a"`},
		{"large", (*U256)(uint256.NewInt(0xDEADBEEF)), `"0xdeadbeef"`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := json.Marshal(tt.input)
			if err != nil {
				t.Errorf("Marshal(%s) unexpected error: %v", tt.input.String(), err)
			}
			if string(result) != tt.expected {
				t.Errorf("Marshal(%s) = %s, want %s", tt.input.String(), result, tt.expected)
			}
		})
	}
}

func TestU256Unmarshaling(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		expected  *uint256.Int
		shouldErr bool
	}{
		{"zero", `"0x0"`, uint256.NewInt(0), false},
		{"empty", `""`, uint256.NewInt(0), false},
		{"small", `"0x2a"`, uint256.NewInt(42), false},
		{"large", `"0xdeadbeef"`, uint256.NewInt(0xDEADBEEF), false},
		{"invalid non-string", `42`, nil, true},
		{"invalid hex", `"0xZZ"`, nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result U256
			err := json.Unmarshal([]byte(tt.input), &result)
			if tt.shouldErr {
				if err == nil {
					t.Errorf("Unmarshal(%s) expected error, got nil", tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("Unmarshal(%s) unexpected error: %v", tt.input, err)
				}
				if (*uint256.Int)(&result).Cmp(tt.expected) != 0 {
					t.Errorf("Unmarshal(%s) = %s, want %s", tt.input, result.String(), tt.expected.Hex())
				}
			}
		})
	}
}

// TestUnmarshalFixedJSON tests UnmarshalFixedJSON function
func TestUnmarshalFixedJSON(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		outLen    int
		shouldErr bool
	}{
		{"valid 4 bytes", `"0x01234567"`, 4, false},
		{"valid 8 bytes", `"0x0123456789abcdef"`, 8, false},
		{"invalid non-string", `42`, 4, true},
		{"invalid length", `"0x0123"`, 4, true},
		{"invalid hex", `"0xZZ"`, 1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := make([]byte, tt.outLen)
			err := UnmarshalFixedJSON(reflect.TypeOf(out), []byte(tt.input), out)
			if tt.shouldErr {
				if err == nil {
					t.Errorf("UnmarshalFixedJSON(%s) expected error, got nil", tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("UnmarshalFixedJSON(%s) unexpected error: %v", tt.input, err)
				}
			}
		})
	}
}

// TestHelperFunctions tests utility functions
func TestIsString(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected bool
	}{
		{"valid string", []byte(`"hello"`), true},
		{"empty string", []byte(`""`), true},
		{"not string", []byte(`42`), false},
		{"partial quote", []byte(`"hello`), false},
		{"empty", []byte{}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isString(tt.input)
			if result != tt.expected {
				t.Errorf("isString(%s) = %t, want %t", tt.input, result, tt.expected)
			}
		})
	}
}

func TestBytesHave0xPrefix(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected bool
	}{
		{"with 0x", []byte("0x123"), true},
		{"with 0X", []byte("0X123"), true},
		{"without prefix", []byte("123"), false},
		{"empty", []byte{}, false},
		{"only 0", []byte("0"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := bytesHave0xPrefix(tt.input)
			if result != tt.expected {
				t.Errorf("bytesHave0xPrefix(%s) = %t, want %t", tt.input, result, tt.expected)
			}
		})
	}
}

// Benchmark tests
func BenchmarkBytesMarshaling(b *testing.B) {
	data := Bytes("hello world this is a test string for benchmarking")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = json.Marshal(data)
	}
}

func BenchmarkBytesUnmarshaling(b *testing.B) {
	jsonData := []byte(`"0x68656c6c6f20776f726c642074686973206973206120746573742073747269"`)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var result Bytes
		_ = json.Unmarshal(jsonData, &result)
	}
}

func BenchmarkUint64Marshaling(b *testing.B) {
	data := Uint64(0xDEADBEEFCAFEBABE)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = json.Marshal(data)
	}
}

func BenchmarkUint64Unmarshaling(b *testing.B) {
	jsonData := []byte(`"0xdeadbeefcafebabe"`)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var result Uint64
		_ = json.Unmarshal(jsonData, &result)
	}
}
