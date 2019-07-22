/*
 * Copyright (c) 2019 uplus.io
 */

package udb

import (
	"fmt"
	"strconv"
	"testing"
	"uplus.io/udb/hash"
)

func TestBaseStream_WriteAndRead(t *testing.T) {
	//base := Base{DataType: BASE_TYPE_DATETIME}
	//fmt.Println(base.StorageLength())
	//stream := NewBaseStream("data")
	//stream.Write(int(1))
	//stream.Write(int64(9999))
	//stream.Write(bool(true))
	//stream.Write([]byte("abcdefg"))
	//stream.Write(string("abcdefghijklmnopqrstuvwxyz"))
}

func TestConsistentHash(t *testing.T) {
	consistent := hash.NewConsistent()

	for i := 0; i < 10; i++ {
		si := fmt.Sprintf("%d", i)
		consistent.Add(hash.NewNode(i, "172.18.1."+si, 1))
	}

	for k, v := range consistent.Nodes {
		fmt.Println("recoverMeta:", k, " Value:", v.Value)
	}

	ipMap := make(map[string]int, 0)
	for i := 0; i < 1000; i++ {
		si := fmt.Sprintf("key%d", i)
		k := consistent.Get(si)
		if _, ok := ipMap[k.Value]; ok {
			ipMap[k.Value] += 1
		} else {
			ipMap[k.Value] = 1
		}
	}

	for k, v := range ipMap {
		fmt.Println("Node IP:", k, " count:", v)
	}
}

func Test32Hash(t *testing.T) {
	consistent := hash.NewConsistentOf(2, 1)

	//balanceMap := make(map[string]int, 0)
	//for i := 0; i < 10000000; i++ {
	//	key := fmt.Sprintf("%d", i)
	//	//key := fmt.Sprintf("%d", time.Now().UnixNano())
	//	node := consistent.Get(key)
	//	//fmt.Println("key:", key, " id:", node.GetId, " value:", node.Value, )
	//
	//	if _, ok := balanceMap[node.Value]; ok {
	//		balanceMap[node.Value] += 1
	//	} else {
	//		balanceMap[node.Value] = 1
	//	}
	//}
	//fmt.Println("hash balance--------")
	//for k1, v := range balanceMap {
	//	fmt.Println("node:", k1, " count:", v)
	//}

	k1 := strconv.Itoa(1)
	node := consistent.Get(k1)
	prevNode := consistent.Prev(k1)
	nextNode := consistent.Next(k1)

	consistent.Next(strconv.Itoa(2))
	consistent.Next(strconv.Itoa(3))
	consistent.Next(strconv.Itoa(4))
	consistent.Next(strconv.Itoa(5))
	consistent.Next(strconv.Itoa(6))
	consistent.Next(strconv.Itoa(7))
	consistent.Next(strconv.Itoa(8))
	consistent.Next(strconv.Itoa(9))
	consistent.Next(strconv.Itoa(10))
	consistent.Next(strconv.Itoa(11))

	fmt.Println(node)
	fmt.Println(prevNode)
	fmt.Println(nextNode)


}
