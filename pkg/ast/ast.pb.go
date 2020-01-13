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
	Value_OUT_POINT   Value_Type = 18
	Value_CELL_DEP    Value_Type = 19
	Value_SCRIPT      Value_Type = 20
	Value_CELL        Value_Type = 21
	Value_TRANSACTION Value_Type = 22
	// Compound fields
	Value_APPLY  Value_Type = 23
	Value_REDUCE Value_Type = 24
	// List fields
	Value_LIST        Value_Type = 25
	Value_QUERY_CELLS Value_Type = 26
	Value_MAP         Value_Type = 27
	Value_FILTER      Value_Type = 28
	// Cell get operations
	Value_GET_CAPACITY  Value_Type = 48
	Value_GET_DATA      Value_Type = 49
	Value_GET_LOCK      Value_Type = 50
	Value_GET_TYPE      Value_Type = 51
	Value_GET_DATA_HASH Value_Type = 52
	// Script get operations
	Value_GET_CODE_HASH Value_Type = 53
	Value_GET_HASH_TYPE Value_Type = 54
	Value_GET_ARGS      Value_Type = 55
	// Operations
	Value_HASH  Value_Type = 73
	Value_NOT   Value_Type = 74
	Value_AND   Value_Type = 75
	Value_OR    Value_Type = 76
	Value_XOR   Value_Type = 77
	Value_EQUAL Value_Type = 78
	Value_LESS  Value_Type = 79
	Value_LEN   Value_Type = 80
	Value_SLICE Value_Type = 81
	Value_INDEX Value_Type = 82
	Value_PLUS  Value_Type = 83
	Value_MINUS Value_Type = 84
)

var Value_Type_name = map[int32]string{
	0:  "UINT64",
	1:  "NIL",
	2:  "BOOL",
	3:  "BYTES",
	16: "ARG",
	17: "PARAM",
	18: "OUT_POINT",
	19: "CELL_DEP",
	20: "SCRIPT",
	21: "CELL",
	22: "TRANSACTION",
	23: "APPLY",
	24: "REDUCE",
	25: "LIST",
	26: "QUERY_CELLS",
	27: "MAP",
	28: "FILTER",
	48: "GET_CAPACITY",
	49: "GET_DATA",
	50: "GET_LOCK",
	51: "GET_TYPE",
	52: "GET_DATA_HASH",
	53: "GET_CODE_HASH",
	54: "GET_HASH_TYPE",
	55: "GET_ARGS",
	73: "HASH",
	74: "NOT",
	75: "AND",
	76: "OR",
	77: "XOR",
	78: "EQUAL",
	79: "LESS",
	80: "LEN",
	81: "SLICE",
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
	"OUT_POINT":     18,
	"CELL_DEP":      19,
	"SCRIPT":        20,
	"CELL":          21,
	"TRANSACTION":   22,
	"APPLY":         23,
	"REDUCE":        24,
	"LIST":          25,
	"QUERY_CELLS":   26,
	"MAP":           27,
	"FILTER":        28,
	"GET_CAPACITY":  48,
	"GET_DATA":      49,
	"GET_LOCK":      50,
	"GET_TYPE":      51,
	"GET_DATA_HASH": 52,
	"GET_CODE_HASH": 53,
	"GET_HASH_TYPE": 54,
	"GET_ARGS":      55,
	"HASH":          73,
	"NOT":           74,
	"AND":           75,
	"OR":            76,
	"XOR":           77,
	"EQUAL":         78,
	"LESS":          79,
	"LEN":           80,
	"SLICE":         81,
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
	// 554 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x54, 0x52, 0x5b, 0x6f, 0xd3, 0x30,
	0x14, 0x5e, 0x96, 0xac, 0xb4, 0x67, 0xb7, 0x33, 0x73, 0x0b, 0x37, 0x51, 0x15, 0x01, 0x7d, 0x6a,
	0x61, 0x1b, 0xe3, 0x0d, 0xc9, 0x4d, 0xcd, 0x1a, 0xe6, 0x26, 0x99, 0xed, 0xa0, 0x95, 0x97, 0x2a,
	0x1d, 0xd1, 0x16, 0x48, 0xd7, 0xaa, 0x4d, 0x60, 0xfc, 0x43, 0x7e, 0x03, 0xbf, 0x06, 0x9d, 0x8c,
	0x82, 0x78, 0xf3, 0x77, 0x39, 0x9f, 0x8f, 0xad, 0x0f, 0x1a, 0xc9, 0xb2, 0xe8, 0xcc, 0x17, 0xb3,
	0x62, 0xc6, 0xec, 0x64, 0x59, 0xb4, 0x7e, 0x39, 0xb0, 0xf1, 0x31, 0xc9, 0xcb, 0x94, 0x3d, 0x01,
	0xab, 0x70, 0xad, 0xa6, 0xd5, 0xde, 0xd9, 0xdf, 0xed, 0x90, 0xab, 0xa2, 0x3b, 0xe6, 0xc7, 0x3c,
	0x55, 0x56, 0xc1, 0x76, 0xc0, 0x9a, 0xb8, 0xeb, 0x4d, 0xab, 0x5d, 0x1f, 0xac, 0x29, 0x6b, 0x42,
	0xb8, 0x74, 0xed, 0xa6, 0xd5, 0x76, 0x08, 0x97, 0x8c, 0x81, 0xbd, 0x48, 0xbe, 0xbb, 0x4e, 0xd3,
	0x6a, 0x6f, 0x0d, 0xd6, 0x14, 0x01, 0xf6, 0x02, 0xea, 0xe7, 0x97, 0x59, 0xfe, 0x79, 0x91, 0x5e,
	0xb9, 0xf5, 0xa6, 0xdd, 0xde, 0xdc, 0x87, 0x7f, 0xc9, 0xea, 0xaf, 0xd6, 0xfa, 0x69, 0x83, 0x43,
	0xf7, 0x30, 0x80, 0x5a, 0xec, 0x07, 0xe6, 0xe8, 0x10, 0xd7, 0xd8, 0x2d, 0xb0, 0x03, 0x5f, 0xa2,
	0xc5, 0xea, 0xe0, 0xf4, 0xc2, 0x50, 0xe2, 0x3a, 0x6b, 0xc0, 0x46, 0x6f, 0x64, 0x84, 0x46, 0x9b,
	0x54, 0xae, 0x8e, 0x11, 0x89, 0x8b, 0xb8, 0xe2, 0x43, 0xdc, 0x63, 0xdb, 0xd0, 0x08, 0x63, 0x33,
	0x8e, 0x42, 0x3f, 0x30, 0xc8, 0xd8, 0x16, 0xd4, 0x3d, 0x21, 0xe5, 0xb8, 0x2f, 0x22, 0xbc, 0x4d,
	0xd1, 0xda, 0x53, 0x7e, 0x64, 0xf0, 0x0e, 0x25, 0x92, 0x82, 0x77, 0xd9, 0x2e, 0x6c, 0x1a, 0xc5,
	0x03, 0xcd, 0x3d, 0xe3, 0x87, 0x01, 0xde, 0xa3, 0x38, 0x1e, 0x45, 0x72, 0x84, 0xf7, 0x69, 0x42,
	0x89, 0x7e, 0xec, 0x09, 0x74, 0x69, 0x42, 0xfa, 0xda, 0xe0, 0x03, 0x9a, 0x38, 0x8d, 0x85, 0x1a,
	0x8d, 0x29, 0x41, 0xe3, 0x43, 0xda, 0x64, 0xc8, 0x23, 0x7c, 0x44, 0xfe, 0xf7, 0xbe, 0x34, 0x42,
	0xe1, 0x63, 0x86, 0xb0, 0x75, 0x2c, 0xcc, 0xd8, 0xe3, 0x11, 0xf7, 0x7c, 0x33, 0xc2, 0x57, 0xb4,
	0x0d, 0x31, 0x7d, 0x6e, 0x38, 0xbe, 0x5e, 0x21, 0x19, 0x7a, 0x27, 0xb8, 0xbf, 0x42, 0x66, 0x14,
	0x09, 0x3c, 0x60, 0x7b, 0xb0, 0xbd, 0x72, 0x8e, 0x07, 0x5c, 0x0f, 0xf0, 0x70, 0x45, 0x79, 0x61,
	0x5f, 0xdc, 0x50, 0x6f, 0x56, 0x14, 0xa1, 0x9b, 0xc1, 0xa3, 0x55, 0x0c, 0x57, 0xc7, 0x1a, 0xdf,
	0xd2, 0xca, 0x95, 0xd5, 0xaf, 0x7e, 0x32, 0x34, 0xf8, 0xa1, 0xfa, 0xb4, 0xa0, 0x8f, 0x27, 0xac,
	0x06, 0xeb, 0xa1, 0x42, 0x49, 0xc4, 0x59, 0xa8, 0x70, 0x48, 0xcf, 0x16, 0xa7, 0x31, 0x97, 0x18,
	0x54, 0x4f, 0x15, 0x5a, 0x63, 0x48, 0xaa, 0x14, 0x01, 0x46, 0xa4, 0x6a, 0xe9, 0x7b, 0x02, 0x4f,
	0xe9, 0xe8, 0x07, 0x7d, 0x71, 0x86, 0x8a, 0x8c, 0x91, 0x8c, 0x35, 0x6a, 0x22, 0x87, 0x7e, 0x10,
	0x6b, 0x34, 0xbd, 0x4d, 0x68, 0xcc, 0x17, 0xd9, 0x34, 0x2b, 0xb2, 0x6f, 0x69, 0xeb, 0x1d, 0x38,
	0x5e, 0x92, 0xe7, 0x8c, 0x81, 0x73, 0x95, 0x4c, 0xd3, 0xaa, 0x5d, 0x0d, 0x55, 0x9d, 0x59, 0x0b,
	0x6a, 0x8b, 0x74, 0x59, 0xe6, 0x45, 0x55, 0xa2, 0xff, 0x9b, 0xf1, 0x47, 0x69, 0xbd, 0x04, 0x47,
	0xcd, 0x66, 0x05, 0x7b, 0x0a, 0x1b, 0xe7, 0x49, 0x9e, 0x2f, 0x5d, 0xab, 0x2a, 0x51, 0xa3, 0xb2,
	0x52, 0xb2, 0xba, 0xe1, 0x7b, 0xcf, 0x3f, 0x3d, 0xbb, 0xc8, 0x8a, 0xcb, 0x72, 0xd2, 0x39, 0x9f,
	0x4d, 0xbb, 0xd7, 0xd7, 0x65, 0xfa, 0x25, 0x4b, 0xbb, 0xc9, 0x55, 0x36, 0x4d, 0x2e, 0xca, 0x65,
	0x77, 0xfe, 0xf5, 0xa2, 0x9b, 0x2c, 0x8b, 0x49, 0xad, 0x2a, 0xfe, 0xc1, 0xef, 0x00, 0x00, 0x00,
	0xff, 0xff, 0x71, 0xfe, 0x48, 0x01, 0x05, 0x03, 0x00, 0x00,
}
