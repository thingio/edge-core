// Code generated by protoc-gen-go. DO NOT EDIT.
// source: define.proto

package talk

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

type TMessage struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Sender               string   `protobuf:"bytes,2,opt,name=sender,proto3" json:"sender,omitempty"`
	Method               string   `protobuf:"bytes,3,opt,name=method,proto3" json:"method,omitempty"`
	Key                  string   `protobuf:"bytes,4,opt,name=key,proto3" json:"key,omitempty"`
	Payload              []byte   `protobuf:"bytes,5,opt,name=payload,proto3" json:"payload,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TMessage) Reset()         { *m = TMessage{} }
func (m *TMessage) String() string { return proto.CompactTextString(m) }
func (*TMessage) ProtoMessage()    {}
func (*TMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_f7e38f9743a0f071, []int{0}
}

func (m *TMessage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TMessage.Unmarshal(m, b)
}
func (m *TMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TMessage.Marshal(b, m, deterministic)
}
func (m *TMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TMessage.Merge(m, src)
}
func (m *TMessage) XXX_Size() int {
	return xxx_messageInfo_TMessage.Size(m)
}
func (m *TMessage) XXX_DiscardUnknown() {
	xxx_messageInfo_TMessage.DiscardUnknown(m)
}

var xxx_messageInfo_TMessage proto.InternalMessageInfo

func (m *TMessage) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *TMessage) GetSender() string {
	if m != nil {
		return m.Sender
	}
	return ""
}

func (m *TMessage) GetMethod() string {
	if m != nil {
		return m.Method
	}
	return ""
}

func (m *TMessage) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

func (m *TMessage) GetPayload() []byte {
	if m != nil {
		return m.Payload
	}
	return nil
}

func init() {
	proto.RegisterType((*TMessage)(nil), "talk.TMessage")
}

func init() { proto.RegisterFile("define.proto", fileDescriptor_f7e38f9743a0f071) }

var fileDescriptor_f7e38f9743a0f071 = []byte{
	// 139 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x49, 0x49, 0x4d, 0xcb,
	0xcc, 0x4b, 0xd5, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x29, 0x49, 0xcc, 0xc9, 0x56, 0x2a,
	0xe3, 0xe2, 0x08, 0xf1, 0x4d, 0x2d, 0x2e, 0x4e, 0x4c, 0x4f, 0x15, 0xe2, 0xe3, 0x62, 0xca, 0x4c,
	0x91, 0x60, 0x54, 0x60, 0xd4, 0xe0, 0x0c, 0x62, 0xca, 0x4c, 0x11, 0x12, 0xe3, 0x62, 0x2b, 0x4e,
	0xcd, 0x4b, 0x49, 0x2d, 0x92, 0x60, 0x02, 0x8b, 0x41, 0x79, 0x20, 0xf1, 0xdc, 0xd4, 0x92, 0x8c,
	0xfc, 0x14, 0x09, 0x66, 0x88, 0x38, 0x84, 0x27, 0x24, 0xc0, 0xc5, 0x9c, 0x9d, 0x5a, 0x29, 0xc1,
	0x02, 0x16, 0x04, 0x31, 0x85, 0x24, 0xb8, 0xd8, 0x0b, 0x12, 0x2b, 0x73, 0xf2, 0x13, 0x53, 0x24,
	0x58, 0x15, 0x18, 0x35, 0x78, 0x82, 0x60, 0xdc, 0x24, 0x36, 0xb0, 0x23, 0x8c, 0x01, 0x01, 0x00,
	0x00, 0xff, 0xff, 0xd0, 0x5c, 0xfc, 0xc6, 0x94, 0x00, 0x00, 0x00,
}
