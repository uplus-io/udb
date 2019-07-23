/*
 * Copyright (c) 2019 uplus.io
 */

package store

import (
	"fmt"
	"sync/atomic"
	"testing"
	"time"
	"uplus.io/udb/utils"
)

type StoreBenchmarkTestConfig struct {
	name       string
	concurrent int
	max        int

	keyLength  int
	dataLength int
}

func storeWrite(store Store, cfg StoreBenchmarkTestConfig) {
	benchmark := utils.NewBenchmark(cfg.name, cfg.concurrent, cfg.max)
	var start time.Time

	dat1K := make([]byte, cfg.dataLength, cfg.dataLength)

	durations := make(chan time.Duration, benchmark.Max)
	couter := make(chan int64, benchmark.Max)

	var exitCount int32 = 0

	benchmark.Start()

	pageSize := benchmark.Max / benchmark.Concurrent

	for c := 0; c < benchmark.Concurrent; c++ {
		go func(id int) {
			var times time.Duration
			var total int64

			var indexStart, indexEnd int
			indexStart = id * pageSize
			indexEnd = (id + 1) * pageSize
			if id+1 == benchmark.Concurrent {
				indexEnd = benchmark.Max
			}
			fmt.Printf("gountine %d start:%d end:%d\n", id, indexStart, indexEnd)
			for j := indexStart; j < indexEnd; j++ {
				key := benchmark.GenerateKey(j)
				start = time.Now()
				store.Set(NewIdentity("ns", "tab", key).IdBytes(), dat1K)
				//fmt.Printf("[%d]write %s \n",id,string(key))
				times += time.Since(start)
				total += 1
			}
			durations <- times
			couter <- total
			atomic.AddInt32(&exitCount, 1)
		}(c)

	}
	benchmark.Finish()

	for ; ; {
		d := <-durations
		c := <-couter
		benchmark.Put(d, c)
		if int(exitCount) == benchmark.Concurrent {
			break
		}
	}
	benchmark.Print()
}

func storeRead(store Store, cfg StoreBenchmarkTestConfig) {
	benchmark := utils.NewBenchmark(cfg.name, cfg.concurrent, cfg.max)
	var start time.Time

	durations := make(chan time.Duration, benchmark.Max)
	couter := make(chan int64, benchmark.Max)

	var exitCount int32 = 0

	benchmark.Start()

	pageSize := benchmark.Max / benchmark.Concurrent

	for c := 0; c < benchmark.Concurrent; c++ {
		go func(id int) {
			var times time.Duration
			var total int64

			var indexStart, indexEnd int
			indexStart = id * pageSize
			indexEnd = (id + 1) * pageSize
			if id+1 == benchmark.Concurrent {
				indexEnd = benchmark.Max
			}
			fmt.Printf("gountine %d start:%d end:%d\n", id, indexStart, indexEnd)
			for j := indexStart; j < indexEnd; j++ {
				key := benchmark.RandomInt16()
				start = time.Now()
				_, err := store.Get(NewIdentity("ns", "tab", key).IdBytes())
				if err != nil {
					fmt.Printf("%s[%d]read error %s\n", cfg.name, id, string(key))
				}
				//fmt.Printf("%s[%d]read %s/%s \n",cfg.name,id,string(key),string(bytes))
				times += time.Since(start)
				total += 1
			}
			durations <- times
			couter <- total
			atomic.AddInt32(&exitCount, 1)
		}(c)

	}
	benchmark.Finish()

	for ; ; {
		d := <-durations
		c := <-couter
		benchmark.Put(d, c)
		if int(exitCount) == benchmark.Concurrent {
			break
		}
	}
	benchmark.Print()
}

func TestStore(t *testing.T) {
	concurrent := 8
	//max := 100000000
	max := 100
	keyLength := 16
	dataLength := 1024
	var store Store

	//store = NewStore(StoreConfig{Category: STORE_TYPE_BOLT, Path: "test-data/bolt.db"})
	//if store != nil {
	//	storeWrite(store, StoreBenchmarkTestConfig{name: "boltdb-w", concurrent: concurrent, max: max, keyLength: keyLength, dataLength: dataLength})
	//	storeRead(store, StoreBenchmarkTestConfig{name: "boltdb-r", concurrent: concurrent, max: max, keyLength: keyLength, dataLength: dataLength})
	//	store.Close()
	//}

	store = NewStore(StoreConfig{Type: STORE_TYPE_LEVELDB, Path: "test-data/leveldb"})
	if store != nil {
		storeWrite(store, StoreBenchmarkTestConfig{name: "leveldb-w", concurrent: concurrent, max: max, keyLength: keyLength, dataLength: dataLength})
		storeRead(store, StoreBenchmarkTestConfig{name: "leveldb-r", concurrent: concurrent, max: max, keyLength: keyLength, dataLength: dataLength})
		store.Close()
	}

	store = NewStore(StoreConfig{Type: STORE_TYPE_BADGER, Path: "test-data/badger"})
	if store != nil {
		storeWrite(store, StoreBenchmarkTestConfig{name: "badger-w", concurrent: concurrent, max: max, keyLength: keyLength, dataLength: dataLength})
		storeRead(store, StoreBenchmarkTestConfig{name: "badger-r", concurrent: concurrent, max: max, keyLength: keyLength, dataLength: dataLength})
		store.Close()
	}
}
