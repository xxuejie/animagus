package indexer

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/gomodule/redigo/redis"
	"github.com/machinebox/graphql"
	blake2b "github.com/minio/blake2b-simd"
	"github.com/xxuejie/animagus/pkg/ast"
	"github.com/xxuejie/animagus/pkg/executor"
	"github.com/xxuejie/animagus/pkg/rpctypes"
)

const Version string = "0.0.1"

type Indexer struct {
	hash          []byte
	values        []ValueContext
	streams       []*ast.Stream
	redisPool     *redis.Pool
	graphqlClient *graphql.Client
}

func NewIndexer(astContent []byte, redisPool *redis.Pool, graphqlUrl string) (*Indexer, error) {
	root := &ast.Root{}
	err := proto.Unmarshal(astContent, root)
	if err != nil {
		return nil, err
	}
	blake2bHash := blake2b.New256()
	_, err = blake2bHash.Write(astContent)
	if err != nil {
		return nil, err
	}
	_, err = blake2bHash.Write([]byte(Version))
	if err != nil {
		return nil, err
	}
	hash := blake2bHash.Sum(nil)
	values := make([]ValueContext, len(root.GetCalls()))
	for i, call := range root.GetCalls() {
		valueContext, err := NewValueContext(call.GetName(), call.GetResult())
		if err != nil {
			return nil, err
		}
		values[i] = valueContext
	}
	// Test GraphQL query
	client := graphql.NewClient(graphqlUrl)
	// client.Log = func(s string) {
	// 	fmt.Printf("GraphQL log: %s\n", s)
	// }
	err = client.Run(context.Background(), graphql.NewRequest(`
query {
  apiVersion
}
`), nil)
	if err != nil {
		return nil, err
	}
	return &Indexer{
		values:        values,
		hash:          hash,
		redisPool:     redisPool,
		graphqlClient: client,
		streams:       root.GetStreams(),
	}, nil
}

type getBlockResponse struct {
	GetBlock *rpctypes.BlockView
}

func (i *Indexer) Run() error {
	redisConn := i.redisPool.Get()
	defer redisConn.Close()

	dbHash, err := redis.Bytes(redisConn.Do("GET", "AST_HASH"))
	if err != nil && err != redis.ErrNil {
		return err
	}
	if len(dbHash) == 0 {
		dbHash = i.hash
		_, err = redisConn.Do("SET", "AST_HASH", dbHash)
		if err != nil {
			return err
		}
	}
	if !bytes.Equal(dbHash, i.hash) {
		return fmt.Errorf("Invalid AST Hash: %x, expected: %x", dbHash, i.hash)
	}
	for {
		var blockToFetch uint64
		var lastBlockHash []byte
		lastBlock, err := redis.Bytes(redisConn.Do("GET", "LAST_BLOCK"))
		if err != nil && err != redis.ErrNil {
			return err
		}
		if len(lastBlock) == 40 {
			lastBlockNumber := binary.LittleEndian.Uint64(lastBlock)
			blockToFetch = lastBlockNumber + 1
			lastBlockHash = lastBlock[8:]
		}

		block, err := i.queryBlock(blockToFetch)
		if err != nil {
			return err
		}
		if block == nil {
			time.Sleep(time.Second)
			continue
		}

		revert := lastBlockHash != nil && (!bytes.Equal(block.Header.ParentHash[:], lastBlockHash))
		if revert {
			if blockToFetch == 0 {
				return fmt.Errorf("Nowhere to revert!")
			}
			err = i.revertBlock(blockToFetch - 1)
			if err != nil {
				return err
			}
			log.Printf("Reverted block number %d", blockToFetch-1)
			continue
		}

		// In the optimal path, we keep only one Redis connection
		commands := &commandBuffer{}
		err = i.indexBlock(*block, commands)
		if err != nil {
			return err
		}
		err = commands.execute(redisConn)
		if err != nil {
			return err
		}
		log.Printf("Indexed block %x, block number %d", block.Header.Hash, block.Header.Number)
	}
}

func (i *Indexer) queryBlock(blockNumber uint64) (*rpctypes.BlockView, error) {
	req := graphql.NewRequest(`
query($blockNumber: String) {
  getBlock(number: $blockNumber) {
    header {
      parent_hash
      hash
      number
    }
    transactions {
      hash
      inputs {
        previous_output {
          cell {
            capacity
            lock {
              code_hash
              hash_type
              args
            }
            type {
              code_hash
              hash_type
              args
            }
          }
          cell_data {
            content
          }
          tx_hash
          index
        }
      }
      outputs {
        capacity
        lock {
          code_hash
          hash_type
          args
        }
        type {
          code_hash
          hash_type
          args
        }
      }
      cells_data {
        content
      }
    }
  }
}
`)
	req.Var("blockNumber", rpctypes.Uint64(blockNumber).EncodeToString())
	var response getBlockResponse
	err := i.graphqlClient.Run(context.Background(), req, &response)
	if err != nil {
		return nil, err
	}
	return response.GetBlock, nil
}

func (i *Indexer) indexBlock(block rpctypes.BlockView, commands *commandBuffer) error {
	var err error
	for _, tx := range block.Transactions {
		for _, input := range tx.RawTransaction.Inputs {
			if input.PreviousOutput.GraphqlCell != nil &&
				input.PreviousOutput.GraphqlCellData != nil {
				err = i.processCell(*input.PreviousOutput.GraphqlCell,
					*input.PreviousOutput.GraphqlCellData.Content,
					input.PreviousOutput,
					false,
					commands)
				if err != nil {
					return err
				}
			}
		}

		for outputIndex, output := range tx.RawTransaction.Outputs {
			err = i.processCell(output,
				*tx.RawTransaction.GraphqlCellsData[outputIndex].Content,
				rpctypes.OutPoint{
					TxHash: tx.Hash,
					Index:  rpctypes.Uint32(outputIndex),
				},
				true,
				commands)
			if err != nil {
				return err
			}
		}
	}
	blockNumber := uint64(block.Header.Number)
	blockHashKey := fmt.Sprintf("BLOCK:%d:HASH", blockNumber)
	commands.do("SET", blockHashKey, block.Header.Hash[:])
	lastBlock := make([]byte, 40)
	binary.LittleEndian.PutUint64(lastBlock, blockNumber)
	copy(lastBlock[8:], block.Header.Hash[:])
	commands.do("SET", "LAST_BLOCK", lastBlock)

	revertKey := fmt.Sprintf("BLOCK:%d:REVERT_COMMANDS", blockNumber)
	commands.setRevertKey(revertKey)
	commands.revertDo("DEL", blockHashKey)
	if blockNumber > 0 {
		previousBlock := make([]byte, 40)
		binary.LittleEndian.PutUint64(previousBlock, blockNumber-1)
		copy(previousBlock[8:], block.Header.ParentHash[:])
		commands.revertDo("SET", "LAST_BLOCK", previousBlock)
	} else {
		commands.revertDo("DEL", "LAST_BLOCK")
	}
	commands.revertDo("DEL", revertKey)

	return nil
}

func (i *Indexer) revertBlock(blockNumber uint64) error {
	conn := i.redisPool.Get()
	defer conn.Close()

	revertKey := fmt.Sprintf("BLOCK:%d:REVERT_COMMANDS", blockNumber)
	revertData, err := redis.Bytes(conn.Do("GET", revertKey))
	if err != nil {
		return err
	}

	var revertCommands []command
	gzipReader, err := gzip.NewReader(bytes.NewReader(revertData))
	if err != nil {
		return err
	}
	decoder := json.NewDecoder(gzipReader)
	err = decoder.Decode(&revertCommands)
	if err != nil {
		return err
	}

	conn.Send("MULTI")
	for _, command := range revertCommands {
		conn.Send(command.Name, command.Args...)
	}
	_, err = conn.Do("EXEC")
	return err
}

func (i *Indexer) processCell(cell rpctypes.CellOutput, cellData rpctypes.Raw, outPoint rpctypes.OutPoint, insert bool, commands *commandBuffer) error {
	// TODO: To maintain reasonable cell set in CKB, it might not be possible to grab
	// meta-data of very old spent cells. Hence for now, we are excluding all cell
	// headers in indexer mode, and only include headers when executing an AST
	// (since executor only uses live cells). A different solution might be that
	// animagus can do its own indexing to cache all header info, but that will
	// be a quite big change so we will leave it till a future time when it is
	// really needed.
	astCell := ast.ConvertCell(cell, cellData, outPoint, nil)
	for _, valueContext := range i.values {
		for queryIndex, query := range valueContext.Queries {
			indexedValues, err := executeIndexingQuery(query, astCell)
			if err != nil {
				return err
			}
			if indexedValues != nil {
				key, err := valueContext.IndexKey(queryIndex, indexedValues)
				if err != nil {
					return err
				}
				if insert {
					commands.insert(key, outPoint)
				} else {
					commands.remove(key, outPoint)
				}
			}
		}
	}
	for _, stream := range i.streams {
		value, err := executeStreamingFilter(stream.GetFilter(), astCell, insert, true)
		if err != nil {
			return err
		}
		if value != nil {
			commands.streamValue(stream.GetName(), value)
		}
		// Prepare revert value
		value, err = executeStreamingFilter(stream.GetFilter(), astCell, !insert, false)
		if err != nil {
			return err
		}
		if value != nil {
			commands.revertStreamValue(stream.GetName(), value)
		}
	}
	return nil
}

type command struct {
	Name string        `json:"n"`
	Args []interface{} `json:"a"`
}

type commandBuffer struct {
	commands       []command
	revertCommands []command
	// Those are kept separated since they will be reversed.
	streamRevertCommands []command
	revertKey            string
	err                  error
}

func (c *commandBuffer) do(commandName string, args ...interface{}) {
	if c.err != nil {
		return
	}
	c.commands = append(c.commands, command{
		Name: commandName,
		Args: args,
	})
}

func (c *commandBuffer) revertDo(commandName string, args ...interface{}) {
	if c.err != nil {
		return
	}
	c.revertCommands = append(c.revertCommands, command{
		Name: commandName,
		Args: args,
	})
}

func (c *commandBuffer) streamRevertDo(commandName string, args ...interface{}) {
	if c.err != nil {
		return
	}
	c.streamRevertCommands = append(c.streamRevertCommands, command{
		Name: commandName,
		Args: args,
	})
}

func (c *commandBuffer) setRevertKey(key string) {
	if c.err != nil {
		return
	}
	c.revertKey = key
}

func (c *commandBuffer) insert(key string, outPoint rpctypes.OutPoint) {
	if c.err != nil {
		return
	}
	var buffer bytes.Buffer
	c.err = outPoint.SerializeToCore(&buffer)
	c.do("SADD", key, buffer.Bytes())
	c.revertDo("SREM", key, buffer.Bytes())
}

func (c *commandBuffer) remove(key string, outPoint rpctypes.OutPoint) {
	if c.err != nil {
		return
	}
	var buffer bytes.Buffer
	c.err = outPoint.SerializeToCore(&buffer)
	c.do("SREM", key, buffer.Bytes())
	c.revertDo("SADD", key, buffer.Bytes())
}

func (c *commandBuffer) streamValue(name string, value []byte) {
	if c.err != nil {
		return
	}
	key := fmt.Sprintf("STREAM:%s", name)
	c.do("PUBLISH", key, value)
}

func (c *commandBuffer) revertStreamValue(name string, value []byte) {
	if c.err != nil {
		return
	}
	key := fmt.Sprintf("STREAM:%s", name)
	c.streamRevertDo("PUBLISH", key, value)
}

func (c *commandBuffer) execute(conn redis.Conn) error {
	if c.err != nil {
		return c.err
	}

	if len(c.revertKey) == 0 {
		return fmt.Errorf("Revert key is missing!")
	}
	var buf bytes.Buffer
	gzipWriter := gzip.NewWriter(&buf)
	encoder := json.NewEncoder(gzipWriter)
	revertCommands := make([]command, len(c.revertCommands)+len(c.streamRevertCommands))
	copy(revertCommands, c.revertCommands)
	for i, c := range c.streamRevertCommands {
		revertCommands[len(revertCommands)-1-i] = c
	}
	err := encoder.Encode(c.revertCommands)
	if err != nil {
		return err
	}
	err = gzipWriter.Close()
	if err != nil {
		return err
	}

	conn.Send("MULTI")
	for _, command := range c.commands {
		conn.Send(command.Name, command.Args...)
	}
	conn.Send("SET", c.revertKey, buf.Bytes())
	_, err = conn.Do("EXEC")
	return err
}

type indexingEnvironment struct {
	cell          *ast.Value
	indexedValues map[int]*ast.Value
}

func (e *indexingEnvironment) Arg(i int) *ast.Value {
	if i == 0 {
		return e.cell
	} else {
		return nil
	}
}

func (e *indexingEnvironment) Param(i int) *ast.Value {
	return &ast.Value{
		T: ast.Value_PARAM,
		Primitive: &ast.Value_U{
			U: uint64(i),
		},
	}
}

func (e *indexingEnvironment) IndexParam(i int, value *ast.Value) error {
	_, found := e.indexedValues[i]
	if found {
		return fmt.Errorf("Param %d is already indexed!", i)
	}
	e.indexedValues[i] = value
	return nil
}

func (e *indexingEnvironment) QueryCell(query *ast.Value) ([]*ast.Value, error) {
	return nil, fmt.Errorf("QueryCell is not allowed in indexer!")
}

func executeIndexingQuery(query *ast.Value, cell *ast.Value) (map[int]*ast.Value, error) {
	if len(query.GetChildren()) != 1 {
		return nil, fmt.Errorf("Invalid number of values to query cell: %d", len(query.GetChildren()))
	}
	environment := &indexingEnvironment{
		cell:          cell,
		indexedValues: make(map[int]*ast.Value),
	}
	value, err := executor.Execute(query.GetChildren()[0], environment)
	if err != nil {
		return nil, err
	}
	if value.GetT() != ast.Value_BOOL {
		return nil, fmt.Errorf("Invalid result value type: %s", value.GetT().String())
	}
	if !value.GetB() {
		return nil, nil
	}
	return environment.indexedValues, nil
}

type streamExecutingEnvironment struct {
	args []*ast.Value
}

func (e streamExecutingEnvironment) Arg(i int) *ast.Value {
	if i < 0 || i >= len(e.args) {
		return nil
	}
	return e.args[i]
}

func (e streamExecutingEnvironment) Param(i int) *ast.Value {
	return nil
}

func (e streamExecutingEnvironment) IndexParam(i int, value *ast.Value) error {
	return fmt.Errorf("Indexing param is not allowed!")
}

func (e streamExecutingEnvironment) QueryCell(query *ast.Value) ([]*ast.Value, error) {
	return nil, fmt.Errorf("Querying cell is not allowed!")
}

func executeStreamingFilter(filter *ast.Value, cell *ast.Value, insert bool, index bool) ([]byte, error) {
	var t string
	if insert {
		t = "insert"
	} else {
		t = "remove"
	}
	var i string
	if index {
		i = "index"
	} else {
		i = "revert"
	}
	e := streamExecutingEnvironment{
		args: []*ast.Value{
			cell,
			&ast.Value{
				T: ast.Value_BYTES,
				Primitive: &ast.Value_Raw{
					Raw: []byte(t),
				},
			},
			&ast.Value{
				T: ast.Value_BYTES,
				Primitive: &ast.Value_Raw{
					Raw: []byte(i),
				},
			},
		},
	}
	value, err := executor.Execute(filter, e)
	if err != nil {
		return nil, err
	}
	if value.GetT() == ast.Value_NIL {
		return nil, nil
	}
	return proto.Marshal(value)
}
