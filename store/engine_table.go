/*
 * Copyright (c) 2019 uplus.io 
 */

package store

import (
	"uplus.io/udb/proto"
)

const (
	ENGINE_NAMESPACE_SYSTEM = "_sys"  // 系统命名空间
	ENGINE_NAMESPACE_USER   = "_user" // 用户命名空间

	ENGINE_TABLE_NAMESPACES = "_ns"   //系统命名空间表名
	ENGINE_TABLE_TABLES     = "_tab"  //系统表名
	ENGINE_TABLE_PARTITIONS = "_part" //系统分区表名
	ENGINE_TABLE_METAS      = "_meta" //系统元数据表名

	ENGINE_KEY_META_STORAGE = "storage" //存储元数据主键
	ENGINE_KEY_META_PART    = "part"    //分区元数据主键
)

var (
	EMPTY_KEY = []byte{}
)

type EngineTable struct {
	engine *Engine
	parts  map[int32]*proto.Partition
}

func NewEngineTable(engine *Engine) *EngineTable {
	table := &EngineTable{engine: engine, parts: make(map[int32]*proto.Partition)}
	table.recoverPartition()
	return table
}

func NewIdOfNs(key []byte) *Identity {
	return NewIdentity(ENGINE_NAMESPACE_SYSTEM, ENGINE_TABLE_NAMESPACES, key)
}

func NewIdOfPart(key []byte) *Identity {
	return NewIdentity(ENGINE_NAMESPACE_SYSTEM, ENGINE_TABLE_PARTITIONS, key)
}

func NewIdOfTab(namespace string, key []byte) *Identity {
	return NewIdentity(namespace, ENGINE_TABLE_TABLES, key)
}

func NewIdOfData(namespace, table string, key []byte) *Identity {
	return NewIdentity(namespace, table, key)
}

func (p *EngineTable) recoverPartition() {
	partSize := len(p.engine.config.Partitions)
	for i := 0; i < partSize; i++ {
		operator := p.engine.partitions[i]
		partition := operator.Part()
		if partition != nil {
			p.parts[partition.Id] = partition
		}
	}
}

func (p *EngineTable) Partition(partId int32) *proto.Partition {
	partition := p.parts[partId]
	return partition
}

func (p *EngineTable) AddPartition(part proto.Partition) error {
	newPart, err := p.engine.part(part.Index).PartIfAbsent(part)
	if err != nil {
		return err
	}
	p.parts[newPart.Id] = newPart
	return nil
}
