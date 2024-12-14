package eval

import (
	"github.com/kvalv/monkey/ast"
	"github.com/kvalv/monkey/object"
)

func evalAssignExpression(expr *ast.AssignExpression, env *object.Environment) object.Object {
	ai, ok := expr.Lhs.(*ast.ArrayIndex)
	if !ok {
		return object.Errorf("not implemented")
	}
	ident, ok := ai.Array.(*ast.Identifier)
	if !ok {
		return object.Errorf("only support identifiers")
	}
	arr, ok := env.Get(ident.String())
	if !ok {
		return object.Errorf("identifier '%s' not defined", arr.String())
	}
	hm := arr.(*object.Hash)

	hm.Set(Eval(ai.Index, env), Eval(expr.Rhs, env))

	return object.NULL
}
