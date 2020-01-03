package executor

import (
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/xxuejie/animagus/pkg/ast"
)

type Environment interface {
	Arg(i int) *ast.Value
	Param(i int) *ast.Value
	IndexParam(i int, value *ast.Value) error
	QueryCell(query *ast.List) ([]*ast.Value, error)
}

func Execute(expr *ast.Value, e Environment) (*ast.Value, error) {
	return evaluateValue(expr, e)
}

func evaluateValue(expr *ast.Value, e Environment) (*ast.Value, error) {
	// Primitive value
	if expr.GetT() < ast.Value_ARG {
		return expr, nil
	}
	switch expr.GetT() {
	case ast.Value_ARG:
		index := int(expr.GetU())
		arg := e.Arg(index)
		if arg == nil {
			return nil, fmt.Errorf("Cannot find arg index %d!", index)
		}
		return arg, nil
	case ast.Value_PARAM:
		index := int(expr.GetU())
		param := e.Param(index)
		if param == nil {
			return nil, fmt.Errorf("Cannot find param index %d!", index)
		}
		return param, nil
	case ast.Value_APPLY:
		children, err := evaluateAstValues(expr.GetChildren(), e)
		if err != nil {
			return nil, err
		}
		if len(children) < 1 {
			return nil, fmt.Errorf("Not enough arguments for apply!")
		}
		if children[0].GetT() != ast.Value_OP {
			return nil, fmt.Errorf("Invalid apply type: %s", children[0].GetT().String())
		}
		return evaluateOp(children[0].GetOp(), children[1:], e)
	case ast.Value_REDUCE:
		if len(expr.GetChildren()) != 2 {
			return nil, fmt.Errorf("Invalid number of arguments for reduce!")
		}
		if expr.GetL() == nil {
			return nil, fmt.Errorf("List must be provided for reduce!")
		}
		f := expr.GetChildren()[0]
		currentValue, err := evaluateValue(expr.GetChildren()[1], e)
		if err != nil {
			return nil, err
		}
		list, err := evaluateList(expr.GetL(), e)
		if err != nil {
			return nil, err
		}
		for _, value := range list {
			currentValue, err = evaluateValue(f, &prependEnvironment{
				e: e,
				args: []*ast.Value{
					currentValue,
					value,
				},
			})
			if err != nil {
				return nil, err
			}
		}
		return currentValue, nil
	}
	return nil, fmt.Errorf("Invalid value type: %s", expr.GetT().String())
}

func evaluateAstValues(values []*ast.Value, e Environment) ([]*ast.Value, error) {
	evaluatedValues := make([]*ast.Value, len(values))
	var err error
	for i, value := range values {
		evaluatedValues[i], err = evaluateValue(value, e)
		if err != nil {
			return nil, err
		}
	}
	return evaluatedValues, nil
}

func evaluateOp(op ast.Op, operands []*ast.Value, e Environment) (*ast.Value, error) {
	switch op {
	case ast.Op_GET:
		if len(operands) != 2 {
			return nil, fmt.Errorf("Invalid number of operands to GET")
		}
		if operands[0].GetT() != ast.Value_FIELD {
			return nil, fmt.Errorf("First operand to GET must be a field!")
		}
		return evaluateOpGet(operands[0].GetField(), operands[1], e)
	case ast.Op_EQUAL:
		if len(operands) != 2 {
			return nil, fmt.Errorf("Invalid number of operands to EQUAL")
		}
		var result bool
		if operands[0].GetT() == ast.Value_PARAM {
			if err := e.IndexParam(int(operands[0].GetU()), operands[1]); err != nil {
				return nil, err
			}
			result = true
		} else if operands[1].GetT() == ast.Value_PARAM {
			if err := e.IndexParam(int(operands[1].GetU()), operands[0]); err != nil {
				return nil, err
			}
			result = true
		} else {
			result = proto.Equal(operands[0], operands[1])
		}
		return &ast.Value{
			T: ast.Value_BOOL,
			Primitive: &ast.Value_B{
				B: result,
			},
		}, nil
	case ast.Op_PLUS:
		if len(operands) != 2 {
			return nil, fmt.Errorf("Invalid number of operands to PLUS")
		}
		if operands[0].GetT() != ast.Value_UINT64 ||
			operands[1].GetT() != ast.Value_UINT64 {
			return nil, fmt.Errorf("Both operands must be uint64 in PLUS!")
		}
		return &ast.Value{
			T: ast.Value_UINT64,
			Primitive: &ast.Value_U{
				U: operands[0].GetU() + operands[1].GetU(),
			},
		}, nil
	case ast.Op_AND:
		if len(operands) == 0 {
			return nil, fmt.Errorf("Invalid number of operands to AND")
		}
		result := true
		for _, operand := range operands {
			if operand.GetT() != ast.Value_BOOL {
				return nil, fmt.Errorf("Invalid operand type %s to AND!", operand.GetT().String())
			}
			if !operand.GetB() {
				result = false
				break
			}
		}
		return &ast.Value{
			T: ast.Value_BOOL,
			Primitive: &ast.Value_B{
				B: result,
			},
		}, nil
	}
	return nil, fmt.Errorf("Invalid op: %s", op.String())
}

func evaluateOpGet(field ast.Field, value *ast.Value, e Environment) (*ast.Value, error) {
	if value.GetT() == ast.Value_NIL {
		// Running GET on NIL values always results in NIL
		return value, nil
	}
	switch field {
	case ast.Field_CAPACITY:
		if value.GetT() != ast.Value_CELL {
			return nil, fmt.Errorf("Cannot fetch capacity on non-cell value!")
		}
		if len(value.GetChildren()) != 4 {
			return nil, fmt.Errorf("Invalid number of cell items!")
		}
		return value.GetChildren()[0], nil
	case ast.Field_LOCK:
		if value.GetT() != ast.Value_CELL {
			return nil, fmt.Errorf("Cannot fetch lock on non-cell value!")
		}
		if len(value.GetChildren()) != 4 {
			return nil, fmt.Errorf("Invalid number of cell items!")
		}
		return value.GetChildren()[1], nil
	case ast.Field_TYPE:
		if value.GetT() != ast.Value_CELL {
			return nil, fmt.Errorf("Cannot fetch type on non-cell value!")
		}
		if len(value.GetChildren()) != 4 {
			return nil, fmt.Errorf("Invalid number of cell items!")
		}
		return value.GetChildren()[2], nil
	case ast.Field_DATA:
		if value.GetT() != ast.Value_CELL {
			return nil, fmt.Errorf("Cannot fetch data on non-cell value!")
		}
		if len(value.GetChildren()) != 4 {
			return nil, fmt.Errorf("Invalid number of cell items!")
		}
		return value.GetChildren()[3], nil
	case ast.Field_CODE_HASH:
		if value.GetT() != ast.Value_SCRIPT {
			return nil, fmt.Errorf("Cannot fetch code hash on non-script value!")
		}
		if len(value.GetChildren()) != 3 {
			return nil, fmt.Errorf("Invalid number of script items!")
		}
		return value.GetChildren()[0], nil
	case ast.Field_HASH_TYPE:
		if value.GetT() != ast.Value_SCRIPT {
			return nil, fmt.Errorf("Cannot fetch hash type on non-script value!")
		}
		if len(value.GetChildren()) != 3 {
			return nil, fmt.Errorf("Invalid number of script items!")
		}
		return value.GetChildren()[1], nil
	case ast.Field_ARGS:
		if value.GetT() != ast.Value_SCRIPT {
			return nil, fmt.Errorf("Cannot fetch args on non-script value!")
		}
		if len(value.GetChildren()) != 3 {
			return nil, fmt.Errorf("Invalid number of script items!")
		}
		return value.GetChildren()[2], nil
	}
	return nil, fmt.Errorf("Invalid get field: %s", field.String())
}

func evaluateList(list *ast.List, e Environment) ([]*ast.Value, error) {
	switch list.GetT() {
	case ast.List_MAP:
		if len(list.GetChildren()) != 1 {
			return nil, fmt.Errorf("Invalid number of lists for map!")
		}
		if len(list.GetValues()) != 1 {
			return nil, fmt.Errorf("Invalid number of values for map!")
		}
		f := list.GetValues()[0]
		list, err := evaluateList(list.GetChildren()[0], e)
		if err != nil {
			return nil, err
		}
		results := make([]*ast.Value, len(list))
		for i, value := range list {
			results[i], err = evaluateValue(f, &prependEnvironment{
				e:    e,
				args: []*ast.Value{value},
			})
			if err != nil {
				return nil, err
			}
		}
		return results, nil
	case ast.List_QUERY_CELLS:
		return e.QueryCell(list)
	}
	return nil, fmt.Errorf("Invalid list type: %s", list.GetT().String())
}
