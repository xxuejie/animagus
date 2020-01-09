// Code generated by protoc-gen-go. DO NOT EDIT.
// source: ast.proto

package ast

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

type Value_Type int32

const (
	// Primitive fields
	Value_UINT64 Value_Type = 0
	Value_NIL    Value_Type = 1
	Value_BOOL   Value_Type = 2
	Value_BYTES  Value_Type = 3
	// In animagus, we distinguish args and params in the following way:
	// * If a Value struct contains an arg, it will be interpretted as a
	// function, when used in constructs such as REDUCE or MAP, args acts
	// as placeholders for the values to test/transform.
	// * Params, on the other hand, denotes user input when calling RPCs,
	// for example, a user might provide an amount to transfer, or an address
	// to transfer to, those will be represented via parameters
	Value_ARG   Value_Type = 16
	Value_PARAM Value_Type = 17
	// Blockchain data structures
	Value_SCRIPT Value_Type = 18
	Value_CELL   Value_Type = 19
	// Compound fields
	Value_APPLY  Value_Type = 22
	Value_REDUCE Value_Type = 23
	// List fields
	Value_LIST        Value_Type = 24
	Value_MAP         Value_Type = 25
	Value_FILTER      Value_Type = 26
	Value_QUERY_CELLS Value_Type = 27
	// Cell get operations
	Value_GET_CAPACITY Value_Type = 48
	Value_GET_DATA     Value_Type = 49
	Value_GET_LOCK     Value_Type = 50
	Value_GET_TYPE     Value_Type = 51
	// Script get operations
	Value_GET_CODE_HASH Value_Type = 52
	Value_GET_HASH_TYPE Value_Type = 53
	Value_GET_ARGS      Value_Type = 54
	// Operations
	Value_HASH        Value_Type = 73
	Value_NOT         Value_Type = 74
	Value_AND         Value_Type = 75
	Value_OR          Value_Type = 76
	Value_XOR         Value_Type = 77
	Value_EQUAL       Value_Type = 78
	Value_LESS        Value_Type = 79
	Value_LEN         Value_Type = 80
	Value_SLICE_BYTES Value_Type = 81
	Value_INDEX       Value_Type = 82
	Value_PLUS        Value_Type = 83
	Value_MINUS       Value_Type = 84
)

var Value_Type_name = map[int32]string{
	0:  "UINT64",
	1:  "NIL",
	2:  "BOOL",
	3:  "BYTES",
	16: "ARG",
	17: "PARAM",
	18: "SCRIPT",
	19: "CELL",
	22: "APPLY",
	23: "REDUCE",
	24: "LIST",
	25: "MAP",
	26: "FILTER",
	27: "QUERY_CELLS",
	48: "GET_CAPACITY",
	49: "GET_DATA",
	50: "GET_LOCK",
	51: "GET_TYPE",
	52: "GET_CODE_HASH",
	53: "GET_HASH_TYPE",
	54: "GET_ARGS",
	73: "HASH",
	74: "NOT",
	75: "AND",
	76: "OR",
	77: "XOR",
	78: "EQUAL",
	79: "LESS",
	80: "LEN",
	81: "SLICE_BYTES",
	82: "INDEX",
	83: "PLUS",
	84: "MINUS",
}

var Value_Type_value = map[string]int32{
	"UINT64":        0,
	"NIL":           1,
	"BOOL":          2,
	"BYTES":         3,
	"ARG":           16,
	"PARAM":         17,
	"SCRIPT":        18,
	"CELL":          19,
	"APPLY":         22,
	"REDUCE":        23,
	"LIST":          24,
	"MAP":           25,
	"FILTER":        26,
	"QUERY_CELLS":   27,
	"GET_CAPACITY":  48,
	"GET_DATA":      49,
	"GET_LOCK":      50,
	"GET_TYPE":      51,
	"GET_CODE_HASH": 52,
	"GET_HASH_TYPE": 53,
	"GET_ARGS":      54,
	"HASH":          73,
	"NOT":           74,
	"AND":           75,
	"OR":            76,
	"XOR":           77,
	"EQUAL":         78,
	"LESS":          79,
	"LEN":           80,
	"SLICE_BYTES":   81,
	"INDEX":         82,
	"PLUS":          83,
	"MINUS":         84,
}

func (x Value_Type) String() string {
	return proto.EnumName(Value_Type_name, int32(x))
}

func (Value_Type) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_37b5b141da493253, []int{0, 0}
}

type Value struct {
	T Value_Type `protobuf:"varint,1,opt,name=t,proto3,enum=ast.Value_Type" json:"t,omitempty"`
	// Types that are valid to be assigned to Primitive:
	//	*Value_B
	//	*Value_U
	//	*Value_Raw
	Primitive            isValue_Primitive `protobuf_oneof:"primitive"`
	Children             []*Value          `protobuf:"bytes,8,rep,name=children,proto3" json:"children,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *Value) Reset()         { *m = Value{} }
func (m *Value) String() string { return proto.CompactTextString(m) }
func (*Value) ProtoMessage()    {}
func (*Value) Descriptor() ([]byte, []int) {
	return fileDescriptor_37b5b141da493253, []int{0}
}

func (m *Value) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Value.Unmarshal(m, b)
}
func (m *Value) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Value.Marshal(b, m, deterministic)
}
func (m *Value) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Value.Merge(m, src)
}
func (m *Value) XXX_Size() int {
	return xxx_messageInfo_Value.Size(m)
}
func (m *Value) XXX_DiscardUnknown() {
	xxx_messageInfo_Value.DiscardUnknown(m)
}

var xxx_messageInfo_Value proto.InternalMessageInfo

func (m *Value) GetT() Value_Type {
	if m != nil {
		return m.T
	}
	return Value_UINT64
}

type isValue_Primitive interface {
	isValue_Primitive()
}

type Value_B struct {
	B bool `protobuf:"varint,2,opt,name=b,proto3,oneof"`
}

type Value_U struct {
	U uint64 `protobuf:"varint,3,opt,name=u,proto3,oneof"`
}

type Value_Raw struct {
	Raw []byte `protobuf:"bytes,4,opt,name=raw,proto3,oneof"`
}

func (*Value_B) isValue_Primitive() {}

func (*Value_U) isValue_Primitive() {}

func (*Value_Raw) isValue_Primitive() {}

func (m *Value) GetPrimitive() isValue_Primitive {
	if m != nil {
		return m.Primitive
	}
	return nil
}

func (m *Value) GetB() bool {
	if x, ok := m.GetPrimitive().(*Value_B); ok {
		return x.B
	}
	return false
}

func (m *Value) GetU() uint64 {
	if x, ok := m.GetPrimitive().(*Value_U); ok {
		return x.U
	}
	return 0
}

func (m *Value) GetRaw() []byte {
	if x, ok := m.GetPrimitive().(*Value_Raw); ok {
		return x.Raw
	}
	return nil
}

func (m *Value) GetChildren() []*Value {
	if m != nil {
		return m.Children
	}
	return nil
}

// XXX_OneofWrappers is for the internal use of the proto package.
func (*Value) XXX_OneofWrappers() []interface{} {
	return []interface{}{
		(*Value_B)(nil),
		(*Value_U)(nil),
		(*Value_Raw)(nil),
	}
}

type Call struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Result               *Value   `protobuf:"bytes,3,opt,name=result,proto3" json:"result,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Call) Reset()         { *m = Call{} }
func (m *Call) String() string { return proto.CompactTextString(m) }
func (*Call) ProtoMessage()    {}
func (*Call) Descriptor() ([]byte, []int) {
	return fileDescriptor_37b5b141da493253, []int{1}
}

func (m *Call) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Call.Unmarshal(m, b)
}
func (m *Call) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Call.Marshal(b, m, deterministic)
}
func (m *Call) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Call.Merge(m, src)
}
func (m *Call) XXX_Size() int {
	return xxx_messageInfo_Call.Size(m)
}
func (m *Call) XXX_DiscardUnknown() {
	xxx_messageInfo_Call.DiscardUnknown(m)
}

var xxx_messageInfo_Call proto.InternalMessageInfo

func (m *Call) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Call) GetResult() *Value {
	if m != nil {
		return m.Result
	}
	return nil
}

type Root struct {
	Calls                []*Call  `protobuf:"bytes,1,rep,name=calls,proto3" json:"calls,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Root) Reset()         { *m = Root{} }
func (m *Root) String() string { return proto.CompactTextString(m) }
func (*Root) ProtoMessage()    {}
func (*Root) Descriptor() ([]byte, []int) {
	return fileDescriptor_37b5b141da493253, []int{2}
}

func (m *Root) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Root.Unmarshal(m, b)
}
func (m *Root) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Root.Marshal(b, m, deterministic)
}
func (m *Root) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Root.Merge(m, src)
}
func (m *Root) XXX_Size() int {
	return xxx_messageInfo_Root.Size(m)
}
func (m *Root) XXX_DiscardUnknown() {
	xxx_messageInfo_Root.DiscardUnknown(m)
}

var xxx_messageInfo_Root proto.InternalMessageInfo

func (m *Root) GetCalls() []*Call {
	if m != nil {
		return m.Calls
	}
	return nil
}

func init() {
	proto.RegisterEnum("ast.Value_Type", Value_Type_name, Value_Type_value)
	proto.RegisterType((*Value)(nil), "ast.Value")
	proto.RegisterType((*Call)(nil), "ast.Call")
	proto.RegisterType((*Root)(nil), "ast.Root")
}

func init() { proto.RegisterFile("ast.proto", fileDescriptor_37b5b141da493253) }

var fileDescriptor_37b5b141da493253 = []byte{
	// 524 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x54, 0x92, 0xd1, 0x6e, 0xd3, 0x4c,
	0x10, 0x85, 0xeb, 0xda, 0xcd, 0x9f, 0x4c, 0xfb, 0xb7, 0xd3, 0x45, 0x02, 0x03, 0x42, 0x44, 0x41,
	0x40, 0xae, 0x12, 0x68, 0x4b, 0x2f, 0x91, 0x36, 0xce, 0xd2, 0x98, 0x6e, 0x6c, 0x77, 0x77, 0x83,
	0x1a, 0x6e, 0x22, 0xa7, 0x58, 0xad, 0xc1, 0x69, 0xa2, 0xc4, 0x86, 0xf2, 0x10, 0x3c, 0x03, 0xaf,
	0x8a, 0x66, 0x43, 0x40, 0xdc, 0xf9, 0x9c, 0xf3, 0xcd, 0xd8, 0x33, 0x1e, 0x68, 0xa4, 0xab, 0xb2,
	0xb3, 0x58, 0xce, 0xcb, 0x39, 0x73, 0xd3, 0x55, 0xd9, 0xfa, 0xe9, 0xc1, 0xce, 0x87, 0xb4, 0xa8,
	0x32, 0xf6, 0x04, 0x9c, 0xd2, 0x77, 0x9a, 0x4e, 0x7b, 0xff, 0xe8, 0xa0, 0x43, 0x94, 0xb5, 0x3b,
	0xe6, 0xfb, 0x22, 0x53, 0x4e, 0xc9, 0xf6, 0xc1, 0x99, 0xfa, 0xdb, 0x4d, 0xa7, 0x5d, 0x1f, 0x6c,
	0x29, 0x67, 0x4a, 0xba, 0xf2, 0xdd, 0xa6, 0xd3, 0xf6, 0x48, 0x57, 0x8c, 0x81, 0xbb, 0x4c, 0xbf,
	0xf9, 0x5e, 0xd3, 0x69, 0xef, 0x0d, 0xb6, 0x14, 0x09, 0xf6, 0x02, 0xea, 0x57, 0x37, 0x79, 0xf1,
	0x69, 0x99, 0xdd, 0xfa, 0xf5, 0xa6, 0xdb, 0xde, 0x3d, 0x82, 0xbf, 0x9d, 0xd5, 0x9f, 0xac, 0xf5,
	0xc3, 0x05, 0x8f, 0xde, 0xc3, 0x00, 0x6a, 0xa3, 0x30, 0x32, 0xa7, 0x27, 0xb8, 0xc5, 0xfe, 0x03,
	0x37, 0x0a, 0x25, 0x3a, 0xac, 0x0e, 0x5e, 0x2f, 0x8e, 0x25, 0x6e, 0xb3, 0x06, 0xec, 0xf4, 0xc6,
	0x46, 0x68, 0x74, 0x29, 0xe5, 0xea, 0x0c, 0x91, 0xbc, 0x84, 0x2b, 0x3e, 0xc4, 0x43, 0xaa, 0xd6,
	0x81, 0x0a, 0x13, 0x83, 0x8c, 0x8a, 0x02, 0x21, 0x25, 0xde, 0x23, 0x80, 0x27, 0x89, 0x1c, 0xe3,
	0x7d, 0x02, 0x94, 0xe8, 0x8f, 0x02, 0x81, 0x0f, 0x08, 0x90, 0xa1, 0x36, 0xe8, 0x53, 0xab, 0x21,
	0x4f, 0xf0, 0x21, 0xc5, 0xef, 0x42, 0x69, 0x84, 0xc2, 0x47, 0xec, 0x00, 0x76, 0x2f, 0x46, 0x42,
	0x8d, 0x27, 0xd4, 0x45, 0xe3, 0x63, 0x86, 0xb0, 0x77, 0x26, 0xcc, 0x24, 0xe0, 0x09, 0x0f, 0x42,
	0x33, 0xc6, 0x57, 0x6c, 0x0f, 0xea, 0xe4, 0xf4, 0xb9, 0xe1, 0xf8, 0x7a, 0xa3, 0x64, 0x1c, 0x9c,
	0xe3, 0xd1, 0x46, 0x99, 0x71, 0x22, 0xf0, 0x98, 0x1d, 0xc2, 0xff, 0xb6, 0x36, 0xee, 0x8b, 0xc9,
	0x80, 0xeb, 0x01, 0x9e, 0x6c, 0x2c, 0x52, 0x6b, 0xea, 0xcd, 0xa6, 0x86, 0xab, 0x33, 0x8d, 0xa7,
	0xf4, 0x7d, 0x16, 0x0d, 0xed, 0x22, 0x62, 0x83, 0xef, 0xed, 0xcc, 0x51, 0x1f, 0xcf, 0x59, 0x0d,
	0xb6, 0x63, 0x85, 0x92, 0x8c, 0xcb, 0x58, 0xe1, 0x90, 0x66, 0x14, 0x17, 0x23, 0x2e, 0x31, 0xb2,
	0x73, 0x09, 0xad, 0x31, 0xa6, 0x54, 0x8a, 0x08, 0x13, 0x9a, 0x45, 0xcb, 0x30, 0x10, 0x93, 0xf5,
	0xf2, 0x2e, 0x08, 0x0f, 0xa3, 0xbe, 0xb8, 0x44, 0x45, 0x78, 0x22, 0x47, 0x1a, 0x35, 0x99, 0xc3,
	0x30, 0x1a, 0x69, 0x34, 0xbd, 0x5d, 0x68, 0x2c, 0x96, 0xf9, 0x2c, 0x2f, 0xf3, 0xaf, 0x59, 0xeb,
	0x2d, 0x78, 0x41, 0x5a, 0x14, 0x8c, 0x81, 0x77, 0x9b, 0xce, 0x32, 0x7b, 0x22, 0x0d, 0x65, 0x9f,
	0x59, 0x0b, 0x6a, 0xcb, 0x6c, 0x55, 0x15, 0xa5, 0xbd, 0x84, 0x7f, 0x7f, 0xef, 0xef, 0xa4, 0xf5,
	0x12, 0x3c, 0x35, 0x9f, 0x97, 0xec, 0x29, 0xec, 0x5c, 0xa5, 0x45, 0xb1, 0xf2, 0x1d, 0x7b, 0x09,
	0x0d, 0x8b, 0x52, 0x67, 0xb5, 0xf6, 0x7b, 0xcf, 0x3f, 0x3e, 0xbb, 0xce, 0xcb, 0x9b, 0x6a, 0xda,
	0xb9, 0x9a, 0xcf, 0xba, 0x77, 0x77, 0x55, 0xf6, 0x39, 0xcf, 0xba, 0xe9, 0x6d, 0x3e, 0x4b, 0xaf,
	0xab, 0x55, 0x77, 0xf1, 0xe5, 0xba, 0x9b, 0xae, 0xca, 0x69, 0xcd, 0x5e, 0xef, 0xf1, 0xaf, 0x00,
	0x00, 0x00, 0xff, 0xff, 0x99, 0x9a, 0xf4, 0xd6, 0xca, 0x02, 0x00, 0x00,
}
