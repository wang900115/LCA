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

import "cmp"

const blockSize = 4096

type element[P cmp.Ordered, V any] struct {
	value    V
	priority P
}

type SetIndexCallback[V any] func(data V, index int)

type sstack[P cmp.Ordered, V any] struct {
	setIndex SetIndexCallback[V]
	size     int // 目前堆疊裡實際存了多少元素
	capacity int // 目前堆疊裡總共可以容納多少元素
	offset   int

	blocks [][]*element[P, V]
	active []*element[P, V]
}

func newSstack[P cmp.Ordered, V any](setIndex SetIndexCallback[V]) *sstack[P, V] {
	result := new(sstack[P, V])
	result.setIndex = setIndex
	result.active = make([]*element[P, V], blockSize)
	result.blocks = [][]*element[P, V]{result.active}
	result.capacity = blockSize
	return result
}

func (s *sstack[P, V]) Push(data any) {
	// 達到容量上限 -> 新增一個block
	if s.size == s.capacity {
		s.active = make([]*element[P, V], blockSize)
		s.blocks = append(s.blocks, s.active)
		s.capacity += blockSize
		s.offset = 0
		// offset到底(block滿ㄌ) -> 換到下一個block的起點
	} else if s.offset == blockSize {
		s.active = s.blocks[s.size/blockSize]
		s.offset = 0
	}

	if s.setIndex != nil {
		s.setIndex(data.(*element[P, V]).value, s.size)
	}
	s.active[s.offset] = data.(*element[P, V])
	s.offset++
	s.size++
}

func (s *sstack[P, V]) Pop() (data any) {
	s.size--
	s.offset--
	if s.offset < 0 {
		s.offset = blockSize - 1
		s.active = s.blocks[s.size/blockSize]
	}
	data, s.active[s.offset] = s.active[s.offset], nil
	if s.setIndex != nil {
		s.setIndex(data.(*element[P, V]).value, -1)
	}
	return
}

func (s *sstack[P, V]) Len() int {
	return s.size
}

func (s *sstack[P, V]) Less(i, j int) bool {
	return s.blocks[i/blockSize][i%blockSize].priority > s.blocks[j/blockSize][j%blockSize].priority
}

func (s *sstack[P, V]) Swap(i, j int) {
	ib, io, jb, jo := i/blockSize, i%blockSize, j/blockSize, j%blockSize
	a, b := s.blocks[jb][jo], s.blocks[ib][io]
	if s.setIndex != nil {
		s.setIndex(a.value, i)
		s.setIndex(b.value, j)
	}
	s.blocks[ib][io], s.blocks[jb][jo] = a, b
}

func (s *sstack[P, V]) Reset() {
	*s = *newSstack[P, V](s.setIndex)
}
