package indexer

import (
	"bytes"
	"fmt"

	"github.com/xxuejie/animagus/pkg/ast"
)

type ValueContext struct {
	Name    string
	Value   *ast.Value
	Queries []*ast.Value
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
		if query == q {
			return i
		}
	}
	return -1
}

func (c ValueContext) IndexKey(queryIndex int, params []*ast.Value) (string, error) {
	var buffer bytes.Buffer
	_, err := buffer.WriteString(fmt.Sprintf("%d", len(params)))
	if err != nil {
		return "", err
	}
	for _, value := range params {
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
		context.Queries = append(context.Queries, value)
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
