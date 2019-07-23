/*
 * Copyright (c) 2019 uplus.io
 */

package store

import (
	"github.com/dgraph-io/badger"
	"uplus.io/udb"
	log "uplus.io/udb/logger"
)

type StoreBadger struct {
	db *badger.DB
}

func OpenStoreBadger(cfg StoreConfig) (Store, error) {
	storeBolt := &StoreBadger{}
	options := badger.LSMOnlyOptions(cfg.Path)
	options.WithTruncate(true)
	options.Logger = nil
	db, err := badger.Open(options)
	//db, err := badger.Open(badger.DefaultOptions(cfg.Path))
	if err != nil {
		return nil, err
	}
	storeBolt.db = db
	return storeBolt, nil
}

func (p *StoreBadger) Close() error {
	return p.db.Close()
}

func (p *StoreBadger) Set(key, value []byte) error {
	return p.db.Update(func(txn *badger.Txn) error {
		return txn.Set(key, value)
	})
}

func (p *StoreBadger) Get(key []byte) (dat []byte, err error) {
	err = p.db.View(func(txn *badger.Txn) error {
		item, e := txn.Get(key)
		if e == badger.ErrKeyNotFound {
			return udb.ErrDbKeyNotFound
		}
		if e != nil {
			return e
		}
		dat, e = item.ValueCopy(nil)
		return e
	})
	return
}

func (p *StoreBadger) Seek(key []byte, iter StoreIterator) (err error) {
	err = p.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		for it.Seek(key); it.ValidForPrefix(key); it.Next() {
			item := it.Item()
			content, err := item.ValueCopy(nil)
			if err != nil {
				return err
			}
			log.Debugf("seek key:%s", item.Key())
			notBreak := iter(item.Key(), content)
			if !notBreak {
				break
			}
		}
		return nil
	})
	return
}

func (p *StoreBadger) ForEach(iter StoreIterator) (err error) {
	err = p.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			content, err := item.ValueCopy(nil)
			if err != nil {
				return err
			}
			log.Debugf("foreach key:%v", string(item.Key()))
			notBreak := iter(item.Key(), content)
			if !notBreak {
				break
			}
		}
		return nil
	})
	return
}
