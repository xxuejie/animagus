package executor

import (
	"fmt"

	internal_ast "github.com/xxuejie/animagus/internal/ast"
	"github.com/xxuejie/animagus/pkg/ast"
)

type Environment interface {
	Arg(i int) *internal_ast.Value
	Param(i int) *internal_ast.Value
	QueryCell(query *ast.List) ([]*internal_ast.Value, error)
}

func Execute(expr *ast.Value, e Environment) (*ast.Value, error) {
	internalValue, err := evaluateValue(&internal_ast.Value{
		Value: expr,
	}, e)
	if err != nil {
		return nil, err
	}
	if internalValue.Value == nil {
		return nil, fmt.Errorf("Execution result cannot be expressed by AST value: %v", *internalValue)
	}
	return internalValue.Value, nil
}

func evaluateValue(expr *internal_ast.Value, e Environment) (*internal_ast.Value, error) {
	if expr.IsTerminal() {
		return expr, nil
	}
	switch expr.Value.GetT() {
	case ast.Value_ARG:
		index := expr.GetIndex()
		arg := e.Arg(index)
		if arg == nil {
			return nil, fmt.Errorf("Cannot find arg index %d!", index)
		}
		return arg, nil
	case ast.Value_PARAM:
		index := expr.GetIndex()
		param := e.Param(index)
		if param == nil {
			return nil, fmt.Errorf("Cannot find param index %d!", index)
		}
		return param, nil
	case ast.Value_APPLY:
		children, err := evaluateAstValues(expr.Value.GetChildren(), e)
		if err != nil {
			return nil, err
		}
		if len(children) < 1 {
			return nil, fmt.Errorf("Not enough arguments for apply!")
		}
		if children[0].Value.GetT() != ast.Value_OP {
			return nil, fmt.Errorf("Invalid apply type: %s", children[0].Value.GetT().String())
		}
		return evaluateOp(children[0].Value.GetOp(), children[1:], e)
	case ast.Value_REDUCE:
		if len(expr.Value.GetChildren()) != 2 {
			return nil, fmt.Errorf("Invalid number of arguments for reduce!")
		}
		if expr.Value.GetL() == nil {
			return nil, fmt.Errorf("List must be provided for reduce!")
		}
		f := &internal_ast.Value{
			Value: expr.Value.GetChildren()[0],
		}
		currentValue, err := evaluateValue(&internal_ast.Value{
			Value: expr.Value.GetChildren()[1],
		}, e)
		if err != nil {
			return nil, err
		}
		list, err := evaluateList(expr.Value.GetL(), e)
		if err != nil {
			return nil, err
		}
		for _, value := range list {
			currentValue, err = evaluateValue(f, &prependEnvironment{
				e: e,
				args: []*internal_ast.Value{
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
	return nil, fmt.Errorf("Invalid value type: %s", expr.Value.GetT().String())
}

func evaluateAstValues(values []*ast.Value, e Environment) ([]*internal_ast.Value, error) {
	evaluatedValues := make([]*internal_ast.Value, len(values))
	var err error
	for i, value := range values {
		evaluatedValues[i], err = evaluateValue(&internal_ast.Value{
			Value: value,
		}, e)
		if err != nil {
			return nil, err
		}
	}
	return evaluatedValues, nil
}

func evaluateOp(op ast.Op, operands []*internal_ast.Value, e Environment) (*internal_ast.Value, error) {
	switch op {
	case ast.Op_GET:
		if len(operands) != 2 {
			return nil, fmt.Errorf("Invalid number of operands to GET")
		}
		if operands[0].Value.GetT() != ast.Value_FIELD {
			return nil, fmt.Errorf("First operand to GET must be a field!")
		}
		return evaluateOpGet(operands[0].Value.GetField(), operands[1], e)
	case ast.Op_EQUAL:
		if len(operands) != 2 {
			return nil, fmt.Errorf("Invalid number of operands to EQUAL")
		}
		result, err := operands[0].Equal(operands[1])
		if err != nil {
			return nil, err
		}
		return &internal_ast.Value{
			Value: &ast.Value{
				T: ast.Value_BOOL,
				Primitive: &ast.Value_B{
					B: result,
				},
			},
		}, nil
	case ast.Op_PLUS:
		if len(operands) != 2 {
			return nil, fmt.Errorf("Invalid number of operands to PLUS")
		}
		if operands[0].Value == nil || operands[1].Value == nil ||
			operands[0].Value.GetT() != ast.Value_UINT64 ||
			operands[1].Value.GetT() != ast.Value_UINT64 {
			return nil, fmt.Errorf("Both operands must be uint64 in PLUS!")
		}
		return &internal_ast.Value{
			Value: &ast.Value{
				T: ast.Value_UINT64,
				Primitive: &ast.Value_U{
					U: operands[0].Value.GetU() + operands[1].Value.GetU(),
				},
			},
		}, nil
	case ast.Op_AND:
		if len(operands) == 0 {
			return nil, fmt.Errorf("Invalid number of operands to AND")
		}
		result := true
		for _, operand := range operands {
			if operand.Value.GetT() != ast.Value_BOOL {
				return nil, fmt.Errorf("Invalid operand type %s to AND!", operand.Value.GetT().String())
			}
			if !operand.Value.GetB() {
				result = false
				break
			}
		}
		return &internal_ast.Value{
			Value: &ast.Value{
				T: ast.Value_BOOL,
				Primitive: &ast.Value_B{
					B: result,
				},
			},
		}, nil
	}
	return nil, fmt.Errorf("Invalid op: %s", op.String())
}

func evaluateOpGet(field ast.Field, value *internal_ast.Value, e Environment) (*internal_ast.Value, error) {
	if value.Value != nil && value.Value.GetT() == ast.Value_NIL {
		// Running GET on NIL values always results in NIL
		return value, nil
	}
	switch field {
	case ast.Field_CAPACITY:
		if value.Cell == nil {
			return nil, fmt.Errorf("Cannot fetch capacity on non-cell value!")
		}
		return &internal_ast.Value{
			Value: &ast.Value{
				T: ast.Value_UINT64,
				Primitive: &ast.Value_U{
					U: uint64(value.Cell.Capacity),
				},
			},
		}, nil
	case ast.Field_DATA:
		if value.CellData == nil {
			return nil, fmt.Errorf("Cannot fetch cell data on non-cell value!")
		}
		return &internal_ast.Value{
			Value: &ast.Value{
				T: ast.Value_BYTES,
				Primitive: &ast.Value_Raw{
					Raw: *value.CellData,
				},
			},
		}, nil
	case ast.Field_LOCK:
		if value.Cell == nil {
			return nil, fmt.Errorf("Cannot fetch lock on non-cell value!")
		}
		return &internal_ast.Value{
			Script: &value.Cell.Lock,
		}, nil
	case ast.Field_TYPE:
		if value.Cell == nil {
			return nil, fmt.Errorf("Cannot fetch type on non-cell value!")
		}
		if value.Cell.Type != nil {
			return &internal_ast.Value{
				Script: value.Cell.Type,
			}, nil
		} else {
			return &internal_ast.Value{
				Value: &ast.Value{
					T: ast.Value_NIL,
				},
			}, nil
		}
	case ast.Field_CODE_HASH:
		if value.Script == nil {
			return nil, fmt.Errorf("Cannot fetch code hash on non-script value!")
		}
		return &internal_ast.Value{
			Value: &ast.Value{
				T: ast.Value_BYTES,
				Primitive: &ast.Value_Raw{
					Raw: value.Script.CodeHash[:],
				},
			},
		}, nil
	case ast.Field_HASH_TYPE:
		if value.Script == nil {
			return nil, fmt.Errorf("Cannot fetch hash type on non-script value!")
		}
		return &internal_ast.Value{
			Value: &ast.Value{
				T: ast.Value_UINT64,
				Primitive: &ast.Value_U{
					U: uint64(value.Script.HashType),
				},
			},
		}, nil
	case ast.Field_ARGS:
		if value.Script == nil {
			return nil, fmt.Errorf("Cannot fetch args on non-script value!")
		}
		return &internal_ast.Value{
			Value: &ast.Value{
				T: ast.Value_BYTES,
				Primitive: &ast.Value_Raw{
					Raw: value.Script.Args,
				},
			},
		}, nil
	}
	return nil, fmt.Errorf("Invalid get field: %s", field.String())
}

func evaluateList(list *ast.List, e Environment) ([]*internal_ast.Value, error) {
	switch list.GetT() {
	case ast.List_MAP:
		if len(list.GetChildren()) != 1 {
			return nil, fmt.Errorf("Invalid number of lists for map!")
		}
		if len(list.GetValues()) != 1 {
			return nil, fmt.Errorf("Invalid number of values for map!")
		}
		f := &internal_ast.Value{
			Value: list.GetValues()[0],
		}
		list, err := evaluateList(list.GetChildren()[0], e)
		if err != nil {
			return nil, err
		}
		results := make([]*internal_ast.Value, len(list))
		for i, value := range list {
			results[i], err = evaluateValue(f, &prependEnvironment{
				e:    e,
				args: []*internal_ast.Value{value},
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
