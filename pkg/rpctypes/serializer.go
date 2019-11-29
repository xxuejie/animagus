package rpctypes

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

type errorWriter struct {
	writer io.Writer
	err    error
}

func newErrorWriter(writer io.Writer) errorWriter {
	return errorWriter{
		writer: writer,
		err:    nil,
	}
}

func (w *errorWriter) Write(p []byte) (int, error) {
	if w.err != nil {
		return 0, w.err
	}
	var n int
	n, w.err = w.writer.Write(p)
	return n, w.err
}

func (w *errorWriter) writeLittleEndianUint32(u uint32) {
	if w.err != nil {
		return
	}
	w.err = binary.Write(w.writer, binary.LittleEndian, u)
}

type CoreSerializer interface {
	SerializeToCore(writer io.Writer) error
}

func (b Bytes) SerializeToCore(writer io.Writer) error {
	_, err := writer.Write(b[:])
	return err
}

func (b VarBytes) SerializeToCore(writer io.Writer) error {
	err := binary.Write(writer, binary.LittleEndian, uint32(len(b)))
	if err != nil {
		return err
	}
	_, err = writer.Write(b[:])
	return err
}

func (u Uint128) SerializeToCore(writer io.Writer) error {
	intBytes := u.V.Bytes()
	if len(intBytes) > 16 {
		return fmt.Errorf("Uint128 value is too big!")
	}
	data := make([]byte, 16)
	for i, b := range intBytes {
		data[len(intBytes)-i-1] = b
	}
	_, err := writer.Write(data)
	return err
}

func (u Uint64) SerializeToCore(writer io.Writer) error {
	return binary.Write(writer, binary.LittleEndian, u)
}

func (u Uint32) SerializeToCore(writer io.Writer) error {
	return binary.Write(writer, binary.LittleEndian, u)
}

func (h Hash) SerializeToCore(writer io.Writer) error {
	_, err := writer.Write(h[:])
	return err
}

func (h ProposalShortId) SerializeToCore(writer io.Writer) error {
	_, err := writer.Write(h[:])
	return err
}

func (s ScriptHashType) SerializeToCore(writer io.Writer) error {
	_, err := writer.Write([]byte{byte(s)})
	return err
}

func (d DepType) SerializeToCore(writer io.Writer) error {
	_, err := writer.Write([]byte{byte(d)})
	return err
}

func (o OutPoint) SerializeToCore(writer io.Writer) error {
	err := o.TxHash.SerializeToCore(writer)
	if err != nil {
		return err
	}
	return o.Index.SerializeToCore(writer)
}

func (c CellInput) SerializeToCore(writer io.Writer) error {
	err := c.Since.SerializeToCore(writer)
	if err != nil {
		return err
	}
	return c.PreviousOutput.SerializeToCore(writer)
}

func serializeTable(items []CoreSerializer, writer io.Writer) error {
	item_count := len(items)
	total_size := 4 * (item_count + 1)
	buffers := make([]bytes.Buffer, item_count)
	offsets := make([]int, item_count)

	for i := 0; i < len(offsets); i++ {
		offsets[i] = total_size
		err := items[i].SerializeToCore(&buffers[i])
		if err != nil {
			return err
		}
		total_size += buffers[i].Len()
	}

	errorWriter := newErrorWriter(writer)
	errorWriter.writeLittleEndianUint32(uint32(total_size))
	for _, v := range offsets {
		errorWriter.writeLittleEndianUint32(uint32(v))
	}
	for _, b := range buffers {
		errorWriter.Write(b.Bytes())
	}
	return errorWriter.err
}

func (s Script) SerializeToCore(writer io.Writer) error {
	return serializeTable([]CoreSerializer{
		s.CodeHash,
		s.HashType,
		s.Args,
	}, writer)
}

func (c CellOutput) SerializeToCore(writer io.Writer) error {
	items := make([]CoreSerializer, 3)
	items[0] = c.Capacity
	items[1] = c.Lock
	if c.Type != nil {
		items[2] = *c.Type
	} else {
		// Empty slice here
		items[2] = Bytes([]byte{})
	}
	return serializeTable(items, writer)
}

func (c CellDep) SerializeToCore(writer io.Writer) error {
	err := c.OutPoint.SerializeToCore(writer)
	if err != nil {
		return err
	}
	return c.DepType.SerializeToCore(writer)
}

type cellDepFixVec []CellDep

func (v cellDepFixVec) SerializeToCore(writer io.Writer) error {
	errorWriter := newErrorWriter(writer)
	errorWriter.writeLittleEndianUint32(uint32(len(v)))
	for _, item := range v {
		item.SerializeToCore(&errorWriter)
	}
	return errorWriter.err
}

type hashFixVec []Hash

func (v hashFixVec) SerializeToCore(writer io.Writer) error {
	errorWriter := newErrorWriter(writer)
	errorWriter.writeLittleEndianUint32(uint32(len(v)))
	for _, item := range v {
		item.SerializeToCore(&errorWriter)
	}
	return errorWriter.err
}

type cellInputFixVec []CellInput

func (v cellInputFixVec) SerializeToCore(writer io.Writer) error {
	errorWriter := newErrorWriter(writer)
	errorWriter.writeLittleEndianUint32(uint32(len(v)))
	for _, item := range v {
		item.SerializeToCore(&errorWriter)
	}
	return errorWriter.err
}

type cellOutputDynVec []CellOutput

func (v cellOutputDynVec) SerializeToCore(writer io.Writer) error {
	items := make([]CoreSerializer, len(v))
	for i, item := range v {
		items[i] = item
	}
	return serializeTable(items, writer)
}

type varbytesDynVec []VarBytes

func (v varbytesDynVec) SerializeToCore(writer io.Writer) error {
	items := make([]CoreSerializer, len(v))
	for i, item := range v {
		items[i] = item
	}
	return serializeTable(items, writer)
}

func (t RawTransaction) SerializeToCore(writer io.Writer) error {
	buffers := make([]bytes.Buffer, 6)
	err := t.Version.SerializeToCore(&buffers[0])
	if err != nil {
		return err
	}
	err = cellDepFixVec(t.CellDeps).SerializeToCore(&buffers[1])
	if err != nil {
		return err
	}
	err = hashFixVec(t.HeaderDeps).SerializeToCore(&buffers[2])
	if err != nil {
		return err
	}
	err = cellInputFixVec(t.Inputs).SerializeToCore(&buffers[3])
	if err != nil {
		return err
	}
	err = cellOutputDynVec(t.Outputs).SerializeToCore(&buffers[4])
	if err != nil {
		return err
	}
	err = varbytesDynVec(t.OutputsData).SerializeToCore(&buffers[5])
	if err != nil {
		return err
	}
	items := make([]CoreSerializer, len(buffers))
	for i, b := range buffers {
		items[i] = Bytes(b.Bytes())
	}
	return serializeTable(items, writer)
}

func (t Transaction) SerializeToCore(writer io.Writer) error {
	var witnessBuffer bytes.Buffer
	err := varbytesDynVec(t.Witnesses).SerializeToCore(&witnessBuffer)
	if err != nil {
		return err
	}
	return serializeTable([]CoreSerializer{
		t.RawTransaction,
		Bytes(witnessBuffer.Bytes()),
	}, writer)
}

func (h RawHeader) SerializeToCore(writer io.Writer) error {
	errorWriter := newErrorWriter(writer)
	h.Version.SerializeToCore(&errorWriter)
	h.CompactTarget.SerializeToCore(&errorWriter)
	h.Timestamp.SerializeToCore(&errorWriter)
	h.Number.SerializeToCore(&errorWriter)
	h.Epoch.SerializeToCore(&errorWriter)
	h.ParentHash.SerializeToCore(&errorWriter)
	h.TransactionsRoot.SerializeToCore(&errorWriter)
	h.ProposalsHash.SerializeToCore(&errorWriter)
	h.UnclesHash.SerializeToCore(&errorWriter)
	h.Dao.SerializeToCore(&errorWriter)
	return errorWriter.err
}

func (h Header) SerializeToCore(writer io.Writer) error {
	err := h.RawHeader.SerializeToCore(writer)
	if err != nil {
		return err
	}
	return h.Nonce.SerializeToCore(writer)
}

type proposalShortIdFixVec []ProposalShortId

func (v proposalShortIdFixVec) SerializeToCore(writer io.Writer) error {
	errorWriter := newErrorWriter(writer)
	errorWriter.writeLittleEndianUint32(uint32(len(v)))
	for _, item := range v {
		item.SerializeToCore(&errorWriter)
	}
	return errorWriter.err
}

func (b UncleBlock) SerializeToCore(writer io.Writer) error {
	var shortIdBuffer bytes.Buffer
	err := proposalShortIdFixVec(b.Proposals).SerializeToCore(&shortIdBuffer)
	if err != nil {
		return err
	}
	return serializeTable([]CoreSerializer{
		b.Header,
		Bytes(shortIdBuffer.Bytes()),
	}, writer)
}

type uncleBlockDynVec []UncleBlock

func (v uncleBlockDynVec) SerializeToCore(writer io.Writer) error {
	items := make([]CoreSerializer, len(v))
	for i, item := range v {
		items[i] = item
	}
	return serializeTable(items, writer)
}

type transactionDynVec []Transaction

func (v transactionDynVec) SerializeToCore(writer io.Writer) error {
	items := make([]CoreSerializer, len(v))
	for i, item := range v {
		items[i] = item
	}
	return serializeTable(items, writer)
}

func (b Block) SerializeToCore(writer io.Writer) error {
	buffers := make([]bytes.Buffer, 4)
	err := b.Header.SerializeToCore(&buffers[0])
	if err != nil {
		return err
	}
	err = uncleBlockDynVec(b.Uncles).SerializeToCore(&buffers[1])
	if err != nil {
		return err
	}
	err = transactionDynVec(b.Transactions).SerializeToCore(&buffers[2])
	if err != nil {
		return err
	}
	err = proposalShortIdFixVec(b.Proposals).SerializeToCore(&buffers[3])
	if err != nil {
		return err
	}
	items := make([]CoreSerializer, len(buffers))
	for i, b := range buffers {
		items[i] = Bytes(b.Bytes())
	}
	return serializeTable(items, writer)
}

func (c CellbaseWitness) SerializeToCore(writer io.Writer) error {
	return serializeTable([]CoreSerializer{
		c.Lock,
		c.Message,
	}, writer)
}

func (c WitnessArgs) SerializeToCore(writer io.Writer) error {
	items := make([]CoreSerializer, 3)
	if c.Lock != nil {
		items[0] = c.Lock
	} else {
		items[0] = Bytes([]byte{})
	}
	if c.InputType != nil {
		items[1] = c.InputType
	} else {
		items[1] = Bytes([]byte{})
	}
	if c.OutputType != nil {
		items[2] = c.OutputType
	} else {
		items[2] = Bytes([]byte{})
	}
	return serializeTable(items, writer)
}
