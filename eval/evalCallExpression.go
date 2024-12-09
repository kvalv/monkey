package eval

import (
	"github.com/kvalv/monkey/ast"
	"github.com/kvalv/monkey/object"
)

func evalCallExpression(node *ast.CallExpression, env *object.Environment) object.Object {
	fn := Eval(node.Function, env).(*object.Function)
	if exp, got := len(fn.Params), len(node.Params); exp != got {
		return object.Errorf("Error invoking function: expected %d arguments but received %d", exp, got)
	}
	scoped := fn.Env.NewScope()
	for i, p := range fn.Params {
		value := Eval(node.Params[i], env)
		if object.IsError(value) {
			return value
		}
		scoped.Set(p.Literal, value)
	}
	return Eval(fn.Body, scoped)
}
