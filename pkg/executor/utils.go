package executor

import (
	internal_ast "github.com/xxuejie/animagus/internal/ast"
	"github.com/xxuejie/animagus/pkg/ast"
)

type prependEnvironment struct {
	e    Environment
	args []*internal_ast.Value
}

func (e *prependEnvironment) Arg(i int) *internal_ast.Value {
	if i < len(e.args) {
		return e.args[i]
	}
	return e.e.Arg(i - len(e.args))
}

func (e *prependEnvironment) Param(i int) *internal_ast.Value {
	return e.e.Param(i)
}

func (e *prependEnvironment) QueryCell(query *ast.List) ([]*internal_ast.Value, error) {
	return e.e.QueryCell(query)
}
