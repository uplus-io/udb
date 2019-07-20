/*
 * Copyright (c) 2019 uplus.io
 */

package hash

import (
	"fmt"
	"strconv"
	"testing"
)

func TestConsistent(t *testing.T) {
	consistent := NewConsistent()

	consistent.Add(NewNode(1, "192.168.0.1:1107", 1, 3))
	consistent.Add(NewNode(2, "192.168.0.1:1108", 1, 3))
	consistent.Add(NewNode(3, "192.168.0.1:1109", 1, 3))
	consistent.Add(NewNode(4, "192.168.0.1:1110", 1, 3))
	consistent.Add(NewNode(5, "192.168.0.1:1111", 1, 3))
	consistent.Add(NewNode(6, "192.168.0.1:1112", 1, 3))

	consistent.Add(NewNode(7, "192.168.0.2:1107", 1, 3))
	consistent.Add(NewNode(8, "192.168.0.2:1108", 1, 3))
	consistent.Add(NewNode(9, "192.168.0.3:1107", 1, 3))
	consistent.Add(NewNode(10, "192.168.0.3:1108", 1, 3))
	consistent.Add(NewNode(11, "192.168.0.4:1107", 1, 3))
	consistent.Add(NewNode(12, "192.168.0.4:1107", 1, 3))

	consistent.Add(NewNode(13, "192.168.1.1:1107", 1, 3))
	consistent.Add(NewNode(14, "192.168.1.2:1107", 1, 3))
	consistent.Add(NewNode(15, "192.168.1.3:1107", 1, 3))
	consistent.Add(NewNode(16, "192.168.2.1:1107", 1, 3))
	consistent.Add(NewNode(17, "192.168.3.1:1107", 1, 3))
	consistent.Distribution()
	consistent.PrintRings()

	testNextPartitions(consistent, 1, 0, 3)
	testNextPartitions(consistent, 12, 2, 3)

	var key string
	var node Node
	var ring1, ring2 int
	for i := 0; i < 100; i++ {
		key = strconv.Itoa(i)
		node = consistent.Get(key)
		ring1 = consistent.PrevRing(key)
		ring2 = consistent.NextRing(key)
		nextNodes := consistent.Next(key, 3)
		prevNodes := consistent.Prev(key, 3)

		fmt.Printf("before key:%s node:%v\n", key, node)
		fmt.Printf("before key:%s prev Ring:%d node:%v\n", key, ring1, consistent.GetNodeByRing(ring1))
		fmt.Printf("before key:%s next Ring:%d node:%v\n", key, ring2, consistent.GetNodeByRing(ring2))
		fmt.Printf("before key:%s nextNodes:%v\n", key, nextNodes)
		fmt.Printf("before key:%s prevNodes:%v\n", key, prevNodes)
		consistent.Remove(NewNode(5, "192.168.0.1:1111", 1, 3))

		key = strconv.Itoa(i)
		node = consistent.Get(key)
		//next = consistent.Next(key)

		fmt.Printf("after key:%s node:%s\n", key, node.Value)

		consistent.Add(NewNode(5, "192.168.0.1:1111", 1, 3))

		fmt.Println(len(consistent.Nodes))

	}
}

func testNextPartitions(consistent *Consistent, currentNode int, currentPart int, size int) []Node {
	nextPartition := consistent.NextPartition(currentNode, currentPart, size)
	fmt.Printf("node[%d-%d] next %d partition:%v\n", currentNode, currentPart, size, nextPartition)
	return nextPartition
}
