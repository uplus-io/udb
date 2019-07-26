package proto

import "uplus.io/udb/utils"

func NewPacket(mode PacketMode, typ PacketType, from int32, to []int32, content []byte) *Packet {
	packet := &Packet{
		Version: PROTO_VERSION,
		Id:      utils.GenId(),
		Mode:    mode,
		Type:    typ,
		From:    from,
		Content: content,
	}
	if mode == PacketMode_Multicast {
		packet.Receivers = to
	} else {
		packet.To = to[0]
	}
	return packet
}

func NewTCPPacket(typ PacketType, from, to int32, content []byte) *Packet {
	return NewPacket(PacketMode_TCP, typ, from, []int32{to}, content)
}

func NewUDPPacket(typ PacketType, from, to int32, content []byte) *Packet {
	return NewPacket(PacketMode_UDP, typ, from, []int32{to}, content)
}
