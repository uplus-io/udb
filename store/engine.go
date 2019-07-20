/*
 * Copyright (c) 2019 uplus.io
 */

package store

import (
	"uplus.io/udb/config"
)

type Engine struct {
	config     config.StorageConfig
	meta       Store
	partitions []Store
	partSize   int
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
	return &Engine{config: config, meta: meta, partitions: stores, partSize: partSize}
}

func (p *Engine) Close() {
	p.meta.Close()
	for _, store := range p.partitions {
		store.Close()
	}
}

func (v Engine) dataPart(id Identity) Store {
	return v.partitions[id.Part(v.partSize)]
}

func (p *Engine) SetMeta(data Data) error {
	return p.meta.Set(data)
}

func (p *Engine) GetMeta(id Identity) (*Data, error) {
	return p.meta.Get(id)
}

func (p *Engine) SetData(data Data) error {
	return p.dataPart(data.Id).Set(data)
}

func (p *Engine) GetData(id Identity) (*Data, error) {
	return p.dataPart(id).Get(id)
}
