package cluster

import "uplus.io/udb/proto"

type ClusterDataCommunicationImplementor struct {
	cluster *Cluster
}

func NewClusterDataCommunicationImplementor(cluster *Cluster) *ClusterDataCommunicationImplementor {
	return &ClusterDataCommunicationImplementor{cluster: cluster}
}

func (p *ClusterDataCommunicationImplementor) Push(to int32, request *proto.PushRequest) *proto.PushResponse {
	content, _ := proto.Marshal(request)
	packet := proto.NewTCPPacket(proto.PacketType_DataPush, p.cluster.id, to, content)
	responseData := p.cluster.SendSyncPacket(packet)
	packetReply := &proto.Packet{}
	proto.Unmarshal(responseData.Content, packetReply)
	response := &proto.PushResponse{}
	proto.Unmarshal(packetReply.Content, response)
	return response
}

func (p *ClusterDataCommunicationImplementor) PushReply(to int32, request *proto.PushResponse) error {
	content, _ := proto.Marshal(request)
	packet := proto.NewTCPPacket(proto.PacketType_DataPushReply, p.cluster.id, to, content)
	p.cluster.SendAsyncPacket(packet)
	return nil
}

func (p *ClusterDataCommunicationImplementor) Pull(to int32, request *proto.PullRequest) *proto.PullResponse {
	content, _ := proto.Marshal(request)
	responseData := proto.NewTCPPacket(proto.PacketType_DataPull, p.cluster.id, to, content)
	packetReply := &proto.Packet{}
	proto.Unmarshal(responseData.Content, packetReply)
	response := &proto.PullResponse{}
	proto.Unmarshal(packetReply.Content, response)
	return response
}

func (p *ClusterDataCommunicationImplementor) PullReply(to int32, request *proto.PullResponse) error {
	content, _ := proto.Marshal(request)
	packet := proto.NewTCPPacket(proto.PacketType_DataPullReply, p.cluster.id, to, content)
	p.cluster.SendAsyncPacket(packet)
	return nil
}

func (p *ClusterDataCommunicationImplementor) MigrateRequest(to int32, request *proto.DataMigrateRequest) error {
	bytes, _ := proto.Marshal(request)
	packet := proto.NewTCPPacket(proto.PacketType_DataMigrate, p.cluster.id, to, bytes)
	p.cluster.SendAsyncPacket(packet)
	return nil
}

func (p *ClusterDataCommunicationImplementor) MigrateResponse(to int32, request *proto.DataMigrateResponse) error {
	bytes, _ := proto.Marshal(request)
	packet := proto.NewTCPPacket(proto.PacketType_DataMigrateReply, p.cluster.id, to, bytes)
	p.cluster.SendAsyncPacket(packet)
	return nil
}
