/*
 * Copyright (c) 2019 uplus.io
 */

package hash

import (
	"hash/crc32"
	"hash/crc64"
)

var (
	//CRC64_Table_DEFAULT = crc64.MakeTable(crc64.ISO)
	CRC64_Table_DEFAULT = crc64.MakeTable(crc64.ECMA)
)

// MurMurHash算法 :https://github.com/spaolacci/murmur3

//byte array hash int32
func Int32(bytes []byte) int32 {
	return int32(Uint32(bytes))
}

//string hash int32
func Int32Of(str string) int32 {
	return int32(UInt32Of(str))
}

//byte array hash uint32
func Uint32(bytes []byte) uint32 {
	return crc32.ChecksumIEEE(bytes)
}

//string hash uint32
func UInt32Of(str string) uint32 {
	return Uint32([]byte(str))
}

func Int64(bytes []byte) int64 {
	return int64(UInt64(bytes))
}

func Int64Of(str string) int64 {
	return int64(UInt64Of(str))
}

func UInt64(bytes []byte) uint64 {
	return crc64.Checksum(bytes, CRC64_Table_DEFAULT)
}

func UInt64Of(str string) uint64 {
	return UInt64([]byte(str))
}

//dat哈希后取size的余数
func Remainder(dat []byte, size int) int {
	hash := Uint32(dat)
	remainder := hash & 0x7fffffff % uint32(size)
	return int(remainder)
}
