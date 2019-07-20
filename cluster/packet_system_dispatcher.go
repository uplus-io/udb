/*
 * Copyright (c) 2019 uplus.io
 */

package cluster

import (
	"uplus.io/udb"
	log "uplus.io/udb/logger"
	"uplus.io/udb/proto"
)

type PacketSystemDispatcher struct {
	cluster *Cluster

	handlerMap map[proto.PacketType]PacketHandler
}

func NewPacketSystemDispatcher(cluster *Cluster) *PacketSystemDispatcher {
	dispatcher := &PacketSystemDispatcher{cluster: cluster, handlerMap: make(map[proto.PacketType]PacketHandler)}
	dispatcher.register(proto.PacketType_CLOCK, dispatcher.handleClock)
	return dispatcher
}

func (p *PacketSystemDispatcher) Dispatch(message proto.Packet) {
	handler, exist := p.handlerMap[message.Type]
	if !exist {
		log.Warnf("message handler[%d] not found", message.Type)
	} else {
		handler(message)
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

func (p *PacketSystemDispatcher) handleClock(message proto.Packet) () {
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
