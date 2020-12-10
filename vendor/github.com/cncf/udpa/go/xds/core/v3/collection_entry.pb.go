// Code generated by protoc-gen-go. DO NOT EDIT.
// source: xds/core/v3/collection_entry.proto

package xds_core_v3

import (
	fmt "fmt"
	_ "github.com/cncf/udpa/go/udpa/annotations"
	_ "github.com/envoyproxy/protoc-gen-validate/validate"
	proto "github.com/golang/protobuf/proto"
	any "github.com/golang/protobuf/ptypes/any"
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

type CollectionEntry struct {
	// Types that are valid to be assigned to ResourceSpecifier:
	//	*CollectionEntry_Locator
	//	*CollectionEntry_InlineEntry_
	ResourceSpecifier    isCollectionEntry_ResourceSpecifier `protobuf_oneof:"resource_specifier"`
	XXX_NoUnkeyedLiteral struct{}                            `json:"-"`
	XXX_unrecognized     []byte                              `json:"-"`
	XXX_sizecache        int32                               `json:"-"`
}

func (m *CollectionEntry) Reset()         { *m = CollectionEntry{} }
func (m *CollectionEntry) String() string { return proto.CompactTextString(m) }
func (*CollectionEntry) ProtoMessage()    {}
func (*CollectionEntry) Descriptor() ([]byte, []int) {
	return fileDescriptor_5b15c821e5994c90, []int{0}
}

func (m *CollectionEntry) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CollectionEntry.Unmarshal(m, b)
}
func (m *CollectionEntry) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CollectionEntry.Marshal(b, m, deterministic)
}
func (m *CollectionEntry) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CollectionEntry.Merge(m, src)
}
func (m *CollectionEntry) XXX_Size() int {
	return xxx_messageInfo_CollectionEntry.Size(m)
}
func (m *CollectionEntry) XXX_DiscardUnknown() {
	xxx_messageInfo_CollectionEntry.DiscardUnknown(m)
}

var xxx_messageInfo_CollectionEntry proto.InternalMessageInfo

type isCollectionEntry_ResourceSpecifier interface {
	isCollectionEntry_ResourceSpecifier()
}

type CollectionEntry_Locator struct {
	Locator *ResourceLocator `protobuf:"bytes,1,opt,name=locator,proto3,oneof"`
}

type CollectionEntry_InlineEntry_ struct {
	InlineEntry *CollectionEntry_InlineEntry `protobuf:"bytes,2,opt,name=inline_entry,json=inlineEntry,proto3,oneof"`
}

func (*CollectionEntry_Locator) isCollectionEntry_ResourceSpecifier() {}

func (*CollectionEntry_InlineEntry_) isCollectionEntry_ResourceSpecifier() {}

func (m *CollectionEntry) GetResourceSpecifier() isCollectionEntry_ResourceSpecifier {
	if m != nil {
		return m.ResourceSpecifier
	}
	return nil
}

func (m *CollectionEntry) GetLocator() *ResourceLocator {
	if x, ok := m.GetResourceSpecifier().(*CollectionEntry_Locator); ok {
		return x.Locator
	}
	return nil
}

func (m *CollectionEntry) GetInlineEntry() *CollectionEntry_InlineEntry {
	if x, ok := m.GetResourceSpecifier().(*CollectionEntry_InlineEntry_); ok {
		return x.InlineEntry
	}
	return nil
}

// XXX_OneofWrappers is for the internal use of the proto package.
func (*CollectionEntry) XXX_OneofWrappers() []interface{} {
	return []interface{}{
		(*CollectionEntry_Locator)(nil),
		(*CollectionEntry_InlineEntry_)(nil),
	}
}

type CollectionEntry_InlineEntry struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Version              string   `protobuf:"bytes,2,opt,name=version,proto3" json:"version,omitempty"`
	Resource             *any.Any `protobuf:"bytes,3,opt,name=resource,proto3" json:"resource,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CollectionEntry_InlineEntry) Reset()         { *m = CollectionEntry_InlineEntry{} }
func (m *CollectionEntry_InlineEntry) String() string { return proto.CompactTextString(m) }
func (*CollectionEntry_InlineEntry) ProtoMessage()    {}
func (*CollectionEntry_InlineEntry) Descriptor() ([]byte, []int) {
	return fileDescriptor_5b15c821e5994c90, []int{0, 0}
}

func (m *CollectionEntry_InlineEntry) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CollectionEntry_InlineEntry.Unmarshal(m, b)
}
func (m *CollectionEntry_InlineEntry) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CollectionEntry_InlineEntry.Marshal(b, m, deterministic)
}
func (m *CollectionEntry_InlineEntry) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CollectionEntry_InlineEntry.Merge(m, src)
}
func (m *CollectionEntry_InlineEntry) XXX_Size() int {
	return xxx_messageInfo_CollectionEntry_InlineEntry.Size(m)
}
func (m *CollectionEntry_InlineEntry) XXX_DiscardUnknown() {
	xxx_messageInfo_CollectionEntry_InlineEntry.DiscardUnknown(m)
}

var xxx_messageInfo_CollectionEntry_InlineEntry proto.InternalMessageInfo

func (m *CollectionEntry_InlineEntry) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *CollectionEntry_InlineEntry) GetVersion() string {
	if m != nil {
		return m.Version
	}
	return ""
}

func (m *CollectionEntry_InlineEntry) GetResource() *any.Any {
	if m != nil {
		return m.Resource
	}
	return nil
}

func init() {
	proto.RegisterType((*CollectionEntry)(nil), "xds.core.v3.CollectionEntry")
	proto.RegisterType((*CollectionEntry_InlineEntry)(nil), "xds.core.v3.CollectionEntry.InlineEntry")
}

func init() { proto.RegisterFile("xds/core/v3/collection_entry.proto", fileDescriptor_5b15c821e5994c90) }

var fileDescriptor_5b15c821e5994c90 = []byte{
	// 375 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x5c, 0x91, 0xcf, 0x6b, 0xdb, 0x30,
	0x1c, 0xc5, 0xe3, 0x64, 0xe4, 0x87, 0x3c, 0x18, 0x88, 0x8c, 0x38, 0x59, 0x06, 0x23, 0xec, 0x10,
	0x18, 0x96, 0x42, 0x72, 0xd9, 0x06, 0x3b, 0xc4, 0x63, 0x90, 0xc1, 0x06, 0xc1, 0xc7, 0x36, 0x6d,
	0x50, 0x6c, 0x25, 0x15, 0x38, 0x52, 0x90, 0x64, 0x93, 0xf4, 0x50, 0x7a, 0xef, 0x7f, 0xd4, 0x6b,
	0x2f, 0xbd, 0xf6, 0xbf, 0x29, 0x3d, 0x15, 0xcb, 0x76, 0x9a, 0xe4, 0x26, 0xf1, 0x3e, 0x4f, 0xef,
	0x7d, 0xf5, 0x05, 0xbd, 0x6d, 0xa8, 0x70, 0x20, 0x24, 0xc5, 0xc9, 0x08, 0x07, 0x22, 0x8a, 0x68,
	0xa0, 0x99, 0xe0, 0x73, 0xca, 0xb5, 0xdc, 0xa1, 0x8d, 0x14, 0x5a, 0x40, 0x7b, 0x1b, 0x2a, 0x94,
	0x32, 0x28, 0x19, 0x75, 0xda, 0x2b, 0x21, 0x56, 0x11, 0xc5, 0x46, 0x5a, 0xc4, 0x4b, 0x4c, 0x78,
	0xce, 0x75, 0x3e, 0xc7, 0xe1, 0x86, 0x60, 0xc2, 0xb9, 0xd0, 0x24, 0x7d, 0x44, 0x61, 0xa5, 0x89,
	0x8e, 0x55, 0x2e, 0x1f, 0x45, 0x49, 0xaa, 0x44, 0x2c, 0x03, 0x3a, 0x8f, 0x44, 0x40, 0xb4, 0x90,
	0x39, 0xd3, 0x4a, 0x48, 0xc4, 0x42, 0xa2, 0x29, 0x2e, 0x0e, 0x99, 0xd0, 0x7b, 0x28, 0x83, 0x0f,
	0xbf, 0xf7, 0xf5, 0xfe, 0xa4, 0xed, 0xe0, 0x77, 0x50, 0xcb, 0xdd, 0x8e, 0xf5, 0xc5, 0xea, 0xdb,
	0xc3, 0x2e, 0x3a, 0x68, 0x8a, 0xfc, 0x3c, 0xe2, 0x5f, 0xc6, 0x4c, 0x4a, 0x7e, 0x81, 0xc3, 0xff,
	0xe0, 0x3d, 0xe3, 0x11, 0xe3, 0x34, 0x9b, 0xd3, 0x29, 0x1b, 0x7b, 0xff, 0xc8, 0x7e, 0x92, 0x86,
	0xfe, 0x1a, 0x83, 0x39, 0x4f, 0x4a, 0xbe, 0xcd, 0xde, 0xae, 0x9d, 0x3b, 0x0b, 0xd8, 0x07, 0x32,
	0x1c, 0x80, 0x77, 0x9c, 0xac, 0xa9, 0x69, 0xd5, 0xf0, 0xba, 0x2f, 0x5e, 0x5b, 0xb6, 0x86, 0x1f,
	0x2f, 0xcf, 0x07, 0xee, 0x0f, 0xe2, 0x5e, 0x8f, 0xdd, 0xb3, 0xf9, 0xcc, 0x9d, 0xa1, 0x9b, 0x9f,
	0x17, 0xdf, 0xbe, 0xfa, 0x86, 0x84, 0x0e, 0xa8, 0x25, 0x54, 0x2a, 0x26, 0xb8, 0xe9, 0xd2, 0xf0,
	0x8b, 0x2b, 0x1c, 0x80, 0x7a, 0xf1, 0x57, 0x4e, 0xc5, 0xd4, 0x6c, 0xa2, 0x6c, 0x05, 0xa8, 0x58,
	0x01, 0x1a, 0xf3, 0x9d, 0xbf, 0xa7, 0xbc, 0x36, 0x80, 0xfb, 0xdf, 0x55, 0x1b, 0x1a, 0xb0, 0x25,
	0xa3, 0x12, 0x56, 0x9e, 0x3d, 0xcb, 0xfb, 0x75, 0x7f, 0xfb, 0xf8, 0x54, 0x2d, 0xd7, 0x2d, 0xf0,
	0x29, 0x10, 0x6b, 0xb4, 0x62, 0xfa, 0x2a, 0x5e, 0xa0, 0x74, 0x6d, 0x87, 0xa3, 0x7b, 0xcd, 0x93,
	0xd9, 0xa7, 0x69, 0xd0, 0xd4, 0x5a, 0x54, 0x4d, 0xe2, 0xe8, 0x35, 0x00, 0x00, 0xff, 0xff, 0x8a,
	0x35, 0xc0, 0xf2, 0x35, 0x02, 0x00, 0x00,
}
