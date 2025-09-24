// Copyright 2022 The go-ethereum Authors
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

package que

import (
	"math/rand"
	"sort"
	"testing"
)

func TestSstack(t *testing.T) {
	// Create some initial data
	size := 16 * blockSize
	data := make([]*element[int64, int], size)
	for i := 0; i < size; i++ {
		data[i] = &element[int64, int]{rand.Int(), rand.Int63()}
	}
	stack := newSstack[int64, int](nil)
	for rep := 0; rep < 2; rep++ {
		secs := []*element[int64, int]{}
		for i := 0; i < size; i++ {
			stack.Push(data[i])
			if i%2 == 0 {
				secs = append(secs, stack.Pop().(*element[int64, int]))
			}
		}
		rest := []*element[int64, int]{}
		for stack.Len() > 0 {
			rest = append(rest, stack.Pop().(*element[int64, int]))
		}
		for i := 0; i < size; i++ {
			if i%2 == 0 && data[i] != secs[i/2] {
				t.Errorf("push/pop mismatch: have %v, want %v.", secs[i/2], data[i])
			}
			if i%2 == 1 && data[i] != rest[len(rest)-i/2-1] {
				t.Errorf("push/pop mismatch: have %v, want %v.", rest[len(rest)-i/2-1], data[i])
			}
		}
	}
}

func TestSstackSort(t *testing.T) {
	size := 16 * blockSize
	data := make([]*element[int64, int], size)
	for i := 0; i < size; i++ {
		data[i] = &element[int64, int]{rand.Int(), int64(i)}
	}
	stack := newSstack[int64, int](nil)
	for _, val := range data {
		stack.Push(val)
	}
	sort.Sort(stack)
	for _, val := range data {
		out := stack.Pop()
		if out != val {
			t.Errorf("push/pop mismatch after sort: have %v, want %v.", out, val)
		}
	}
}

func TestSstackReset(t *testing.T) {
	size := 16 * blockSize
	data := make([]*element[int64, int], size)
	for i := 0; i < size; i++ {
		data[i] = &element[int64, int]{rand.Int(), rand.Int63()}
	}
	stack := newSstack[int64, int](nil)
	for rep := 0; rep < 2; rep++ {
		secs := []*element[int64, int]{}
		for i := 0; i < size; i++ {
			stack.Push(data[i])
			if i%2 == 0 {
				secs = append(secs, stack.Pop().(*element[int64, int]))
			}
		}
		stack.Reset()
		if stack.Len() != 0 {
			t.Errorf("stack not empty after reset: %v", stack)
		}
		for i := 0; i < size; i++ {
			if i%2 == 0 && data[i] != secs[i/2] {
				t.Errorf("push/pop mismatch: have %v, want %v.", secs[i/2], data[i])
			}
		}
	}
}
