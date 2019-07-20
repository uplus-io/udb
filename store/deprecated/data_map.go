/*
 * Copyright (c) 2019 uplus.io
 */

package deprecated

import (
	"sync"
)

type DataMapConfig struct {
}

type DataMapItem interface {
}
type DataMap struct {
	sync.RWMutex
}

func NewDataMap(config DataMapConfig) *DataMap {
	return &DataMap{}
}

func init() {
	//var max = 100;
	//tree := btree.NewWithIntComparator(3)
	//for i := 0; i < max; i++ {
	//	tree.Put(i+1, fmt.Sprintf("%d", i+1))
	//}
	//fmt.Println(tree)

	//value, found := tree.Get(1)
	//if found {
	//	fmt.Println(value)
	//}

	//fmt.Println(tree.Height())
	//json,_ := tree.ToJSON()
	//fmt.Println(string(json))

}
