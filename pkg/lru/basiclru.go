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

type BasicLRU[K comparable, V any] struct {
	list  *list[K]
	items map[K]cacheItem[K, V]
	cap   int
}

type cacheItem[K any, V any] struct {
	elem  *listElem[K]
	value V
}

func NewBasicLRU[K comparable, V any](capacity int) BasicLRU[K, V] {
	if capacity <= 0 {
		capacity = 1
	}
	c := BasicLRU[K, V]{
		list:  newList[K](),
		items: make(map[K]cacheItem[K, V]),
		cap:   capacity,
	}
	return c
}

func (c *BasicLRU[K, V]) Add(key K, value V) (evicted bool) {
	_, _, evicted = c.Add3(key, value)
	return evicted
}

func (c *BasicLRU[K, V]) Add3(key K, value V) (ek K, ev V, evicted bool) {
	item, ok := c.items[key]
	if ok {
		item.value = value
		c.items[key] = item
		c.list.moveToFront(item.elem)
		return ek, ev, false
	}

	var elem *listElem[K]
	if c.Len() >= c.cap {
		elem = c.list.removeLast()
		evicted = true
		ek = elem.v
		ev = c.items[ek].value
		delete(c.items, ek)
	} else {
		elem = new(listElem[K])
	}
	elem.v = key
	c.items[key] = cacheItem[K, V]{elem, value}
	c.list.pushElem(elem)
	return ek, ev, evicted
}

func (c *BasicLRU[K, V]) Contains(key K) bool {
	_, ok := c.items[key]
	return ok
}

func (c *BasicLRU[K, V]) Get(key K) (value V, ok bool) {
	item, ok := c.items[key]
	if !ok {
		return value, false
	}
	c.list.moveToFront(item.elem)
	return item.value, true
}

func (c *BasicLRU[K, V]) GetOldest() (key K, value V, ok bool) {
	lastElem := c.list.last()
	if lastElem == nil {
		return key, value, false
	}
	key = lastElem.v
	item := c.items[key]
	return key, item.value, true
}

func (c *BasicLRU[K, V]) Len() int {
	return len(c.items)
}

func (c *BasicLRU[K, V]) Peek(key K) (value V, ok bool) {
	item, ok := c.items[key]
	return item.value, ok
}

func (c *BasicLRU[K, V]) Purge() {
	c.list.init()
	clear(c.items)
}

func (c *BasicLRU[K, V]) Remove(key K) bool {
	item, ok := c.items[key]
	if ok {
		delete(c.items, key)
		c.list.remove(item.elem)
	}
	return ok
}

func (c *BasicLRU[K, V]) RemoveOldest() (key K, value V, ok bool) {
	lastElem := c.list.last()
	if lastElem == nil {
		return key, value, false
	}

	key = lastElem.v
	item := c.items[key]
	delete(c.items, key)
	c.list.remove(lastElem)
	return key, item.value, true
}

func (c *BasicLRU[K, V]) Keys() []K {
	keys := make([]K, 0, len(c.items))
	return c.list.appendTo(keys)
}
