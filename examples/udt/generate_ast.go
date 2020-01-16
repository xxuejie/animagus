package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/golang/protobuf/proto"
	"github.com/xxuejie/animagus/pkg/ast"
)

var (
	SecpTypeHash = []byte{0x9b, 0xd7, 0xe0, 0x6f, 0x3e, 0xcf, 0x4b, 0xe0, 0xf2, 0xfc, 0xd2, 0x18, 0x8b, 0x23, 0xf1, 0xb9, 0xfc, 0xc8, 0x8e, 0x5d, 0x4b, 0x65, 0xa8, 0x63, 0x7b, 0x17, 0x72, 0x3b, 0xbd, 0xa3, 0xcc, 0xe8}
	// Mainnet
	// SecpCellDep  = []byte{0x71, 0xa7, 0xba, 0x8f, 0xc9, 0x63, 0x49, 0xfe, 0xa0, 0xed, 0x3a, 0x5c, 0x47, 0x99, 0x2e, 0x3b, 0x40, 0x84, 0xb0, 0x31, 0xa4, 0x22, 0x64, 0xa0, 0x18, 0xe0, 0x07, 0x2e, 0x81, 0x72, 0xe4, 0x6c}
	// Devnet
	SecpCellDep = []byte{0xac, 0xe5, 0xea, 0x83, 0xc4, 0x78, 0xbb, 0x86, 0x6e, 0xdf, 0x12, 0x2f, 0xf8, 0x62, 0x08, 0x57, 0x89, 0x15, 0x8f, 0x5c, 0xbf, 0xf1, 0x55, 0xb7, 0xbb, 0x5f, 0x13, 0x05, 0x85, 0x55, 0xb7, 0x08}

	UdtCodeHash = []byte{0x57, 0xdd, 0x00, 0x67, 0x81, 0x4d, 0xab, 0x35, 0x6e, 0x05, 0xc6, 0xde, 0xf0, 0xd0, 0x94, 0xbb, 0x79, 0x77, 0x67, 0x11, 0xe6, 0x8f, 0xfd, 0xfa, 0xd2, 0xdf, 0x6a, 0x7f, 0x87, 0x7f, 0x7d, 0xb6}
)

func fetch_field(field ast.Value_Type, value *ast.Value) *ast.Value {
	return &ast.Value{
		T:        field,
		Children: []*ast.Value{value},
	}
}

func equal(a *ast.Value, b *ast.Value) *ast.Value {
	return &ast.Value{
		T: ast.Value_EQUAL,
		Children: []*ast.Value{
			a,
			b,
		},
	}
}

func bytes_value(b []byte) *ast.Value {
	return &ast.Value{
		T: ast.Value_BYTES,
		Primitive: &ast.Value_Raw{
			Raw: b,
		},
	}
}

func uint_value(u uint64) *ast.Value {
	return &ast.Value{
		T: ast.Value_UINT64,
		Primitive: &ast.Value_U{
			U: u,
		},
	}
}

func and(values ...*ast.Value) *ast.Value {
	return &ast.Value{
		T:        ast.Value_AND,
		Children: values,
	}
}

func arg(i uint64) *ast.Value {
	return &ast.Value{
		T: ast.Value_ARG,
		Primitive: &ast.Value_U{
			U: i,
		},
	}
}

func param(i uint64) *ast.Value {
	return &ast.Value{
		T: ast.Value_PARAM,
		Primitive: &ast.Value_U{
			U: i,
		},
	}
}

func map_funcs(list *ast.Value, funcs ...*ast.Value) *ast.Value {
	for _, f := range funcs {
		list = &ast.Value{
			T:        ast.Value_MAP,
			Children: []*ast.Value{f, list},
		}
	}
	return list
}

func isDefaultSecpCell(argIndex uint64) *ast.Value {
	lock := fetch_field(ast.Value_GET_LOCK, arg(argIndex))

	code_hash_test := equal(
		fetch_field(ast.Value_GET_CODE_HASH, lock),
		bytes_value(SecpTypeHash))

	hash_type_test := equal(
		fetch_field(ast.Value_GET_HASH_TYPE, lock),
		uint_value(1))

	args_test := equal(fetch_field(ast.Value_GET_ARGS, lock), param(1))

	return and(code_hash_test, hash_type_test, args_test)
}

func isSimpleUdtCell(argIndex uint64, paramIndex uint64) *ast.Value {
	t := fetch_field(ast.Value_GET_TYPE, arg(argIndex))

	code_hash_test := equal(
		fetch_field(ast.Value_GET_CODE_HASH, t),
		bytes_value(UdtCodeHash))

	hash_type_test := equal(
		fetch_field(ast.Value_GET_HASH_TYPE, t),
		uint_value(0))

	args_test := equal(fetch_field(ast.Value_GET_ARGS, t), param(paramIndex))

	return and(code_hash_test, hash_type_test, args_test)
}

func assembleSecpLock(paramIndex uint64) *ast.Value {
	return &ast.Value{
		T: ast.Value_SCRIPT,
		Children: []*ast.Value{
			bytes_value(SecpTypeHash),
			uint_value(1),
			param(paramIndex),
		},
	}
}

func assembleSecpCellDep() *ast.Value {
	return &ast.Value{
		T: ast.Value_CELL_DEP,
		Children: []*ast.Value{
			&ast.Value{
				T: ast.Value_OUT_POINT,
				Children: []*ast.Value{
					bytes_value(SecpCellDep),
					uint_value(0),
				},
			},
			uint_value(1),
		},
	}
}

func assembleUdtType(paramIndex uint64) *ast.Value {
	return &ast.Value{
		T: ast.Value_SCRIPT,
		Children: []*ast.Value{
			bytes_value(UdtCodeHash),
			uint_value(0),
			param(paramIndex),
		},
	}
}

func adjustFee(tx *ast.Value) *ast.Value {
	length := &ast.Value{
		T: ast.Value_LEN,
		Children: []*ast.Value{
			&ast.Value{
				T:        ast.Value_SERIALIZE_TO_CORE,
				Children: []*ast.Value{tx},
			},
		},
	}
	fee := &ast.Value{
		T: ast.Value_MULTIPLY,
		Children: []*ast.Value{
			// Adding extra bytes here to set aside for signatures
			&ast.Value{
				T: ast.Value_ADD,
				Children: []*ast.Value{
					length,
					uint_value(100),
				},
			},
			uint_value(1),
		},
	}
	changeCell := tx.GetChildren()[1].GetChildren()[1]
	adjustedChangeCell := &ast.Value{
		T: ast.Value_CELL,
		Children: []*ast.Value{
			&ast.Value{
				T: ast.Value_SUBTRACT,
				Children: []*ast.Value{
					changeCell.GetChildren()[0],
					fee,
				},
			},
			changeCell.GetChildren()[1],
			changeCell.GetChildren()[2],
			changeCell.GetChildren()[3],
		},
	}
	return &ast.Value{
		T: ast.Value_TRANSACTION,
		Children: []*ast.Value{
			tx.GetChildren()[0],
			&ast.Value{
				T: ast.Value_LIST,
				Children: []*ast.Value{
					tx.GetChildren()[1].GetChildren()[0],
					adjustedChangeCell,
				},
			},
			tx.GetChildren()[2],
		},
	}
}

func main() {
	typeCells := &ast.Value{
		T: ast.Value_QUERY_CELLS,
		Children: []*ast.Value{
			equal(
				&ast.Value{
					T:        ast.Value_GET_DATA_HASH,
					Children: []*ast.Value{arg(0)},
				},
				bytes_value(UdtCodeHash),
			),
		},
	}

	ready := equal(
		&ast.Value{
			T:        ast.Value_LEN,
			Children: []*ast.Value{typeCells},
		},
		uint_value(1),
	)

	cells := &ast.Value{
		T: ast.Value_QUERY_CELLS,
		Children: []*ast.Value{
			and(isDefaultSecpCell(0), isSimpleUdtCell(0, 0)),
		},
	}

	tokens := map_funcs(
		cells,
		fetch_field(ast.Value_GET_DATA, arg(0)),
		&ast.Value{
			T: ast.Value_SLICE,
			Children: []*ast.Value{
				uint_value(0),
				uint_value(16),
				arg(0),
			},
		},
	)

	totalCapacities := &ast.Value{
		T: ast.Value_REDUCE,
		Children: []*ast.Value{
			&ast.Value{
				T: ast.Value_ADD,
				Children: []*ast.Value{
					arg(0),
					arg(1),
				},
			},
			uint_value(0),
			&ast.Value{
				T:        ast.Value_MAP,
				Children: []*ast.Value{fetch_field(ast.Value_GET_CAPACITY, arg(0)), cells},
			},
		},
	}

	balance := &ast.Value{
		T: ast.Value_REDUCE,
		Children: []*ast.Value{
			&ast.Value{
				T: ast.Value_ADD,
				Children: []*ast.Value{
					arg(0),
					arg(1),
				},
			},
			bytes_value([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}),
			tokens,
		},
	}

	balance = &ast.Value{
		T: ast.Value_SLICE,
		Children: []*ast.Value{
			uint_value(0),
			uint_value(16),
			balance,
		},
	}

	// This helps cast uint64 values to bytes to make it handy.
	transferTokens := &ast.Value{
		T: ast.Value_SLICE,
		Children: []*ast.Value{
			uint_value(0),
			uint_value(16),
			&ast.Value{
				T: ast.Value_ADD,
				Children: []*ast.Value{
					bytes_value([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}),
					param(3),
				},
			},
		},
	}

	changeTokens := &ast.Value{
		T: ast.Value_SLICE,
		Children: []*ast.Value{
			uint_value(0),
			uint_value(16),
			&ast.Value{
				T: ast.Value_SUBTRACT,
				Children: []*ast.Value{
					balance,
					transferTokens,
				},
			},
		},
	}

	changeCapacities := &ast.Value{
		T: ast.Value_SUBTRACT,
		Children: []*ast.Value{
			totalCapacities,
			uint_value(142 * 100000000),
		},
	}

	transferCell := &ast.Value{
		T: ast.Value_CELL,
		Children: []*ast.Value{
			uint_value(142 * 100000000),
			assembleSecpLock(2),
			assembleUdtType(0),
			transferTokens,
		},
	}

	changeCell := &ast.Value{
		T: ast.Value_CELL,
		Children: []*ast.Value{
			changeCapacities,
			assembleSecpLock(1),
			assembleUdtType(0),
			changeTokens,
		},
	}

	// TODO: witness support
	transaction := &ast.Value{
		T: ast.Value_TRANSACTION,
		Children: []*ast.Value{
			cells,
			&ast.Value{
				T: ast.Value_LIST,
				Children: []*ast.Value{
					transferCell,
					changeCell,
				},
			},
			&ast.Value{
				T: ast.Value_LIST,
				Children: []*ast.Value{
					assembleSecpCellDep(),
					&ast.Value{
						T: ast.Value_INDEX,
						Children: []*ast.Value{
							uint_value(0),
							typeCells,
						},
					},
				},
			},
		},
	}
	transaction = adjustFee(transaction)

	serializedTransaction := &ast.Value{
		T:        ast.Value_SERIALIZE_TO_JSON,
		Children: []*ast.Value{transaction},
	}

	root := &ast.Root{
		Calls: []*ast.Call{
			&ast.Call{
				Name:   "ready",
				Result: ready,
			},
			&ast.Call{
				Name:   "balance",
				Result: balance,
			},
			&ast.Call{
				Name:   "transfer",
				Result: serializedTransaction,
			},
		},
	}

	fmt.Println(proto.MarshalTextString(root))

	bytes, err := proto.Marshal(root)
	if err != nil {
		log.Fatal(err)
	}
	ioutil.WriteFile("simple_udt.bin", bytes, 0644)
}
