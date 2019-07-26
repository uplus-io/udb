/*
 * Copyright (c) 2019 uplus.io
 */

package store

import (
	"fmt"
	"log"
	"strings"
	"uplus.io/udb/hash"
	"uplus.io/udb/proto"
	"uplus.io/udb/utils"
)

type StoreType uint8
type StoreIterator func(key, data []byte) bool
type StoreMetaIterator func(key []byte, meta proto.DataMeta) bool

const (
	VERSION = 1
)

const (
	STORE_TYPE_UDB StoreType = iota
	STORE_TYPE_BOLT
	STORE_TYPE_LEVELDB
	STORE_TYPE_BADGER
)

var (
	StoreType_Value = map[string]uint8{
		"UDB":     0,
		"BOLG":    1,
		"LEVELDB": 2,
		"BADGER":  3,
	}

	FLAG_META = []byte{'0'}
	FLAG_DATA = []byte{'1'}
)

func StoreTypeOfValue(val string) StoreType {
	return StoreType(StoreType_Value[val])
}

type StoreConfig struct {
	Path string
	Type StoreType
}

type Store interface {
	Close() error

	Set(key, value []byte) error

	Get(key []byte) ([]byte, error)

	Seek(key []byte, iter StoreIterator) error

	ForEach(iter StoreIterator) error

	//WALName() string
}

func NewStore(cfg StoreConfig) (s Store) {
	var err error
	if STORE_TYPE_BOLT == cfg.Type {
		//s, err = OpenStoreBolt(cfg)
	} else if STORE_TYPE_LEVELDB == cfg.Type {
		//s, err = OpenStoreLeveldb(cfg)
	} else if STORE_TYPE_BADGER == cfg.Type {
		s, err = OpenStoreBadger(cfg)
	}
	if err != nil {
		log.Fatal(err)
	}
	return
}

type Identity struct {
	NamespaceId    int32  //命名空间id
	Namespace      string //命名空间名称
	TableId        int32  //表id
	Table          string //表名称
	PartitionIndex int32  //分区索引
	PartitionId    int32  //分区id
	Key            []byte //用户主键
	keyValue       int64  //用户主键系统id值
	id             []byte //系统id
	ring           int32  //系统id ring值
}

func NsTabBytes(ns, tab int32) []byte {
	bytes := []byte(fmt.Sprintf("/%d/%d/", ns, tab))
	//val := hash.UInt64(bytes)
	//return utils.LUInt64ToBytes(val)
	return bytes
}

func KeyBytes(key []byte) []byte {
	//val := hash.UInt64(key)
	//return utils.LUInt64ToBytes(val)
	return key
}

func NsTabKeyBytes(ns, tab int32, key []byte) []byte {
	return append(NsTabBytes(ns, tab), KeyBytes(key)...)
}

func IdentityVersionId(identity Identity, version int32) []byte {
	bytes := []byte(fmt.Sprintf("/%d", version))
	id := append(identity.IdBytes(), bytes...)
	return append(FLAG_DATA, id...)
}

func IdentityMetaId(identity Identity) []byte {
	return append(FLAG_META, identity.IdBytes()...)
}

func IdentityDataId(identity Identity) []byte {
	return append(FLAG_DATA, identity.IdBytes()...)
}

func NewIdentity(ns, tab string, key []byte) *Identity {
	identity := &Identity{Namespace: ns, Table: tab, Key: key}
	identity.NamespaceId = hash.Int32Of(ns)
	identity.TableId = hash.Int32Of(tab)
	identity.keyValue = hash.Int64(identity.Key)
	identity.id = NsTabKeyBytes(identity.NamespaceId, identity.TableId, identity.Key)
	identity.ring = hash.Int32(identity.id)
	return identity
}

func NewIdentityWithValue(idBytes []byte) *Identity {
	bits := strings.Split(string(idBytes), "/")
	ns := bits[0]
	tab := bits[1]
	key := bits[2]
	identity := &Identity{}
	identity.NamespaceId = utils.StringToInt32(ns)
	identity.TableId = utils.StringToInt32(tab)
	identity.Key = []byte(key)
	identity.keyValue = hash.Int64(identity.Key)
	identity.id = NsTabKeyBytes(identity.NamespaceId, identity.TableId, identity.Key)
	return NewIdentity(ns, tab, []byte(key))
}

//func (v Identity) IdValue() int64 {
//	return v.idValue
//}
//
//func (v Identity) KeyValue() int64 {
//	return v.keyValue
//}

func (v Identity) IdBytes() []byte {
	return v.id
}

//func (v Identity) KeyBytes() []byte {
//	return utils.LInt64ToBytes(v.keyValue)
//}

func (v Identity) Part(partSize int) int32 {
	i := v.keyValue & 0x7FFFFFFFFFFFFFFF
	return int32(i % int64(partSize))
}

type Data struct {
	Id      Identity
	Content []byte
	Version int32
}

func NewData(id Identity, con []byte) *Data {
	return &Data{Id: id, Content: con}
}
