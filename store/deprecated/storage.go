/*
 * Copyright (c) 2019 uplus.io 
 */

package deprecated

import (
	"path/filepath"
	"uplus.io/udb/config"
	"uplus.io/udb/proto"
)

type Storage struct {
	config config.EngineConfig

	schema *SchemaStorage
}

func NewStorage(config config.EngineConfig) (*Storage, error) {
	var err error = nil
	storage := &Storage{config: config}
	storage.init()
	return storage, err
}

func(p *Storage) metaFile(filename string)string  {
	return p.config.StorePath + string(filepath.Separator) + "meta" + string(filepath.Separator) + filename
}

func (p *Storage) init() error {
	schema, err := NewSchemaStorage(p.metaFile("schemas"))
	if err != nil {
		return err
	}
	p.schema = schema
	return nil
}

func (p *Storage) Open() {

}

func (p *Storage) Close() {

}

func (p *Storage) UpdateSchema(schema *Schema) Schema {
	return p.schema.UpdateSchema(schema)
}

func (p *Storage) Put(schema string, descriptor string, key interface{}, value interface{}, version uint32) (*proto.Data, error) {
	return nil, nil
}

func (p *Storage) Get(schema string, descriptor string, key interface{}) (*proto.Data, error) {
	return nil, nil
}
