/*
 * Copyright (c) 2019 uplus.io
 */

package store

import (
	"fmt"
	"testing"
	"uplus.io/udb/config"
	"uplus.io/udb/proto"
	"uplus.io/udb/utils"
)

func TestEngine_AddPart(t *testing.T) {
	engine := createEngine()
	defer engine.Close()
	engine.AddPart(proto.Partition{Id: 1000, Version: VERSION, Index: 0})
	engine.AddPart(proto.Partition{Id: 1001, Version: VERSION, Index: 1})
	fmt.Printf("part0 info - %s\n", engine.GetPart(1000).String())
	fmt.Printf("part1 info - %s\n", engine.GetPart(1001).String())
}

func TestEngine_SetData(t *testing.T) {
	engine := createEngine()
	defer engine.Close()
	var err error
	identity1 := NewIdOfData("udb", "user", []byte("1"))
	identity2 := NewIdOfData("udb", "user", []byte("2"))
	err = engine.SetData(1000, *NewData(*identity1, utils.LInt32ToBytes(101)))
	printError(err, "set data error")
	err = engine.SetData(1001, *NewData(*identity2, utils.LInt32ToBytes(102)))
	printError(err, "set data error")

	val1, _ := engine.GetData(1000, *identity1)
	printDataContent(identity1.Key, val1)
	val2, _ := engine.GetData(1001, *identity2)
	printDataContent(identity2.Key, val2)
	printAllData(engine)
}

func printError(err error, format string, args ...interface{}) {
	if err != nil {
		msg := fmt.Sprintf(format, args...)
		fmt.Printf("%s - %v\n", msg, err)
	}
}

func printDataContent(key []byte, data *proto.DataContent) {
	fmt.Printf("data key:%s content[%s]\n", key, data.String())
}

func printAllData(engine *Engine) {
	fmt.Printf("meta part|meta---------------------------\n")
	metaOperator := engine.meta
	printMetaStorage(metaOperator)
	fmt.Printf("meta part|data---------------------------\n")
	printDataStorage(metaOperator)
	fmt.Printf("meta ns tab---------------------------\n")
	printNSAndTAB(metaOperator)

	for i := 0; i < len(engine.partitions); i++ {
		fmt.Printf("data part[%d]|meta---------------------------\n", i)
		operator := engine.partitions[i]
		printMetaStorage(operator)
		fmt.Printf("data part[%d]|data---------------------------\n", i)
		printDataStorage(operator)
		fmt.Printf("data part[%d]|ns tab---------------------------\n", i)
		printNSAndTAB(operator)
	}
}

func printNSAndTAB(operator StoreOperator) {
	nss := operator.NSs()
	for _, ns := range nss {
		fmt.Printf("ns[%s]---------------------------\n", ns.String())
		tables := operator.TABs(ns.Name)
		for _, tab := range tables {
			fmt.Printf("tab[%s]---------------------------\n", tab.String())
		}
	}
}

func printMetaStorage(operator StoreOperator) {
	operator.MetaForEach(func(key []byte, meta proto.DataMeta) bool {
		printMetaInfo(key,&meta)
		return true
	})
}

func printDataStorage(operator StoreOperator) {
	operator.DataForEach(func(key, data []byte) bool {
		printDataInfo(key, data)
		return true
	})
}

func printDataInfo(key []byte, data []byte) {
	fmt.Printf("data key:%s valLen:%d\n", string(key), len(data))
}

func printMetaInfo(key []byte, data *proto.DataMeta) {
	fmt.Printf("meta key:%s val:%s\n", string(key), data.String())
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
