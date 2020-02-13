package validator

import (
	"fmt"

	"github.com/xxuejie/animagus/pkg/ast"
)

type syscallContext struct {
	c      *context
	method string
	field  string
	source string
	start  uint64
	end    uint64

	remainingFuncs []*ast.Value
}

func createSyscallContext(c *context, current *ast.Value, map_funcs []*ast.Value) (*syscallContext, error) {
	sc := &syscallContext{
		c:              c,
		remainingFuncs: map_funcs,
	}
	err := sc.visit(current)
	if err != nil {
		return nil, err
	}
	return sc, nil
}

func (sc *syscallContext) visit(current *ast.Value) error {
	for current != nil {
		switch current.GetT() {
		case ast.Value_GET_INPUTS:
			sc.method = "ckb_load_cell_by_field"
			if len(current.GetChildren()) == 1 &&
				current.GetChildren()[0].GetT() == ast.Value_UINT64 &&
				current.GetChildren()[0].GetU() == 1 {
				sc.source = "CKB_SOURCE_GROUP_INPUT"
			} else {
				sc.source = "CKB_SOURCE_INPUT"
			}
			current = sc.shift()
		case ast.Value_GET_OUTPUTS:
			sc.method = "ckb_load_cell_by_field"
			if len(current.GetChildren()) == 1 &&
				current.GetChildren()[0].GetT() == ast.Value_UINT64 &&
				current.GetChildren()[0].GetU() == 1 {
				sc.source = "CKB_SOURCE_GROUP_OUTPUT"
			} else {
				sc.source = "CKB_SOURCE_OUTPUT"
			}
			current = sc.shift()
		case ast.Value_GET_DATA:
			if len(current.GetChildren()) != 1 {
				return fmt.Errorf("Invalid number of operands to GET_DATA")
			}
			sc.method = "ckb_load_cell_data"
			current = sc.checkArgOrShift(current.GetChildren()[0])
		case ast.Value_SLICE:
			children := current.GetChildren()
			if len(children) != 3 {
				return fmt.Errorf("Invalid number of operands to SLICE")
			}
			if children[0].GetT() != ast.Value_UINT64 ||
				children[1].GetT() != ast.Value_UINT64 {
				return fmt.Errorf("Invalid operand type to SLICE!")
			}
			sc.start = children[0].GetU()
			sc.end = children[1].GetU()
			current = sc.checkArgOrShift(children[2])
		default:
			sc.remainingFuncs = append([]*ast.Value{current}, sc.remainingFuncs...)
			current = nil
		}
	}
	return nil
}

func (sc *syscallContext) checkArgOrShift(expr *ast.Value) *ast.Value {
	if expr.GetT() == ast.Value_ARG && expr.GetU() == 0 {
		return sc.shift()
	} else {
		return expr
	}
}

func (sc *syscallContext) shift() *ast.Value {
	if len(sc.remainingFuncs) > 0 {
		v := sc.remainingFuncs[0]
		sc.remainingFuncs = sc.remainingFuncs[1:]
		return v
	}
	return nil
}

func (sc *syscallContext) generate() (int, variableType, error) {
	idx := sc.c.newVariable(varTypeUint64)
	t := variableType{
		t:      varBytes,
		length: sc.end - sc.start,
	}
	i := sc.c.newVariable(t)
	sc.c.printfln("uint64_t v%d_length = %d;", i, t.length)
	sc.c.printfln("uint8_t v%d[%d];", i, t.length)
	sc.c.printfln("uint64_t v%d = 0;", idx)
	sc.c.printfln("while (1) {")
	sc.c.indent()
	sc.c.printfln("v%d_length = %d;", i, t.length)
	sc.c.printfln("memset(v%d, 0, %d);", i, t.length)
	if sc.method == "ckb_load_cell_data" {
		sc.c.printfln("int ret = ckb_load_cell_data(v%d, &v%d_length, %d, v%d, %s);",
			i, i, sc.start, idx, sc.source)
	} else {
		sc.c.printfln("int ret = %s(v%d, &v%d_length, %d, v%d, %s, %s);",
			sc.method, i, i, sc.start, idx, sc.source, sc.field)
	}
	sc.c.printfln("if (ret == CKB_INDEX_OUT_OF_BOUND) { break; }")
	sc.c.printfln("if (ret != 0) { return %d; }", sc.c.newErrorCode())
	sc.c.printfln("if (v%d_length != %d) { return %d; }", i, t.length, sc.c.newErrorCode())
	sc.c.printfln("v%d += 1;", idx)
	return i, t, nil
}
