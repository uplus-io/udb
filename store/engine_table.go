/*
 * Copyright (c) 2019 uplus.io 
 */

package store

import (
	"uplus.io/udb"
	"uplus.io/udb/hash"
	log "uplus.io/udb/logger"
	"uplus.io/udb/proto"
)

const (
	ENGINE_SYSTEM_NAMESPACE        = "_system"     // 系统命名空间
	ENGINE_SYSTEM_TABLE_NAMESPACES = "_namespaces" //系统命名空间表名
	ENGINE_SYSTEM_TABLE_TABLES     = "_tables"     //系统表名
	ENGINE_SYSTEM_TABLE_METAS      = "_metas"      //系统元数据表名
	ENGINE_KEY_META_STORAGE        = "storage"     //存储元数据主键
	ENGINE_KEY_META_PART           = "partitions"  //分区元数据主键
)

var (
	EMPTY_KEY = []byte{}
)

func NewIdSysNs(key []byte) *Identity {
	return NewIdentity(ENGINE_SYSTEM_NAMESPACE, ENGINE_SYSTEM_TABLE_NAMESPACES, key)
}

func NewIdSysTab(key []byte) *Identity {
	return NewIdentity(ENGINE_SYSTEM_NAMESPACE, ENGINE_SYSTEM_TABLE_TABLES, key)
}

func (p *Engine) SystemNamespaces() []proto.Namespace {
	namespaces := []proto.Namespace{}
	err := p.meta.Seek(*NewIdSysNs(EMPTY_KEY), func(data Data) bool {
		namespace := proto.Namespace{}
		err := proto.Unmarshal(data.Content, &namespace)
		if err != nil {
			return false
		}
		namespaces = append(namespaces, []proto.Namespace{namespace}...)
		return true
	})
	if err != nil {
		return nil
	}
	return namespaces
}

func (p *Engine) IfAbsentCreateNamespace(namespace string) *proto.Namespace {
	pns := &proto.Namespace{}
	ns := []byte(namespace)
	data, err := p.meta.Get(*NewIdSysNs(ns))
	if err == udb.ErrDbKeyNotFound {
		pns.Id = hash.Int64Of(namespace)
		pns.Name = namespace
		nsData, _ := proto.Marshal(pns)
		data := NewData(*NewIdSysNs(ns), nsData)
		err := p.meta.Set(*data)
		if err != nil {
			log.Fatalf("db:create namespace fail - %v", err)
			return nil
		}
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

func (p *Engine) GetPartitionMeta() *proto.PartitionMeta {
	return nil
}

func (p *Engine) SetPartitionMeta(*proto.PartitionMeta) error {
	return nil
}
