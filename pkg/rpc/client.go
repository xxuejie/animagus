package rpc

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

type responseBody struct {
	ID      int         `json:"id"`
	Jsonrpc string      `json:"jsonrpc"`
	Result  interface{} `json:"result"`
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

func RpcRequest(client *http.Client, url string, params RequestParams, target interface{}) error {
	b, _ := json.Marshal(params)
	resp, err := client.Post(url, "application/json", strings.NewReader(string(b)))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	result, _ := ioutil.ReadAll(resp.Body)

	hr := &responseBody{Result: &target}

	e := json.Unmarshal(result, hr)
	return e
}
