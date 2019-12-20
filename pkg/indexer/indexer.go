package indexer

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"log"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/gomodule/redigo/redis"
	"github.com/machinebox/graphql"
	blake2b "github.com/minio/blake2b-simd"
	internal_ast "github.com/xxuejie/animagus/internal/ast"
	"github.com/xxuejie/animagus/pkg/ast"
	"github.com/xxuejie/animagus/pkg/executor"
	"github.com/xxuejie/animagus/pkg/rpctypes"
)

const Version string = "0.0.1"

type Indexer struct {
	hash          []byte
	values        []ValueContext
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
	}, nil
}

type getBlockResponse struct {
	GetBlock *rpctypes.Block
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

		req := graphql.NewRequest(`
query($blockNumber: String) {
  getBlock(number: $blockNumber) {
    header {
      parent_hash
      hash
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
          }
          cell_data {
            content
          }
        }
      }
      outputs {
        capacity
        lock {
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
		req.Var("blockNumber", rpctypes.Uint64(blockToFetch).EncodeToString())
		var response getBlockResponse
		err = i.graphqlClient.Run(context.Background(), req, &response)
		if err != nil {
			return err
		}
		if response.GetBlock == nil {
			time.Sleep(time.Second)
			continue
		}
		block := *response.GetBlock

		if lastBlockHash != nil && (!bytes.Equal(block.Header.ParentHash[:], lastBlockHash)) {
			return fmt.Errorf("TODO: implement forking")
		}

		commands := &commandBuffer{}
		for _, tx := range block.Transactions {
			for _, input := range tx.RawTransaction.Inputs {
				if input.PreviousOutput.GraphqlCell != nil &&
					input.PreviousOutput.GraphqlCellData != nil {
					for _, valueContext := range i.values {
						for queryIndex, query := range valueContext.Queries {
							// TODO: cell data support
							params, err := executeIndexingQuery(query,
								input.PreviousOutput.GraphqlCell,
								input.PreviousOutput.GraphqlCellData.Content)
							if err != nil {
								return err
							}
							if params != nil {
								key, err := valueContext.IndexKey(queryIndex, params)
								if err != nil {
									return nil
								}
								commands.remove(key, input.PreviousOutput)
							}
						}
					}
				}
			}

			for outputIndex, output := range tx.RawTransaction.Outputs {
				for _, valueContext := range i.values {
					for queryIndex, query := range valueContext.Queries {
						params, err := executeIndexingQuery(query, &output,
							tx.RawTransaction.GraphqlCellsData[outputIndex].Content)
						if err != nil {
							return err
						}
						if params != nil {
							key, err := valueContext.IndexKey(queryIndex, params)
							if err != nil {
								return err
							}
							commands.insert(key, rpctypes.OutPoint{
								TxHash: *tx.RawTransaction.GraphqlHash,
								Index:  rpctypes.Uint32(outputIndex),
							})
						}
					}
				}
			}
		}
		commands.do("SET", fmt.Sprintf("BLOCK:%d:HASH", blockToFetch), block.Header.GraphqlHash[:])
		lastBlock = make([]byte, 40)
		binary.LittleEndian.PutUint64(lastBlock, blockToFetch)
		copy(lastBlock[8:], block.Header.GraphqlHash[:])
		commands.do("SET", "LAST_BLOCK", lastBlock)
		err = commands.execute(redisConn)
		if err != nil {
			return err
		}
		log.Printf("Indexed block %x, block number %d", *block.Header.GraphqlHash, blockToFetch)
	}
}

type command struct {
	name string
	args []interface{}
}

type commandBuffer struct {
	commands []command
	err      error
}

func (c *commandBuffer) do(commandName string, args ...interface{}) {
	if c.err != nil {
		return
	}
	c.commands = append(c.commands, command{
		name: commandName,
		args: args,
	})
}

func (c *commandBuffer) insert(key string, outPoint rpctypes.OutPoint) {
	if c.err != nil {
		return
	}
	var buffer bytes.Buffer
	c.err = outPoint.SerializeToCore(&buffer)
	c.do("SADD", key, buffer.Bytes())
}

func (c *commandBuffer) remove(key string, outPoint rpctypes.OutPoint) {
	if c.err != nil {
		return
	}
	var buffer bytes.Buffer
	c.err = outPoint.SerializeToCore(&buffer)
	c.do("SREM", key, buffer.Bytes())
}

func (c *commandBuffer) execute(conn redis.Conn) error {
	if c.err != nil {
		return c.err
	}
	conn.Send("MULTI")
	for _, command := range c.commands {
		conn.Send(command.name, command.args...)
	}
	_, err := conn.Do("EXEC")
	return err
}

type indexingEnvironment struct {
	cell   *internal_ast.Value
	params map[int]*internal_ast.Value
}

func (e *indexingEnvironment) Arg(i int) *internal_ast.Value {
	if i == 0 {
		return e.cell
	} else {
		return nil
	}
}

func (e *indexingEnvironment) Param(i int) *internal_ast.Value {
	param, found := e.params[i]
	// TODO: for better security we can visit the AST first, gather the number
	// of params required, then come back here and allocate all params at first.
	// That will work around problems where you have shortcuts in the AST.
	if !found {
		param = &internal_ast.Value{
			Indexing: true,
		}
		e.params[i] = param
	}
	return param
}

func (e *indexingEnvironment) QueryCell(query *ast.List) ([]*internal_ast.Value, error) {
	return nil, fmt.Errorf("QueryCell is not allowed in indexer!")
}

func executeIndexingQuery(query *ast.List, cell *rpctypes.CellOutput, cellData *rpctypes.Raw) ([]*internal_ast.Value, error) {
	if len(query.GetValues()) != 1 {
		return nil, fmt.Errorf("Invalid number of values to query cell: %d", len(query.GetValues()))
	}
	environment := &indexingEnvironment{
		cell: &internal_ast.Value{
			Cell:     cell,
			CellData: cellData,
		},
		params: make(map[int]*internal_ast.Value),
	}
	value, err := executor.Execute(query.GetValues()[0], environment)
	if err != nil {
		return nil, err
	}
	if value.GetT() != ast.Value_BOOL {
		return nil, fmt.Errorf("Invalid result value type: %s", value.GetT().String())
	}
	if !value.GetB() {
		return nil, nil
	}
	sortedParams := make([]*internal_ast.Value, len(environment.params))
	for i, param := range environment.params {
		if i >= len(sortedParams) {
			return nil, fmt.Errorf("Params are not all used!")
		}
		sortedParams[i] = param
	}
	return sortedParams, nil
}
