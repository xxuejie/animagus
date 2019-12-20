package main

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/golang/protobuf/proto"
	"github.com/xxuejie/animagus/pkg/ast"
)

func fetch_field(field ast.Field, value *ast.Value) *ast.Value {
	return &ast.Value{
		T: ast.Value_APPLY,
		Children: []*ast.Value{
			&ast.Value{
				T: ast.Value_OP,
				Primitive: &ast.Value_Op{
					Op: ast.Op_GET,
				},
			},
			&ast.Value{
				T: ast.Value_FIELD,
				Primitive: &ast.Value_Field{
					Field: field,
				},
			},
			value,
		},
	}
}

func equal(a *ast.Value, b *ast.Value) *ast.Value {
	return &ast.Value{
		T: ast.Value_APPLY,
		Children: []*ast.Value{
			&ast.Value{
				T: ast.Value_OP,
				Primitive: &ast.Value_Op{
					Op: ast.Op_EQUAL,
				},
			},
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
	children := []*ast.Value{
		&ast.Value{
			T: ast.Value_OP,
			Primitive: &ast.Value_Op{
				Op: ast.Op_AND,
			},
		},
	}
	children = append(children, values...)
	return &ast.Value{
		T:        ast.Value_APPLY,
		Children: children,
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

func main() {
	lock := fetch_field(ast.Field_LOCK, arg(0))

	expected_code_hash, err := hex.DecodeString("9bd7e06f3ecf4be0f2fcd2188b23f1b9fcc88e5d4b65a8637b17723bbda3cce8")
	if err != nil {
		log.Fatal(err)
	}
	code_hash_test := equal(
		fetch_field(ast.Field_CODE_HASH, lock),
		bytes_value(expected_code_hash))

	hash_type_test := equal(
		fetch_field(ast.Field_HASH_TYPE, lock),
		uint_value(1))

	args_test := equal(fetch_field(ast.Field_ARGS, lock), param(0))

	test := and(code_hash_test, hash_type_test, args_test)

	cells := &ast.List{
		T:      ast.List_QUERY_CELLS,
		Values: []*ast.Value{test},
	}

	capacities := &ast.List{
		T:        ast.List_MAP,
		Values:   []*ast.Value{fetch_field(ast.Field_CAPACITY, arg(0))},
		Children: []*ast.List{cells},
	}

	balance := &ast.Value{
		T: ast.Value_REDUCE,
		L: capacities,
		Children: []*ast.Value{
			&ast.Value{
				T: ast.Value_APPLY,
				Children: []*ast.Value{
					&ast.Value{
						T: ast.Value_OP,
						Primitive: &ast.Value_Op{
							Op: ast.Op_PLUS,
						},
					},
					arg(0),
					arg(1),
				},
			},
			uint_value(0),
		},
	}

	root := &ast.Root{
		Calls: []*ast.Call{
			&ast.Call{
				Name:   "balance",
				Result: balance,
			},
		},
	}

	fmt.Println(proto.MarshalTextString(root))

	bytes, err := proto.Marshal(root)
	if err != nil {
		log.Fatal(err)
	}
	ioutil.WriteFile("balance.bin", bytes, 0644)
}
