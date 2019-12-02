package coretypes

import (
	"encoding/binary"
	"math/big"
)

const (
	Data byte = 0
	Type byte = 1

	Code     byte = 0
	DepGroup byte = 1

	Byte32Size          = 32
	HashSize            = 32
	OutPointSize        = 36
	CellInputSize       = 44
	CellDepSize         = 37
	RawHeaderSize       = 192
	HeaderSize          = 208
	ProposalShortIdSize = 10
)

type CoreVerifier interface {
	Verify(compatible bool) bool
}

type Byte32 []byte

func (b Byte32) Verify(_compatible bool) bool {
	return len(b) == Byte32Size
}

type Hash []byte

func (h Hash) Verify(_compatible bool) bool {
	return len(h) == HashSize
}

type ScriptHashType []byte

func (t ScriptHashType) Verify(_compatible bool) bool {
	if len(t) != 1 {
		return false
	}
	return t[0] == Data || t[0] == Type
}

func (t ScriptHashType) Value() byte {
	return t[0]
}

type DepType []byte

func (t DepType) Verify(_compatible bool) bool {
	if len(t) != 1 {
		return false
	}
	return t[0] == Code || t[0] == DepGroup
}

func (t DepType) Value() byte {
	return t[0]
}

type Bytes []byte

func (b Bytes) Verify(_compatible bool) bool {
	if len(b) < 4 {
		return false
	}
	count := int(binary.LittleEndian.Uint32(b[0:4]))
	return len(b) == 4+count
}

func (b Bytes) Value() []byte {
	return b[4:]
}

type OutPoint []byte

func (o OutPoint) Verify(compatible bool) bool {
	return len(o) == OutPointSize && o.TxHash().Verify(compatible)
}

func (o OutPoint) TxHash() Hash {
	return Hash(o[0:32])
}

func (o OutPoint) Index() uint32 {
	return binary.LittleEndian.Uint32(o[32:36])
}

type CellInput []byte

func (c CellInput) Verify(compatible bool) bool {
	return len(c) == CellInputSize && c.PreviousOutput().Verify(compatible)
}

func (c CellInput) Since() uint64 {
	return binary.LittleEndian.Uint64(c[0:8])
}

func (c CellInput) PreviousOutput() OutPoint {
	return OutPoint(c[8:44])
}

func extractOffsets(b []byte) ([]int, bool) {
	if (len(b)) < 4 {
		return nil, false
	}
	slice_len := int(binary.LittleEndian.Uint32(b[0:4]))
	if slice_len != len(b) {
		return nil, false
	}
	if slice_len == 4 {
		return []int{len(b)}, true
	}
	if slice_len < 8 {
		return nil, false
	}
	first_offset := int(binary.LittleEndian.Uint32(b[4:8]))
	if first_offset%4 != 0 || first_offset < 8 {
		return nil, false
	}
	field_count := first_offset/4 - 1
	if slice_len < first_offset {
		return nil, false
	}
	offsets := make([]int, field_count+1)
	for i := 0; i < field_count; i++ {
		start := 4 + i*4
		offsets[i] = int(binary.LittleEndian.Uint32(b[start : start+4]))
	}
	offsets[field_count] = slice_len
	return offsets, true
}

func verifyAndExtractOffsets(b []byte, expected_field_count int, compatible bool) ([]int, bool) {
	offsets, success := extractOffsets(b)
	if !success {
		return nil, false
	}
	field_count := len(offsets) - 1
	if field_count < expected_field_count {
		return nil, false
	} else if (!compatible) && field_count > expected_field_count {
		return nil, false
	}
	return offsets, true
}

func uncheckedField(b []byte, index int, last bool) []byte {
	start := 4 + index*4
	offset := int(binary.LittleEndian.Uint32(b[start : start+4]))
	var offset_end int
	if !last {
		offset_end = int(binary.LittleEndian.Uint32(b[start+4 : start+8]))
	} else {
		field_count := int(binary.LittleEndian.Uint32(b[4:8]))/4 - 1
		if index+1 < field_count {
			offset_end = int(binary.LittleEndian.Uint32(b[start+4 : start+8]))
		} else {
			offset_end = len(b)
		}
	}
	return b[offset:offset_end]
}

type Script []byte

func (s Script) Verify(compatible bool) bool {
	offsets, success := verifyAndExtractOffsets(s, 3, compatible)
	if !success {
		return false
	}
	if !Hash(s[offsets[0]:offsets[1]]).Verify(compatible) {
		return false
	}
	if !ScriptHashType(s[offsets[1]:offsets[2]]).Verify(compatible) {
		return false
	}
	if !Bytes(s[offsets[2]:offsets[3]]).Verify(compatible) {
		return false
	}
	return true
}

func (s Script) CodeHash() Hash {
	return Hash(uncheckedField(s, 0, false))
}

func (s Script) HashType() ScriptHashType {
	return ScriptHashType(uncheckedField(s, 1, false))
}

func (s Script) Args() Bytes {
	return Bytes(uncheckedField(s, 2, true))
}

type CellOutput []byte

func (c CellOutput) Verify(compatible bool) bool {
	offsets, success := verifyAndExtractOffsets(c, 3, compatible)
	if !success {
		return false
	}
	if len(c[offsets[0]:offsets[1]]) != 8 {
		return false
	}
	if !Script(c[offsets[1]:offsets[2]]).Verify(compatible) {
		return false
	}
	typeSlice := c[offsets[2]:offsets[3]]
	if len(typeSlice) > 0 {
		if !Script(typeSlice).Verify(compatible) {
			return false
		}
	}
	return true
}

func (c CellOutput) Capacity() uint64 {
	return binary.LittleEndian.Uint64(uncheckedField(c, 0, false))
}

func (c CellOutput) Lock() Script {
	return Script(uncheckedField(c, 1, false))
}

func (c CellOutput) HasType() bool {
	return len(uncheckedField(c, 2, true)) > 0
}

func (c CellOutput) Type() Script {
	return Script(uncheckedField(c, 2, true))
}

func (c CellOutput) MaybeType() *Script {
	field := uncheckedField(c, 2, true)
	if len(field) > 0 {
		s := Script(field)
		return &s
	} else {
		return nil
	}
}

type CellDep []byte

func (c CellDep) Verify(compatible bool) bool {
	return len(c) == CellDepSize &&
		c.OutPoint().Verify(compatible) &&
		c.DepType().Verify(compatible)
}

func (c CellDep) OutPoint() OutPoint {
	return OutPoint(c[0:36])
}

func (c CellDep) DepType() DepType {
	return DepType(c[36:37])
}

type CellDepFixVec []byte

func (c CellDepFixVec) Verify(compatible bool) bool {
	if (len(c)) < 4 {
		return false
	}
	count := int(binary.LittleEndian.Uint32(c[0:4]))
	if len(c) != 4+CellDepSize*count {
		return false
	}
	for i := 0; i < count; i++ {
		if !c.Get(i).Verify(compatible) {
			return false
		}
	}
	return true
}

func (c CellDepFixVec) Len() int {
	return int(binary.LittleEndian.Uint32(c[0:4]))
}

func (c CellDepFixVec) Get(index int) CellDep {
	start := 4 + index*CellDepSize
	return CellDep(c[start : start+CellDepSize])
}

type HashFixVec []byte

func (h HashFixVec) Verify(compatible bool) bool {
	if (len(h)) < 4 {
		return false
	}
	count := int(binary.LittleEndian.Uint32(h[0:4]))
	if len(h) != 4+HashSize*count {
		return false
	}
	for i := 0; i < count; i++ {
		if !h.Get(i).Verify(compatible) {
			return false
		}
	}
	return true
}

func (h HashFixVec) Len() int {
	return int(binary.LittleEndian.Uint32(h[0:4]))
}

func (h HashFixVec) Get(index int) Hash {
	start := 4 + index*HashSize
	return Hash(h[start : start+HashSize])
}

type CellInputFixVec []byte

func (c CellInputFixVec) Verify(compatible bool) bool {
	if (len(c)) < 4 {
		return false
	}
	count := int(binary.LittleEndian.Uint32(c[0:4]))
	if len(c) != 4+CellInputSize*count {
		return false
	}
	for i := 0; i < count; i++ {
		if !c.Get(i).Verify(compatible) {
			return false
		}
	}
	return true
}

func (c CellInputFixVec) Len() int {
	return int(binary.LittleEndian.Uint32(c[0:4]))
}

func (c CellInputFixVec) Get(index int) CellInput {
	start := 4 + index*CellInputSize
	return CellInput(c[start : start+CellInputSize])
}

type CellOutputDynVec []byte

func (c CellOutputDynVec) Verify(compatible bool) bool {
	offsets, success := extractOffsets(c)
	if !success {
		return false
	}
	for i := 0; i < len(offsets)-1; i++ {
		if !CellOutput(c[offsets[i]:offsets[i+1]]).Verify(compatible) {
			return false
		}
	}
	return true
}

func (c CellOutputDynVec) Len() int {
	if len(c) < 8 {
		return 0
	} else {
		return int(binary.LittleEndian.Uint32(c[4:8]))/4 - 1
	}
}

func (c CellOutputDynVec) Get(i int) CellOutput {
	return CellOutput(uncheckedField(c, i, true))
}

type BytesDynVec []byte

func (b BytesDynVec) Verify(compatible bool) bool {
	offsets, success := extractOffsets(b)
	if !success {
		return false
	}
	for i := 0; i < len(offsets)-1; i++ {
		if !Bytes(b[offsets[i]:offsets[i+1]]).Verify(compatible) {
			return false
		}
	}
	return true
}

func (b BytesDynVec) Len() int {
	if len(b) < 8 {
		return 0
	} else {
		return int(binary.LittleEndian.Uint32(b[4:8]))/4 - 1
	}
}

func (b BytesDynVec) Get(i int) Bytes {
	return Bytes(uncheckedField(b, i, true))
}

type RawTransaction []byte

func (t RawTransaction) Verify(compatible bool) bool {
	offsets, success := verifyAndExtractOffsets(t, 6, compatible)
	if !success {
		return false
	}
	if len(t[offsets[0]:offsets[1]]) != 4 {
		return false
	}
	if !CellDepFixVec(t[offsets[1]:offsets[2]]).Verify(compatible) {
		return false
	}
	if !HashFixVec(t[offsets[2]:offsets[3]]).Verify(compatible) {
		return false
	}
	if !CellInputFixVec(t[offsets[3]:offsets[4]]).Verify(compatible) {
		return false
	}
	if !CellOutputDynVec(t[offsets[4]:offsets[5]]).Verify(compatible) {
		return false
	}
	if !BytesDynVec(t[offsets[5]:offsets[6]]).Verify(compatible) {
		return false
	}
	return true
}

func (t RawTransaction) Version() uint32 {
	return binary.LittleEndian.Uint32(uncheckedField(t, 0, false))
}

func (t RawTransaction) CellDeps() CellDepFixVec {
	return CellDepFixVec(uncheckedField(t, 1, false))
}

func (t RawTransaction) HeaderDeps() HashFixVec {
	return HashFixVec(uncheckedField(t, 2, false))
}

func (t RawTransaction) Inputs() CellInputFixVec {
	return CellInputFixVec(uncheckedField(t, 3, false))
}

func (t RawTransaction) Outputs() CellOutputDynVec {
	return CellOutputDynVec(uncheckedField(t, 4, false))
}

func (t RawTransaction) OutputsData() BytesDynVec {
	return BytesDynVec(uncheckedField(t, 5, true))
}

type Transaction []byte

func (t Transaction) Verify(compatible bool) bool {
	offsets, success := verifyAndExtractOffsets(t, 2, compatible)
	if !success {
		return false
	}
	if !RawTransaction(t[offsets[0]:offsets[1]]).Verify(compatible) {
		return false
	}
	if !BytesDynVec(t[offsets[1]:offsets[2]]).Verify(compatible) {
		return false
	}
	return true
}

func (t Transaction) RawTransaction() RawTransaction {
	return RawTransaction(uncheckedField(t, 0, false))
}

func (t Transaction) Witnesses() BytesDynVec {
	return BytesDynVec(uncheckedField(t, 1, true))
}

type RawHeader []byte

func (h RawHeader) Verify(compatible bool) bool {
	return len(h) == RawHeaderSize &&
		h.ParentHash().Verify(compatible) &&
		h.TransactionsRoot().Verify(compatible) &&
		h.ProposalsHash().Verify(compatible) &&
		h.UnclesHash().Verify(compatible) &&
		h.Dao().Verify(compatible)
}

func (h RawHeader) Version() uint32 {
	return binary.LittleEndian.Uint32(h[0:4])
}

func (h RawHeader) CompactTarget() uint32 {
	return binary.LittleEndian.Uint32(h[4:8])
}

func (h RawHeader) Timestamp() uint64 {
	return binary.LittleEndian.Uint64(h[8:16])
}

func (h RawHeader) Number() uint64 {
	return binary.LittleEndian.Uint64(h[16:24])
}

func (h RawHeader) Epoch() uint64 {
	return binary.LittleEndian.Uint64(h[24:32])
}

func (h RawHeader) ParentHash() Hash {
	return Hash(h[32:64])
}

func (h RawHeader) TransactionsRoot() Hash {
	return Hash(h[64:96])
}

func (h RawHeader) ProposalsHash() Hash {
	return Hash(h[96:128])
}

func (h RawHeader) UnclesHash() Hash {
	return Hash(h[128:160])
}

func (h RawHeader) Dao() Byte32 {
	return Byte32(h[160:196])
}

type Header []byte

func (h Header) Verify(compatible bool) bool {
	return len(h) == HeaderSize && h.RawHeader().Verify(compatible)
}

func (h Header) RawHeader() RawHeader {
	return RawHeader(h[0:192])
}

func (h Header) Nonce() *big.Int {
	b := make([]byte, 16)
	for i, v := range h[192:208] {
		b[16-i-1] = v
	}
	i := new(big.Int)
	return i.SetBytes(b)
}

type ProposalShortId []byte

func (p ProposalShortId) Verify(_compatible bool) bool {
	return len(p) == ProposalShortIdSize
}

type ProposalShortIdFixVec []byte

func (p ProposalShortIdFixVec) Verify(compatible bool) bool {
	if (len(p)) < 4 {
		return false
	}
	count := int(binary.LittleEndian.Uint32(p[0:4]))
	if len(p) != 4+ProposalShortIdSize*count {
		return false
	}
	for i := 0; i < count; i++ {
		if !p.Get(i).Verify(compatible) {
			return false
		}
	}
	return true
}

func (p ProposalShortIdFixVec) Len() int {
	return int(binary.LittleEndian.Uint32(p[0:4]))
}

func (p ProposalShortIdFixVec) Get(index int) ProposalShortId {
	start := 4 + index*ProposalShortIdSize
	return ProposalShortId(p[start : start+ProposalShortIdSize])
}

type UncleBlock []byte

func (b UncleBlock) Verify(compatible bool) bool {
	offsets, success := verifyAndExtractOffsets(b, 2, compatible)
	if !success {
		return false
	}
	if !Header(b[offsets[0]:offsets[1]]).Verify(compatible) {
		return false
	}
	if !ProposalShortIdFixVec(b[offsets[1]:offsets[2]]).Verify(compatible) {
		return false
	}
	return true
}

func (b UncleBlock) Header() Header {
	return Header(uncheckedField(b, 0, false))
}

func (b UncleBlock) Proposals() ProposalShortIdFixVec {
	return ProposalShortIdFixVec(uncheckedField(b, 1, true))
}

type UncleBlockDynVec []byte

func (b UncleBlockDynVec) Verify(compatible bool) bool {
	offsets, success := extractOffsets(b)
	if !success {
		return false
	}
	for i := 0; i < len(offsets)-1; i++ {
		if !UncleBlock(b[offsets[i]:offsets[i+1]]).Verify(compatible) {
			return false
		}
	}
	return true
}

func (b UncleBlockDynVec) Len() int {
	if len(b) < 8 {
		return 0
	} else {
		return int(binary.LittleEndian.Uint32(b[4:8]))/4 - 1
	}
}

func (b UncleBlockDynVec) Get(i int) UncleBlock {
	return UncleBlock(uncheckedField(b, i, true))
}

type TransactionDynVec []byte

func (t TransactionDynVec) Verify(compatible bool) bool {
	offsets, success := extractOffsets(t)
	if !success {
		return false
	}
	for i := 0; i < len(offsets)-1; i++ {
		if !Transaction(t[offsets[i]:offsets[i+1]]).Verify(compatible) {
			return false
		}
	}
	return true
}

func (t TransactionDynVec) Len() int {
	if len(t) < 8 {
		return 0
	} else {
		return int(binary.LittleEndian.Uint32(t[4:8]))/4 - 1
	}
}

func (t TransactionDynVec) Get(i int) Transaction {
	return Transaction(uncheckedField(t, i, true))
}

type Block []byte

func (b Block) Verify(compatible bool) bool {
	offsets, success := verifyAndExtractOffsets(b, 4, compatible)
	if !success {
		return false
	}
	if !Header(b[offsets[0]:offsets[1]]).Verify(compatible) {
		return false
	}
	if !UncleBlockDynVec(b[offsets[1]:offsets[2]]).Verify(compatible) {
		return false
	}
	if !TransactionDynVec(b[offsets[2]:offsets[3]]).Verify(compatible) {
		return false
	}
	if !ProposalShortIdFixVec(b[offsets[3]:offsets[4]]).Verify(compatible) {
		return false
	}
	return true
}

func (b Block) Header() Header {
	return Header(uncheckedField(b, 0, false))
}

func (b Block) Uncles() UncleBlockDynVec {
	return UncleBlockDynVec(uncheckedField(b, 1, false))
}

func (b Block) Transactions() TransactionDynVec {
	return TransactionDynVec(uncheckedField(b, 2, false))
}

func (b Block) Proposals() ProposalShortIdFixVec {
	return ProposalShortIdFixVec(uncheckedField(b, 3, true))
}

type CellbaseWitness []byte

func (w CellbaseWitness) Verify(compatible bool) bool {
	offsets, success := verifyAndExtractOffsets(w, 2, compatible)
	if !success {
		return false
	}
	if !Script(w[offsets[0]:offsets[1]]).Verify(compatible) {
		return false
	}
	if !Bytes(w[offsets[1]:offsets[2]]).Verify(compatible) {
		return false
	}
	return true
}

func (w CellbaseWitness) Lock() Script {
	return Script(uncheckedField(w, 0, false))
}

func (w CellbaseWitness) Message() Bytes {
	return Bytes(uncheckedField(w, 1, true))
}

type WitnessArgs []byte

func (w WitnessArgs) Verify(compatible bool) bool {
	offsets, success := verifyAndExtractOffsets(w, 3, compatible)
	if !success {
		return false
	}
	lockSlice := w[offsets[0]:offsets[1]]
	if len(lockSlice) > 0 {
		if !Bytes(lockSlice).Verify(compatible) {
			return false
		}
	}
	inputTypeSlice := w[offsets[1]:offsets[2]]
	if len(inputTypeSlice) > 0 {
		if !Bytes(inputTypeSlice).Verify(compatible) {
			return false
		}
	}
	outputTypeSlice := w[offsets[2]:offsets[3]]
	if len(outputTypeSlice) > 0 {
		if !Bytes(outputTypeSlice).Verify(compatible) {
			return false
		}
	}
	return true
}

func (w WitnessArgs) HasLock() bool {
	return len(uncheckedField(w, 0, false)) > 0
}

func (w WitnessArgs) Lock() Bytes {
	return Bytes(uncheckedField(w, 0, false))
}

func (w WitnessArgs) MaybeLock() *Bytes {
	field := uncheckedField(w, 0, false)
	if len(field) > 0 {
		s := Bytes(field)
		return &s
	} else {
		return nil
	}
}

func (w WitnessArgs) HasInputType() bool {
	return len(uncheckedField(w, 1, false)) > 0
}

func (w WitnessArgs) InputType() Bytes {
	return Bytes(uncheckedField(w, 1, false))
}

func (w WitnessArgs) MaybeInputType() *Bytes {
	field := uncheckedField(w, 1, false)
	if len(field) > 0 {
		s := Bytes(field)
		return &s
	} else {
		return nil
	}
}

func (w WitnessArgs) HasOutputType() bool {
	return len(uncheckedField(w, 2, true)) > 0
}

func (w WitnessArgs) OutputType() Bytes {
	return Bytes(uncheckedField(w, 2, true))
}

func (w WitnessArgs) MaybeOutputType() *Bytes {
	field := uncheckedField(w, 2, true)
	if len(field) > 0 {
		s := Bytes(field)
		return &s
	} else {
		return nil
	}
}
