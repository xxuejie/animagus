package ast

import (
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/xxuejie/animagus/pkg/ast"
	"github.com/xxuejie/animagus/pkg/rpctypes"
)

// TODO: merge the special values into ast.Value when we add transaction
// building support.
type Value struct {
	Indexing      bool
	IndexingValue *Value

	Value    *ast.Value
	Cell     *rpctypes.CellOutput
	CellData *rpctypes.Raw
	Script   *rpctypes.Script
	List     []*Value
}

func (v *Value) IsTerminal() bool {
	t := v.Value.GetT()
	return t != ast.Value_ARG && t != ast.Value_PARAM &&
		t != ast.Value_LIST && t != ast.Value_APPLY && t != ast.Value_REDUCE
}

func (v *Value) GetIndex() int {
	return int(v.Value.GetU())
}

func (v *Value) indexValue(target *Value) (bool, error) {
	if v.IndexingValue != nil {
		return false, fmt.Errorf("Indexing value is already indexed!")
	}
	v.IndexingValue = target
	return true, nil
}

func (v *Value) Equal(b *Value) (bool, error) {
	if v == nil || b == nil {
		return v == b, nil
	}
	if v.Indexing {
		return v.indexValue(b)
	}
	if b.Indexing {
		return b.indexValue(v)
	}
	if v.Value != nil {
		return proto.Equal(v.Value, b.Value), nil
	}
	if v.Cell != nil {
		return deriveEqualCell(v.Cell, b.Cell) &&
			deriveEqualCellData(v.CellData, b.CellData), nil
	}
	return deriveEqualScript(v.Script, b.Script), nil
}
