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
	SecpCellDep  = []byte{0x71, 0xa7, 0xba, 0x8f, 0xc9, 0x63, 0x49, 0xfe, 0xa0, 0xed, 0x3a, 0x5c, 0x47, 0x99, 0x2e, 0x3b, 0x40, 0x84, 0xb0, 0x31, 0xa4, 0x22, 0x64, 0xa0, 0x18, 0xe0, 0x07, 0x2e, 0x81, 0x72, 0xe4, 0x6c}

	UdtCodeHash   = []byte{0x57, 0xdd, 0x00, 0x67, 0x81, 0x4d, 0xab, 0x35, 0x6e, 0x05, 0xc6, 0xde, 0xf0, 0xd0, 0x94, 0xbb, 0x79, 0x77, 0x67, 0x11, 0xe6, 0x8f, 0xfd, 0xfa, 0xd2, 0xdf, 0x6a, 0x7f, 0x87, 0x7f, 0x7d, 0xb6}
	UdtScriptHash = []byte{0x41, 0x37, 0x0c, 0x40, 0xe3, 0xb3, 0xcf, 0x76, 0xab, 0x9a, 0x8a, 0x25, 0xe3, 0xd7, 0xad, 0x50, 0xb1, 0xc8, 0x28, 0xb9, 0x16, 0x3f, 0xfb, 0x3e, 0xab, 0x8f, 0x28, 0x1b, 0x39, 0x3e, 0xd5, 0xe6}
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

	args_test := equal(fetch_field(ast.Value_GET_ARGS, lock), param(0))

	return and(code_hash_test, hash_type_test, args_test)
}

func isSimpleUdtCell(argIndex uint64) *ast.Value {
	t := fetch_field(ast.Value_GET_TYPE, arg(argIndex))
	return equal(
		&ast.Value{
			T: ast.Value_HASH,
			Children: []*ast.Value{
				t,
			},
		},
		bytes_value(UdtScriptHash),
	)
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

func assembleUdtType() *ast.Value {
	return &ast.Value{
		T: ast.Value_SCRIPT,
		Children: []*ast.Value{
			bytes_value(UdtCodeHash),
			uint_value(0),
			bytes_value([]byte{}),
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
			and(isDefaultSecpCell(0), isSimpleUdtCell(0)),
		},
	}

	tokens := map_funcs(
		cells,
		fetch_field(ast.Value_GET_DATA, arg(0)),
		&ast.Value{
			T: ast.Value_SLICE,
			Children: []*ast.Value{
				arg(0),
				&ast.Value{
					T: ast.Value_UINT64,
					Primitive: &ast.Value_U{
						U: 16,
					},
				},
			},
		},
	)

	totalCapacities := &ast.Value{
		T: ast.Value_REDUCE,
		Children: []*ast.Value{
			&ast.Value{
				T: ast.Value_PLUS,
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
				T: ast.Value_PLUS,
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
		T: ast.Value_PLUS,
		Children: []*ast.Value{
			bytes_value([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}),
			param(2),
		},
	}

	changeTokens := &ast.Value{
		T: ast.Value_MINUS,
		Children: []*ast.Value{
			balance,
			transferTokens,
		},
	}

	changeCapacities := &ast.Value{
		T: ast.Value_MINUS,
		Children: []*ast.Value{
			totalCapacities,
			uint_value(94),
		},
	}

	transferCell := &ast.Value{
		T: ast.Value_CELL,
		Children: []*ast.Value{
			uint_value(94),
			assembleSecpLock(1),
			assembleUdtType(),
			transferTokens,
		},
	}

	changeCell := &ast.Value{
		T: ast.Value_CELL,
		Children: []*ast.Value{
			changeCapacities,
			assembleSecpLock(0),
			assembleUdtType(),
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
				Result: transaction,
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
