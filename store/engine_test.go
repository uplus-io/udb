/*
 * Copyright (c) 2019 uplus.io
 */

package store

import (
	"fmt"
	"testing"
	"uplus.io/udb/config"
)

func TestEngine_Namespaces(t *testing.T) {
	engine := createEngine()
	defer engine.Close()
	table := engine.Table()
	table.IfAbsentNS("uplus")
	for _, ns := range table.NSs() {
		fmt.Printf("ns[%s]\n", ns.String())
	}
}

func TestEngine_Tables(t *testing.T) {
	engine := createEngine()
	defer engine.Close()
	table := engine.Table()
	table.IfAbsentTab("_user")
	for _, tab := range table.Tabs() {
		fmt.Printf("tab[%s]\n", tab.String())
	}
}

func TestEngine_MetaIterator(t *testing.T) {
	engine := createEngine()
	defer engine.Close()
	engine.MetaForEach(func(data Data) bool {
		fmt.Printf("key[%s/%s/%s] val len:%d\n", data.Id.Namespace, data.Id.Table, string(data.Id.Key), len(data.Content))
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
