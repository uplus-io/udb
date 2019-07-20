/*
 * Copyright (c) 2019 uplus.io
 */

package store

//import (
//	"github.com/syndtr/goleveldb/leveldb"
//)
//
//type StoreLeveldb struct {
//	db *leveldb.DB
//}
//
//func OpenStoreLeveldb(cfg StoreConfig) (Store, error) {
//	storeBolt := &StoreLeveldb{}
//	db, err := leveldb.OpenFile(cfg.Path, nil)
//	if err != nil {
//		return nil, err
//	}
//	if err != nil {
//		return nil, err
//	}
//	storeBolt.db = db
//	return storeBolt, nil
//}
//
//func (p *StoreLeveldb) Set(k, v []byte) error {
//	return p.db.Put(k, v, nil)
//}
//
//func (p *StoreLeveldb) Get(k []byte) (dat []byte, err error) {
//	return p.db.Get(k, nil)
//}
//
//func (p *StoreLeveldb) Close() error {
//	return p.db.Close()
//}
