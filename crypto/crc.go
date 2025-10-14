package crypto

import (
	"hash/crc32"
	"hash/crc64"
)

func CRC32(data []byte) uint32 {
	return crc32.ChecksumIEEE(data)
}

func VerifyCRC32(data []byte, exceptedCRC uint32) bool {
	return crc32.ChecksumIEEE(data) == exceptedCRC
}

func CRC64(data []byte) uint64 {
	table := crc64.MakeTable(crc64.ECMA)
	return crc64.Checksum(data, table)
}

func VerifyCRC64(data []byte, exceptedCRC uint64) bool {
	table := crc64.MakeTable(crc64.ECMA)
	return crc64.Checksum(data, table) == exceptedCRC
}
