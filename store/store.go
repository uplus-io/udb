/*
 * Copyright (c) 2019 uplus.io
 */

package store

import (
	"fmt"
	"log"
	"strings"
	"uplus.io/udb/hash"
)

type StoreType uint8
type StoreIterator func(data Data) bool

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
)

const (
	DEFAULT_KV_DATABASE_NAME = "udb-kv"
)

func StoreTypeOfValue(val string) StoreType {
	return StoreType(StoreType_Value[val])
}

type StoreConfig struct {
	Path string
	Type StoreType
}

type Store interface {
	Set(data Data) error

	Get(id Identity) (*Data, error)

	//Delete(k []byte) error
	//
	//Exits(k []byte) (bool, error)
	//
	Seek(id Identity, iter StoreIterator) error

	ForEach(iter StoreIterator) error
	//
	Close() error
	//
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
	Namespace string
	Table     string
	Key       []byte
	keyValue  int64
	id        []byte
}

func NsTabBytes(ns, tab string) []byte {
	bytes := []byte(fmt.Sprintf("%s/%s/", ns, tab))
	//val := hash.UInt64(bytes)
	//return utils.LUInt64ToBytes(val)
	return bytes
}

func KeyBytes(key []byte) []byte {
	//val := hash.UInt64(key)
	//return utils.LUInt64ToBytes(val)
	return key
}

func NsTabKeyBytes(ns, tab string, key []byte) []byte {
	return append(NsTabBytes(ns, tab), KeyBytes(key)...)
}

func NewIdentity(ns, tab string, key []byte) *Identity {
	identity := &Identity{Namespace: ns, Table: tab, Key: key}
	identity.keyValue = hash.Int64(key)
	identity.id = NsTabKeyBytes(ns, tab, key)
	//identity.idValue = hash.Int64(identity.id)
	return identity
}

func NewIdentityWithValue(idBytes []byte) *Identity {
	bits := strings.Split(string(idBytes), "/")
	ns := bits[0]
	tab := bits[1]
	key := bits[2]
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

func (v Identity) Part(partSize int) int {
	i := v.keyValue & 0x7FFFFFFFFFFFFFFF
	return int(i % int64(partSize))
}

type Data struct {
	Id      Identity
	Content []byte
}

func NewData(id Identity, con []byte) *Data {
	return &Data{Id: id, Content: con}
}
