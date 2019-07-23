/*
 * Copyright (c) 2019 uplus.io
 */

package cluster

import "uplus.io/udb/proto"

type PacketDispatcher interface {
	Dispatch(packet proto.Packet) error
	register(packetType proto.PacketType, handler PacketHandler) error
}

type PacketHandler func(packet proto.Packet) error
