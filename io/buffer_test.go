/*
 * Copyright (c) 2019 uplus.io
 */

package io

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"testing"
	"time"
	"uplus.io/udb/core"
	"uplus.io/udb/proto"
)

func TestNativeBuffer(t *testing.T) {
	//buffer := NewBigBuffer(16)
	//
	//buffer.writeInt(1)
	//buffer.writeInt(2)
	//
	//fmt.Println(buffer.ReadInt())
	//fmt.Println(buffer.ReadInt())

	//tmp := make([]byte, 16, 16)
	//buffer := bytes.NewBigBuffer(tmp)
	buffer := bytes.NewBuffer([]byte{})

	//binary.BigEndian.PutUint32(int32(1))
	binary.Write(buffer, binary.BigEndian, int32(1))
	binary.Write(buffer, binary.BigEndian, int32(2))
	binary.Write(buffer, binary.BigEndian, int32(3))
	binary.Write(buffer, binary.BigEndian, int32(4))

	var first, second, third int
	binary.Read(buffer, binary.BigEndian, &first)
	binary.Read(buffer, binary.BigEndian, &second)
	binary.Read(buffer, binary.BigEndian, &third)

	fmt.Println(first)
	fmt.Println(second)
	fmt.Println(third)

	temp := make([]byte, 4, 4)
	var five int32 = -5
	binary.BigEndian.PutUint32(temp, uint32(five))

	u := binary.BigEndian.Uint32(temp)
	fmt.Println(int32(u))
}

func TestBuffer(t *testing.T) {
	s := "abcdefghijklmnopqrstuvwxyz"

	buffer := NewBigBuffer()
	buffer.WriteInt(1)
	buffer.WriteInt(2)
	buffer.WriteInt(3)
	buffer.WriteInt(4)
	buffer.WriteInt32(9)
	buffer.WriteLong(1000000000000000)
	buffer.WriteBytes([]byte(s))

	fmt.Println(buffer.ReadInt())
	fmt.Println(buffer.ReadInt())
	fmt.Println(buffer.ReadInt())
	fmt.Println(buffer.ReadInt())
	fmt.Println(buffer.ReadInt32())
	fmt.Println(buffer.ReadLong())
	fmt.Println(buffer.ReadBytes(26))

	node := core.Node{Id: 123, Timestamp: *proto.NewTimestamp(time.Now().UnixNano(), time.Now().UnixNano())}
	buf := NewBigBuffer()
	buf.writeBuffer(node)
	dat := buf.Bytes()
	fmt.Println(dat)

}
