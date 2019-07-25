package cluster

import "uplus.io/udb/proto"

type ClusterDataCommunicationImplementor struct {
	cluster *Cluster
}

func (p *ClusterDataCommunicationImplementor) Push(to int32, request *proto.PushRequest) *proto.PushResponse {
	content, _ := proto.Marshal(request)
	packet := proto.NewPacket(proto.PacketType_DataPush, p.cluster.id, to, content)
	bytes, _ := proto.Marshal(packet)
	responseData, err := p.cluster.transport.SyncSendToTCP(to, bytes)
	if err != nil {
		return nil
	}
	packetReply := &proto.Packet{}
	proto.Unmarshal(responseData, packetReply)
	response := &proto.PushResponse{}
	proto.Unmarshal(packetReply.Content, response)
	return response
}

func (p *ClusterDataCommunicationImplementor) Pull(to int32, request *proto.PullRequest) *proto.PullRequest {
	return nil
}

func (p *ClusterDataCommunicationImplementor) MigrateRequest(to int32, request *proto.DataMigrateRequest) error {
	return nil
}

func (p *ClusterDataCommunicationImplementor) MigrateResponse(to int32, request *proto.DataMigrateResponse) error {
	return nil
}
