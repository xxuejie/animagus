package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/golang/protobuf/proto"
	"github.com/xxuejie/animagus/pkg/ast"
)

var (
	DaoTypeHash = []byte{0x82, 0xd7, 0x6d, 0x1b, 0x75, 0xfe, 0x2f, 0xd9, 0xa2, 0x7d, 0xfb, 0xaa, 0x65, 0xa0, 0x39, 0x22, 0x1a, 0x38, 0x0d, 0x76, 0xc9, 0x26, 0xf3, 0x78, 0xd3, 0xf8, 0x1c, 0xf3, 0xe7, 0xe1, 0x3f, 0x2e}
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

func bytes_value(b []byte) *ast.Value {
	return &ast.Value{
		T: ast.Value_BYTES,
		Primitive: &ast.Value_Raw{
			Raw: b,
		},
	}
}

func string_value(s string) *ast.Value {
	return &ast.Value{
		T: ast.Value_BYTES,
		Primitive: &ast.Value_Raw{
			Raw: []byte(s),
		},
	}
}

func main() {
	script := fetch_field(ast.Value_GET_TYPE, arg(0))
	code_hash_test := equal(
		fetch_field(ast.Value_GET_CODE_HASH, script),
		bytes_value(DaoTypeHash))
	type_test := equal(arg(1), string_value("insert"))
	index_test := equal(arg(2), string_value("index"))

	tests := and(code_hash_test, type_test, index_test)

	out_point := fetch_field(ast.Value_GET_OUT_POINT, arg(0))

	filter := &ast.Value{
		T: ast.Value_COND,
		Children: []*ast.Value{
			tests,
			out_point,
			&ast.Value{
				T: ast.Value_NIL,
			},
		},
	}

	root := &ast.Root{
		Streams: []*ast.Stream{
			&ast.Stream{
				Name:   "nervosdao_deposits",
				Filter: filter,
			},
		},
	}

	fmt.Println(proto.MarshalTextString(root))

	bytes, err := proto.Marshal(root)
	if err != nil {
		log.Fatal(err)
	}
	ioutil.WriteFile("nervosdao_events.bin", bytes, 0644)
}
