package main

import (
	"encoding/json"
	"fmt"

	"github.com/xxuejie/animagus/pkg/ast"
)

type AstValueType interface {
	ToValue() (*ast.Value, error)
	// TODO: FromValue
}

type Root struct {
	Value AstValueType
}

func (r *Root) UnmarshalJSON(data []byte) error {
	typeStruct := struct {
		Type string `json:"type"`
	}{}
	err := json.Unmarshal(data, &typeStruct)
	if err != nil {
		return err
	}
	switch typeStruct.Type {
	case "arg":
		r.Value = &Arg{}
	case "equal":
		r.Value = &Equal{}
	case "get":
		r.Value = &Get{}
	default:
		return fmt.Errorf("Invalid type: %s", typeStruct.Type)
	}
	return json.Unmarshal(data, &r.Value)
}

type Arg struct {
	U int `json:"u"`
}

func (a *Arg) MarshalJSON() ([]byte, error) {
	type Alias Arg
	return json.Marshal(struct {
		Type string `json:"type"`
		*Alias
	}{
		Type:  "arg",
		Alias: (*Alias)(a),
	})
}

func (a *Arg) ToValue() (*ast.Value, error) {
	return &ast.Value{
		T: ast.Value_ARG,
		Primitive: &ast.Value_U{
			U: uint64(a.U),
		},
	}, nil
}

type Param struct {
	U int `json:"u"`
}

func (a *Param) MarshalJSON() ([]byte, error) {
	type Alias Param
	return json.Marshal(struct {
		Type string `json:"type"`
		*Alias
	}{
		Type:  "param",
		Alias: (*Alias)(a),
	})
}

func (a *Param) ToValue() (*ast.Value, error) {
	return &ast.Value{
		T: ast.Value_PARAM,
		Primitive: &ast.Value_U{
			U: uint64(a.U),
		},
	}, nil
}

type Get struct {
	Field string       `json:"field"`
	Value AstValueType `json:"value"`
}

func (a Get) MarshalJSON() ([]byte, error) {
	type Alias Get
	return json.Marshal(struct {
		Type string `json:"type"`
		Alias
	}{
		Type:  "get",
		Alias: Alias(a),
	})
}

func (a *Get) UnmarshalJSON(data []byte) error {
	type Alias Get
	aux := struct {
		Value Root `json:"value"`
		*Alias
	}{
		Alias: (*Alias)(a),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	a.Value = aux.Value.Value
	return nil
}

func (a *Get) ToValue() (*ast.Value, error) {
	var field ast.Field
	switch a.Field {
	case "capacity":
		field = ast.Field_CAPACITY
	default:
		return nil, fmt.Errorf("Invalid get field: %s", a.Field)
	}
	childValue, err := a.Value.ToValue()
	if err != nil {
		return nil, err
	}
	return &ast.Value{
		T: ast.Value_APPLY,
		Children: []*ast.Value{
			&ast.Value{
				T: ast.Value_OP,
				Primitive: &ast.Value_Op{
					Op: ast.Op_GET,
				},
			},
			&ast.Value{
				T: ast.Value_FIELD,
				Primitive: &ast.Value_Field{
					Field: field,
				},
			},
			childValue,
		},
	}, nil
}

type Equal struct {
	Args []AstValueType `json:"args"`
}

func (a Equal) MarshalJSON() ([]byte, error) {
	type Alias Equal
	return json.Marshal(struct {
		Type string `json:"type"`
		Alias
	}{
		Type:  "equal",
		Alias: Alias(a),
	})
}

func (a *Equal) UnmarshalJSON(data []byte) error {
	type Alias Equal
	aux := struct {
		Args []Root `json:"args"`
		*Alias
	}{
		Alias: (*Alias)(a),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	a.Args = make([]AstValueType, len(aux.Args))
	for i, arg := range aux.Args {
		a.Args[i] = arg.Value
	}
	return nil
}

func (a *Equal) ToValue() (*ast.Value, error) {
	children := make([]*ast.Value, len(a.Args)+1)
	children[0] = &ast.Value{
		T: ast.Value_OP,
		Primitive: &ast.Value_Op{
			Op: ast.Op_EQUAL,
		},
	}
	for i, arg := range a.Args {
		value, err := arg.ToValue()
		if err != nil {
			return nil, err
		}
		children[i+1] = value
	}
	return &ast.Value{
		T:        ast.Value_APPLY,
		Children: children,
	}, nil
}
