// Code generated by protoc-gen-go. DO NOT EDIT.
// source: fetch.proto

package bean

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

type NeighbourReq struct {
	Req                  string   `protobuf:"bytes,1,opt,name=req,proto3" json:"req,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *NeighbourReq) Reset()         { *m = NeighbourReq{} }
func (m *NeighbourReq) String() string { return proto.CompactTextString(m) }
func (*NeighbourReq) ProtoMessage()    {}
func (*NeighbourReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_fetch_fae93ed012df61ac, []int{0}
}
func (m *NeighbourReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NeighbourReq.Unmarshal(m, b)
}
func (m *NeighbourReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NeighbourReq.Marshal(b, m, deterministic)
}
func (dst *NeighbourReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NeighbourReq.Merge(dst, src)
}
func (m *NeighbourReq) XXX_Size() int {
	return xxx_messageInfo_NeighbourReq.Size(m)
}
func (m *NeighbourReq) XXX_DiscardUnknown() {
	xxx_messageInfo_NeighbourReq.DiscardUnknown(m)
}

var xxx_messageInfo_NeighbourReq proto.InternalMessageInfo

func (m *NeighbourReq) GetReq() string {
	if m != nil {
		return m.Req
	}
	return ""
}

type BlockReq struct {
	Req                  string   `protobuf:"bytes,1,opt,name=req,proto3" json:"req,omitempty"`
	MinHeight            int64    `protobuf:"varint,2,opt,name=min_height,json=minHeight,proto3" json:"min_height,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *BlockReq) Reset()         { *m = BlockReq{} }
func (m *BlockReq) String() string { return proto.CompactTextString(m) }
func (*BlockReq) ProtoMessage()    {}
func (*BlockReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_fetch_fae93ed012df61ac, []int{1}
}
func (m *BlockReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_BlockReq.Unmarshal(m, b)
}
func (m *BlockReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_BlockReq.Marshal(b, m, deterministic)
}
func (dst *BlockReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BlockReq.Merge(dst, src)
}
func (m *BlockReq) XXX_Size() int {
	return xxx_messageInfo_BlockReq.Size(m)
}
func (m *BlockReq) XXX_DiscardUnknown() {
	xxx_messageInfo_BlockReq.DiscardUnknown(m)
}

var xxx_messageInfo_BlockReq proto.InternalMessageInfo

func (m *BlockReq) GetReq() string {
	if m != nil {
		return m.Req
	}
	return ""
}

func (m *BlockReq) GetMinHeight() int64 {
	if m != nil {
		return m.MinHeight
	}
	return 0
}

type BlockResp struct {
	Resp                 string   `protobuf:"bytes,1,opt,name=resp,proto3" json:"resp,omitempty"`
	MaxHeight            int64    `protobuf:"varint,2,opt,name=max_height,json=maxHeight,proto3" json:"max_height,omitempty"`
	NewBlock             *Block   `protobuf:"bytes,3,opt,name=new_block,json=newBlock,proto3" json:"new_block,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *BlockResp) Reset()         { *m = BlockResp{} }
func (m *BlockResp) String() string { return proto.CompactTextString(m) }
func (*BlockResp) ProtoMessage()    {}
func (*BlockResp) Descriptor() ([]byte, []int) {
	return fileDescriptor_fetch_fae93ed012df61ac, []int{2}
}
func (m *BlockResp) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_BlockResp.Unmarshal(m, b)
}
func (m *BlockResp) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_BlockResp.Marshal(b, m, deterministic)
}
func (dst *BlockResp) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BlockResp.Merge(dst, src)
}
func (m *BlockResp) XXX_Size() int {
	return xxx_messageInfo_BlockResp.Size(m)
}
func (m *BlockResp) XXX_DiscardUnknown() {
	xxx_messageInfo_BlockResp.DiscardUnknown(m)
}

var xxx_messageInfo_BlockResp proto.InternalMessageInfo

func (m *BlockResp) GetResp() string {
	if m != nil {
		return m.Resp
	}
	return ""
}

func (m *BlockResp) GetMaxHeight() int64 {
	if m != nil {
		return m.MaxHeight
	}
	return 0
}

func (m *BlockResp) GetNewBlock() *Block {
	if m != nil {
		return m.NewBlock
	}
	return nil
}

func init() {
	proto.RegisterType((*NeighbourReq)(nil), "bean.neighbour_req")
	proto.RegisterType((*BlockReq)(nil), "bean.block_req")
	proto.RegisterType((*BlockResp)(nil), "bean.block_resp")
}

func init() { proto.RegisterFile("fetch.proto", fileDescriptor_fetch_fae93ed012df61ac) }

var fileDescriptor_fetch_fae93ed012df61ac = []byte{
	// 186 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x8f, 0xb1, 0xca, 0xc2, 0x30,
	0x14, 0x85, 0xc9, 0xdf, 0xf2, 0x63, 0x6e, 0x11, 0x24, 0x53, 0x11, 0x84, 0xda, 0x29, 0x53, 0x07,
	0x5d, 0x9d, 0x9c, 0x9c, 0xfb, 0x02, 0x25, 0xa9, 0x57, 0x13, 0xb4, 0x49, 0x1b, 0x23, 0xed, 0xe3,
	0x4b, 0x52, 0x5d, 0xc4, 0xe9, 0x1e, 0xbe, 0xcb, 0xf9, 0xe0, 0x40, 0x76, 0x41, 0xdf, 0xaa, 0xaa,
	0x77, 0xd6, 0x5b, 0x96, 0x4a, 0x14, 0x66, 0x0d, 0x67, 0xe1, 0xc5, 0x4c, 0xca, 0x2d, 0x2c, 0x0d,
	0xea, 0xab, 0x92, 0xf6, 0xe9, 0x1a, 0x87, 0x03, 0x5b, 0x41, 0xe2, 0x70, 0xc8, 0x49, 0x41, 0x38,
	0xad, 0x43, 0x2c, 0x0f, 0x40, 0xe5, 0xdd, 0xb6, 0xb7, 0xdf, 0x6f, 0xb6, 0x01, 0xe8, 0xb4, 0x69,
	0x54, 0xb0, 0xf8, 0xfc, 0xaf, 0x20, 0x3c, 0xa9, 0x69, 0xa7, 0xcd, 0x29, 0x82, 0x52, 0x03, 0x7c,
	0xda, 0x8f, 0x9e, 0x31, 0x48, 0xc3, 0x7d, 0xf7, 0x63, 0x8e, 0x02, 0x31, 0x7d, 0x0b, 0xc4, 0x34,
	0x0b, 0x18, 0x07, 0x6a, 0x70, 0x6c, 0xa2, 0x24, 0x4f, 0x0a, 0xc2, 0xb3, 0x5d, 0x56, 0x85, 0x1d,
	0x55, 0x44, 0xf5, 0xc2, 0xe0, 0x78, 0x0c, 0x49, 0xfe, 0xc7, 0x49, 0xfb, 0x57, 0x00, 0x00, 0x00,
	0xff, 0xff, 0x87, 0x91, 0x7b, 0x21, 0xf3, 0x00, 0x00, 0x00,
}
