package main

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/golang/protobuf/proto"
	"github.com/xxuejie/animagus/pkg/ast"
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

	expected_code_hash, err := hex.DecodeString("9bd7e06f3ecf4be0f2fcd2188b23f1b9fcc88e5d4b65a8637b17723bbda3cce8")
	if err != nil {
		log.Fatal(err)
	}
	code_hash_test := equal(
		fetch_field(ast.Value_GET_CODE_HASH, lock),
		bytes_value(expected_code_hash))

	hash_type_test := equal(
		fetch_field(ast.Value_GET_HASH_TYPE, lock),
		uint_value(1))

	args_test := equal(fetch_field(ast.Value_GET_ARGS, lock), param(0))

	return and(code_hash_test, hash_type_test, args_test)
}

func isSimpleUdtCell(argIndex uint64) *ast.Value {
	expected_script_hash, err := hex.DecodeString("41370c40e3b3cf76ab9a8a25e3d7ad50b1c828b9163ffb3eab8f281b393ed5e6")
	if err != nil {
		log.Fatal(err)
	}

	t := fetch_field(ast.Value_GET_TYPE, arg(argIndex))
	return equal(
		&ast.Value{
			T: ast.Value_HASH,
			Children: []*ast.Value{
				t,
			},
		},
		bytes_value(expected_script_hash),
	)
}

func main() {
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
			T: ast.Value_SLICE_BYTES,
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
	ioutil.WriteFile("simple_udt.bin", bytes, 0644)
}
