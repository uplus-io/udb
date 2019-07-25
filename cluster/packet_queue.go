package cluster

import (
	"uplus.io/udb/proto"
)

type PacketQueue struct {
	messages map[string]*PacketMessage
}

func Put(packetId string, message PacketMessage) {
	//todo: socket sync message
}

type PacketMessage struct {
	messageCh chan *proto.Packet
}
