/*
 * Copyright (c) 2019 uplus.io
 */

package data

import (
	"strings"
	"sync"
	"uplus.io/udb/core"
	log "uplus.io/udb/logger"
	"uplus.io/udb/proto"
	"uplus.io/udb/utils"
)

type Warehouse struct {
	Centers       *core.Array
	applicants    *core.Array
	communication proto.ClusterCommunication

	readying bool
	sync.RWMutex
}

func NewWarehouse(communication proto.ClusterCommunication) *Warehouse {
	return &Warehouse{Centers: core.NewArray(), applicants: core.NewArray(), communication: communication}
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

func (p *Warehouse) GetNode(dc int32, nodeId int32) *Node {
	p.RLock()
	defer p.RUnlock()
	center := p.Centers.Id(dc).(*DataCenter)
	if center != nil {
		return center.nodes.Id(nodeId).(*Node)
	}
	return nil
}

func (p *Warehouse) JoinNode(ip string, port int) *Node {
	node := NewNode(ip, port, 0)
	if p.readying {
		p.communication.SendNodeInfoTo(node.Id)
	} else {
		p.applicants.Add(node)
		log.Infof("cluster applicant[%d:%s:%d] join", node.Id, node.Ip, node.Port)
	}
	return node
}

func (p *Warehouse) LeaveNode(ip string, port int) *Node {
	node := NewNode(ip, port, 0)
	if p.readying {
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

func BuildNode(ip string, port int) *Node {
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

func (p *Warehouse) Readying() {
	p.Lock()
	defer p.Unlock()
	p.readying = true
}

func (p *Warehouse) print() {
	for i := 0; i < p.Centers.Len(); i++ {
		center := p.Centers.Index(i).(*DataCenter)
		log.Debugf("%d dataCenter[%d] has %d areas", i, center.Id, center.Area.Len())
		for j := 0; j < center.Area.Len(); j++ {
			area := center.Area.Index(j).(*Area)
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
