package indexer

import (
	"bytes"
	"fmt"
	"sort"

	"github.com/golang/protobuf/proto"
	"github.com/xxuejie/animagus/pkg/ast"
)

type ValueContext struct {
	Name        string
	Value       *ast.Value
	Queries     []*ast.Value
	QueryParams [][]int
}

func NewValueContext(name string, value *ast.Value) (ValueContext, error) {
	context := ValueContext{
		Name:    name,
		Value:   value,
		Queries: make([]*ast.Value, 0),
	}
	if err := visitValue(value, &context); err != nil {
		return ValueContext{}, err
	}
	return context, nil
}

func (c ValueContext) QueryIndex(query *ast.Value) int {
	for i, q := range c.Queries {
		if proto.Equal(query, q) {
			return i
		}
	}
	return -1
}

func (c ValueContext) IndexKey(queryIndex int, paramValues map[int]*ast.Value) (string, error) {
	params := c.QueryParams[queryIndex]
	var buffer bytes.Buffer
	_, err := buffer.WriteString(fmt.Sprintf("%d", len(params)))
	if err != nil {
		return "", err
	}
	for i := range params {
		value, found := paramValues[i]
		if !found {
			return "", fmt.Errorf("Requested param index %d is not provided!", i)
		}
		switch value.GetT() {
		case ast.Value_UINT64:
			_, err = buffer.WriteString(fmt.Sprintf("n%d", value.GetU()))
		case ast.Value_BOOL:
			_, err = buffer.WriteString(fmt.Sprintf("o%t", value.GetB()))
		case ast.Value_BYTES:
			_, err = buffer.WriteString(fmt.Sprintf("x%x", value.GetRaw()))
		default:
			err = fmt.Errorf("Invalid param value type: %s", value.GetT().String())
		}
		if err != nil {
			return "", err
		}
	}
	paramKey := string(buffer.Bytes())
	return fmt.Sprintf("CALL:%s:QUERY:%d:PARAM:%s:CELLS", c.Name, queryIndex, paramKey), nil
}

func visitValue(value *ast.Value, context *ValueContext) error {
	if value.GetT() == ast.Value_QUERY_CELLS {
		for _, q := range context.Queries {
			if proto.Equal(q, value) {
				return nil
			}
		}
		paramSet := make(map[int]bool)
		gatherQueryParams(value, &paramSet)
		params := make([]int, 0, len(paramSet))
		for k := range paramSet {
			params = append(params, k)
		}
		sort.Ints(params)
		context.Queries = append(context.Queries, value)
		context.QueryParams = append(context.QueryParams, params)
		return nil
	}
	for _, child := range value.GetChildren() {
		err := visitValue(child, context)
		if err != nil {
			return err
		}
	}
	return nil
}

func gatherQueryParams(value *ast.Value, paramSet *map[int]bool) {
	if value.GetT() == ast.Value_PARAM {
		(*paramSet)[int(value.GetU())] = true
		return
	}
	for _, child := range value.GetChildren() {
		gatherQueryParams(child, paramSet)
	}
}
