/*
 * Copyright (c) 2019 uplus.io
 */

package cluster

import "uplus.io/udb/proto"

type Pipeline interface {
	InSyncWrite(packet *proto.Packet) *SyncChannel
	InWrite(packet *proto.Packet)
	InRead() <-chan *proto.Packet
	OutWrite(packet *proto.Packet)
	OutRead() <-chan *proto.Packet
}

type PipelinePacket struct {
	packetInCh  chan *proto.Packet
	packetOutCh chan *proto.Packet

	syncChannelMap map[string]*SyncChannel
}

func NewPipelinePacket() *PipelinePacket {
	return &PipelinePacket{
		packetInCh:     make(chan *proto.Packet),
		packetOutCh:    make(chan *proto.Packet),
		syncChannelMap: make(map[string]*SyncChannel),
	}
}

//func (p *PipelinePacket) In() chan *proto.Packet {
//	return p.packetInCh
//}
//
//func (p *PipelinePacket) Out() chan *proto.Packet {
//	return p.packetOutCh
//}

func (p *PipelinePacket) InSyncWrite(packet *proto.Packet) *SyncChannel {
	syncChannel := NewSyncChannel(packet.Id)
	p.syncChannelMap[packet.Id] = syncChannel
	p.InWrite(packet)
	return syncChannel
}

func (p *PipelinePacket) InWrite(packet *proto.Packet) {
	p.packetInCh <- packet
}

func (p *PipelinePacket) InRead() <-chan *proto.Packet {
	return p.packetInCh
}

func (p *PipelinePacket) OutWrite(packet *proto.Packet) {
	syncChannel, exist := p.syncChannelMap[packet.Id]
	if exist {
		syncChannel.Write(packet)
	} else {
		p.packetOutCh <- packet
	}
}

func (p *PipelinePacket) OutRead() <-chan *proto.Packet {
	return p.packetOutCh
}

type SyncChannel struct {
	PacketId string
	packetCh chan *proto.Packet
}

func NewSyncChannel(packetId string) *SyncChannel {
	return &SyncChannel{PacketId: packetId, packetCh: make(chan *proto.Packet)}
}
func (p *SyncChannel) Write(packet *proto.Packet) {
	p.packetCh <- packet
}
func (p *SyncChannel) Read() <-chan *proto.Packet {
	return p.packetCh
}
