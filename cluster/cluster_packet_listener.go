/*
 * Copyright (c) 2019 uplus.io
 */

package cluster

import "uplus.io/udb/proto"

type PacketListener interface {
	OnReceive(packet *proto.Packet)
}

type ClusterPacketListener struct {
	Pipeline Pipeline
}

func NewClusterPacketListener(pipeline Pipeline) *ClusterPacketListener {
	return &ClusterPacketListener{Pipeline: pipeline}
}

func (p *ClusterPacketListener) OnReceive(packet *proto.Packet) {
	p.Pipeline.OutWrite(packet)
}
