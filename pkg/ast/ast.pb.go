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
	Value_CELL_INPUT  Value_Type = 19
	Value_CELL_DEP    Value_Type = 20
	Value_SCRIPT      Value_Type = 21
	Value_CELL        Value_Type = 22
	Value_TRANSACTION Value_Type = 23
	// Compound fields
	Value_APPLY  Value_Type = 24
	Value_REDUCE Value_Type = 25
	// List fields
	Value_LIST        Value_Type = 26
	Value_QUERY_CELLS Value_Type = 27
	Value_MAP         Value_Type = 28
	Value_FILTER      Value_Type = 29
	// Cell get operations
	Value_GET_CAPACITY  Value_Type = 48
	Value_GET_DATA      Value_Type = 49
	Value_GET_LOCK      Value_Type = 50
	Value_GET_TYPE      Value_Type = 51
	Value_GET_DATA_HASH Value_Type = 52
	Value_GET_OUT_POINT Value_Type = 53
	// Script get operations
	Value_GET_CODE_HASH Value_Type = 54
	Value_GET_HASH_TYPE Value_Type = 55
	Value_GET_ARGS      Value_Type = 56
	// Operations
	Value_HASH      Value_Type = 73
	Value_SERIALIZE Value_Type = 74
	Value_NOT       Value_Type = 75
	Value_AND       Value_Type = 76
	Value_OR        Value_Type = 77
	Value_XOR       Value_Type = 78
	Value_EQUAL     Value_Type = 79
	Value_LESS      Value_Type = 80
	Value_LEN       Value_Type = 81
	Value_SLICE     Value_Type = 82
	Value_INDEX     Value_Type = 83
	Value_PLUS      Value_Type = 84
	Value_MINUS     Value_Type = 85
)

var Value_Type_name = map[int32]string{
	0:  "UINT64",
	1:  "NIL",
	2:  "BOOL",
	3:  "BYTES",
	16: "ARG",
	17: "PARAM",
	18: "OUT_POINT",
	19: "CELL_INPUT",
	20: "CELL_DEP",
	21: "SCRIPT",
	22: "CELL",
	23: "TRANSACTION",
	24: "APPLY",
	25: "REDUCE",
	26: "LIST",
	27: "QUERY_CELLS",
	28: "MAP",
	29: "FILTER",
	48: "GET_CAPACITY",
	49: "GET_DATA",
	50: "GET_LOCK",
	51: "GET_TYPE",
	52: "GET_DATA_HASH",
	53: "GET_OUT_POINT",
	54: "GET_CODE_HASH",
	55: "GET_HASH_TYPE",
	56: "GET_ARGS",
	73: "HASH",
	74: "SERIALIZE",
	75: "NOT",
	76: "AND",
	77: "OR",
	78: "XOR",
	79: "EQUAL",
	80: "LESS",
	81: "LEN",
	82: "SLICE",
	83: "INDEX",
	84: "PLUS",
	85: "MINUS",
}

var Value_Type_value = map[string]int32{
	"UINT64":        0,
	"NIL":           1,
	"BOOL":          2,
	"BYTES":         3,
	"ARG":           16,
	"PARAM":         17,
	"OUT_POINT":     18,
	"CELL_INPUT":    19,
	"CELL_DEP":      20,
	"SCRIPT":        21,
	"CELL":          22,
	"TRANSACTION":   23,
	"APPLY":         24,
	"REDUCE":        25,
	"LIST":          26,
	"QUERY_CELLS":   27,
	"MAP":           28,
	"FILTER":        29,
	"GET_CAPACITY":  48,
	"GET_DATA":      49,
	"GET_LOCK":      50,
	"GET_TYPE":      51,
	"GET_DATA_HASH": 52,
	"GET_OUT_POINT": 53,
	"GET_CODE_HASH": 54,
	"GET_HASH_TYPE": 55,
	"GET_ARGS":      56,
	"HASH":          73,
	"SERIALIZE":     74,
	"NOT":           75,
	"AND":           76,
	"OR":            77,
	"XOR":           78,
	"EQUAL":         79,
	"LESS":          80,
	"LEN":           81,
	"SLICE":         82,
	"INDEX":         83,
	"PLUS":          84,
	"MINUS":         85,
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
	// 582 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x54, 0x53, 0x5b, 0x73, 0xd3, 0x3a,
	0x10, 0xae, 0x6b, 0x27, 0x27, 0x51, 0x6f, 0x5b, 0x9d, 0x73, 0xc0, 0x5c, 0x3a, 0x64, 0xc2, 0x00,
	0x79, 0x4a, 0xa0, 0x2d, 0x85, 0x27, 0x66, 0x14, 0x47, 0x34, 0xa2, 0x8a, 0xed, 0x4a, 0x32, 0xd3,
	0xf4, 0x25, 0xe3, 0x14, 0x4f, 0x6b, 0x70, 0x9a, 0x4c, 0x62, 0x43, 0xf9, 0x0f, 0xfc, 0x62, 0x9e,
	0x98, 0x75, 0x31, 0x1d, 0xde, 0xf4, 0x5d, 0xf6, 0x5b, 0x69, 0x67, 0x45, 0x9a, 0xf1, 0x2a, 0xef,
	0x2e, 0x96, 0xf3, 0x7c, 0x4e, 0xed, 0x78, 0x95, 0xb7, 0x7f, 0xd4, 0x48, 0xed, 0x63, 0x9c, 0x15,
	0x09, 0xdd, 0x23, 0x56, 0xee, 0x5a, 0x2d, 0xab, 0xb3, 0xbd, 0xbf, 0xd3, 0x45, 0x57, 0x49, 0x77,
	0xcd, 0xf7, 0x45, 0xa2, 0xac, 0x9c, 0x6e, 0x13, 0x6b, 0xea, 0xae, 0xb7, 0xac, 0x4e, 0x63, 0xb8,
	0xa6, 0xac, 0x29, 0xe2, 0xc2, 0xb5, 0x5b, 0x56, 0xc7, 0x41, 0x5c, 0x50, 0x4a, 0xec, 0x65, 0xfc,
	0xcd, 0x75, 0x5a, 0x56, 0x67, 0x73, 0xb8, 0xa6, 0x10, 0xd0, 0xe7, 0xa4, 0x71, 0x71, 0x95, 0x66,
	0x9f, 0x96, 0xc9, 0xb5, 0xdb, 0x68, 0xd9, 0x9d, 0x8d, 0x7d, 0x72, 0x97, 0xac, 0xfe, 0x68, 0xed,
	0x9f, 0x36, 0x71, 0xb0, 0x0f, 0x25, 0xa4, 0x1e, 0x09, 0xdf, 0x1c, 0x1d, 0xc2, 0x1a, 0xfd, 0x87,
	0xd8, 0xbe, 0x90, 0x60, 0xd1, 0x06, 0x71, 0xfa, 0x41, 0x20, 0x61, 0x9d, 0x36, 0x49, 0xad, 0x3f,
	0x36, 0x5c, 0x83, 0x8d, 0x2a, 0x53, 0xc7, 0x00, 0xc8, 0x85, 0x4c, 0xb1, 0x11, 0xec, 0xd2, 0x2d,
	0xd2, 0x0c, 0x22, 0x33, 0x09, 0x03, 0xe1, 0x1b, 0xa0, 0x74, 0x9b, 0x10, 0x8f, 0x4b, 0x39, 0x11,
	0x7e, 0x18, 0x19, 0xf8, 0x97, 0x6e, 0x92, 0x46, 0x89, 0x07, 0x3c, 0x84, 0xff, 0xb0, 0x95, 0xf6,
	0x94, 0x08, 0x0d, 0xfc, 0x8f, 0x1d, 0x50, 0x81, 0x7b, 0x74, 0x87, 0x6c, 0x18, 0xc5, 0x7c, 0xcd,
	0x3c, 0x23, 0x02, 0x1f, 0xee, 0x63, 0x3c, 0x0b, 0x43, 0x39, 0x06, 0x17, 0x2b, 0x14, 0x1f, 0x44,
	0x1e, 0x87, 0x07, 0x58, 0x21, 0x85, 0x36, 0xf0, 0x10, 0x2b, 0x4e, 0x23, 0xae, 0xc6, 0x13, 0x4c,
	0xd0, 0xf0, 0x08, 0x6f, 0x36, 0x62, 0x21, 0x3c, 0x46, 0xff, 0x7b, 0x21, 0x0d, 0x57, 0xb0, 0x47,
	0x81, 0x6c, 0x1e, 0x73, 0x33, 0xf1, 0x58, 0xc8, 0x3c, 0x61, 0xc6, 0xf0, 0x12, 0x6f, 0x83, 0xcc,
	0x80, 0x19, 0x06, 0xaf, 0x2a, 0x24, 0x03, 0xef, 0x04, 0xf6, 0x2b, 0x64, 0xc6, 0x21, 0x87, 0x03,
	0xba, 0x4b, 0xb6, 0x2a, 0xe7, 0x64, 0xc8, 0xf4, 0x10, 0x0e, 0x2b, 0xea, 0xee, 0xb5, 0xaf, 0x2b,
	0xca, 0x0b, 0x06, 0xfc, 0xd6, 0x75, 0x54, 0x51, 0x88, 0x6e, 0xb3, 0xde, 0x54, 0xc9, 0x4c, 0x1d,
	0x6b, 0x78, 0x8b, 0xaf, 0x28, 0xad, 0x02, 0x47, 0xa7, 0xb9, 0x12, 0x4c, 0x8a, 0x73, 0x0e, 0x1f,
	0xca, 0xd9, 0x07, 0x06, 0x4e, 0xca, 0x31, 0xfb, 0x03, 0x90, 0xb4, 0x4e, 0xd6, 0x03, 0x05, 0x23,
	0x24, 0xce, 0x02, 0x05, 0x3e, 0x0e, 0x86, 0x9f, 0x46, 0x4c, 0x42, 0x50, 0x0e, 0x83, 0x6b, 0x0d,
	0x21, 0xaa, 0x92, 0xfb, 0x70, 0x8a, 0xaa, 0x96, 0xc2, 0xe3, 0xa0, 0xf0, 0x28, 0xfc, 0x01, 0x3f,
	0x03, 0x8d, 0xc6, 0x50, 0x46, 0x1a, 0x0c, 0x92, 0x23, 0xe1, 0x47, 0x1a, 0xa2, 0xfe, 0x06, 0x69,
	0x2e, 0x96, 0xe9, 0x2c, 0xcd, 0xd3, 0xaf, 0x49, 0xfb, 0x1d, 0x71, 0xbc, 0x38, 0xcb, 0x28, 0x25,
	0xce, 0x75, 0x3c, 0x4b, 0xca, 0x7d, 0x6c, 0xaa, 0xf2, 0x4c, 0xdb, 0xa4, 0xbe, 0x4c, 0x56, 0x45,
	0x96, 0x97, 0x6b, 0xf7, 0xf7, 0x2e, 0xfd, 0x56, 0xda, 0x2f, 0x88, 0xa3, 0xe6, 0xf3, 0x9c, 0x3e,
	0x21, 0xb5, 0x8b, 0x38, 0xcb, 0x56, 0xae, 0x55, 0xae, 0x5d, 0xb3, 0xb4, 0x62, 0xb2, 0xba, 0xe5,
	0xfb, 0xcf, 0xce, 0x9f, 0x5e, 0xa6, 0xf9, 0x55, 0x31, 0xed, 0x5e, 0xcc, 0x67, 0xbd, 0x9b, 0x9b,
	0x22, 0xf9, 0x9c, 0x26, 0xbd, 0xf8, 0x3a, 0x9d, 0xc5, 0x97, 0xc5, 0xaa, 0xb7, 0xf8, 0x72, 0xd9,
	0x8b, 0x57, 0xf9, 0xb4, 0x5e, 0x7e, 0x95, 0x83, 0x5f, 0x01, 0x00, 0x00, 0xff, 0xff, 0x3e, 0xb1,
	0xa2, 0x77, 0x37, 0x03, 0x00, 0x00,
}
