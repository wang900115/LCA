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

import "errors"

var (
	errMissingData      = errors.New("missing bytes on input")
	errUnreferencedData = errors.New("extra bytes on input")
	errExceededTarget   = errors.New("target data size exceeded")
	errZeroContent      = errors.New("zero byte in input content")
)

func DecompressBytes(data []byte, target int) ([]byte, error) {
	if len(data) > target {
		return nil, errExceededTarget
	}
	if len(data) == target {
		cpy := make([]byte, len(data))
		copy(cpy, data)
		return cpy, nil
	}
	return bitsetDecodeBytes(data, target)
}

func bitsetDecodeBytes(data []byte, target int) ([]byte, error) {
	out, size, err := bitsetDecodePartialBytes(data, target)
	if err != nil {
		return nil, err
	}
	if size != len(data) {
		return nil, errUnreferencedData
	}
	return out, nil
}

func bitsetDecodePartialBytes(data []byte, target int) ([]byte, int, error) {
	if target == 0 {
		return nil, 0, nil
	}

	decomp := make([]byte, target)
	if len(data) == 0 {
		return decomp, 0, nil
	}

	if target == 1 {
		decomp[0] = data[0]
		if data[0] != 0 {
			return decomp, 1, nil
		}
		return decomp, 0, nil
	}

	nonZeroBitset, ptr, err := bitsetDecodePartialBytes(data, (target+7)/8)
	if err != nil {
		return nil, ptr, err
	}

	for i := 0; i < 8*len(nonZeroBitset); i++ {
		if nonZeroBitset[i/8]&(1<<byte(7-i%8)) != 0 {
			if ptr > len(data) {
				return nil, 0, errMissingData
			}
			if i >= len(decomp) {
				return nil, 0, errExceededTarget
			}

			if data[ptr] == 0 {
				return nil, 0, errZeroContent
			}

			decomp[i] = data[ptr]
			ptr++
		}
	}
	return decomp, ptr, nil
}
