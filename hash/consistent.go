/*
 * Copyright (c) 2019 uplus.io
 */

package hash

import (
	"fmt"
	"sort"
	"sync"
	log "uplus.io/udb/logger"
)

type HashRing []uint32

func (c HashRing) Len() int {
	return len(c)
}

func (c HashRing) Less(i, j int) bool {
	return c[i] < c[j]
}

func (c HashRing) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

type Node struct {
	Id            int
	Value         string
	Weight        float32
	Partition     int
	partitionSize int
}

func NewNode(id int, value string, weight float32, partitionSize int) *Node {
	return &Node{
		Id:            id,
		Value:         value,
		Weight:        weight,
		partitionSize: partitionSize,
	}
}

func (p *Node) Clone() *Node {
	return NewNode(p.Id, p.Value, p.Weight, p.partitionSize)
}

type Consistent struct {
	Nodes     map[uint32]Node
	Resources map[int]bool
	ring      HashRing
	nodeMap   map[int][]int
	sync.RWMutex
}

func NewConsistent() *Consistent {
	return &Consistent{
		Nodes:     make(map[uint32]Node),
		Resources: make(map[int]bool),
		ring:      HashRing{},
		nodeMap:   make(map[int][]int),
	}
}

func (p *Consistent) Add(node *Node) bool {
	p.Lock()
	defer p.Unlock()

	if _, ok := p.Resources[node.Id]; ok {
		return false
	}

	count := int(float32(node.partitionSize) * node.Weight)
	for i := 0; i < count; i++ {
		str := p.joinStr(i, node)
		clone := node.Clone()
		clone.Partition = i
		//p.Nodes[UInt32Of(str)] = *(node)
		p.Nodes[UInt32Of(str)] = *(clone)
	}
	p.Resources[node.Id] = true
	return true
}

func (p *Consistent) Distribution() {
	p.sortHashRing()
	p.initNodeMap()
}

func (p *Consistent) sortHashRing() {
	p.ring = HashRing{}
	for k := range p.Nodes {
		p.ring = append(p.ring, k)
	}
	sort.Sort(p.ring)
}

func (p *Consistent) initNodeMap() {
	p.nodeMap = make(map[int][]int)
	for i, r := range p.ring {
		node := p.Nodes[r]
		rings, exist := p.nodeMap[node.Id]
		if !exist {
			rings = make([]int, 1)
			rings[0] = i
		} else {
			rings = append(rings, i)
		}
		p.nodeMap[node.Id] = rings
	}
}

func (p *Consistent) joinStr(i int, node *Node) string {
	str := fmt.Sprintf("%s-%f-%d", node.Value, node.Weight, i)
	return str
	//return node.Value + "*" + strconv.Itoa(node.Weight) +
	//	"-" + strconv.Itoa(i)
}

func (p *Consistent) Get(key string) Node {
	p.RLock()
	defer p.RUnlock()
	i := p.GetRing(key)
	return p.GetNodeByRing(i)
}

func (p *Consistent) GetRing(key string) int {
	p.RLock()
	defer p.RUnlock()
	hash := UInt32Of(key)
	return p.search(hash)
}

func (p *Consistent) GetNodeByRing(ring int) Node {
	p.RLock()
	defer p.RUnlock()
	return p.Nodes[p.ring[ring]]
}

func (p *Consistent) NextPartition(nodeId int, partition int, size int) []Node {
	for ring, node := range p.Nodes {
		if node.Id == nodeId && node.Partition == partition {
			ringIndex := p.search(ring)

			nodes := make([]Node, size)
			node := p.GetNodeByRing(ringIndex)

			var lastRing int
			var lastNode Node
			lastRing = ringIndex
			var loopTimes = 0
			for i := 0; i < size; i++ {
				for {
					lastRing++
					if lastRing > len(p.ring)-1 {
						lastRing = 0
						loopTimes++
					}
					if loopTimes > 1{
						return nil
					}
					lastNode = p.GetNodeByRing(lastRing)
					if lastNode.Id != node.Id {
						nodes[i] = lastNode
						node = lastNode
						break
					}
				}
			}
			return nodes

		}
	}
	return nil
}

func (p *Consistent) Prev(key string, size int) []Node {
	p.RLock()
	defer p.RUnlock()
	nodes := make([]Node, size)
	ring := p.GetRing(key)
	node := p.GetNodeByRing(ring)
	var lastRing int
	var lastNode Node
	lastRing = ring
	for i := 0; i < size; i++ {
		for {
			lastRing--
			if lastRing < 0 {
				lastRing = len(p.ring) - 1
			}
			lastNode = p.GetNodeByRing(lastRing)
			if lastNode.Id != node.Id {
				nodes[i] = lastNode
				node = lastNode
				break
			}
		}
	}
	return nodes
}

func (p *Consistent) Next(key string, size int) []Node {
	p.RLock()
	defer p.RUnlock()
	nodes := make([]Node, size)
	ring := p.GetRing(key)
	node := p.GetNodeByRing(ring)
	var lastRing int
	var lastNode Node
	lastRing = ring
	for i := 0; i < size; i++ {
		for {
			lastRing++
			if lastRing > len(p.ring)-1 {
				lastRing = 0
			}
			lastNode = p.GetNodeByRing(lastRing)
			if lastNode.Id != node.Id {
				nodes[i] = lastNode
				node = lastNode
				break
			}
		}
	}
	return nodes
}

func (p *Consistent) PrevRing(key string) int {
	p.RLock()
	defer p.RUnlock()

	i := p.GetRing(key)
	if i == 0 {
		i = len(p.ring) - 1
	} else {
		i -= 1
	}
	return i
}

func (p *Consistent) NextRing(key string) int {
	p.RLock()
	defer p.RUnlock()

	i := p.GetRing(key)
	if i == len(p.ring)-1 {
		i = 0
	} else {
		i += 1
	}
	return i
}

func (p *Consistent) search(hash uint32) int {
	//顺时针查找哈希环中接近最大值的hash
	ringLength := len(p.ring)
	i := sort.Search(ringLength, func(i int) bool { return p.ring[i] >= hash })
	//判断查找的接近最大值 超出索引和匹配索引
	if i < ringLength {
		//最后一个哈希，则返回首个以形成回环；否则返回索引值
		return i
	} else {
		return 0
	}
}

func (p *Consistent) Remove(node *Node) {
	p.Lock()
	defer p.Unlock()

	if _, ok := p.Resources[node.Id]; !ok {
		return
	}

	delete(p.Resources, node.Id)

	count := int(float32(node.partitionSize) * node.Weight)
	for i := 0; i < count; i++ {
		str := p.joinStr(i, node)
		delete(p.Nodes, UInt32Of(str))
	}
	p.sortHashRing()
}

func (p *Consistent) PrintRings() {
	var last = 0
	for i, r := range p.ring {
		node := p.Nodes[r]
		log.Infof("index:%d ring:%d node:[%d-%d value:%s] lastDifference:%dK",
			i, r,
			node.Id, node.Partition, node.Value,
			(int(r)-last)/1000)
		last = int(r)
	}

	//for id, rings := range p.nodeMap {
	//	log.Infof("node:%d rings:%v", id, rings)
	//}
}
