/*
 * Copyright (c) 2019 uplus.io
 */

package store

import (
	"uplus.io/udb/config"
	"uplus.io/udb/proto"
)

type Engine struct {
	config     config.StorageConfig
	meta       Store
	partitions []Store
	partSize   int

	table *EngineTable
}

func NewEngine(config config.StorageConfig) *Engine {
	storeType := StoreTypeOfValue(config.Engine)
	partSize := len(config.Partitions)
	stores := make([]Store, partSize)
	meta := NewStore(StoreConfig{Path: config.Meta, Type: storeType})
	for i, path := range config.Partitions {
		store := NewStore(StoreConfig{Path: path, Type: storeType})
		stores[i] = store
	}
	engine := &Engine{config: config, meta: meta, partitions: stores, partSize: partSize}
	engine.table = NewEngineTable(engine)
	return engine
}

func (p *Engine) Close() {
	p.meta.Close()
	for _, store := range p.partitions {
		store.Close()
	}
}

func (p *Engine) Table() *EngineTable {
	return p.table
}

func (p *Engine) MetaSeek(id Identity, iter StoreIterator) error {
	return p.meta.Seek(id, iter)
}

func (p *Engine) MetaForEach(iter StoreIterator) error {
	return p.meta.ForEach(iter)
}

func (p *Engine) SetMeta(data Data) error {
	return p.meta.Set(data)
}

func (p *Engine) GetMeta(id Identity) (*Data, error) {
	return p.meta.Get(id)
}

func (p *Engine) DataSeek(partIndex int32, id Identity, iter StoreIterator) error {
	return p.part(partIndex).Seek(id, iter)
}

func (p *Engine) part(partIndex int32) Store {
	return p.partitions[partIndex]
}

func (p *Engine) SetData(partId int32, data Data) error {
	//return p.dataPart(&data.Id).Set(data)
	return nil
}

func (p *Engine) GetData(partId int32, id Identity) (*Data, error) {
	//return p.dataPart(&id).Get(id)
	return nil, nil
}

func (p *Engine) AddPart(part proto.Partition) error {
	return p.Table().AddPartition(part)
}
