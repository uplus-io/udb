/*
 * Copyright (c) 2019 uplus.io
 */

package store

import (
	"uplus.io/udb"
	"uplus.io/udb/config"
	"uplus.io/udb/hash"
	log "uplus.io/udb/logger"
	"uplus.io/udb/proto"
)

type Engine struct {
	config     config.StorageConfig
	meta       StoreOperator
	partitions []StoreOperator
	partSize   int

	table *EngineTable
}

func NewEngine(config config.StorageConfig) *Engine {
	storeType := StoreTypeOfValue(config.Engine)
	partSize := len(config.Partitions)
	stores := make([]StoreOperator, partSize)
	meta := NewStore(StoreConfig{Path: config.Meta, Type: storeType})
	for i, path := range config.Partitions {
		store := NewStore(StoreConfig{Path: path, Type: storeType})
		stores[i] = NewStoreOperatorKV(store)
	}
	engine := &Engine{config: config, meta: NewStoreOperatorKV(meta), partitions: stores, partSize: partSize}
	engine.table = NewEngineTable(engine)
	return engine
}

func (p *Engine) ValidatePartition(nodeId int32) error {
	for i, store := range p.partitions {
		ring := int32(hash.GenerateConsistentRing(nodeId, i))
		partition := store.Part()
		if partition == nil {
			part := proto.Partition{}
			part.Version = VERSION
			part.Id = ring
			part.Index = int32(i)
			_, err := store.PartIfAbsent(part)
			if err != nil {
				return err
			}
		} else {
			if partition.Id != ring {
				return udb.ErrPartRingVerifyFailed
			}
		}
		log.Debugf("validate partition[%d] ring[%d]", i, ring)
	}
	return nil
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

func (p *Engine) SetData(partId int32, data Data) error {
	partition := p.Table().Partition(partId)
	if partition == nil {
		return udb.ErrPartNotFound
	}
	p.meta.NSIfAbsent(data.Id.Namespace)
	p.meta.TABIfAbsent(data.Id.Namespace, data.Id.Table)
	return p.part(partition.Index).SetData(data)
}

func (p *Engine) GetData(partId int32, id Identity) (*proto.DataContent, error) {
	partition := p.Table().Partition(partId)
	if partition == nil {
		return nil, udb.ErrPartNotFound
	}
	return p.part(partition.Index).GetData(id)
}

func (p *Engine) part(partIndex int32) StoreOperator {
	return p.partitions[partIndex]
}

func (p *Engine) PartSize() int {
	return p.partSize
}

func (p *Engine) GetPart(partId int32) *proto.Partition {
	return p.Table().Partition(partId)
}

func (p *Engine) AddPart(part proto.Partition) error {
	return p.Table().AddPartition(part)
}
