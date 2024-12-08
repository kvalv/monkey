package eval

import (
	"github.com/kvalv/monkey/ast"
	"github.com/kvalv/monkey/object"
)

func evalLetStatement(node *ast.LetStatement, env *object.Environment) object.Object {
	value := Eval(node.Rhs, env)
	if object.IsError(value) {
		return value
	}
	env.Set(node.Lhs.Literal, value)
	return value
}
