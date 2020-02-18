package main

import (
	"bytes"
	"fmt"
	"log"

	"github.com/xxuejie/animagus/pkg/ast"
	"github.com/xxuejie/animagus/pkg/validator"
)

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

func fetch_field(field ast.Value_Type, value *ast.Value) *ast.Value {
	return &ast.Value{
		T:        field,
		Children: []*ast.Value{value},
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

func map_funcs(list *ast.Value, funcs ...*ast.Value) *ast.Value {
	for _, f := range funcs {
		list = &ast.Value{
			T:        ast.Value_MAP,
			Children: []*ast.Value{f, list},
		}
	}
	return list
}

func main() {
	inputs := &ast.Value{
		T:        ast.Value_GET_INPUTS,
		Children: []*ast.Value{uint_value(1)},
	}
	inputTokens := map_funcs(
		inputs,
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
	inputSum := &ast.Value{
		T: ast.Value_REDUCE,
		Children: []*ast.Value{
			&ast.Value{
				T: ast.Value_ADD,
				Children: []*ast.Value{
					arg(0),
					arg(1),
				},
			},
			&ast.Value{
				T: ast.Value_ADD,
				Children: []*ast.Value{
					bytes_value([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}),
					uint_value(0),
				},
			},
			inputTokens,
		},
	}

	outputs := &ast.Value{
		T:        ast.Value_GET_OUTPUTS,
		Children: []*ast.Value{uint_value(1)},
	}
	outputTokens := map_funcs(
		outputs,
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
	outputSum := &ast.Value{
		T: ast.Value_REDUCE,
		Children: []*ast.Value{
			&ast.Value{
				T: ast.Value_ADD,
				Children: []*ast.Value{
					arg(0),
					arg(1),
				},
			},
			&ast.Value{
				T: ast.Value_ADD,
				Children: []*ast.Value{
					bytes_value([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}),
					uint_value(0),
				},
			},
			outputTokens,
		},
	}

	root := &ast.Value{
		T: ast.Value_EQUAL,
		Children: []*ast.Value{
			inputSum,
			outputSum,
		},
	}

	var source bytes.Buffer
	err := validator.Generate(root, &source)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(source.String())
}
