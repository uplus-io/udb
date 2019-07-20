/*
 * Copyright (c) 2019 uplus.io
 */

package cluster

import "uplus.io/udb/proto"

type Pipeline interface {
	In() chan *proto.Packet

	Out() chan *proto.Packet
}

type PipelinePacket struct {
	packetInCh  chan *proto.Packet
	packetOutCh chan *proto.Packet
}

func NewPipelinePacket() *PipelinePacket {
	return &PipelinePacket{packetInCh: make(chan *proto.Packet), packetOutCh: make(chan *proto.Packet)}
}

func (p *PipelinePacket) In() chan *proto.Packet {
	return p.packetInCh
}

func (p *PipelinePacket) Out() chan *proto.Packet {
	return p.packetOutCh
}

