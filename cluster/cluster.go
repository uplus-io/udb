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

	engine            *store.Engine
	warehouse         *data.Warehouse
	communication     proto.ClusterCommunication
	dataCommunication proto.ClusterDataCommunication

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
	p.dataCommunication = NewClusterDataCommunicationImplementor(p)
	p.startEngine()
	go p.launchGossip()
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

	//server launch
	transportInfo := p.transport.Serving()
	p.id = transportInfo.Id
	//storage partition validate
	p.engine.ValidatePartition(p.id)

	p.launched = true
	log.Debugf("cluster node[%d] started %v", p.id, p.launched)

	localInfo := p.collectLocalInfo()
	p.JoinNode(p.id, int(localInfo.PartitionSize), int(localInfo.ReplicaSize))
	p.checkWarehouse()
}

func (p *Cluster) checkWarehouse() {
	p.contactCluster()
	p.warehouse.Readying(func(status data.WarehouseStatus) {
		transportInfo := p.transport.Me()
		repository := proto.ParseRepository(transportInfo.Addr.String())
		parts := p.engine.Parts()
		for _, part := range parts {
			center := p.warehouse.GetCenter(repository.DataCenter)
			next := center.NextOfRing(uint32(part.Id))
			request := &proto.DataMigrateRequest{StartRing: part.Id, EndRing: int32(next)}

			for _, node := range center.Nodes() {
				if node.Id != transportInfo.Id {
					p.dataCommunication.MigrateRequest(node.Id, request)
				}
			}
		}
	})
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

}

func (p *Cluster) JoinNode(nodeId int32, partitionSize int, replicaSize int) {
	node := p.transport.Node(nodeId)

	p.warehouse.AddNode(data.NewNode(node.Addr.String(), int(node.Port), 1), partitionSize, replicaSize)
	p.warehouse.Group()

}

func (p *Cluster) SendAsyncPacket(packet *proto.Packet) {
	p.pipeline.InWrite(packet)
}

func (p *Cluster) SendSyncPacket(packet *proto.Packet) *proto.Packet {
	channel := p.pipeline.InSyncWrite(packet)
	return <-channel.Read()
}

func (p *Cluster) packetInLoop() {
	for {
		packet := <-p.pipeline.InRead()
		log.Debugf("send packet[%s]", packet.String())
		bytes, err := proto.Marshal(packet)
		if err != nil {
			log.Errorf("waiting to send packet marshal error[%s]", packet.String())
			continue
		}
		if packet.Mode == proto.PacketMode_TCP {
			err := p.transport.SendToTCP(packet.To, bytes)
			if err != nil {
				log.Errorf("sending tcp packet error[%s]", packet.String())
				continue
			}
		} else if packet.Mode == proto.PacketMode_UDP {
			err := p.transport.SendToUDP(packet.To, bytes)
			if err != nil {
				log.Errorf("sending udp packet error[%s]", packet.String())
				continue
			}
		}
		//node := p.warehouse.GetNode(packet.GetDataCenter(), packet.GetTo())
		//if node != nil {
		//}
	}
}

func (p *Cluster) packetOutLoop() {
	for {
		packet := <-p.pipeline.OutRead()
		log.Debugf("received packet[%s]", packet.String())
		go func() {
			err := p.packetDispatcher.Dispatch(*packet)
			if err != nil {
				p.pipeline.OutWrite(packet)
			}
		}()

		//node := p.warehouse.GetNode(packet.GetDataCenter(), packet.GetTo())
		//if node != nil {
		//}
	}
}
