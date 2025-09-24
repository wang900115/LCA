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
	"cmp"
	"container/heap"
)

type Prque[P cmp.Ordered, V any] struct {
	cont *sstack[P, V]
}

func New[P cmp.Ordered, V any](setIndex SetIndexCallback[V]) *Prque[P, V] {
	return &Prque[P, V]{newSstack[P, V](setIndex)}
}

func (p *Prque[P, V]) Push(data V, priority P) {
	heap.Push(p.cont, &element[P, V]{data, priority})
}

func (p *Prque[P, V]) Pop() (V, P) {
	element := heap.Pop(p.cont).(*element[P, V])
	return element.value, element.priority
}

func (p *Prque[P, V]) Peek() (V, P) {
	element := p.cont.blocks[0][0]
	return element.value, element.priority
}

func (p *Prque[P, V]) PopItem() V {
	return heap.Pop(p.cont).(*element[P, V]).value
}

func (p *Prque[P, V]) Remove(i int) V {
	return heap.Remove(p.cont, i).(*element[P, V]).value
}

func (p *Prque[P, V]) Empty() bool {
	return p.cont.Len() == 0
}

func (p *Prque[P, V]) Size() int {
	return p.cont.Len()
}

func (p *Prque[P, V]) Reset() {
	*p = *New[P, V](p.cont.setIndex)
}
