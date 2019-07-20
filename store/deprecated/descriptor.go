/*
 * Copyright (c) 2019 uplus.io
 */

package deprecated

import (
	"uplus.io/udb/io"
)

const CONTENT_LENGTH_DESCRIPTOR uint32 = 266

type Descriptor struct {
	//编号 4
	Id uint32
	//删除标记 1
	Deleted bool
	//字段名 256
	Name string
	//字段数据类型 1
	BaseType BaseType
	//字段数据体长度 4
	Length uint32
}

func NewDescriptorOf(dat []byte) *Descriptor {
	buffer := io.NewBigBuffer(nil)
	buffer.WriteBytes(dat)
	descriptor := &Descriptor{}
	descriptor.Id = buffer.ReadUint32()
	descriptor.Deleted = buffer.ReadBool()
	descriptor.Name = buffer.ReadString(256)
	descriptor.BaseType = BaseType(buffer.ReadUint8())
	descriptor.Length = buffer.ReadUint32()
	return descriptor
}

func (d Descriptor) Bytes() []byte {
	buffer := io.NewBigBuffer(nil)
	buffer.WriteUint32(d.Id)
	buffer.WriteBool(d.Deleted)
	buffer.WriteString(d.Name, 256)
	buffer.WriteUint8(uint8(d.BaseType))
	buffer.WriteUint32(d.Length)
	return buffer.Bytes()
}
