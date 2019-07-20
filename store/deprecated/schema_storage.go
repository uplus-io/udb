/*
 * Copyright (c) 2019 uplus.io
 */

package deprecated

import (
	io2 "io"
	"log"
	"strings"
	"sync"
	"uplus.io/udb/io"
	"uplus.io/udb/version"
)

type SchemaStorage struct {
	Path    string
	tail    uint32
	schemas map[uint32]Schema

	schemaStream *io.Stream

	sync.RWMutex
}

func NewSchemaStorage(path string) (*SchemaStorage, error) {
	storage := &SchemaStorage{Path: path, schemas: make(map[uint32]Schema)}
	storage.init()
	err := storage.recoverSchema()
	if err != nil {
		return nil, err
	}
	return storage, nil
}

func (p *SchemaStorage) init() {
	p.schemaStream = io.NewStream(p.Path)
}

func (p *SchemaStorage) Schema(name string) (*Schema) {
	for _, s := range p.schemas {
		if strings.EqualFold(s.Name, name) {
			return &s
		}
	}
	return nil
}

func (p *SchemaStorage) UpdateSchema(schema *Schema) Schema {
	p.Lock()
	defer p.Unlock()
	exist := p.Schema(schema.Name)
	if exist != nil {
		p.schemas[exist.Id] = *schema
		p.storeSchema()
		return p.schemas[exist.Id]
	}

	p.schemas[p.tail] = *schema
	schema.Id = p.tail
	p.tail++
	p.storeSchema()
	return *schema
}

func (p *SchemaStorage) recoverSchema() error {
	ver, err := p.schemaStream.ReadLine(0)
	if err == io2.EOF {
		return nil
	}
	if err != nil {
		return err
	}
	if len(ver) == 0 {
		return nil
	}
	log.Printf("schema version:%s\n", ver)
	tail, err := p.schemaStream.ReadLine(7)
	if err != nil {
		return err
	}
	buffer := io.NewBigBuffer(nil)
	buffer.WriteBytes(tail)
	p.tail = buffer.ReadUint32()

	//偏移版本号 tail
	var offset int64 = 12
	for i := 0; i < int(p.tail); i++ {
		bytes, e := p.schemaStream.ReadLine(offset)
		if e != nil {
			return e
		}
		p.schemas[uint32(i)] = *NewSchemaOf(bytes)
	}
	return nil
}

func (p *SchemaStorage) storeSchema() {
	ver := io.NewBigBuffer(nil)
	ver.WriteString(version.STORAGE_VERSION, 6)
	ver.WriteEndLine()
	p.schemaStream.Write(0, *ver)

	tail := io.NewBigBuffer(nil)
	tail.WriteUint32(p.tail)
	tail.WriteEndLine()
	p.schemaStream.Write(7, *tail)
	var offset int64 = 12
	for _, s := range p.schemas {
		bytes := s.Bytes()
		length := len(bytes)
		buffer := io.NewBigBuffer(nil)
		buffer.WriteBytes(bytes)
		buffer.WriteEndLine()
		p.schemaStream.Write(offset, *buffer)
		offset += int64(length)
	}
	p.schemaStream.Flush()
}
