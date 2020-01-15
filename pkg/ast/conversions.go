package ast

import (
	"github.com/xxuejie/animagus/pkg/rpctypes"
)

func ConvertScript(script rpctypes.Script) *Value {
	return &Value{
		T: Value_SCRIPT,
		Children: []*Value{
			&Value{
				T: Value_BYTES,
				Primitive: &Value_Raw{
					Raw: script.CodeHash[:],
				},
			},
			&Value{
				T: Value_UINT64,
				Primitive: &Value_U{
					U: uint64(script.HashType),
				},
			},
			&Value{
				T: Value_BYTES,
				Primitive: &Value_Raw{
					Raw: script.Args,
				},
			},
		},
	}
}

func ConvertOutPoint(outPoint rpctypes.OutPoint) *Value {
	return &Value{
		T: Value_OUT_POINT,
		Children: []*Value{
			&Value{
				T: Value_BYTES,
				Primitive: &Value_Raw{
					Raw: outPoint.TxHash[:],
				},
			},
			&Value{
				T: Value_UINT64,
				Primitive: &Value_U{
					U: uint64(outPoint.Index),
				},
			},
		},
	}
}

func ConvertCell(cell rpctypes.CellOutput, cellData rpctypes.Raw, outPoint rpctypes.OutPoint) *Value {
	typeScript := &Value{T: Value_NIL}
	if cell.Type != nil {
		typeScript = ConvertScript(*cell.Type)
	}
	return &Value{
		T: Value_CELL,
		Children: []*Value{
			&Value{
				T: Value_UINT64,
				Primitive: &Value_U{
					U: uint64(cell.Capacity),
				},
			},
			ConvertScript(cell.Lock),
			typeScript,
			&Value{
				T: Value_BYTES,
				Primitive: &Value_Raw{
					Raw: cellData,
				},
			},
			ConvertOutPoint(outPoint),
		},
	}
}
