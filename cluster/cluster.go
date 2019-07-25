/*
 * Copyright (c) 2019 uplus.io
 */

package cluster

import (
	"sync"
	"uplus.io/udb/config"
	"uplus.io/udb/data"
	"uplus.io/udb/logger"
	log "uplus.io/udb/logger"
	"uplus.io/udb/proto"
	"uplus.io/udb/store"
)

type Cluster struct {
	id     int32
	config config.ClusterConfig

	engine        *store.Engine
	warehouse     *data.Warehouse
	communication proto.ClusterCommunication

	packetDispatcher PacketDispatcher

	exit chan bool

	transport Transport
	pipeline  Pipeline

	clusterHealth proto.ClusterHealth
	nodeHealth    proto.NodeHealth
	nodeStatus    proto.NodeStatus

	launched bool

	sync.RWMutex
}

func NewCluster(config config.ClusterConfig) *Cluster {
	//init log
	logger.DebugLoggerEnable(true)
	logConfig := config.LogConfig
	level := logger.LoggerLevel(logConfig.Level)
	logger.NewLogger(level, logConfig.Path, logConfig.Filename)
	return &Cluster{
		config:        config,
		pipeline:      NewPipelinePacket(),
		exit:          make(chan bool),
		clusterHealth: proto.ClusterHealth_CH_Unknown,
		nodeHealth:    proto.NodeHealth_Suspect,
		nodeStatus:    proto.NodeStatus_Unknown,
	}
}

func (p *Cluster) Listen() {
	p.communication = NewClusterCommunicationImplementor(p)
	p.startEngine()
	p.launchGossip()
	go p.packetInLoop()
	go p.packetOutLoop()
	p.nodeHealth = proto.NodeHealth_Alive
	<-p.exit
}

func (p *Cluster) startEngine() {
	p.engine = store.NewEngine(p.config.StorageConfig)
	log.Infof("cluster storage engine launched")
	p.collectLocalInfo()
}

func (p *Cluster) getClusterHealth() proto.ClusterHealth {
	return p.clusterHealth
}

func (p *Cluster) getLocalNodeHealth() proto.NodeHealth {
	return p.nodeHealth
}

func (p *Cluster) getLocalNodeStatus() proto.NodeStatus {
	return p.nodeStatus
}

//collect local nodeInfo
func (p *Cluster) collectLocalInfo() *proto.NodeInfo {
	repository := p.engine.Table().Repository()
	if repository.DataCenter == 0 {
		p.clusterHealth = proto.ClusterHealth_CH_NotInitialize
	}
	partSize := len(p.config.StorageConfig.Partitions)
	replicaSize := 2
	partitions := make([]*proto.Partition, partSize)
	for i := 0; i < partSize; i++ {
		partitions[i] = p.engine.Table().PartitionOfIndex(int32(i))
	}
	log.Infof("cluster local info loaded")
	return &proto.NodeInfo{
		Repository:    repository,
		PartitionSize: int32(partSize),
		ReplicaSize:   int32(replicaSize),
		Health:        p.getLocalNodeHealth(),
		Status:        p.getLocalNodeStatus(),
		Partitions:    partitions,
	}
}

func (p *Cluster) launchGossip() {
	p.Lock()
	defer p.Unlock()
	p.packetDispatcher = NewPacketSystemDispatcher(p)
	p.warehouse = data.NewWarehouse(p.communication)
	transportConfig := &TransportConfig{
		Seeds:          p.config.Seeds,
		Secret:         p.config.SecurityConfig.Secret,
		BindIp:         p.config.BindIp,
		BindPort:       p.config.BindPort,
		AdvertisePort:  p.config.AdvertisePort,
		EventListener:  NewClusterEventListener(p.warehouse),
		PacketListener: NewClusterPacketListener(p.pipeline)}

	p.transport = NewTransportGossip(transportConfig)
	transportInfo := p.transport.Serving()
	p.id = transportInfo.Id
	p.engine.ValidatePartition(p.id)
	p.launched = true
	log.Debugf("cluster node[%d] started %v", p.id, p.launched)

	localInfo := p.collectLocalInfo()
	p.JoinNode(p.id, int(localInfo.PartitionSize), int(localInfo.ReplicaSize))
	p.contactCluster()
}

func (p *Cluster) contactCluster() {
	applicants := p.warehouse.Applicants()
	for i := 0; i < applicants.Len(); i++ {
		node := applicants.Index(i).(*data.Node)
		if p.id != node.Id {
			err := p.communication.SendNodeInfoTo(node.Id)
			if err != nil {
				log.Errorf("contact cluster[%d->%d] error", p.id, node.Id)
			}
		}
	}
	p.warehouse.Readying()
}

func (p *Cluster) JoinNode(nodeId int32, partitionSize int, replicaSize int) {
	node := p.transport.Node(nodeId)

	p.warehouse.AddNode(data.NewNode(node.Addr.String(), int(node.Port), 1), partitionSize, replicaSize)
	p.warehouse.Group()
}

func (p *Cluster) packetInLoop() {
	for {
		packet := <-p.pipeline.In()
		log.Debugf("send packet[%s]", packet.String())
		//node := p.warehouse.GetNode(packet.GetDataCenter(), packet.GetTo())
		//if node != nil {
		//}
	}
}

func (p *Cluster) packetOutLoop() {
	for {
		packet := <-p.pipeline.Out()
		log.Debugf("received packet[%s]", packet.String())
		go func() {
			err := p.packetDispatcher.Dispatch(*packet)
			if err != nil {
				p.pipeline.Out() <- packet
			}
		}()

		//node := p.warehouse.GetNode(packet.GetDataCenter(), packet.GetTo())
		//if node != nil {
		//}
	}
}
