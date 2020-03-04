package executor

import (
	"fmt"
	"testing"

	"github.com/xxuejie/animagus/pkg/ast"
)

type testEnvironment struct {
	args   []*ast.Value
	params []*ast.Value
}

func (e *testEnvironment) ReplaceArgs(args []*ast.Value) error {
	if len(args) > len(e.args) {
		return fmt.Errorf("Too many args provided")
	}
	l := len(args)
	if len(e.args) < l {
		l = len(e.args)
	}
	for i := 0; i < l; i++ {
		e.args[i] = args[i]
	}
	return nil
}

func (e *testEnvironment) Arg(i int) *ast.Value {
	return e.args[i]
}

func (e *testEnvironment) Param(i int) *ast.Value {
	return e.params[i]
}

func (e *testEnvironment) IndexParam(i int, value *ast.Value) error {
	return fmt.Errorf("Index param is not expected!")
}

func (e *testEnvironment) QueryCell(query *ast.Value) ([]*ast.Value, error) {
	return nil, fmt.Errorf("Query cell is not expected!")
}

func uint_value(u uint64) *ast.Value {
	return &ast.Value{
		T: ast.Value_UINT64,
		Primitive: &ast.Value_U{
			U: u,
		},
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

func TestTailRecursion(t *testing.T) {
	n := uint_value(10)
	loop_i := arg(0)
	loop_a := arg(1)
	loop_b := arg(2)

	test := &ast.Value{
		T: ast.Value_LESS,
		Children: []*ast.Value{
			loop_i,
			n,
		},
	}

	next_i := &ast.Value{
		T: ast.Value_ADD,
		Children: []*ast.Value{
			loop_i,
			uint_value(1),
		},
	}
	next_a := loop_b
	next_b := &ast.Value{
		T: ast.Value_ADD,
		Children: []*ast.Value{
			loop_a,
			loop_b,
		},
	}

	next := &ast.Value{
		T: ast.Value_TAIL_RECURSION,
		Children: []*ast.Value{
			next_i,
			next_a,
			next_b,
		},
	}

	f := &ast.Value{
		T: ast.Value_COND,
		Children: []*ast.Value{
			test,
			next,
			loop_b,
		},
	}

	e := &testEnvironment{
		args: []*ast.Value{
			uint_value(0),
			uint_value(0),
			uint_value(1),
		},
		params: nil,
	}

	value, err := Execute(f, e)
	if err != nil {
		t.Fatal(err)
	}
	if value.GetT() != ast.Value_UINT64 {
		t.Errorf("Invalid value type: %s", value.GetT().String())
	}
	if value.GetU() != 89 {
		t.Errorf("Invalid result: %d, expected: 89", value.GetU())
	}
}
