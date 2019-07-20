/*
 * Copyright (c) 2019 uplus.io
 */

package store

import (
	"fmt"
	"testing"
	"uplus.io/udb/config"
)

func TestEngine_SystemNamespace(t *testing.T) {
	engine := createEngine()
	defer engine.Close()
	engine.IfAbsentCreateNamespace("uplus")
	for _, ns := range engine.SystemNamespaces() {
		fmt.Printf("ns[id:%d name:%s]\n", ns.Id, ns.Name)
	}
}

func TestEngine_MetaIterator(t *testing.T) {
	engine := createEngine()
	defer engine.Close()
	engine.meta.ForEach(func(data Data) bool {
		fmt.Printf("key[%s/%s/%s] val:[%s]\n", data.Id.Namespace, data.Id.Table, string(data.Id.Key), string(data.Content))
		return true
	})
}

func createEngine() *Engine {
	cfg := config.StorageConfig{
		Engine: "BADGER",
		Meta:   "./test-data/engine-test/meta",
		Wal:    "./test-data/engine-test/wal",
		Partitions: []string{
			"./test-data/engine-test/data/0",
			"./test-data/engine-test/data/1",
		},
	}
	return NewEngine(cfg)
}
