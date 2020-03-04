package executor

import (
	"github.com/xxuejie/animagus/pkg/ast"
)

type prependEnvironment struct {
	e    Environment
	args []*ast.Value
}

func (e *prependEnvironment) ReplaceArgs(args []*ast.Value) error {
	l := len(args)
	if len(e.args) < l {
		l = len(e.args)
	}
	for i := 0; i < l; i++ {
		e.args[i] = args[i]
	}
	if len(args) > len(e.args) {
		if err := e.ReplaceArgs(args[len(e.args):]); err != nil {
			return err
		}
	}
	return nil
}

func (e *prependEnvironment) Arg(i int) *ast.Value {
	if i < len(e.args) {
		return e.args[i]
	}
	return e.e.Arg(i - len(e.args))
}

func (e *prependEnvironment) Param(i int) *ast.Value {
	return e.e.Param(i)
}

func (e *prependEnvironment) IndexParam(i int, value *ast.Value) error {
	return e.e.IndexParam(i, value)
}

func (e *prependEnvironment) QueryCell(query *ast.Value) ([]*ast.Value, error) {
	return e.e.QueryCell(query)
}
