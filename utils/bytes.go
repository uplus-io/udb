/*
 * Copyright (c) 2019 uplus.io
 */

package utils

import (
	"encoding/binary"
	"math"
)

func LIntToBytes(val int) ([]byte) {
	return LUInt32ToBytes(uint32(val))
}

func LBytesToInt(bytes []byte) int {
	val := LBytesToUInt32(bytes)
	return int(val)
}

func BIntToBytes(val int) ([]byte) {
	return BUInt32ToBytes(uint32(val))
}

func BBytesToInt(bytes []byte) int {
	val := BBytesToUInt32(bytes)
	return int(val)
}

func LUIntToBytes(val uint) ([]byte) {
	return LUInt32ToBytes(uint32(val))
}

func LBytesToUInt(bytes []byte) uint {
	val := LBytesToUInt32(bytes)
	return uint(val)
}

func BUIntToBytes(val uint) ([]byte) {
	return BUInt32ToBytes(uint32(val))
}

func BBytesToUInt(bytes []byte) uint {
	val := BBytesToUInt32(bytes)
	return uint(val)
}

func StringToBytes(val string) []byte {
	return []byte(val)
}

func BytesToString(dat []byte) string {
	return string(dat)
}

//单字节数据类型直接转换

func BoolToBytes(val bool) (dat []byte) {
	dat = make([]byte, 1)
	if val {
		dat[0] = 0x1
	} else {
		dat[0] = 0x0
	}
	return
}

func BytesToBool(dat []byte) bool {
	if dat[0] == 0x1 {
		return true
	} else {
		return false
	}
}

func Int8ToBytes(val int8) (dat []byte) {
	dat = make([]byte, 1)
	dat[0] = byte(val)
	return
}

func BytesToInt8(dat []byte) (val int8) {
	return int8(dat[0])
}

func UInt8ToBytes(val uint8) (dat []byte) {
	dat = make([]byte, 1)
	dat[0] = byte(val)
	return
}

func BytesToUInt8(dat []byte) (val uint8) {
	return uint8(dat[0])
}

//小端序

func LInt16ToBytes(val int16) (dat []byte) {
	dat = make([]byte, 2)
	binary.LittleEndian.PutUint16(dat, uint16(val))
	return
}

func LBytesToInt16(bytes []byte) (int16) {
	val := binary.LittleEndian.Uint16(bytes)
	return int16(val)
}

func LUInt16ToBytes(val uint16) (dat []byte) {
	dat = make([]byte, 2)
	binary.LittleEndian.PutUint16(dat, val)
	return
}

func LBytesToUInt16(bytes []byte) (uint16) {
	val := binary.LittleEndian.Uint16(bytes)
	return val
}

func LInt32ToBytes(val int32) (dat []byte) {
	dat = make([]byte, 4)
	binary.LittleEndian.PutUint32(dat, uint32(val))
	return
}

func LBytesToInt32(bytes []byte) (int32) {
	val := binary.LittleEndian.Uint32(bytes)
	return int32(val)
}

func LUInt32ToBytes(val uint32) (dat []byte) {
	dat = make([]byte, 4)
	binary.LittleEndian.PutUint32(dat, val)
	return
}

func LBytesToUInt32(bytes []byte) (uint32) {
	val := binary.LittleEndian.Uint32(bytes)
	return val
}

func LInt64ToBytes(val int64) (dat []byte) {
	dat = make([]byte, 8)
	binary.LittleEndian.PutUint64(dat, uint64(val))
	return
}

func LBytesToInt64(bytes []byte) (int64) {
	val := binary.LittleEndian.Uint64(bytes)
	return int64(val)
}

func LUInt64ToBytes(val uint64) (dat []byte) {
	dat = make([]byte, 8)
	binary.LittleEndian.PutUint64(dat, val)
	return
}

func LBytesToUInt64(bytes []byte) (uint64) {
	val := binary.LittleEndian.Uint64(bytes)
	return val
}

//大端序

func BInt16ToBytes(val int16) (dat []byte) {
	dat = make([]byte, 2)
	binary.BigEndian.PutUint16(dat, uint16(val))
	return
}

func BBytesToInt16(bytes []byte) (int16) {
	val := binary.BigEndian.Uint16(bytes)
	return int16(val)
}

func BUInt16ToBytes(val uint16) (dat []byte) {
	dat = make([]byte, 2)
	binary.BigEndian.PutUint16(dat, val)
	return
}

func BBytesToUInt16(bytes []byte) (uint16) {
	val := binary.BigEndian.Uint16(bytes)
	return val
}

func BInt32ToBytes(val int32) (dat []byte) {
	dat = make([]byte, 4)
	binary.BigEndian.PutUint32(dat, uint32(val))
	return
}

func BBytesToInt32(bytes []byte) (int32) {
	val := binary.BigEndian.Uint32(bytes)
	return int32(val)
}

func BUInt32ToBytes(val uint32) (dat []byte) {
	dat = make([]byte, 4)
	binary.BigEndian.PutUint32(dat, val)
	return
}

func BBytesToUInt32(bytes []byte) (uint32) {
	val := binary.BigEndian.Uint32(bytes)
	return val
}

func BInt64ToBytes(val int64) (dat []byte) {
	dat = make([]byte, 8)
	binary.BigEndian.PutUint64(dat, uint64(val))
	return
}

func BBytesToInt64(bytes []byte) (int64) {
	val := binary.BigEndian.Uint64(bytes)
	return int64(val)
}

func BUInt64ToBytes(val uint64) (dat []byte) {
	dat = make([]byte, 8)
	binary.BigEndian.PutUint64(dat, val)
	return
}

func BBytesToUInt64(bytes []byte) (uint64) {
	val := binary.BigEndian.Uint64(bytes)
	return val
}

//大小端序 浮点数

func LFloat32ToBytes(val float32) []byte {
	bits := math.Float32bits(val)
	return LUInt32ToBytes(bits)
}

func LBytesToFloat32(dat []byte) float32 {
	bits := LBytesToUInt32(dat)
	return math.Float32frombits(bits)
}

func LFloat64ToBytes(val float64) []byte {
	bits := math.Float64bits(val)
	return LUInt64ToBytes(bits)
}

func LBytesToFloat64(dat []byte) float64 {
	bits := LBytesToUInt64(dat)
	return math.Float64frombits(bits)
}

func BFloat32ToBytes(val float32) []byte {
	bits := math.Float32bits(val)
	return BUInt32ToBytes(bits)
}

func BBytesToFloat32(dat []byte) float32 {
	bits := BBytesToUInt32(dat)
	return math.Float32frombits(bits)
}

func BFloat64ToBytes(val float64) []byte {
	bits := math.Float64bits(val)
	return BUInt64ToBytes(bits)
}

func BBytesToFloat64(dat []byte) float64 {
	bits := BBytesToUInt64(dat)
	return math.Float64frombits(bits)
}
