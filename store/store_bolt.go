/*
 * Copyright (c) 2019 uplus.io
 */

package store

//import (
//	"github.com/boltdb/bolt"
//	"time"
//)
//
//type StoreBolt struct {
//	db       *bolt.DB
//	database []byte
//
//
//}
//
//func OpenStoreBolt(cfg StoreConfig) (Store, error) {
//	storeBolt := &StoreBolt{database: []byte(DEFAULT_KV_DATABASE_NAME)}
//	db, err := bolt.Open(cfg.Path, 0600, &bolt.Options{Timeout: 3600 * time.Second})
//	if err != nil {
//		return nil, err
//	}
//	err = db.Update(func(tx *bolt.Tx) error {
//		_, err = tx.CreateBucketIfNotExists(storeBolt.database)
//		return err
//	})
//
//	if err != nil {
//		return nil, err
//	}
//	storeBolt.db = db
//	return storeBolt, nil
//}
//
//func (p *StoreBolt) Set(k, v []byte) error {
//	return p.db.Update(func(tx *bolt.Tx) error {
//		return tx.Bucket(p.database).Put(k, v)
//	})
//}
//
//func (p *StoreBolt) Get(k []byte) (dat []byte, err error) {
//	err = p.db.View(func(tx *bolt.Tx) error {
//		dat = tx.Bucket(p.database).Get(k)
//		return nil
//	})
//	return
//}
//
//func (p *StoreBolt) Close() error {
//	return p.db.Close()
//}
