/*
 * Copyright (c) 2019 uplus.io
 */

package store

import (
	"github.com/dgraph-io/badger"
	"uplus.io/udb"
)

type StoreBadger struct {
	db *badger.DB
}

func OpenStoreBadger(cfg StoreConfig) (Store, error) {
	storeBolt := &StoreBadger{}
	options := badger.LSMOnlyOptions(cfg.Path)
	options.WithTruncate(true)
	db, err := badger.Open(options)
	//db, err := badger.Open(badger.DefaultOptions(cfg.Path))
	if err != nil {
		return nil, err
	}
	storeBolt.db = db
	return storeBolt, nil
}

func (p *StoreBadger) Set(data Data) error {
	return p.db.Update(func(txn *badger.Txn) error {
		return txn.Set(data.Id.IdBytes(), data.Content)
	})
}

func (p *StoreBadger) Get(id Identity) (dat *Data, err error) {
	err = p.db.View(func(txn *badger.Txn) error {
		item, e := txn.Get(id.IdBytes())
		if e == badger.ErrKeyNotFound {
			return udb.ErrDbKeyNotFound
		}
		if e != nil {
			return e
		}
		content, e := item.ValueCopy(nil)
		dat = NewData(id, content)
		return e
	})
	return
}

func (p *StoreBadger) Seek(id Identity, iter StoreIterator) (err error) {
	err = p.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		for it.Seek(id.IdBytes()); it.ValidForPrefix(id.IdBytes()); it.Next() {
			item := it.Item()
			content, err := item.ValueCopy(nil)
			if err != nil {
				return err
			}
			notBreak := iter(*NewData(id, content))
			if !notBreak {
				break
			}
		}
		return nil
	})
	return
}

func (p *StoreBadger) ForEach(iter StoreIterator) (err error)  {
	err = p.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			id := item.Key()
			identity := NewIdentityWithValue(id)
			content, err := item.ValueCopy(nil)
			if err != nil {
				return err
			}
			notBreak := iter(*NewData(*identity, content))
			if !notBreak {
				break
			}
		}
		return nil
	})
	return
}

func (p *StoreBadger) Close() error {
	return p.db.Close()
}
