// Code generated by protoc-gen-go. DO NOT EDIT.
// source: proto/api.proto

/*
Package rpc is a generated protocol buffer package.

It is generated from these files:
	proto/api.proto

It has these top-level messages:
	Msg
*/
package rpc

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

type Msg struct {
	Data string `protobuf:"bytes,1,opt,name=Data" json:"Data,omitempty"`
}

func (m *Msg) Reset()                    { *m = Msg{} }
func (m *Msg) String() string            { return proto.CompactTextString(m) }
func (*Msg) ProtoMessage()               {}
func (*Msg) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Msg) GetData() string {
	if m != nil {
		return m.Data
	}
	return ""
}

func init() {
	proto.RegisterType((*Msg)(nil), "protoFile.Msg")
}

func init() { proto.RegisterFile("proto/api.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 97 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x2f, 0x28, 0xca, 0x2f,
	0xc9, 0xd7, 0x4f, 0x2c, 0xc8, 0xd4, 0x03, 0xb3, 0x84, 0x38, 0xc1, 0x94, 0x5b, 0x66, 0x4e, 0xaa,
	0x92, 0x24, 0x17, 0xb3, 0x6f, 0x71, 0xba, 0x90, 0x10, 0x17, 0x8b, 0x4b, 0x62, 0x49, 0xa2, 0x04,
	0xa3, 0x02, 0xa3, 0x06, 0x67, 0x10, 0x98, 0xed, 0x24, 0x1a, 0x25, 0xac, 0xa7, 0x9f, 0x5b, 0x59,
	0x9c, 0x5a, 0x54, 0x96, 0x5a, 0xa4, 0x9f, 0x9c, 0x5f, 0x94, 0xaa, 0x5f, 0x54, 0x90, 0x9c, 0xc4,
	0x06, 0xd6, 0x6c, 0x0c, 0x08, 0x00, 0x00, 0xff, 0xff, 0xcc, 0x95, 0xec, 0x7e, 0x56, 0x00, 0x00,
	0x00,
}
