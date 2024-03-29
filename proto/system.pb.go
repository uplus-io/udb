// Code generated by protoc-gen-go. DO NOT EDIT.
// source: system.proto

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

type StorageMeta struct {
	Version              int32    `protobuf:"varint,1,opt,name=version,proto3" json:"version,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *StorageMeta) Reset()         { *m = StorageMeta{} }
func (m *StorageMeta) String() string { return proto.CompactTextString(m) }
func (*StorageMeta) ProtoMessage()    {}
func (*StorageMeta) Descriptor() ([]byte, []int) {
	return fileDescriptor_86a7260ebdc12f47, []int{0}
}

func (m *StorageMeta) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StorageMeta.Unmarshal(m, b)
}
func (m *StorageMeta) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StorageMeta.Marshal(b, m, deterministic)
}
func (m *StorageMeta) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StorageMeta.Merge(m, src)
}
func (m *StorageMeta) XXX_Size() int {
	return xxx_messageInfo_StorageMeta.Size(m)
}
func (m *StorageMeta) XXX_DiscardUnknown() {
	xxx_messageInfo_StorageMeta.DiscardUnknown(m)
}

var xxx_messageInfo_StorageMeta proto.InternalMessageInfo

func (m *StorageMeta) GetVersion() int32 {
	if m != nil {
		return m.Version
	}
	return 0
}

type PartitionMeta struct {
	Repo                 *Repository  `protobuf:"bytes,1,opt,name=repo,proto3" json:"repo,omitempty"`
	Node                 int32        `protobuf:"varint,2,opt,name=node,proto3" json:"node,omitempty"`
	Parts                []*Partition `protobuf:"bytes,3,rep,name=parts,proto3" json:"parts,omitempty"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *PartitionMeta) Reset()         { *m = PartitionMeta{} }
func (m *PartitionMeta) String() string { return proto.CompactTextString(m) }
func (*PartitionMeta) ProtoMessage()    {}
func (*PartitionMeta) Descriptor() ([]byte, []int) {
	return fileDescriptor_86a7260ebdc12f47, []int{1}
}

func (m *PartitionMeta) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PartitionMeta.Unmarshal(m, b)
}
func (m *PartitionMeta) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PartitionMeta.Marshal(b, m, deterministic)
}
func (m *PartitionMeta) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PartitionMeta.Merge(m, src)
}
func (m *PartitionMeta) XXX_Size() int {
	return xxx_messageInfo_PartitionMeta.Size(m)
}
func (m *PartitionMeta) XXX_DiscardUnknown() {
	xxx_messageInfo_PartitionMeta.DiscardUnknown(m)
}

var xxx_messageInfo_PartitionMeta proto.InternalMessageInfo

func (m *PartitionMeta) GetRepo() *Repository {
	if m != nil {
		return m.Repo
	}
	return nil
}

func (m *PartitionMeta) GetNode() int32 {
	if m != nil {
		return m.Node
	}
	return 0
}

func (m *PartitionMeta) GetParts() []*Partition {
	if m != nil {
		return m.Parts
	}
	return nil
}

func init() {
	proto.RegisterType((*StorageMeta)(nil), "proto.StorageMeta")
	proto.RegisterType((*PartitionMeta)(nil), "proto.PartitionMeta")
}

func init() { proto.RegisterFile("system.proto", fileDescriptor_86a7260ebdc12f47) }

var fileDescriptor_86a7260ebdc12f47 = []byte{
	// 169 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x3c, 0x8c, 0xc1, 0x8a, 0xc2, 0x40,
	0x0c, 0x86, 0xe9, 0xb6, 0xdd, 0x85, 0xb4, 0x0b, 0xbb, 0x73, 0x1a, 0x3c, 0x95, 0x82, 0xda, 0x53,
	0x0f, 0xf5, 0x39, 0x04, 0x19, 0x9f, 0x60, 0xd4, 0x20, 0x73, 0xe8, 0x64, 0xc8, 0x04, 0xa1, 0x6f,
	0x2f, 0xa4, 0xea, 0x29, 0xc9, 0xff, 0x7f, 0xf9, 0xa0, 0xcd, 0x4b, 0x16, 0x9c, 0xc7, 0xc4, 0x24,
	0x64, 0x6a, 0x1d, 0x9b, 0xf6, 0x4a, 0xf3, 0x4c, 0x71, 0x0d, 0xfb, 0x3d, 0x34, 0x67, 0x21, 0xf6,
	0x77, 0x3c, 0xa2, 0x78, 0x63, 0xe1, 0xe7, 0x81, 0x9c, 0x03, 0x45, 0x5b, 0x74, 0xc5, 0x50, 0xbb,
	0xf7, 0xd9, 0x33, 0xfc, 0x9e, 0x3c, 0x4b, 0x90, 0x40, 0x51, 0xd1, 0x2d, 0x54, 0x8c, 0x89, 0x94,
	0x6b, 0xa6, 0xff, 0xd5, 0x37, 0x3a, 0x4c, 0x94, 0x83, 0x10, 0x2f, 0x4e, 0x6b, 0x63, 0xa0, 0x8a,
	0x74, 0x43, 0xfb, 0xa5, 0x3a, 0xdd, 0xcd, 0x0e, 0xea, 0xe4, 0x59, 0xb2, 0x2d, 0xbb, 0x72, 0x68,
	0xa6, 0xbf, 0xd7, 0xef, 0xc7, 0xef, 0xd6, 0xfa, 0xf2, 0xad, 0xf9, 0xe1, 0x19, 0x00, 0x00, 0xff,
	0xff, 0x0e, 0xf5, 0x81, 0x52, 0xc8, 0x00, 0x00, 0x00,
}
