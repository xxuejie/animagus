package executor

import (
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/golang/protobuf/proto"
	"github.com/xxuejie/animagus/pkg/ast"
	"github.com/xxuejie/animagus/pkg/rpctypes"
)

type Environment interface {
	Arg(i int) *ast.Value
	Param(i int) *ast.Value
	IndexParam(i int, value *ast.Value) error
	QueryCell(query *ast.Value) ([]*ast.Value, error)
}

func Execute(expr *ast.Value, e Environment) (*ast.Value, error) {
	return evaluateValue(expr, e)
}

func isPrimitive(expr *ast.Value) bool {
	return expr.GetT() < ast.Value_ARG
}

func isOp(expr *ast.Value) bool {
	return expr.GetT() >= ast.Value_HASH
}

func isGetOp(expr *ast.Value) bool {
	return expr.GetT() >= ast.Value_GET_CAPACITY && expr.GetT() < ast.Value_HASH
}

func isListOp(expr *ast.Value) bool {
	return expr.GetT() >= ast.Value_LIST && expr.GetT() < ast.Value_GET_CAPACITY
}

func evaluateValue(expr *ast.Value, e Environment) (*ast.Value, error) {
	// Primitive value
	if isPrimitive(expr) {
		return expr, nil
	}
	if isListOp(expr) {
		list, err := evaluateList(expr, e)
		if err != nil {
			return nil, err
		}
		return &ast.Value{
			T:        ast.Value_LIST,
			Children: list,
		}, nil
	}
	if isOp(expr) {
		children, err := evaluateAstValues(expr.GetChildren(), e)
		if err != nil {
			return nil, err
		}
		return evaluateOp(expr.GetT(), children, e)
	}
	if isGetOp(expr) {
		if len(expr.GetChildren()) != 1 {
			return nil, fmt.Errorf("Invalid number of operands to GET")
		}
		operand, err := evaluateValue(expr.GetChildren()[0], e)
		if err != nil {
			return nil, err
		}
		return evaluateOpGet(expr.GetT(), operand, e)
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
	case ast.Value_TRANSACTION:
		if len(expr.GetChildren()) != 3 {
			return nil, fmt.Errorf("Not enough arguments for transaction!")
		}
		inputCells, err := evaluateList(expr.GetChildren()[0], e)
		if err != nil {
			return nil, err
		}
		inputs := make([]*ast.Value, len(inputCells))
		for i, inputCell := range inputCells {
			if inputCell.GetT() != ast.Value_CELL ||
				len(inputCell.GetChildren()) != 5 {
				return nil, fmt.Errorf("Invalid input cell!")
			}
			inputs[i] = &ast.Value{
				T: ast.Value_CELL_INPUT,
				Children: []*ast.Value{
					inputCell.GetChildren()[4],
					&ast.Value{
						T: ast.Value_UINT64,
						Primitive: &ast.Value_U{
							U: 0,
						},
					},
				},
			}
		}
		outputs, err := evaluateValue(expr.GetChildren()[1], e)
		if err != nil {
			return nil, err
		}
		depValues, err := evaluateList(expr.GetChildren()[2], e)
		if err != nil {
			return nil, err
		}
		deps := make([]*ast.Value, len(depValues))
		for i, depValue := range depValues {
			switch depValue.GetT() {
			case ast.Value_CELL_DEP:
				deps[i] = depValue
			case ast.Value_CELL:
				if len(depValue.GetChildren()) != 5 {
					return nil, fmt.Errorf("Invalid dep cell!")
				}
				deps[i] = &ast.Value{
					T: ast.Value_CELL_DEP,
					Children: []*ast.Value{
						depValue.GetChildren()[4],
						&ast.Value{
							T: ast.Value_UINT64,
							Primitive: &ast.Value_U{
								U: 0,
							},
						},
					},
				}
			default:
				return nil, fmt.Errorf("Invalid dep type: %s", depValue.GetT().String())
			}
		}
		return &ast.Value{
			T: ast.Value_TRANSACTION,
			Children: []*ast.Value{
				&ast.Value{
					T:        ast.Value_LIST,
					Children: inputs,
				},
				outputs,
				&ast.Value{
					T:        ast.Value_LIST,
					Children: deps,
				},
			},
		}, nil
	case ast.Value_CELL:
		value, err := evaluateChildren(expr, e)
		if err != nil {
			return nil, err
		}
		err = ast.IsValidCell(value)
		if err != nil {
			return nil, err
		}
		return value, nil
	case ast.Value_SCRIPT:
		value, err := evaluateChildren(expr, e)
		if err != nil {
			return nil, err
		}
		err = ast.IsValidScript(value)
		if err != nil {
			return nil, err
		}
		return value, nil
	case ast.Value_CELL_DEP:
		value, err := evaluateChildren(expr, e)
		if err != nil {
			return nil, err
		}
		err = ast.IsValidCellDep(value)
		if err != nil {
			return nil, err
		}
		return value, nil
	case ast.Value_OUT_POINT:
		value, err := evaluateChildren(expr, e)
		if err != nil {
			return nil, err
		}
		err = ast.IsValidOutPoint(value)
		if err != nil {
			return nil, err
		}
		return value, nil
	case ast.Value_APPLY:
		if len(expr.GetChildren()) < 1 {
			return nil, fmt.Errorf("Not enough arguments for apply!")
		}
		args, err := evaluateAstValues(expr.GetChildren()[1:], e)
		if err != nil {
			return nil, err
		}
		return evaluateValue(expr.GetChildren()[0], &prependEnvironment{
			e:    e,
			args: args,
		})
	case ast.Value_REDUCE:
		if len(expr.GetChildren()) != 3 {
			return nil, fmt.Errorf("Invalid number of arguments for reduce!")
		}
		f := expr.GetChildren()[0]
		currentValue, err := evaluateValue(expr.GetChildren()[1], e)
		if err != nil {
			return nil, err
		}
		list, err := evaluateList(expr.GetChildren()[2], e)
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

func evaluateChildren(value *ast.Value, e Environment) (*ast.Value, error) {
	children, err := evaluateAstValues(value.GetChildren(), e)
	if err != nil {
		return nil, err
	}
	return &ast.Value{
		T:        value.GetT(),
		Children: children,
	}, nil
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

func evaluateOp(op ast.Value_Type, operands []*ast.Value, e Environment) (*ast.Value, error) {
	switch op {
	case ast.Value_HASH:
		if len(operands) != 1 {
			return nil, fmt.Errorf("Invalid number of operands to HASH")
		}
		return evaluateHash(operands[0])
	case ast.Value_SERIALIZE:
		if len(operands) != 1 {
			return nil, fmt.Errorf("Invalid number of operands to HASH")
		}
		return evaluateSerialize(operands[0])
	case ast.Value_EQUAL:
		if len(operands) != 2 {
			return nil, fmt.Errorf("Invalid number of operands to EQUAL")
		}
		var result bool
		if operands[0].GetT() == ast.Value_PARAM && operands[1].GetT() != ast.Value_NIL {
			if err := e.IndexParam(int(operands[0].GetU()), operands[1]); err != nil {
				return nil, err
			}
			result = true
		} else if operands[1].GetT() == ast.Value_PARAM && operands[0].GetT() != ast.Value_NIL {
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
	case ast.Value_PLUS:
		if len(operands) != 2 {
			return nil, fmt.Errorf("Invalid number of operands to PLUS")
		}
		if operands[0].GetT() == ast.Value_UINT64 &&
			operands[1].GetT() == ast.Value_UINT64 {
			return &ast.Value{
				T: ast.Value_UINT64,
				Primitive: &ast.Value_U{
					U: operands[0].GetU() + operands[1].GetU(),
				},
			}, nil
		}
		a, err := valueToBigInt(operands[0])
		if err != nil {
			return nil, err
		}
		b, err := valueToBigInt(operands[1])
		if err != nil {
			return nil, err
		}
		return bigIntToValue(new(big.Int).Add(a, b)), nil
	case ast.Value_MINUS:
		if len(operands) != 2 {
			return nil, fmt.Errorf("Invalid number of operands to MINUS")
		}
		if operands[0].GetT() == ast.Value_UINT64 &&
			operands[1].GetT() == ast.Value_UINT64 {
			return &ast.Value{
				T: ast.Value_UINT64,
				Primitive: &ast.Value_U{
					U: operands[0].GetU() - operands[1].GetU(),
				},
			}, nil
		}
		a, err := valueToBigInt(operands[0])
		if err != nil {
			return nil, err
		}
		b, err := valueToBigInt(operands[1])
		if err != nil {
			return nil, err
		}
		return bigIntToValue(new(big.Int).Sub(a, b)), nil
	case ast.Value_AND:
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
	case ast.Value_SLICE:
		if len(operands) != 3 {
			return nil, fmt.Errorf("Invalid number of operands to SLICE")
		}
		if operands[0].GetT() != ast.Value_UINT64 ||
			operands[1].GetT() != ast.Value_UINT64 ||
			operands[2].GetT() != ast.Value_BYTES {
			return nil, fmt.Errorf("Invalid operand type to SLICE")
		}
		start := int(operands[0].GetU())
		end := int(operands[1].GetU())
		source := operands[2].GetRaw()
		result := make([]byte, end-start)
		if start > len(source) {
			return nil, fmt.Errorf("Invalid slice start: %d!", start)
		}
		if end > len(source) {
			end = len(source)
		}
		copy(result, source[start:end])
		return &ast.Value{
			T: ast.Value_BYTES,
			Primitive: &ast.Value_Raw{
				Raw: result,
			},
		}, nil
	case ast.Value_INDEX:
		if len(operands) != 2 {
			return nil, fmt.Errorf("Invalid number of operands to INDEX")
		}
		if operands[0].GetT() != ast.Value_UINT64 {
			return nil, fmt.Errorf("Invalid operand type to INDEX")
		}
		list, err := evaluateList(operands[1], e)
		if err != nil {
			return nil, err
		}
		i := int(operands[0].GetU())
		if i < 0 || i >= len(list) {
			return nil, fmt.Errorf("Index out of range!")
		}
		return list[i], nil
	}
	return nil, fmt.Errorf("Invalid op: %s", op.String())
}

func evaluateOpGet(field ast.Value_Type, value *ast.Value, e Environment) (*ast.Value, error) {
	if value.GetT() == ast.Value_NIL {
		// Running GET on NIL values always results in NIL
		return value, nil
	}
	switch field {
	case ast.Value_GET_CAPACITY:
		return value.GetChildren()[0], nil
	case ast.Value_GET_LOCK:
		return value.GetChildren()[1], nil
	case ast.Value_GET_TYPE:
		return value.GetChildren()[2], nil
	case ast.Value_GET_DATA:
		return value.GetChildren()[3], nil
	case ast.Value_GET_DATA_HASH:
		data := value.GetChildren()[3]
		if data.GetT() != ast.Value_BYTES {
			return nil, fmt.Errorf("Invalid data value type: %s", data.GetT().String())
		}
		h, err := rpctypes.CalculateHash(rpctypes.Raw(data.GetRaw()))
		if err != nil {
			return nil, err
		}
		return &ast.Value{
			T: ast.Value_BYTES,
			Primitive: &ast.Value_Raw{
				Raw: h,
			},
		}, nil
	case ast.Value_GET_CODE_HASH:
		return value.GetChildren()[0], nil
	case ast.Value_GET_HASH_TYPE:
		return value.GetChildren()[1], nil
	case ast.Value_GET_ARGS:
		return value.GetChildren()[2], nil
	}
	return nil, fmt.Errorf("Invalid get field: %s", field.String())
}

func evaluateList(list *ast.Value, e Environment) ([]*ast.Value, error) {
	switch list.GetT() {
	case ast.Value_LIST:
		return evaluateAstValues(list.GetChildren(), e)
	case ast.Value_MAP:
		if len(list.GetChildren()) != 2 {
			return nil, fmt.Errorf("Invalid number of lists for map!")
		}
		f := list.GetChildren()[0]
		list, err := evaluateList(list.GetChildren()[1], e)
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
	case ast.Value_QUERY_CELLS:
		return e.QueryCell(list)
	}
	return nil, fmt.Errorf("Invalid list type: %s", list.GetT().String())
}

func evaluateHash(value *ast.Value) (*ast.Value, error) {
	if value.GetT() == ast.Value_NIL {
		// TODO: Running HASH on NIL values always results in NIL, this might hit
		// problems in the future, ideally we should change this once conditionals
		// are better supported.
		return value, nil
	}
	switch value.GetT() {
	case ast.Value_SCRIPT:
		if err := ast.IsValidScript(value); err != nil {
			return nil, err
		}
		script := rpctypes.Script{
			HashType: rpctypes.ScriptHashType(value.GetChildren()[1].GetU()),
			Args:     rpctypes.Bytes(value.GetChildren()[2].GetRaw()),
		}
		copy(script.CodeHash[:], value.GetChildren()[0].GetRaw())
		h, err := rpctypes.CalculateHash(script)
		if err != nil {
			return nil, err
		}
		return &ast.Value{
			T: ast.Value_BYTES,
			Primitive: &ast.Value_Raw{
				Raw: h,
			},
		}, nil
	}
	return nil, fmt.Errorf("Invalid value type: %s, cannot calculate hash", value.GetT().String())
}

func evaluateSerialize(value *ast.Value) (*ast.Value, error) {
	switch value.GetT() {
	case ast.Value_TRANSACTION:
		tx, err := ast.RestoreTransaction(value, true)
		if err != nil {
			return nil, err
		}
		data, err := json.Marshal(tx)
		if err != nil {
			return nil, err
		}
		return &ast.Value{
			T: ast.Value_BYTES,
			Primitive: &ast.Value_Raw{
				Raw: data,
			},
		}, nil
	}
	return nil, fmt.Errorf("Invalid value type: %s", value.GetT().String())
}

func valueToBigInt(value *ast.Value) (*big.Int, error) {
	i := new(big.Int)
	if value.GetT() == ast.Value_BYTES {
		a := make([]byte, len(value.GetRaw()))
		copy(a, value.GetRaw())
		for i := len(a)/2 - 1; i >= 0; i-- {
			opp := len(a) - 1 - i
			a[i], a[opp] = a[opp], a[i]
		}
		i.SetBytes(a)
	} else if value.GetT() == ast.Value_UINT64 {
		i.SetUint64(value.GetU())
	} else {
		return nil, fmt.Errorf("Cannot convert value type %s to big int!", value.GetT().String())
	}
	return i, nil
}

func bigIntToValue(i *big.Int) *ast.Value {
	a := i.Bytes()
	for i := len(a)/2 - 1; i >= 0; i-- {
		opp := len(a) - 1 - i
		a[i], a[opp] = a[opp], a[i]
	}
	return &ast.Value{
		T: ast.Value_BYTES,
		Primitive: &ast.Value_Raw{
			Raw: a,
		},
	}
}
