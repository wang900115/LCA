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

package lru

type list[T any] struct {
	root listElem[T]
}

type listElem[T any] struct {
	next *listElem[T]
	prev *listElem[T]
	v    T
}

func newList[T any]() *list[T] {
	l := new(list[T])
	l.init()
	return l
}

func (l *list[T]) init() {
	l.root.next = &l.root
	l.root.prev = &l.root
}

func (l *list[T]) pushElem(e *listElem[T]) {
	e.prev = &l.root
	e.next = l.root.next
	l.root.next = e
	e.next.prev = e
}

func (l *list[T]) moveToFront(e *listElem[T]) {
	e.prev.next = e.next
	e.next.prev = e.prev
	l.pushElem(e)
}

func (l *list[T]) remove(e *listElem[T]) {
	e.prev.next = e.next
	e.next.prev = e.prev
	e.next, e.prev = nil, nil
}

func (l *list[T]) removeLast() *listElem[T] {
	last := l.last()
	if last != nil {
		l.remove(last)
	}
	return last
}

func (l *list[T]) last() *listElem[T] {
	e := l.root.prev
	if e == &l.root {
		return nil
	}
	return e
}

func (l *list[T]) appendTo(slice []T) []T {
	for e := l.root.prev; e != &l.root; e = e.prev {
		slice = append(slice, e.v)
	}
	return slice
}
