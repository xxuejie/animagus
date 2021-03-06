package rpc

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/xxuejie/animagus/pkg/rpctypes"
)

type responseBody struct {
	ID      int         `json:"id"`
	Jsonrpc string      `json:"jsonrpc"`
	Result  interface{} `json:"result"`
}

type blockResponseBody struct {
	responseBody
	Result *rpctypes.BlockView `json:"result"`
}

type headerResponseBody struct {
	responseBody
	Result *rpctypes.HeaderView `json:"result"`
}

type transactionWithStatusResponseBody struct {
	responseBody
	Result *rpctypes.TransactionWithStatusView `json:"result"`
}

type blockNumberResponseBody struct {
	responseBody
	Result *rpctypes.Uint64 `json:"result"`
}

type RequestParams struct {
	ID      int      `json:"id"`
	Jsonrpc string   `json:"jsonrpc"`
	Method  string   `json:"method"`
	Params  []string `json:"params"`
}

func NewRequestParams(method string, params []string) RequestParams {
	return RequestParams{
		ID:      42,
		Jsonrpc: "2.0",
		Method:  method,
		Params:  params,
	}
}

type Client struct {
	HttpClient *http.Client
	Url        string
}

func NewClient(url string) *Client {
	return &Client{
		HttpClient: &http.Client{
			Transport: &http.Transport{
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 100,
			},
		},
		Url: url,
	}
}

func (c *Client) rpcRequest(params RequestParams) ([]byte, error) {
	b, _ := json.Marshal(params)
	bodyReader := strings.NewReader(string(b))
	req, err := http.NewRequest("POST", c.Url, bodyReader)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	result, err := ioutil.ReadAll(resp.Body)

	return result, err
}

func (c *Client) GetTipBlockNumber() (*rpctypes.Uint64, error) {
	params := NewRequestParams(
		"get_tip_block_number",
		[]string{},
	)

	result, err := c.rpcRequest(params)

	if err != nil {
		return nil, err
	}

	rb := &blockNumberResponseBody{}
	err = json.Unmarshal(result, rb)
	if err != nil {
		return nil, err
	}
	return rb.Result, nil
}

func (c *Client) GetBlockByNumber(blockNumber rpctypes.Uint64) (*rpctypes.BlockView, error) {
	params := NewRequestParams(
		"get_block_by_number",
		[]string{blockNumber.EncodeToString()},
	)

	result, err := c.rpcRequest(params)

	if err != nil {
		return nil, err
	}

	rb := &blockResponseBody{}
	err = json.Unmarshal(result, rb)
	if err != nil {
		return nil, err
	}
	return rb.Result, nil
}

func (c *Client) GetTransaction(txHash *rpctypes.Hash) (*rpctypes.TransactionWithStatusView, error) {
	txHashStr := fmt.Sprintf("0x%x", *txHash)
	params := NewRequestParams(
		"get_transaction",
		[]string{txHashStr},
	)

	result, err := c.rpcRequest(params)

	if err != nil {
		return nil, err
	}

	rb := &transactionWithStatusResponseBody{}
	err = json.Unmarshal(result, rb)

	if err != nil {
		return nil, err
	}

	transactionWithStatus := rb.Result

	return transactionWithStatus, nil
}

func (c *Client) getTransactions(txHashes []rpctypes.Hash) ([]*rpctypes.TransactionWithStatusView, error) {
	txHashLength := len(txHashes)
	if txHashLength == 0 {
		return []*rpctypes.TransactionWithStatusView{}, nil
	}

	type transactionWithError struct {
		TransactionWithStatusView *rpctypes.TransactionWithStatusView
		Err                       error
	}

	done := make(chan *transactionWithError, txHashLength)

	for _, txHash := range txHashes {
		go func(txHash rpctypes.Hash) {
			transactionWithStatus, err := c.GetTransaction(&txHash)

			done <- &transactionWithError{
				TransactionWithStatusView: transactionWithStatus,
				Err:                       err,
			}
		}(txHash)
	}

	transactionWithStatusViews := []*rpctypes.TransactionWithStatusView{}
	for i := 0; i < txHashLength; i++ {
		txViewWithError := <-done
		err := txViewWithError.Err
		if err != nil {
			return nil, err
		}
		txWithStatusView := txViewWithError.TransactionWithStatusView
		transactionWithStatusViews = append(transactionWithStatusViews, txWithStatusView)

		if i == txHashLength-1 {
			close(done)
		}
	}

	return transactionWithStatusViews, nil
}

func (c *Client) GetAllTransactions(txHashes []rpctypes.Hash, size int) ([]*rpctypes.TransactionWithStatusView, error) {
	var txHashSlices [][]rpctypes.Hash
	txHashLength := len(txHashes)
	for i := 0; i < txHashLength; i += size {
		rightEdge := i + size
		if rightEdge > txHashLength {
			rightEdge = txHashLength
		}
		hashes := txHashes[i:rightEdge]
		txHashSlices = append(txHashSlices, hashes)
	}

	var transactionWithStatusViews []*rpctypes.TransactionWithStatusView
	for _, slice := range txHashSlices {
		txsWithStatus, err := c.getTransactions(slice)
		if err != nil {
			return txsWithStatus, err
		}
		for _, tx := range txsWithStatus {
			transactionWithStatusViews = append(transactionWithStatusViews, tx)
		}
	}

	return transactionWithStatusViews, nil
}

func (c *Client) GetHeader(blockHash *rpctypes.Hash) (*rpctypes.HeaderView, error) {
	blockHashStr := fmt.Sprintf("0x%x", *blockHash)
	params := NewRequestParams(
		"get_header",
		[]string{blockHashStr},
	)

	result, err := c.rpcRequest(params)
	if err != nil {
		return nil, err
	}

	rb := &headerResponseBody{}
	err = json.Unmarshal(result, rb)
	if err != nil {
		return nil, err
	}

	header := rb.Result

	return header, err
}

func (c *Client) getHeaders(blockHashes []rpctypes.Hash) ([]*rpctypes.HeaderView, error) {
	blockHashLength := len(blockHashes)
	if blockHashLength == 0 {
		return []*rpctypes.HeaderView{}, nil
	}

	type headerWithError struct {
		Header *rpctypes.HeaderView
		Err    error
	}

	done := make(chan *headerWithError, blockHashLength)

	for _, blockHash := range blockHashes {
		go func(blockHash rpctypes.Hash) {
			header, err := c.GetHeader(&blockHash)

			done <- &headerWithError{
				Header: header,
				Err:    err,
			}
		}(blockHash)
	}

	headers := []*rpctypes.HeaderView{}
	for i := 0; i < blockHashLength; i++ {
		blockWithError := <-done
		err := blockWithError.Err
		if err != nil {
			return nil, err
		}
		header := blockWithError.Header
		headers = append(headers, header)

		if i == blockHashLength-1 {
			close(done)
		}
	}

	return headers, nil
}

func (c *Client) GetAllHeaders(blockHashes []rpctypes.Hash, size int) ([]*rpctypes.HeaderView, error) {
	var blockHashSlices [][]rpctypes.Hash
	blockHashLength := len(blockHashes)
	for i := 0; i < blockHashLength; i += size {
		rightEdge := i + size
		if rightEdge > blockHashLength {
			rightEdge = blockHashLength
		}
		hashes := blockHashes[i:rightEdge]
		blockHashSlices = append(blockHashSlices, hashes)
	}

	var headersResult []*rpctypes.HeaderView
	for _, slice := range blockHashSlices {
		headers, err := c.getHeaders(slice)
		if err != nil {
			return headers, err
		}
		for _, header := range headers {
			headersResult = append(headersResult, header)
		}
	}

	return headersResult, nil
}
