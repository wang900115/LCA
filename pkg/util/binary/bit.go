// Copyright 2017 The go-ethereum Authors
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

package binary

import (
	"runtime"
	"unsafe"
)

const wordSize = int(unsafe.Sizeof(uintptr(0)))
const supportsUnaligned = runtime.GOARCH == "386" || runtime.GOARCH == "amd64" || runtime.GOARCH == "ppc64" || runtime.GOARCH == "ppc64le" || runtime.GOARCH == "s390x"

type Opselector func(uintptr, uintptr) uintptr

var (
	AND Opselector = func(u1, u2 uintptr) uintptr { return u1 & u2 }
	OR  Opselector = func(u1, u2 uintptr) uintptr { return u1 | u2 }
	XOR Opselector = func(u1, u2 uintptr) uintptr { return u1 ^ u2 }
)

func OPBytes(dst, a, b []byte, op Opselector) int {
	if supportsUnaligned {
		return fastOPBytes(dst, a, b, op)
	}
	return safeOPBytes(dst, a, b)
}

func fastOPBytes(dst, a, b []byte, op Opselector) int {
	n := len(a)
	if len(b) < n {
		n = len(b)
	}
	w := n / int(wordSize)
	if w > 0 {
		dw := *(*[]uintptr)(unsafe.Pointer(&dst))
		aw := *(*[]uintptr)(unsafe.Pointer(&a))
		bw := *(*[]uintptr)(unsafe.Pointer(&b))
		for i := 0; i < w; i++ {
			dw[i] = op(aw[i], bw[i])
		}
	}
	for i := n - n%wordSize; i < n; i++ {
		dst[i] = a[i] & b[i]
	}
	return n
}

func safeOPBytes(dst, a, b []byte) int {
	n := len(a)
	if len(b) < n {
		n = len(b)
	}
	for i := 0; i < n; i++ {
		dst[i] = a[i] & b[i]
	}
	return n
}
