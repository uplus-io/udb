/*
 * Copyright (c) 2019 uplus.io
 */

package deprecated

import (
	"errors"
	"fmt"
	"strings"
	"uplus.io/udb/io"
)

const CONTENT_LENGTH_SCHEMA_HEAD uint32 = 265

type Schema struct {
	//编号 4
	Id uint32
	//删除标记 1
	Deleted bool
	//结构名 256
	Name string
	//最后一个字段的id值 4
	tail uint32
	//数据结构描述
	Descriptors map[uint32]Descriptor
}

func NewSchema(name string) *Schema {
	schema := &Schema{Name: name}
	schema.Descriptors = make(map[uint32]Descriptor)
	return schema
}

func NewSchemaOf(dat []byte) (*Schema) {
	schema := &Schema{}
	schema.Descriptors = make(map[uint32]Descriptor)
	buffer := io.NewBigBuffer(nil)
	buffer.WriteBytes(dat)
	schema.Id = buffer.ReadUint32()
	schema.Deleted = buffer.ReadBool()
	schema.Name = buffer.ReadString(256)
	schema.tail = buffer.ReadUint32()
	for i := 0; i < int(schema.tail); i++ {
		bytes := buffer.ReadBytes(int(CONTENT_LENGTH_DESCRIPTOR))
		schema.Descriptors[uint32(i)] = *NewDescriptorOf(bytes)
	}
	return schema
}

func (s Schema) ContentLength() uint32 {
	return CONTENT_LENGTH_SCHEMA_HEAD + s.tail*CONTENT_LENGTH_DESCRIPTOR
}

func (d Schema) Bytes() []byte {
	buffer := io.NewBigBuffer(nil)
	buffer.WriteUint32(d.Id)
	buffer.WriteBool(d.Deleted)
	buffer.WriteString(d.Name, 256)
	buffer.WriteUint32(d.tail)
	for _, desc := range d.Descriptors {
		buffer.WriteBytes(desc.Bytes())
	}
	return buffer.Bytes()
}

func (p *Schema) Create(name string, baseType BaseType, length uint32) (Descriptor) {
	descriptor := Descriptor{Id: p.tail, Name: name, BaseType: baseType, Length: length}
	p.Descriptors[p.tail] = descriptor
	p.tail++
	return descriptor
}

func (p *Schema) Remove(name string) error {
	for id, desc := range p.Descriptors {
		if strings.EqualFold(desc.Name, name) {
			desc.Deleted = true
			p.Descriptors[id] = desc
			return nil
		}
	}
	return errors.New(fmt.Sprintf("schema[%s]not found descriptor[%s]", p.Name, name))
}
