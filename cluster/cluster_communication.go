package cluster

import (
	log "uplus.io/udb/logger"
	"uplus.io/udb/proto"
)

type ClusterCommunicationImplementor struct {
	cluster *Cluster
}

func NewClusterCommunicationImplementor(cluster *Cluster) *ClusterCommunicationImplementor {
	return &ClusterCommunicationImplementor{cluster: cluster}
}

func (p *ClusterCommunicationImplementor) SendNodeInfoTo(to int32) error {
	transport := p.cluster.transport

	nodeInfo := p.cluster.collectLocalInfo()
	nodeInfoData, _ := proto.Marshal(nodeInfo)

	clusterStat := proto.NewPacket(proto.PacketType_SystemHi, int32(transport.Me().Id), to, nodeInfoData)

	statData, err := proto.Marshal(clusterStat)
	if err != nil {
		log.Errorf("marshal cluster stat error")
		return err
	}
	return transport.SendToTCP(to, statData)
}
