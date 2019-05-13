// Code generated by protoc-gen-go. DO NOT EDIT.
// source: node.proto

package types // import "github.com/aergoio/aergo/types"

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type PeerAddress struct {
	// address is string representation of ip address or domain name.
	Address              string   `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
	Port                 uint32   `protobuf:"varint,2,opt,name=port,proto3" json:"port,omitempty"`
	PeerID               []byte   `protobuf:"bytes,3,opt,name=peerID,proto3" json:"peerID,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PeerAddress) Reset()         { *m = PeerAddress{} }
func (m *PeerAddress) String() string { return proto.CompactTextString(m) }
func (*PeerAddress) ProtoMessage()    {}
func (*PeerAddress) Descriptor() ([]byte, []int) {
	return fileDescriptor_node_8abca6e41f19b114, []int{0}
}
func (m *PeerAddress) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PeerAddress.Unmarshal(m, b)
}
func (m *PeerAddress) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PeerAddress.Marshal(b, m, deterministic)
}
func (dst *PeerAddress) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PeerAddress.Merge(dst, src)
}
func (m *PeerAddress) XXX_Size() int {
	return xxx_messageInfo_PeerAddress.Size(m)
}
func (m *PeerAddress) XXX_DiscardUnknown() {
	xxx_messageInfo_PeerAddress.DiscardUnknown(m)
}

var xxx_messageInfo_PeerAddress proto.InternalMessageInfo

func (m *PeerAddress) GetAddress() string {
	if m != nil {
		return m.Address
	}
	return ""
}

func (m *PeerAddress) GetPort() uint32 {
	if m != nil {
		return m.Port
	}
	return 0
}

func (m *PeerAddress) GetPeerID() []byte {
	if m != nil {
		return m.PeerID
	}
	return nil
}

func init() {
	proto.RegisterType((*PeerAddress)(nil), "types.PeerAddress")
}

func init() { proto.RegisterFile("node.proto", fileDescriptor_node_8abca6e41f19b114) }

var fileDescriptor_node_8abca6e41f19b114 = []byte{
	// 141 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0xca, 0xcb, 0x4f, 0x49,
	0xd5, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x2d, 0xa9, 0x2c, 0x48, 0x2d, 0x56, 0x0a, 0xe6,
	0xe2, 0x0e, 0x48, 0x4d, 0x2d, 0x72, 0x4c, 0x49, 0x29, 0x4a, 0x2d, 0x2e, 0x16, 0x92, 0xe0, 0x62,
	0x4f, 0x84, 0x30, 0x25, 0x18, 0x15, 0x18, 0x35, 0x38, 0x83, 0x60, 0x5c, 0x21, 0x21, 0x2e, 0x96,
	0x82, 0xfc, 0xa2, 0x12, 0x09, 0x26, 0x05, 0x46, 0x0d, 0xde, 0x20, 0x30, 0x5b, 0x48, 0x8c, 0x8b,
	0xad, 0x20, 0x35, 0xb5, 0xc8, 0xd3, 0x45, 0x82, 0x59, 0x81, 0x51, 0x83, 0x27, 0x08, 0xca, 0x73,
	0x52, 0x88, 0x92, 0x4b, 0xcf, 0x2c, 0xc9, 0x28, 0x4d, 0xd2, 0x4b, 0xce, 0xcf, 0xd5, 0x4f, 0x4c,
	0x2d, 0x4a, 0xcf, 0xcf, 0xcc, 0x87, 0xd0, 0xfa, 0x60, 0x6b, 0x93, 0xd8, 0xc0, 0x8e, 0x30, 0x06,
	0x04, 0x00, 0x00, 0xff, 0xff, 0x92, 0x13, 0x1d, 0xcc, 0x92, 0x00, 0x00, 0x00,
}
