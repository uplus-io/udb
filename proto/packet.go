package proto

import "uplus.io/udb/utils"

func NewPacket(typ PacketType, from, to int32, content []byte) *Packet {
	return &Packet{
		Version: PROTO_VERSION,
		Id:      utils.GenId(),
		Type:    typ,
		From:    from,
		To:      to, Content: content,
	}
}
