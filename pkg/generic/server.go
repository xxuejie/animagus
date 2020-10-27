package generic

import (
	"context"
	"fmt"
	"io"

	"github.com/golang/protobuf/proto"
	"github.com/gomodule/redigo/redis"
	"github.com/xxuejie/animagus/pkg/ast"
	"github.com/xxuejie/animagus/pkg/coretypes"
	"github.com/xxuejie/animagus/pkg/executor"
	"github.com/xxuejie/animagus/pkg/indexer"
	"github.com/xxuejie/animagus/pkg/rpc"
	"github.com/xxuejie/animagus/pkg/rpctypes"
	"github.com/xxuejie/animagus/pkg/verifier"
)

type callInfo struct {
	expr    *ast.Value
	context indexer.ValueContext
}

type Server struct {
	calls     map[string]callInfo
	streams   []*ast.Stream
	redisPool *redis.Pool
	rpcClient *rpc.Client
}

func NewServer(astContent []byte, redisPool *redis.Pool, rpcUrl string) (*Server, error) {
	root := &ast.Root{}
	err := proto.Unmarshal(astContent, root)
	if err != nil {
		return nil, err
	}
	calls := make(map[string]callInfo)
	for _, call := range root.GetCalls() {
		err = verifier.Verify(call.GetResult())
		if err != nil {
			return nil, fmt.Errorf("Verification failure for call %s: %s", call.GetName(), err)
		}
		valueContext, err := indexer.NewValueContext(call.GetName(), call.GetResult())
		if err != nil {
			return nil, err
		}
		calls[call.GetName()] = callInfo{
			expr:    call.GetResult(),
			context: valueContext,
		}
	}
	for _, stream := range root.GetStreams() {
		err = verifier.Verify(stream.GetFilter())
		if err != nil {
			return nil, fmt.Errorf("Verification failure for stream %s: %s", stream.GetName(), err)
		}
	}
	client := rpc.NewClient(rpcUrl)
	return &Server{
		calls:     calls,
		streams:   root.GetStreams(),
		redisPool: redisPool,
		rpcClient: client,
	}, nil
}

type executeEnvironment struct {
	params       *GenericParams
	valueContext indexer.ValueContext
	s            *Server
}

func (e executeEnvironment) ReplaceArgs(args []*ast.Value) error {
	if len(args) > 0 {
		return fmt.Errorf("No arguments provided for replacement!")
	}
	return nil
}

func (e executeEnvironment) Arg(i int) *ast.Value {
	return nil
}

func (e executeEnvironment) Param(i int) *ast.Value {
	if i < 0 || i >= len(e.params.GetParams()) {
		return nil
	}
	return e.params.GetParams()[i]
}

func (e executeEnvironment) IndexParam(i int, value *ast.Value) error {
	return fmt.Errorf("Indexing param is not allowed when executing!")
}

type getCellsResponse struct {
	GetCells []*rpctypes.OutPoint
}

func (s *Server) getCells(coreOutPoints []coretypes.OutPoint) ([]*rpctypes.OutPoint, error) {
	var rpcOutPoints []*rpctypes.OutPoint
	set := make(map[rpctypes.Hash]int)

	for _, coreOutPoint := range coreOutPoints {
		txHash := []byte(coreOutPoint.TxHash())
		var rpcTxHash rpctypes.Hash
		copy(rpcTxHash[:], txHash)
		rpcIndex := rpctypes.Uint32(coreOutPoint.Index())
		rpcOutPoint := rpctypes.OutPoint{
			TxHash: rpcTxHash,
			Index:  rpcIndex,
		}
		rpcOutPoints = append(rpcOutPoints, &rpcOutPoint)
		set[rpcTxHash] = 1
	}

	var txHashes []rpctypes.Hash
	for key := range set {
		txHashes = append(txHashes, key)
	}

	transactionWithStatusViews, err := s.rpcClient.GetAllTransactions(txHashes, 50)
	if err != nil {
		return nil, err
	}

	txWithStatusMap := make(map[rpctypes.Hash]*rpctypes.TransactionWithStatusView)
	blockHashSet := make(map[rpctypes.Hash]int)
	for _, txWithStatusView := range transactionWithStatusViews {
		txWithStatusMap[txWithStatusView.Transaction.Hash] = txWithStatusView
		blockHashSet[*txWithStatusView.TxStatus.BlockHash] = 1
	}
	var blockHashes []rpctypes.Hash
	for key := range blockHashSet {
		blockHashes = append(blockHashes, key)
	}
	headers, err := s.rpcClient.GetAllHeaders(blockHashes, 50)
	if err != nil {
		return nil, err
	}
	headerMap := make(map[rpctypes.Hash]*rpctypes.Header)
	for _, header := range headers {
		headerMap[header.Hash] = &header.Header
	}

	var result []*rpctypes.OutPoint
	for _, outPoint := range rpcOutPoints {
		// get transaction
		transactionWithStatus := txWithStatusMap[outPoint.TxHash]
		transactionView := &transactionWithStatus.Transaction

		index := outPoint.Index
		cell := transactionView.Outputs[index]
		originData := &transactionView.OutputsData[index]
		raw := rpctypes.Raw([]byte(*originData))

		// get header
		header := headerMap[*transactionWithStatus.TxStatus.BlockHash]
		resultOutPoint := rpctypes.OutPoint{
			TxHash:   outPoint.TxHash,
			Index:    outPoint.Index,
			Cell:     &cell,
			CellData: &raw,
			Header:   header,
		}
		result = append(result, &resultOutPoint)
	}
	return result, nil
}

func (e executeEnvironment) QueryCell(query *ast.Value) ([]*ast.Value, error) {
	queryIndex := e.valueContext.QueryIndex(query)
	if queryIndex == -1 {
		return nil, fmt.Errorf("Invalid query cell argument!")
	}
	paramValues := make(map[int]*ast.Value)
	for i, value := range e.params.GetParams() {
		paramValues[i] = value
	}
	indexKey, err := e.valueContext.IndexKey(queryIndex, paramValues)
	if err != nil {
		return nil, err
	}
	conn := e.s.redisPool.Get()
	defer conn.Close()
	slices, err := redis.ByteSlices(conn.Do("SMEMBERS", indexKey))
	if err != nil {
		return nil, err
	}
	if len(slices) == 0 {
		return []*ast.Value{}, nil
	}
	outPoints := make([]coretypes.OutPoint, len(slices))
	for i, slice := range slices {
		outPoints[i] = coretypes.OutPoint(slice)
		if !outPoints[i].Verify(true) {
			return nil, fmt.Errorf("OutPoint %x verification failure!", slice)
		}
	}
	resultCells, err := e.s.getCells(outPoints)
	if err != nil {
		return nil, err
	}
	results := make([]*ast.Value, len(resultCells))
	for i, cell := range resultCells {
		results[i] = ast.ConvertCell(*cell.Cell, *cell.CellData, *cell, cell.Header)
	}
	return results, nil
}

func (s *Server) Call(ctx context.Context, p *GenericParams) (*ast.Value, error) {
	callInfo, found := s.calls[p.GetName()]
	if !found {
		return nil, fmt.Errorf("Calling non-exist function: %s", p.GetName())
	}
	environment := executeEnvironment{
		params:       p,
		valueContext: callInfo.context,
		s:            s,
	}
	return executor.Execute(callInfo.expr, environment)
}

func (s *Server) Stream(p *GenericParams, streamServer GenericService_StreamServer) error {
	var selectedStream *ast.Stream
	for _, aStream := range s.streams {
		if p.GetName() == aStream.GetName() {
			selectedStream = aStream
			break
		}
	}
	if selectedStream == nil {
		return fmt.Errorf("Calling non-exist stream: %s", p.GetName())
	}

	psc := redis.PubSubConn{Conn: s.redisPool.Get()}
	defer psc.Close()

	key := fmt.Sprintf("STREAM:%s", selectedStream.GetName())
	err := psc.Subscribe(key)
	if err != nil {
		return err
	}
	for {
		switch v := psc.Receive().(type) {
		case redis.Message:
			if v.Channel == key {
				value := &ast.Value{}
				err = proto.Unmarshal(v.Data, value)
				if err != nil {
					break
				}
				err = streamServer.Send(value)
				if err == io.EOF {
					return psc.Unsubscribe()
				}
				if err != nil {
					break
				}
			}
		case error:
			return v
		}
	}
	psc.Unsubscribe()
	return err
}
