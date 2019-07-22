/*
 * Copyright (c) 2019 uplus.io 
 */

package store

import (
	"fmt"
	"uplus.io/udb"
	"uplus.io/udb/hash"
	log "uplus.io/udb/logger"
	"uplus.io/udb/proto"
)

const (
	ENGINE_NAMESPACE_SYSTEM = "_sys"  // 系统命名空间
	ENGINE_NAMESPACE_USER   = "_user" // 用户命名空间

	ENGINE_TABLE_NAMESPACES = "_ns"    //系统命名空间表名
	ENGINE_TABLE_TABLES     = "_tab"   //系统表名
	ENGINE_TABLE_PARTITIONS = "_part"  //系统分区表名
	ENGINE_TABLE_METAS      = "_meta" //系统元数据表名

	ENGINE_KEY_META_STORAGE = "storage" //存储元数据主键
	ENGINE_KEY_META_PART    = "part"    //分区元数据主键
)

var (
	EMPTY_KEY = []byte{}
)

type EngineTable struct {
	engine *Engine
	nss    map[int32]*proto.Namespace
	tabs   map[int32]*proto.Table
	parts  map[int32]*proto.Partition
}

func NewEngineTable(engine *Engine) *EngineTable {
	table := &EngineTable{engine: engine, nss: make(map[int32]*proto.Namespace), tabs: make(map[int32]*proto.Table)}
	table.recoverMeta()
	table.recoverPartition()
	return table
}

func (p *EngineTable) recoverMeta() {
	p.engine.MetaSeek(*NewIdSysNs(EMPTY_KEY), func(data Data) bool {
		namespace := &proto.Namespace{}
		err := proto.Unmarshal(data.Content, namespace)
		if err != nil {
			return false
		}
		p.nss[namespace.Id] = namespace
		return true
	})

	p.engine.MetaSeek(*NewIdSysTab(EMPTY_KEY), func(data Data) bool {
		table := &proto.Table{}
		err := proto.Unmarshal(data.Content, table)
		if err != nil {
			return false
		}
		p.tabs[table.Id] = table
		return true
	})

	p.engine.MetaSeek(*NewIdOfPart(EMPTY_KEY), func(data Data) bool {
		part := &proto.Partition{}
		err := proto.Unmarshal(data.Content, part)
		if err != nil {
			return false
		}
		p.parts[part.Id] = part
		return true
	})
}

func (p *EngineTable) recoverPartition() {
	partSize := len(p.engine.config.Partitions)
	for i := 0; i < partSize; i++ {
		part := &proto.Partition{}
		index := int32(i)
		data, err := p.engine.part(index).Get(*NewIdOfPart([]byte(ENGINE_KEY_META_PART)))
		if err != udb.ErrDbKeyNotFound {
			err := proto.Unmarshal(data.Content, part)
			if err == nil {
				//validate part real index
				if part.Index != index {
					panic(fmt.Sprintf("part index inconsistent"))
				}
			}
		}
		p.parts[part.Id] = part
	}
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

func (p *EngineTable) NSs() []proto.Namespace {
	namespaces := []proto.Namespace{}
	for _, ns := range p.nss {
		namespaces = append(namespaces, *ns)
	}
	return namespaces
}

func (p *EngineTable) NS(namespace string) *proto.Namespace {
	return p.NSValue(hash.Int32Of(namespace))
}

func (p *EngineTable) NSValue(namespace int32) *proto.Namespace {
	return p.nss[namespace]
}

func (p *EngineTable) IfAbsentNS(namespace string) *proto.Namespace {
	pns := &proto.Namespace{}
	identity := NewIdSysNs([]byte(namespace))
	data, err := p.engine.GetMeta(*identity)
	if err == udb.ErrDbKeyNotFound {
		pns.Id = hash.Int32Of(namespace)
		pns.Name = namespace
		pns.Desc = proto.NewDescription(identity.NamespaceId, identity.TableId)
		nsData, _ := proto.Marshal(pns)
		data := NewData(*identity, nsData)
		err := p.engine.SetMeta(*data)
		if err != nil {
			log.Fatalf("db:create namespace fail - %v", err)
			return nil
		}
		p.nss[pns.Id] = pns
		return pns
	}
	if err != nil {
		return nil
	}
	err = proto.Unmarshal(data.Content, pns)
	if err != nil {
		return nil
	}
	return pns
}

func (p *EngineTable) Tabs() []proto.Table {
	tables := []proto.Table{}
	for _, tab := range p.tabs {
		tables = append(tables, *tab)
	}
	return tables
}

func (p *EngineTable) Tab(tab string) *proto.Table {
	return p.TabValue(hash.Int32Of(tab))
}

func (p *EngineTable) TabValue(tab int32) *proto.Table {
	return p.tabs[tab]
}

func (p *EngineTable) IfAbsentTab(tab string) *proto.Table {
	ptab := &proto.Table{}
	identity := NewIdSysTab([]byte(tab))
	data, err := p.engine.meta.Get(*identity)
	if err == udb.ErrDbKeyNotFound {
		ptab.Id = hash.Int32Of(tab)
		ptab.Name = tab
		ptab.Desc = proto.NewDescription(identity.NamespaceId, identity.TableId)
		nsData, _ := proto.Marshal(ptab)
		data := NewData(*identity, nsData)
		err := p.engine.SetMeta(*data)
		if err != nil {
			log.Fatalf("db:create namespace fail - %v", err)
			return nil
		}
		p.tabs[ptab.Id] = ptab
		return ptab
	}
	if err != nil {
		return nil
	}
	err = proto.Unmarshal(data.Content, ptab)
	if err != nil {
		return nil
	}
	return ptab
}

func (p *EngineTable) Partition(partId int32) *proto.Partition {
	partition := p.parts[partId]
	return partition
}

func (p *EngineTable) AddPartition(part proto.Partition) error {
	partition := p.Partition(part.Id)
	//version != 0,used partition,throw error
	if partition.Version != 0 {
		return udb.ErrPartAllocated
	}
	partition.Id = part.Id
	partition.Version = VERSION
	partition.Index = part.Index
	partIdentity := NewIdOfPart([]byte(ENGINE_KEY_META_PART))
	bytes, _ := proto.Marshal(partition)
	err := p.engine.part(partition.Index).Set(*NewData(*partIdentity, bytes))
	if err != nil {
		partition.Id = 0
		partition.Version = 0
		partition.Index = -1
		return err
	}
	return nil
}

func (p *EngineTable) UpdateDataMeta(partId int32, identity Identity) (id *Identity, err error) {
	part := p.Partition(partId)
	if part.Version == 0 {
		id = nil
		err = udb.ErrPartNotAllocate
		return
	}
	store := p.engine.part(part.Id)

}
