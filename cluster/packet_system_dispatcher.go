/*
 * Copyright (c) 2019 uplus.io
 */

package cluster

import (
	"uplus.io/udb"
	log "uplus.io/udb/logger"
	"uplus.io/udb/proto"
	"uplus.io/udb/store"
)

type PacketSystemDispatcher struct {
	cluster *Cluster

	handlerMap map[proto.PacketType]PacketHandler
}

func NewPacketSystemDispatcher(cluster *Cluster) *PacketSystemDispatcher {
	dispatcher := &PacketSystemDispatcher{cluster: cluster, handlerMap: make(map[proto.PacketType]PacketHandler)}
	dispatcher.register(proto.PacketType_SystemHi, dispatcher.handleClusterHi)
	dispatcher.register(proto.PacketType_DataMigrate, dispatcher.handleMigrate)
	dispatcher.register(proto.PacketType_DataMigrateReply, dispatcher.handleMigrateReply)
	dispatcher.register(proto.PacketType_DataPush, dispatcher.handlePush)
	dispatcher.register(proto.PacketType_DataPushReply, dispatcher.handlePushReply)
	dispatcher.register(proto.PacketType_DataPull, dispatcher.handlePull)
	return dispatcher
}

func (p *PacketSystemDispatcher) Dispatch(message proto.Packet) error {
	handler, exist := p.handlerMap[message.Type]
	if !exist {
		log.Errorf("message handler[%d] not found", message.Type)
		return nil
	} else {
		return handler(message)
	}
}

func (p *PacketSystemDispatcher) register(messageType proto.PacketType, handler PacketHandler) error {
	_, ok := p.handlerMap[messageType]
	if ok {
		return udb.ErrMessageHandlerExist
	}
	p.handlerMap[messageType] = handler
	return nil
}

func (p *PacketSystemDispatcher) handleClusterHi(packet proto.Packet) error {
	nodeInfo := &proto.NodeInfo{}
	err := proto.Unmarshal(packet.Content, nodeInfo)
	if err != nil {
		log.Warnf("handle systemHi unmarshal packet error")
	}
	log.Debugf("handleClusterHi from[%d] nodeInfo[%s]", packet.From, nodeInfo.String())
	//todo: mmap
	/**
		var buf bytes.Buffer
		enc := gob.NewEncoder(&buf)
		err := enc.Encode(request.data)
		buf.Bytes()

		buf := bytes.NewReader(value)
		dec := gob.NewDecoder(buf)
		var data types.DocData
		err := dec.Decode(&data)

		https://stackoverflow.com/questions/9203526/mapping-an-array-to-a-file-via-mmap-in-go
	 */
	//todo: local partition manager impl
	p.cluster.JoinNode(packet.From, int(nodeInfo.PartitionSize), int(nodeInfo.ReplicaSize))
	return nil
	//node, b := p.center.Cluster.Node(message.From)
	//if !b {
	//	log.Errorf("cannot update clock[%v] of node[%d]", message.Timestamp,message.From)
	//	return
	//}
	//timestamp := message.Timestamp
	//node.Timestamp.Remote = timestamp.Remote
	//node.Timestamp.Local = timestamp.Local
	//log.Debugf("update cluster node[%d] timestamp l:%d r:%d", node.GetId, node.Timestamp.Local, node.Timestamp.Remote)
}

func (p *PacketSystemDispatcher) handleMigrate(packet proto.Packet) error {
	request := proto.DataMigrateRequest{}
	proto.Unmarshal(packet.Content, &request)
	operations := store.NewDataOperations(p.cluster.engine, p.cluster.dataCommunication, packet.From)
	operations.Migrate(request.StartRing, request.EndRing)
	return nil
}
func (p *PacketSystemDispatcher) handleMigrateReply(packet proto.Packet) error {
	migrateResponse := proto.DataMigrateResponse{}
	proto.Unmarshal(packet.Content, &migrateResponse)
	log.Debugf("migrateReply[%s]", migrateResponse.String())
	return nil
}

func (p *PacketSystemDispatcher) handlePush(packet proto.Packet) error {
	pushRequest := proto.PushRequest{}
	proto.Unmarshal(packet.Content, &pushRequest)
	operations := store.NewDataOperations(p.cluster.engine, p.cluster.dataCommunication, packet.From)
	operations.Push(pushRequest.Data)
	return nil
}

func (p *PacketSystemDispatcher) handlePushReply(packet proto.Packet) error {
	pushResponse := proto.PushResponse{}
	proto.Unmarshal(packet.Content, &pushResponse)
	log.Debugf("pushReply[%s]", pushResponse.String())
	return nil
}
func (p *PacketSystemDispatcher) handlePull(packet proto.Packet) error {
	return nil
}
