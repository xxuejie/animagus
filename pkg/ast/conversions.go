package ast

import (
	"bytes"
	"fmt"

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

func ConvertHeader(header rpctypes.Header) *Value {
	var nonceBuffer bytes.Buffer
	err := header.Nonce.SerializeToCore(&nonceBuffer)
	if err != nil {
		return &Value{
			T: Value_ERROR,
			Primitive: &Value_Raw{
				Raw: []byte(fmt.Sprintf("Serializing nonce error: %v", err)),
			},
		}
	}
	return &Value{
		T: Value_HEADER,
		Children: []*Value{
			&Value{
				T: Value_UINT64,
				Primitive: &Value_U{
					U: uint64(header.CompactTarget),
				},
			},
			&Value{
				T: Value_UINT64,
				Primitive: &Value_U{
					U: uint64(header.Timestamp),
				},
			},
			&Value{
				T: Value_UINT64,
				Primitive: &Value_U{
					U: uint64(header.Number),
				},
			},
			&Value{
				T: Value_UINT64,
				Primitive: &Value_U{
					U: uint64(header.Epoch),
				},
			},
			&Value{
				T: Value_BYTES,
				Primitive: &Value_Raw{
					Raw: header.ParentHash[:],
				},
			},
			&Value{
				T: Value_BYTES,
				Primitive: &Value_Raw{
					Raw: header.TransactionsRoot[:],
				},
			},
			&Value{
				T: Value_BYTES,
				Primitive: &Value_Raw{
					Raw: header.ProposalsHash[:],
				},
			},
			&Value{
				T: Value_BYTES,
				Primitive: &Value_Raw{
					Raw: header.UnclesHash[:],
				},
			},
			&Value{
				T: Value_BYTES,
				Primitive: &Value_Raw{
					Raw: header.Dao[:],
				},
			},
			&Value{
				T: Value_BYTES,
				Primitive: &Value_Raw{
					Raw: nonceBuffer.Bytes(),
				},
			},
		},
	}
}

func ConvertCell(cell rpctypes.CellOutput, cellData rpctypes.Raw,
	outPoint rpctypes.OutPoint, header *rpctypes.Header) *Value {
	typeScript := &Value{T: Value_NIL}
	if cell.Type != nil {
		typeScript = ConvertScript(*cell.Type)
	}
	children := []*Value{
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
	}
	return &Value{
		T:        Value_CELL,
		Children: children,
	}
}
