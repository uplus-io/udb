/*
 * Copyright (c) 2019 uplus.io
 */

package data

import (
	"strings"
	"sync"
	"time"
	"uplus.io/udb/core"
	log "uplus.io/udb/logger"
	"uplus.io/udb/proto"
	"uplus.io/udb/utils"
)

const (
	warehouseAvailableInterval = 5 //可用时间间隔 second
)

type WarehouseListener func(status WarehouseStatus)

type WarehouseStatus uint8

const (
	WarehouseStatus_Launching WarehouseStatus = iota
	WarehouseStatus_Normal
	WarehouseStatus_Node_Changed
)

type Warehouse struct {
	Centers       *core.Array
	applicants    *core.Array
	communication proto.ClusterCommunication

	available             bool  //可用
	lastCommunicationTime int64 //最后节点状态变化通讯时间
	clusterReadying       bool  //当前节点服务已就绪标记
	status                WarehouseStatus

	listener WarehouseListener

	sync.RWMutex
}

func NewWarehouse(communication proto.ClusterCommunication) *Warehouse {
	return &Warehouse{Centers: core.NewArray(), applicants: core.NewArray(), communication: communication, status: WarehouseStatus_Launching}
}

func GenerateRepositoryId(group string) int32 {
	return utils.StringToInt32(group)
}

func (p *Warehouse) IfAbsentCreateDataCenter(group string) *DataCenter {
	id := GenerateRepositoryId(group)
	return p.Centers.IfAbsentCreate(NewDataCenter(id)).(*DataCenter)
}

func (p *Warehouse) IfPresent(ipv4 string) *DataCenter {
	return nil
}

func (p *Warehouse) GetCenter(dc int32) *DataCenter {
	p.RLock()
	defer p.RUnlock()
	return p.Centers.Id(dc).(*DataCenter)
}

func (p *Warehouse) GetNode(dc int32, nodeId int32) *Node {
	p.RLock()
	defer p.RUnlock()
	center := p.GetCenter(dc)
	if center != nil {
		return center.nodes.Id(nodeId).(*Node)
	}
	return nil
}

func (p *Warehouse) JoinNode(ip string, port int) *Node {
	p.lastCommunicationTime = time.Now().Unix()
	p.status = WarehouseStatus_Node_Changed
	node := NewNode(ip, port, 0)
	if p.clusterReadying {
		p.communication.SendNodeInfoTo(node.Id)
	} else {
		p.applicants.Add(node)
		log.Infof("cluster applicant[%d:%s:%d] join", node.Id, node.Ip, node.Port)
	}
	return node
}

func (p *Warehouse) LeaveNode(ip string, port int) *Node {
	p.lastCommunicationTime = time.Now().Unix()

	node := NewNode(ip, port, 0)
	if p.clusterReadying {
		//todo:set node invalid
	} else {
		p.applicants.Delete(node.Id)
		log.Infof("cluster applicant[%d:%s:%d] leave", node.Id, node.Ip, node.Port)
	}
	return node
}

func (p *Warehouse) AddNode(node *Node, partitionSize int, replicaSize int) error {
	p.Lock()
	defer p.Unlock()
	bits := strings.Split(node.Ip, ".")
	center := p.IfAbsentCreateDataCenter(bits[0])
	area := center.IfAbsentCreateArea(bits[1])
	rack := area.IfAbsentCreateRack(bits[2])
	newNode := rack.IfAbsentCreateNode(node.Ip, node.Port)
	newNode.DataCenter = center.Id
	newNode.Area = area.Id
	newNode.Rack = rack.Id
	//todo:需要注意 分区数与比重与已存数据不一致问题
	newNode.Weight = node.Weight
	newNode.PartitionSize = partitionSize
	newNode.ReplicaSize = replicaSize

	center.addNode(newNode)

	node = newNode
	return nil
}

func (p *Warehouse) Group() {
	p.Lock()
	defer p.Unlock()
	for i := 0; i < p.Centers.Len(); i++ {
		center := p.Centers.Index(i).(*DataCenter)
		center.Group()
	}
}

func (p *Warehouse) Applicants() *core.Array {
	return p.applicants
}

func (p *Warehouse) Readying(listener WarehouseListener) {
	p.Lock()
	defer p.Unlock()
	p.clusterReadying = true
	p.listener = listener
	go p.selfCheck()
}

func (p *Warehouse) selfCheck() {
	for {
		ts := time.Now().Unix() - p.lastCommunicationTime
		log.Debugf("self check,ts:%d interval:%d", ts, warehouseAvailableInterval)
		if ts >= warehouseAvailableInterval {
			p.status = WarehouseStatus_Normal
			p.listener(p.status)
			break
			//runtime.Gosched()
		}
		time.Sleep(time.Second)
	}
}

func (p *Warehouse) print() {
	for i := 0; i < p.Centers.Len(); i++ {
		center := p.Centers.Index(i).(*DataCenter)
		log.Debugf("%d dataCenter[%d] has %d areas", i, center.Id, center.Areas.Len())
		for j := 0; j < center.Areas.Len(); j++ {
			area := center.Areas.Index(j).(*Area)
			log.Debugf("    %d area[%d] has %d racks", j, area.Id, area.Racks.Len())
			for k := 0; k < area.Racks.Len(); k++ {
				rack := area.Racks.Index(k).(*Rack)
				log.Debugf("        %d rack[%d] has %d nodes", j, rack.Id, rack.Nodes.Len())
				for l := 0; l < rack.Nodes.Len(); l++ {
					node := rack.Nodes.Index(l).(*Node)
					log.Debugf("            %d node[id:%d dataCenter:%d area:%d rack:%d] part[size:%d] replica[size:%d]", k, node.Id, node.DataCenter, node.Area, node.Rack, node.PartitionSize, node.ReplicaSize)
					for m := 0; m < node.Partitions.Len(); m++ {
						part := node.Partitions.Index(m).(*Partition)
						log.Debugf("                %d node[%d-%d] part[id:%d index:%d dataCenter:%d area:%d rack:%d] has replicas:%d",
							m, part.Node, part.Index, part.Id, part.Index, part.DataCenter, part.Area, part.Rack, part.Replicas.Len())
						for n := 0; n < part.Replicas.Len(); n++ {
							replica := part.Replicas.Index(n).(*Partition)
							log.Debugf("                    %d node[%d-%d] replica[id:%d index:%d dataCenter:%d area:%d rack:%d]",
								n, replica.Node, replica.Index, replica.Id, replica.Index, replica.DataCenter, replica.Area, replica.Rack)
						}
					}
				}
			}
		}
	}
}
