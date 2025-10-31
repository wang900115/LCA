package enode

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math/bits"
	"net/netip"
	"strings"
)

var errMissingPrefix = errors.New("missing 'enr:' prefix for base64-encoded record")

type Node struct {
	id       ID
	hostname string
	ip       netip.Addr
	udp      uint16
	tcp      uint16
}

func validIP(ip netip.Addr) bool {
	return ip.IsValid() && !ip.IsMulticast()
}

func localityScore(ip netip.Addr) int {
	switch {
	case ip.IsUnspecified():
		return 0
	case ip.IsLoopback():
		return 1
	case ip.IsLinkLocalMulticast():
		return 2
	case ip.IsPrivate():
		return 3
	default:
		return 4
	}
}

type ID [32]byte

func (n ID) Bytes() []byte {
	return n[:]
}

func (n ID) String() string {
	return fmt.Sprintf("%x", n[:])
}

func (n ID) GoString() string {
	return fmt.Sprintf("encode.HexID(\"%x\")", n[:])
}

func (n ID) TerminalString() string {
	return hex.EncodeToString(n[:8])
}

func (n ID) MarshalText() ([]byte, error) {
	return []byte(hex.EncodeToString(n[:])), nil
}

func HexID(in string) ID {
	id, err := ParseID(in)
	if err != nil {
		panic(err)
	}
	return id
}

func ParseID(in string) (ID, error) {
	var id ID
	b, err := hex.DecodeString(strings.TrimPrefix(in, "0x"))
	if err != nil {
		return id, err
	} else if len(b) != len(id) {
		return id, fmt.Errorf("invalid length %d for ID", len(b))
	}
	copy(id[:], b)
	return id, nil
}

// DistCmp compares the distances of two node IDs from a target ID.
func DistCmp(target, a, b ID) int {
	for i := range target {
		da := a[i] ^ target[i]
		db := b[i] ^ target[i]
		if da > db {
			return 1
		} else if da < db {
			return -1
		}
	}
	return 0
}

// LogDist returns the logarithmic distance between two node IDs.
func LogDist(a, b ID) int {
	lz := 0
	for i := range a {
		x := a[i] ^ b[i]
		if x == 0 {
			lz += 8
		} else {
			lz += bits.LeadingZeros8(x)
			break
		}
	}
	return len(a)*8 - lz
}
