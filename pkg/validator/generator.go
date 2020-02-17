package validator

import (
	"fmt"
	"io"

	"github.com/xxuejie/animagus/pkg/ast"
)

type variableType struct {
	t      int
	length uint64
}

func (t variableType) cType() string {
	switch t.t {
	case varBool:
		return "bool"
	case varUint64:
		return "uint64_t"
	case varUint128:
		return "uint128_t"
	case varBytes:
		return "uint8_t*"
	}
	return ""
}

const (
	varBool = iota + 1
	varUint64
	varUint128
	varBytes
)

var (
	varTypeEmpty = variableType{}
	varTypeBool  = variableType{
		t: varBool,
	}
	varTypeUint64 = variableType{
		t: varUint64,
	}
	varTypeUint128 = variableType{
		t: varUint128,
	}
)

type context struct {
	err           error
	nextErrorCode int8
	currentIndent string
	variables     []variableType
	args          []int
	writer        io.Writer
}

func newContext(writer io.Writer) *context {
	return &context{
		nextErrorCode: -1,
		writer:        writer,
	}
}

func (c *context) newVariable(t variableType) int {
	i := len(c.variables)
	c.variables = append(c.variables, t)
	return i
}

func (c *context) newErrorCode() int8 {
	code := c.nextErrorCode
	c.nextErrorCode -= 1
	// Skip 0
	if c.nextErrorCode == 0 {
		c.nextErrorCode = -1
	}
	return code
}

func (c *context) printf(format string, a ...interface{}) {
	if c.err != nil {
		return
	}
	_, c.err = fmt.Fprintf(c.writer, format, a...)
}

func (c *context) print_indent() {
	c.printf(c.currentIndent)
}

func (c *context) printfln(format string, a ...interface{}) {
	c.printf("%s%s\n", c.currentIndent, fmt.Sprintf(format, a...))
}

func (c *context) prologue() {
	c.printfln("#include \"blockchain.h\"")
	c.printfln("#include \"ckb_syscalls.h\"\n")
	c.printfln("typedef unsigned __int128 uint128_t;")
	c.printfln("\nint main() {")
	c.indent()
}

func (c *context) epilogue() {
	c.shrink()
	c.printfln("}")
}

func (c *context) indent() {
	c.currentIndent += "  "
}

func (c *context) shrink() {
	c.currentIndent = c.currentIndent[2:]
}

func (c *context) generateSlice(expr *ast.Value, start, end uint64) (int, variableType, error) {
	switch expr.GetT() {
	case ast.Value_BYTES:
		a, at, err := c.generateVariable(expr)
		if err != nil {
			return -1, varTypeEmpty, err
		}
		if start >= at.length || end > at.length {
			return -1, varTypeEmpty, fmt.Errorf("Invalid slice range!")
		}
		t := variableType{
			t:      varBytes,
			length: end - start,
		}
		i := c.newVariable(t)
		c.printfln("size_t v%d_length = %d;", i, t.length)
		c.printfln("uint8_t* v%d = &v%d[%d];", i, a, start)
		return i, t, nil
	}
	return -1, varTypeEmpty, fmt.Errorf("Invalid slice item type %s", expr.GetT().String())
}

func (c *context) loopStart(expr *ast.Value) (int, variableType, error) {
	map_funcs := []*ast.Value{}
	for expr.GetT() == ast.Value_MAP {
		map_funcs = append([]*ast.Value{expr.GetChildren()[0]}, map_funcs...)
		expr = expr.GetChildren()[1]
	}
	switch expr.GetT() {
	case ast.Value_GET_INPUTS:
		fallthrough
	case ast.Value_GET_OUTPUTS:
		sc, err := createSyscallContext(c, expr, map_funcs)
		if err != nil {
			return -1, varTypeEmpty, err
		}
		return sc.generate()
	}
	return -1, varTypeEmpty, fmt.Errorf("Invalid loop type: %s", expr.GetT().String())
}

func (c *context) loopEnd(expr *ast.Value) error {
	switch expr.GetT() {
	case ast.Value_MAP:
		return c.loopEnd(expr.GetChildren()[1])
	default:
		c.shrink()
		c.printfln("}")
		return nil
	}
}

func (c *context) castBytesToInteger(i int, t variableType) (int, variableType, error) {
	if t.t == varUint128 || t.t == varUint64 {
		return i, t, nil
	}
	if t.t != varBytes {
		return -1, varTypeEmpty, fmt.Errorf("Requested type %d is not bytes!", t.t)
	}
	switch t.length {
	case 8:
		j := c.newVariable(varTypeUint64)
		c.printfln("uint64_t v%d = *((uint64_t*) v%d);", j, i)
		return j, varTypeUint64, nil
	case 16:
		j := c.newVariable(varTypeUint128)
		c.printfln("uint128_t v%d = *((uint128_t*) v%d);", j, i)
		return j, varTypeUint128, nil
	default:
		return -1, varTypeEmpty, fmt.Errorf("Invalid byte length!")
	}
}

func (c *context) generateVariable(expr *ast.Value) (int, variableType, error) {
	switch expr.GetT() {
	case ast.Value_ARG:
		index := int(expr.GetU())
		if index < 0 || index >= len(c.args) {
			return -1, varTypeEmpty, fmt.Errorf("Invalid arg index!")
		}
		return c.args[index], c.variables[c.args[index]], nil
	case ast.Value_UINT64:
		i := c.newVariable(varTypeUint64)
		c.printfln("uint64_t v%d = %d;", i, expr.GetU())
		return i, varTypeUint64, nil
	case ast.Value_EQUAL:
		a, at, err := c.generateVariable(expr.GetChildren()[0])
		if err != nil {
			return -1, varTypeEmpty, err
		}
		b, bt, err := c.generateVariable(expr.GetChildren()[1])
		if err != nil {
			return -1, varTypeEmpty, err
		}
		if at != bt {
			return -1, varTypeEmpty, fmt.Errorf("Invalid equal types: %d %d", at, bt)
		}
		i := c.newVariable(varTypeBool)
		if at.t == varBytes {
			c.printfln("bool v%d = (v%d_length == v%d_length) && (memcmp(v%d, v%d, v%d_length) == 0);", i, a, b, a, b, a)
		} else {
			c.printfln("bool v%d = (v%d == v%d);", i, a, b)
		}
		return i, varTypeBool, nil
	case ast.Value_BYTES:
		l := len(expr.GetRaw())
		t := variableType{
			t:      varBytes,
			length: uint64(l),
		}
		i := c.newVariable(t)
		c.printfln("size_t v%d_length = %d;", i, l)
		c.print_indent()
		c.printf("uint8_t v%d[%d] = { ", i, l)
		for j, ch := range expr.GetRaw() {
			if j != 0 {
				c.printf(", ")
			}
			c.printf("0x%x", ch)
		}
		c.printf(" };\n")
		return i, t, nil
	case ast.Value_SLICE:
		children := expr.GetChildren()
		if children[0].GetT() != ast.Value_UINT64 ||
			children[1].GetT() != ast.Value_UINT64 {
			return -1, varTypeEmpty, fmt.Errorf("Invalid operand type to SLICE!")
		}
		return c.generateSlice(children[2], children[0].GetU(), children[1].GetU())
	case ast.Value_ADD:
		a, at, err := c.generateVariable(expr.GetChildren()[0])
		if err != nil {
			return -1, varTypeEmpty, err
		}
		a, at, err = c.castBytesToInteger(a, at)
		if err != nil {
			return -1, varTypeEmpty, err
		}
		b, bt, err := c.generateVariable(expr.GetChildren()[1])
		if err != nil {
			return -1, varTypeEmpty, err
		}
		b, bt, err = c.castBytesToInteger(b, bt)
		if err != nil {
			return -1, varTypeEmpty, err
		}
		if at.t == varUint128 || bt.t == varUint128 {
			if at.t != varUint128 {
				newA := c.newVariable(varTypeUint128)
				c.printfln("uint128_t v%d = (uint128_t) v%d;", newA, a)
				a = newA
				at = varTypeUint128
			}
			if bt.t != varUint128 {
				newB := c.newVariable(varTypeUint128)
				c.printfln("uint128_t v%d = (uint128_t) v%d;", newB, b)
				b = newB
				bt = varTypeUint128
			}
		}
		if at.t != bt.t {
			return -1, varTypeEmpty, fmt.Errorf("ADD operand types do not match!")
		}
		i := c.newVariable(at)
		c.printfln("%s v%d = v%d + v%d;", at.cType(), i, a, b)
		return i, at, nil
	case ast.Value_REDUCE:
		initial, initialType, err := c.generateVariable(expr.GetChildren()[1])
		if err != nil {
			return -1, varTypeEmpty, err
		}
		loopVar, _, err := c.loopStart(expr.GetChildren()[2])
		if err != nil {
			return -1, varTypeEmpty, err
		}
		c.args = append([]int{initial, loopVar}, c.args...)
		reducedValue, reducedValueType, err := c.generateVariable(expr.GetChildren()[0])
		c.args = c.args[2:]
		if err != nil {
			return -1, varTypeEmpty, err
		}
		// TODO: see how we can handle this later
		if reducedValueType.t != initialType.t {
			return -1, varTypeEmpty, fmt.Errorf("Cannot assign type %d to type %d!", reducedValueType.t, initialType.t)
		}
		if reducedValueType.t == varBytes {
			if reducedValueType.length != initialType.length {
				return -1, varTypeEmpty, fmt.Errorf("Bytes length mismatch!")
			}
			c.printfln("memcpy(v%d, v%d, v%d_length);", initial, reducedValue, initial)
		} else {
			c.printfln("v%d = v%d;", initial, reducedValue)
		}
		err = c.loopEnd(expr.GetChildren()[2])
		if err != nil {
			return -1, varTypeEmpty, err
		}
		return initial, initialType, nil
	}
	return -1, varTypeEmpty, fmt.Errorf("Invalid value type %s", expr.GetT().String())
}

func Generate(expr *ast.Value, writer io.Writer) error {
	c := newContext(writer)
	c.prologue()
	i, t, err := c.generateVariable(expr)
	if err != nil {
		return err
	}
	if t.t != varBool {
		return fmt.Errorf("Terminated variable must be a bool!")
	}
	c.printfln("if (v%d) { return 0; } else { return %d; }", i, c.newErrorCode())
	c.epilogue()
	return c.err
}
