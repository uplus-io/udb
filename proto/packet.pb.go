// Code generated by protoc-gen-go. DO NOT EDIT.
// source: packet.proto

package proto

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type PacketCategory int32

const (
	PacketCategory_System PacketCategory = 0
	PacketCategory_Data   PacketCategory = 2
)

var PacketCategory_name = map[int32]string{
	0: "System",
	2: "Data",
}

var PacketCategory_value = map[string]int32{
	"System": 0,
	"Data":   2,
}

func (x PacketCategory) String() string {
	return proto.EnumName(PacketCategory_name, int32(x))
}

func (PacketCategory) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_e9ef1a6541f9f9e7, []int{0}
}

type PacketType int32

const (
	PacketType_SystemHi     PacketType = 0
	PacketType_DataPush     PacketType = 100
	PacketType_DataPull     PacketType = 101
	PacketType_DataPushPull PacketType = 102
)

var PacketType_name = map[int32]string{
	0:   "SystemHi",
	100: "DataPush",
	101: "DataPull",
	102: "DataPushPull",
}

var PacketType_value = map[string]int32{
	"SystemHi":     0,
	"DataPush":     100,
	"DataPull":     101,
	"DataPushPull": 102,
}

func (x PacketType) String() string {
	return proto.EnumName(PacketType_name, int32(x))
}

func (PacketType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_e9ef1a6541f9f9e7, []int{1}
}

type Packet struct {
	Version              int32          `protobuf:"varint,1,opt,name=version,proto3" json:"version,omitempty"`
	Id                   string         `protobuf:"bytes,2,opt,name=id,proto3" json:"id,omitempty"`
	Category             PacketCategory `protobuf:"varint,3,opt,name=category,proto3,enum=proto.PacketCategory" json:"category,omitempty"`
	Type                 PacketType     `protobuf:"varint,4,opt,name=type,proto3,enum=proto.PacketType" json:"type,omitempty"`
	DataCenter           int32          `protobuf:"varint,5,opt,name=dataCenter,proto3" json:"dataCenter,omitempty"`
	From                 int32          `protobuf:"varint,6,opt,name=from,proto3" json:"from,omitempty"`
	To                   int32          `protobuf:"varint,7,opt,name=to,proto3" json:"to,omitempty"`
	Content              []byte         `protobuf:"bytes,8,opt,name=content,proto3" json:"content,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *Packet) Reset()         { *m = Packet{} }
func (m *Packet) String() string { return proto.CompactTextString(m) }
func (*Packet) ProtoMessage()    {}
func (*Packet) Descriptor() ([]byte, []int) {
	return fileDescriptor_e9ef1a6541f9f9e7, []int{0}
}

func (m *Packet) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Packet.Unmarshal(m, b)
}
func (m *Packet) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Packet.Marshal(b, m, deterministic)
}
func (m *Packet) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Packet.Merge(m, src)
}
func (m *Packet) XXX_Size() int {
	return xxx_messageInfo_Packet.Size(m)
}
func (m *Packet) XXX_DiscardUnknown() {
	xxx_messageInfo_Packet.DiscardUnknown(m)
}

var xxx_messageInfo_Packet proto.InternalMessageInfo

func (m *Packet) GetVersion() int32 {
	if m != nil {
		return m.Version
	}
	return 0
}

func (m *Packet) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Packet) GetCategory() PacketCategory {
	if m != nil {
		return m.Category
	}
	return PacketCategory_System
}

func (m *Packet) GetType() PacketType {
	if m != nil {
		return m.Type
	}
	return PacketType_SystemHi
}

func (m *Packet) GetDataCenter() int32 {
	if m != nil {
		return m.DataCenter
	}
	return 0
}

func (m *Packet) GetFrom() int32 {
	if m != nil {
		return m.From
	}
	return 0
}

func (m *Packet) GetTo() int32 {
	if m != nil {
		return m.To
	}
	return 0
}

func (m *Packet) GetContent() []byte {
	if m != nil {
		return m.Content
	}
	return nil
}

func init() {
	proto.RegisterEnum("proto.PacketCategory", PacketCategory_name, PacketCategory_value)
	proto.RegisterEnum("proto.PacketType", PacketType_name, PacketType_value)
	proto.RegisterType((*Packet)(nil), "proto.Packet")
}

func init() { proto.RegisterFile("packet.proto", fileDescriptor_e9ef1a6541f9f9e7) }

var fileDescriptor_e9ef1a6541f9f9e7 = []byte{
	// 260 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x54, 0x8f, 0xd1, 0x4a, 0xf3, 0x40,
	0x10, 0x85, 0xbb, 0xf9, 0x93, 0x34, 0xff, 0x10, 0xc2, 0x3a, 0x20, 0xec, 0x95, 0x04, 0x41, 0x09,
	0xbd, 0x28, 0xa8, 0x8f, 0x50, 0x2f, 0x7a, 0x59, 0xa2, 0x2f, 0xb0, 0x26, 0x53, 0x0d, 0xb6, 0xd9,
	0xb0, 0x19, 0x85, 0xbc, 0xb2, 0x4f, 0x21, 0x99, 0x18, 0x6b, 0xaf, 0x76, 0xcf, 0x39, 0x1f, 0xc3,
	0x39, 0x90, 0x76, 0xb6, 0x7a, 0x27, 0x5e, 0x77, 0xde, 0xb1, 0xc3, 0x48, 0x9e, 0xeb, 0x2f, 0x05,
	0xf1, 0x4e, 0x7c, 0x34, 0xb0, 0xfc, 0x24, 0xdf, 0x37, 0xae, 0x35, 0x2a, 0x57, 0x45, 0x54, 0xce,
	0x12, 0x33, 0x08, 0x9a, 0xda, 0x04, 0xb9, 0x2a, 0xfe, 0x97, 0x41, 0x53, 0xe3, 0x1d, 0x24, 0x95,
	0x65, 0x7a, 0x75, 0x7e, 0x30, 0xff, 0x72, 0x55, 0x64, 0xf7, 0x97, 0xd3, 0xd5, 0xf5, 0x74, 0x6a,
	0xf3, 0x13, 0x96, 0xbf, 0x18, 0xde, 0x40, 0xc8, 0x43, 0x47, 0x26, 0x14, 0xfc, 0xe2, 0x0c, 0x7f,
	0x1e, 0x3a, 0x2a, 0x25, 0xc6, 0x2b, 0x80, 0xda, 0xb2, 0xdd, 0x50, 0xcb, 0xe4, 0x4d, 0x24, 0x35,
	0xfe, 0x38, 0x88, 0x10, 0xee, 0xbd, 0x3b, 0x9a, 0x58, 0x12, 0xf9, 0x8f, 0xed, 0xd8, 0x99, 0xa5,
	0x38, 0x01, 0xbb, 0x71, 0x47, 0xe5, 0x5a, 0xa6, 0x96, 0x4d, 0x92, 0xab, 0x22, 0x2d, 0x67, 0xb9,
	0xba, 0x85, 0xec, 0xbc, 0x20, 0x02, 0xc4, 0x4f, 0x43, 0xcf, 0x74, 0xd4, 0x0b, 0x4c, 0x20, 0x7c,
	0xb4, 0x6c, 0x75, 0xb0, 0xda, 0x02, 0x9c, 0x9a, 0x61, 0x0a, 0xc9, 0xc4, 0x6c, 0x1b, 0xbd, 0x18,
	0xd5, 0x48, 0xed, 0x3e, 0xfa, 0x37, 0x5d, 0x9f, 0xd4, 0xe1, 0xa0, 0x09, 0x35, 0xa4, 0x73, 0x26,
	0xce, 0xfe, 0x25, 0x96, 0x9d, 0x0f, 0xdf, 0x01, 0x00, 0x00, 0xff, 0xff, 0x21, 0x3d, 0x5c, 0x1a,
	0x7c, 0x01, 0x00, 0x00,
}