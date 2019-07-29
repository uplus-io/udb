/*
 * Copyright (c) 2019 uplus.io
 */

package data

import (
	"fmt"
	"uplus.io/udb/core"
	"uplus.io/udb/hash"
	log "uplus.io/udb/logger"
	"uplus.io/udb/proto"
)

type PartitionDiskType uint8

const (
	PartitionDiskType_Physical PartitionDiskType = iota
	PartitionDiskType_Cloud
)

type Capacity struct {
	Total uint64
	Used  uint64
	Free  uint64
}

type Partition struct {
	Id         int32 //分区Id
	Index      int32
	DataCenter int32  //数据中心Id
	Area       int32  //区Id
	Rack       int32  //机架Id
	Node       int32  //节点Id
	Ip         string //节点Ip
	Port       int    //节点端口
	Path       string //分区根路径

	Replicas *core.Array

	DiskCapacity Capacity
}

func NewPartition(node int32, partition int32) *Partition {
	id := hash.Int32Of(fmt.Sprintf("%d-%d", node, partition))
	return &Partition{Id: id, Node: node, Index: partition, Replicas: core.NewArray()}
}

func (p *Partition) GetId() int32 {
	return p.Id
}

func (p *Partition) Compare(item core.ArrayItem) int32 {
	return p.GetId() - item.GetId()
}

type Node struct {
	DataCenter int32
	Area       int32
	Rack       int32
	Id         int32
	Ip         string
	Port       int

	Status proto.NodeStatus
	Health proto.NodeHealth

	Partitions *core.Array

	Weight        float32
	ReplicaSize   int
	PartitionSize int
}

func NewNode(ip string, port int, weight float32) *Node {
	node := &Node{Ip: ip, Port: port, Weight: weight, Partitions: core.NewArray()}
	node.Id = hash.Int32Of(fmt.Sprintf("%s:%d", node.Ip, node.Port))
	return node
}

func (p *Node) GetId() int32 {
	return p.Id
}

func (p *Node) Compare(item core.ArrayItem) int32 {
	return p.GetId() - item.GetId()
}

type Rack struct {
	Id    int32  //机架Id
	Name  string //机架名
	Nodes *core.Array
}

func NewRack(id int32) *Rack {
	return &Rack{Id: id, Nodes: core.NewArray()}
}

func (p *Rack) GetId() int32 {
	return p.Id
}

func (p *Rack) Compare(item core.ArrayItem) int32 {
	return p.GetId() - item.GetId()
}

func (p *Rack) Node(id int32) *Node {
	return p.Nodes.Id(id).(*Node)
}

func (p *Rack) IfAbsentCreateNode(ip string, port int) *Node {
	return p.Nodes.IfAbsentCreate(NewNode(ip, port, 0)).(*Node)
}

type Area struct {
	Id    int32
	Racks *core.Array
}

func NewArea(id int32) *Area {
	return &Area{Id: id, Racks: core.NewArray()}
}

func (p *Area) GetId() int32 {
	return p.Id
}

func (p *Area) Compare(item core.ArrayItem) int32 {
	return p.GetId() - item.GetId()
}

func (p *Area) Rack(id int32) *Rack {
	return p.Racks.Id(id).(*Rack)
}

func (p *Area) IfAbsentCreateRack(group string) *Rack {
	id := GenerateRepositoryId(group)
	return p.Racks.IfAbsentCreate(NewRack(id)).(*Rack)
}

type DataCenterReplicaStrategy uint8

const (
	DataCenterReplicaStrategy_Circle DataCenterReplicaStrategy = iota //按照一致性哈希环复制
)

type DataCenter struct {
	Id              int32  //数据中心Id
	Name            string //数据中心名
	Areas           *core.Array
	dataHash        *hash.Consistent
	parts           *core.Array
	nodes           *core.Array
	ReplicaStrategy DataCenterReplicaStrategy
}

func NewDataCenter(id int32) *DataCenter {
	return &DataCenter{
		Id:              id,
		Areas:           core.NewArray(),
		dataHash:        hash.NewConsistent(),
		parts:           core.NewArray(),
		nodes:           core.NewArray(),
		ReplicaStrategy: DataCenterReplicaStrategy_Circle,
	}
}

func (p *DataCenter) GetId() int32 {
	return p.Id
}

func (p *DataCenter) Compare(item core.ArrayItem) int32 {
	return p.GetId() - item.GetId()
}

func (p *DataCenter) Area(areaId int32) *Area {
	return p.Areas.Id(areaId).(*Area)
}

func (p *DataCenter) IfAbsentCreateArea(group string) *Area {
	id := GenerateRepositoryId(group)
	return p.Areas.IfAbsentCreate(NewArea(id)).(*Area)
}

func (p *DataCenter) LeaveNode(nodeId uint32) error {
	return nil
}

func (p *DataCenter) UpdateNodeStatus(nodeId uint32, status proto.NodeStatus) error {
	return nil
}

func (p *DataCenter) NextOfRing(ring uint32) uint32 {
	return p.dataHash.NextOfRing(ring)
}

func (p *DataCenter) Nodes() []Node {
	nodes := make([]Node, p.nodes.Len())
	for i := 0; i < p.nodes.Len(); i++ {
		nodes[i] = *(p.nodes.Index(i).(*Node))
	}
	return nodes
}

func (p *DataCenter) addNode(node *Node) error {
	total := int(float32(node.PartitionSize) * node.Weight)
	for i := 0; i < total; i++ {
		part := node.Partitions.Add(NewPartition(node.Id, int32(i))).(*Partition)
		part.DataCenter = node.DataCenter
		part.Area = node.Area
		part.Rack = node.Rack
		p.parts.Add(part)
	}
	address := fmt.Sprintf("%s:%d", node.Ip, node.Port)
	p.dataHash.Add(hash.NewNode(node.Id, address, node.Weight, node.PartitionSize))
	p.nodes.Add(node)
	return nil
}

func (p *DataCenter) Group() {
	p.dataHash.Distribution()
	log.Debugf("dataCenter[%d] ring list-----------------", p.Id)
	p.dataHash.PrintRings()
	if p.ReplicaStrategy == DataCenterReplicaStrategy_Circle {
		p.groupCircle()
	} else {
		panic(fmt.Sprintf("not support replica strategy:%v", p.ReplicaStrategy))
	}
	log.Debugf("dataCenter[%d] node and ring maps-----------------", p.Id)
	p.dataHash.PrintMaps()
}
func (p *DataCenter) groupCircle() {
	for i := 0; i < p.parts.Len(); i++ {
		part := p.parts.Index(i).(*Partition)
		node := p.nodes.Id(part.Node).(*Node)
		replicas := p.dataHash.NextPartition(part.Node, int(part.Index), node.ReplicaSize)
		for _, replica := range replicas {
			partId := hash.Int32Of(fmt.Sprintf("%d-%d", replica.Id, replica.Partition))
			replicaPart := p.parts.Id(partId).(*Partition)
			part.Replicas.Add(replicaPart)
		}
	}
}
