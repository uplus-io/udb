/*
 * Copyright (c) 2019 uplus.io
 */

package io

import (
	"bytes"
	"encoding/binary"
	"strings"
)

type Buffer struct {
	offset uint64
	buffer *bytes.Buffer
	order  binary.ByteOrder

	errors       []error
	errorsLength int
}

func NewBigBuffer(dat []byte) *Buffer {
	if dat == nil {
		dat = []byte{}
	}
	buf := bytes.NewBuffer(dat)
	buffer := &Buffer{offset: 0, buffer: buf, order: binary.BigEndian}
	buffer.errors = make([]error, 1)
	buffer.errorsLength = 0
	return buffer
}

func NewLittleBuffer(dat []byte) *Buffer {
	if dat == nil {
		dat = []byte{}
	}
	buffer := bytes.NewBuffer(dat)
	return &Buffer{offset: 0, buffer: buffer, order: binary.LittleEndian}
}

func (p *Buffer) HasError() bool {
	return p.errorsLength != 0
}

func (p *Buffer) Errors() []error {
	return p.errors
}

func (p *Buffer) LastError() error {
	return p.errors[p.errorsLength]
}

func (p *Buffer) checkError(err error) {
	if err != nil {
		p.errors[p.errorsLength] = err
		p.errorsLength++
	}
}

func (p *Buffer) Write(data interface{}) error {
	err := binary.Write(p.buffer, p.order, data)
	p.checkError(err)
	return err
}

func (p *Buffer) Read(data interface{}) error {
	err := binary.Read(p.buffer, p.order, data)
	p.checkError(err)
	return err
}

func (p *Buffer) Length() int {
	return len(p.buffer.Bytes())
}

func (p *Buffer) Bytes() []byte {
	return p.buffer.Bytes()
}

func (p *Buffer) WriteInt8(val int8) {
	p.Write(val)
}

func (p *Buffer) ReadInt8() int8 {
	var val int8
	p.Read(&val)
	return val
}

func (p *Buffer) WriteUint8(val uint8) {
	p.Write(val)
}

func (p *Buffer) ReadUint8() uint8 {
	var val uint8
	p.Read(&val)
	return val
}

func (p *Buffer) WriteInt(val int) {
	p.Write(uint32(val))
}

func (p *Buffer) ReadInt() int {
	var i uint32
	p.Read(&i)
	return int(i)
}

func (p *Buffer) WriteInt32(val int32) {
	p.Write(val)
}

func (p *Buffer) ReadInt32() int32 {
	var val int32
	p.Read(&val)
	return val
}

func (p *Buffer) WriteUint32(val uint32) {
	p.Write(val)
}

func (p *Buffer) ReadUint32() uint32 {
	var val uint32
	p.Read(&val)
	return val
}

func (p *Buffer) WriteInt64(val int64) {
	p.Write(val)
}

func (p *Buffer) ReadInt64() int64 {
	var val int64
	p.Read(&val)
	return val
}

func (p *Buffer) WriteBytes(bytes []byte) {
	p.Write(bytes)
}

func (p *Buffer) ReadBytes(len int) []byte {
	bytes := make([]byte, len)
	p.Read(&bytes)
	return bytes
}

func (p *Buffer) WriteString(str string, max int) {
	tmp := make([]byte, max)
	dat := []byte(str)
	len := len(dat)
	for i := 0; i < len; i++ {
		tmp[i] = dat[i]
	}
	for i := max - len; i > 0; i-- {
		tmp[len+i-1] = ' '
	}
	p.Write(tmp)
}

func (p *Buffer) ReadString(len int) string {
	bytes := make([]byte, len)
	binary.Read(p.buffer, binary.BigEndian, &bytes)
	return strings.TrimSpace(string(bytes))
}

func (p *Buffer) ReadBool() bool {
	bytes := make([]byte, 1)
	p.Read(&bytes)
	if bytes[0] == 0x1 {
		return true
	} else {
		return false
	}
}

func (p *Buffer) WriteBool(b bool) {
	tmp := make([]byte, 1)
	if b {
		tmp[0] = 0x1
	} else {
		tmp[0] = 0x0
	}
	p.Write(tmp)
}

func (p *Buffer) WriteEndLine() {
	p.WriteBytes([]byte{'\n'})
}
