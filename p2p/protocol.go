package p2p

import (
	"cmp"
	"fmt"
	"strings"
)

type Protocol struct {
	Name     string
	Version  uint
	Length   uint64
	NodeInfo func() interface{}
}

func (p Protocol) cap() Cap {
	return Cap{Name: p.Name, Version: p.Version}
}

type Cap struct {
	Name    string
	Version uint
}

func (c Cap) String() string {
	return fmt.Sprintf("%s/%d", c.Name, c.Version)
}

func (c Cap) Cmp(other Cap) int {
	if c.Name == other.Name {
		return cmp.Compare(c.Version, other.Version)
	}
	return strings.Compare(c.Name, other.Name)
}
