/*
 * Copyright (c) 2019 uplus.io
 */

package utils

import (
	"bytes"
	"strings"
	"testing"
	"time"
	"uplus.io/udb/proto"
)

func TestStructSerializable(t *testing.T) {
	message := proto.Message{}
	message.Timestamp = *proto.NewTimestamp(time.Now().UnixNano(), time.Now().UnixNano())
	message.Category = proto.MessageCategorySystem
	message.Type = proto.MessageTypeClock
	message.From = uint32(1234567890)
	message.Content = StringToBytes("abcdefgh")

	dat, err := message.Serialize()
	if err != nil {
		t.Fatal(err)
	}
	copy := proto.Message{}
	copy.Deserialize(dat)

	if message.Category != copy.Category && message.Type != copy.Type &&
		bytes.Equal(message.Content, copy.Content) &&
		message.Timestamp.Local != copy.Timestamp.Local && message.Timestamp.Remote != copy.Timestamp.Remote {
		t.Fail()
	}
}

func TestSingleByteSerializable(t *testing.T) {
	var i1 int8 = -1
	var i2 uint8 = 2
	var b bool = true

	if BytesToInt8(Int8ToBytes(i1)) != i1 {
		t.Fail()
	}

	if BytesToUInt8(UInt8ToBytes(i2)) != i2 {
		t.Fail()
	}

	if BytesToBool(BoolToBytes(b)) != b {
		t.Fail()
	}

}
func TestSignedNumberSerializable(t *testing.T) {
	var i1 int = 1
	var i2 int16 = 2
	var i3 int32 = 3
	var i4 int64 = 4

	if LBytesToInt(LIntToBytes(i1)) != i1 {
		t.Fail()
	}

	if BBytesToInt(BIntToBytes(i1)) != i1 {
		t.Fail()
	}

	if LBytesToInt16(LInt16ToBytes(i2)) != i2 {
		t.Fail()
	}

	if BBytesToInt16(BInt16ToBytes(i2)) != i2 {
		t.Fail()
	}

	if LBytesToInt32(LInt32ToBytes(i3)) != i3 {
		t.Fail()
	}

	if BBytesToInt32(BInt32ToBytes(i3)) != i3 {
		t.Fail()
	}
	if LBytesToInt64(LInt64ToBytes(i4)) != i4 {
		t.Fail()
	}
	if BBytesToInt64(BInt64ToBytes(i4)) != i4 {
		t.Fail()
	}
}

func TestUnsignedNumberSerializable(t *testing.T) {
	var i1 uint = 1
	var i2 uint16 = 2
	var i3 uint32 = 3
	var i4 uint64 = 4

	if LBytesToUInt(LUIntToBytes(i1)) != i1 {
		t.Fail()
	}

	if BBytesToUInt(BUIntToBytes(i1)) != i1 {
		t.Fail()
	}

	if LBytesToUInt16(LUInt16ToBytes(i2)) != i2 {
		t.Fail()
	}

	if BBytesToUInt16(BUInt16ToBytes(i2)) != i2 {
		t.Fail()
	}

	if LBytesToUInt32(LUInt32ToBytes(i3)) != i3 {
		t.Fail()
	}

	if BBytesToUInt32(BUInt32ToBytes(i3)) != i3 {
		t.Fail()
	}
	if LBytesToUInt64(LUInt64ToBytes(i4)) != i4 {
		t.Fail()
	}
	if BBytesToUInt64(BUInt64ToBytes(i4)) != i4 {
		t.Fail()
	}
}

func TestFloatSerializable(t *testing.T) {
	var f1 float32 = 3.1415926
	var f2 float64 = 3.1415927
	if LBytesToFloat32(LFloat32ToBytes(f1)) != f1 {
		t.Fail()
	}
	if LBytesToFloat64(LFloat64ToBytes(f2)) != f2 {
		t.Fail()
	}

	if BBytesToFloat32(BFloat32ToBytes(f1)) != f1 {
		t.Fail()
	}
	if BBytesToFloat64(BFloat64ToBytes(f2)) != f2 {
		t.Fail()
	}
}

func TestStringSerializable(t *testing.T) {
	var s = "hello world"
	if !strings.EqualFold(s, BytesToString(StringToBytes(s))) {
		t.Fail()
	}
}
