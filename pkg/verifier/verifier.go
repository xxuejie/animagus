package verifier

import (
	"fmt"
	"unicode/utf8"

	"github.com/xxuejie/animagus/pkg/ast"
)

func Verify(expr *ast.Value) error {
	for i, child := range expr.GetChildren() {
		if err := Verify(child); err != nil {
			return fmt.Errorf("ERROR occured for argument %d in %s: %s", i, expr.GetT().String(), err)
		}
	}

	switch expr.GetT() {
	case ast.Value_NIL:
	case ast.Value_UINT64:
		if _, ok := expr.GetPrimitive().(*ast.Value_U); !ok {
			return fmt.Errorf("UINT64 type must have u set!")
		}
		if len(expr.GetChildren()) > 0 {
			return fmt.Errorf("UINT64 type should not have children!")
		}
	case ast.Value_BOOL:
		if _, ok := expr.GetPrimitive().(*ast.Value_B); !ok {
			return fmt.Errorf("BOOL type must have b set!")
		}
		if len(expr.GetChildren()) > 0 {
			return fmt.Errorf("BOOL type should not have children!")
		}
	case ast.Value_BYTES:
		if _, ok := expr.GetPrimitive().(*ast.Value_Raw); !ok {
			return fmt.Errorf("BYTES type must have raw set!")
		}
		if len(expr.GetChildren()) > 0 {
			return fmt.Errorf("BYTES type should not have children!")
		}
	case ast.Value_ERROR:
		raw, ok := expr.GetPrimitive().(*ast.Value_Raw)
		if !ok {
			return fmt.Errorf("ERROR type must have raw set!")
		}
		if !utf8.Valid(raw.Raw) {
			return fmt.Errorf("ERROR type value must be a utf8 formatted string!")
		}
		if len(expr.GetChildren()) > 0 {
			return fmt.Errorf("ERROR type should not have children!")
		}
	case ast.Value_ARG:
		if _, ok := expr.GetPrimitive().(*ast.Value_U); !ok {
			return fmt.Errorf("ARG type must have u set!")
		}
		if len(expr.GetChildren()) > 0 {
			return fmt.Errorf("ARG type should not have children!")
		}
	case ast.Value_PARAM:
		if _, ok := expr.GetPrimitive().(*ast.Value_U); !ok {
			return fmt.Errorf("PARAM type must have u set!")
		}
		if len(expr.GetChildren()) > 0 {
			return fmt.Errorf("PARAM type should not have children!")
		}
	case ast.Value_OUT_POINT:
		if len(expr.GetChildren()) != 2 {
			return fmt.Errorf("Invalid number of arguments for %s!", expr.GetT().String())
		}
	case ast.Value_CELL_INPUT:
		if len(expr.GetChildren()) != 2 {
			return fmt.Errorf("Invalid number of arguments for %s!", expr.GetT().String())
		}
	case ast.Value_CELL_DEP:
		if len(expr.GetChildren()) != 2 {
			return fmt.Errorf("Invalid number of arguments for %s!", expr.GetT().String())
		}
	case ast.Value_SCRIPT:
		if len(expr.GetChildren()) != 3 {
			return fmt.Errorf("Invalid number of arguments for %s!", expr.GetT().String())
		}
	case ast.Value_CELL:
		if len(expr.GetChildren()) != 4 || len(expr.GetChildren()) != 6 {
			return fmt.Errorf("Invalid number of arguments for %s!", expr.GetT().String())
		}
	case ast.Value_TRANSACTION:
		if len(expr.GetChildren()) != 3 {
			return fmt.Errorf("Invalid number of arguments for %s!", expr.GetT().String())
		}
	case ast.Value_HEADER:
		if len(expr.GetChildren()) != 10 {
			return fmt.Errorf("Invalid number of arguments for %s!", expr.GetT().String())
		}
	case ast.Value_APPLY:
		if len(expr.GetChildren()) < 1 {
			return fmt.Errorf("Invalid number of arguments for %s!", expr.GetT().String())
		}
		if err := verifyFuncArgs(expr.GetChildren()[0], len(expr.GetChildren())-1); err != nil {
			return fmt.Errorf("ERROR occured verifying APPLY argument length: %s", err)
		}
	case ast.Value_REDUCE:
		if len(expr.GetChildren()) != 3 {
			return fmt.Errorf("Invalid number of arguments for %s!", expr.GetT().String())
		}
		if err := verifyFuncArgs(expr.GetChildren()[0], 2); err != nil {
			return fmt.Errorf("ERROR occured verifying REDUCE argument length: %s", err)
		}
	case ast.Value_LIST:
	case ast.Value_QUERY_CELLS:
		if len(expr.GetChildren()) != 1 {
			return fmt.Errorf("Invalid number of arguments for %s!", expr.GetT().String())
		}
		if err := verifyFuncArgs(expr.GetChildren()[0], 1); err != nil {
			return fmt.Errorf("ERROR occured verifying QUERY_CELLS argument length: %s", err)
		}
	case ast.Value_MAP:
		if len(expr.GetChildren()) != 2 {
			return fmt.Errorf("Invalid number of arguments for %s!", expr.GetT().String())
		}
		if err := verifyFuncArgs(expr.GetChildren()[0], 1); err != nil {
			return fmt.Errorf("ERROR occured verifying MAP function: %s", err)
		}
		if !isList(expr.GetChildren()[1]) {
			return fmt.Errorf("Argument 1 of MAP is not a list: %s", expr.GetChildren()[1].GetT().String())
		}
	case ast.Value_FILTER:
		if len(expr.GetChildren()) != 2 {
			return fmt.Errorf("Invalid number of arguments for %s!", expr.GetT().String())
		}
		if err := verifyFuncArgs(expr.GetChildren()[0], 1); err != nil {
			return fmt.Errorf("ERROR occured verifying FILTER function: %s", err)
		}
		if !isList(expr.GetChildren()[1]) {
			return fmt.Errorf("Argument 1 of FILTER is not a list: %s", expr.GetChildren()[1].GetT().String())
		}
	case ast.Value_GET_CAPACITY:
		fallthrough
	case ast.Value_GET_DATA:
		fallthrough
	case ast.Value_GET_LOCK:
		fallthrough
	case ast.Value_GET_TYPE:
		fallthrough
	case ast.Value_GET_DATA_HASH:
		if len(expr.GetChildren()) != 1 {
			return fmt.Errorf("Invalid number of arguments for %s!", expr.GetT().String())
		}
	case ast.Value_GET_OUT_POINT:
		if len(expr.GetChildren()) != 1 {
			return fmt.Errorf("Invalid number of arguments for %s!", expr.GetT().String())
		}
		if len(expr.GetChildren()[0].GetChildren()) < 5 {
			return fmt.Errorf("Specified cell does not provide OutPoint!")
		}
	case ast.Value_GET_CODE_HASH:
		fallthrough
	case ast.Value_GET_HASH_TYPE:
		fallthrough
	case ast.Value_GET_ARGS:
		if len(expr.GetChildren()) != 1 {
			return fmt.Errorf("Invalid number of arguments for %s!", expr.GetT().String())
		}
	case ast.Value_GET_CELL_DEPS:
		fallthrough
	case ast.Value_GET_HEADER_DEPS:
		fallthrough
	case ast.Value_GET_INPUTS:
		fallthrough
	case ast.Value_GET_OUTPUTS:
		fallthrough
	case ast.Value_GET_WITNESSES:
		if len(expr.GetChildren()) != 1 {
			return fmt.Errorf("Invalid number of arguments for %s!", expr.GetT().String())
		}
	case ast.Value_GET_COMPACT_TARGET:
		fallthrough
	case ast.Value_GET_TIMESTAMP:
		fallthrough
	case ast.Value_GET_NUMBER:
		fallthrough
	case ast.Value_GET_EPOCH:
		fallthrough
	case ast.Value_GET_PARENT_HASH:
		fallthrough
	case ast.Value_GET_TRANSACTIONS_ROOT:
		fallthrough
	case ast.Value_GET_PROPOSALS_HASH:
		fallthrough
	case ast.Value_GET_UNCLES_HASH:
		fallthrough
	case ast.Value_GET_DAO:
		fallthrough
	case ast.Value_GET_NONCE:
		if len(expr.GetChildren()) != 1 {
			return fmt.Errorf("Invalid number of arguments for %s!", expr.GetT().String())
		}
	case ast.Value_GET_HEADER:
		if len(expr.GetChildren()) != 1 {
			return fmt.Errorf("Invalid number of arguments for %s!", expr.GetT().String())
		}
		if len(expr.GetChildren()[0].GetChildren()) < 6 {
			return fmt.Errorf("Specified cell does not provide Header!")
		}
	case ast.Value_HASH:
		if len(expr.GetChildren()) != 1 {
			return fmt.Errorf("Invalid number of arguments for %s!", expr.GetT().String())
		}
	case ast.Value_SERIALIZE_TO_CORE:
		fallthrough
	case ast.Value_SERIALIZE_TO_JSON:
		if len(expr.GetChildren()) != 1 {
			return fmt.Errorf("Invalid number of arguments for %s!", expr.GetT().String())
		}
		value := expr.GetChildren()[0]
		switch value.GetT() {
		case ast.Value_SCRIPT:
		case ast.Value_HEADER:
		case ast.Value_TRANSACTION:
		default:
			return fmt.Errorf("Cannot perform %s operation on %s", expr.GetT().String(), value.GetT().String())
		}
	case ast.Value_NOT:
		if len(expr.GetChildren()) != 1 {
			return fmt.Errorf("Invalid number of arguments for %s!", expr.GetT().String())
		}
	case ast.Value_AND:
		fallthrough
	case ast.Value_OR:
		if len(expr.GetChildren()) == 0 {
			return fmt.Errorf("Invalid number of arguments for %s!", expr.GetT().String())
		}
	case ast.Value_EQUAL:
		if len(expr.GetChildren()) != 2 {
			return fmt.Errorf("Invalid number of arguments for %s!", expr.GetT().String())
		}
	case ast.Value_LEN:
		if len(expr.GetChildren()) != 1 {
			return fmt.Errorf("Invalid number of arguments for %s!", expr.GetT().String())
		}
	case ast.Value_SLICE:
		if len(expr.GetChildren()) != 3 {
			return fmt.Errorf("Invalid number of arguments for %s!", expr.GetT().String())
		}
	case ast.Value_INDEX:
		if len(expr.GetChildren()) != 2 {
			return fmt.Errorf("Invalid number of arguments for %s!", expr.GetT().String())
		}
	case ast.Value_LESS:
		fallthrough
	case ast.Value_ADD:
		fallthrough
	case ast.Value_SUBTRACT:
		fallthrough
	case ast.Value_MULTIPLY:
		fallthrough
	case ast.Value_DIVIDE:
		fallthrough
	case ast.Value_MOD:
		if len(expr.GetChildren()) != 2 {
			return fmt.Errorf("Invalid number of arguments for %s!", expr.GetT().String())
		}
	case ast.Value_COND:
		if len(expr.GetChildren()) != 3 {
			return fmt.Errorf("Invalid number of arguments for %s!", expr.GetT().String())
		}
	case ast.Value_TAIL_RECURSION:
		if len(expr.GetChildren()) == 0 {
			return fmt.Errorf("To keep recursion going, at least one argument must be provided!")
		}
	default:
		return fmt.Errorf("Invalid value type: %s", expr.GetT().String())
	}
	return nil
}

func isList(l *ast.Value) bool {
	switch l.GetT() {
	case ast.Value_LIST:
	case ast.Value_MAP:
	case ast.Value_FILTER:
	case ast.Value_QUERY_CELLS:
	default:
		return false
	}
	return true
}

func verifyFuncArgs(f *ast.Value, args int) error {
	remainingValues := []*ast.Value{f}
	for len(remainingValues) > 0 {
		value := remainingValues[0]
		remainingValues = remainingValues[1:]

		if value.GetT() == ast.Value_ARG {
			i := int(f.GetU())
			if i < 0 || i >= args {
				return fmt.Errorf("Invalid argument index: %d", value.GetU())
			}
		} else {
			remainingValues = append(remainingValues, value.GetChildren()...)
		}
	}
	return nil
}
