/*
CRC checksum functions used in the transmit integrity layer.
Ensure packet completeness and detect data corruption.

CRC64 → Used to validate entire packet integrity.
CRC32 → Used to validate integrity of packet's internal RPC section.
*/
package crypto

import (
	"hash/crc32"
	"hash/crc64"
)

var (
	crc32Table = crc32.MakeTable(crc32.IEEE)
	crc64Table = crc64.MakeTable(crc64.ECMA)
)

// CRC32 computes the CRC32 checksum (IEEE polynomial)
func CRC32(data []byte) uint32 {
	return crc32.Checksum(data, crc32Table)
}

// VerifyCRC32 verifies whether the provided checksum matches the data
func VerifyCRC32(data []byte, expectedCRC uint32) bool {
	return crc32.Checksum(data, crc32Table) == expectedCRC
}

// CRC64 computes the CRC64 checksum (ECMA polynomial)
func CRC64(data []byte) uint64 {
	return crc64.Checksum(data, crc64Table)
}

// VerifyCRC64 verifies whether the provided checksum matches the data
func VerifyCRC64(data []byte, exceptedCRC uint64) bool {
	return crc64.Checksum(data, crc64Table) == exceptedCRC
}
