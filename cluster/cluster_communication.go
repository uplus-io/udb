package cluster

import (
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
	nodeInfoData, err := proto.Marshal(nodeInfo)
	if err != nil {
		return err
	}
	clusterStat := proto.NewTCPPacket( proto.PacketType_SystemHi, int32(transport.Me().Id), to, nodeInfoData)
	p.cluster.SendAsyncPacket(clusterStat)
	return nil
}
