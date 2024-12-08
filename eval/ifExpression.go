package eval

import (
	"github.com/kvalv/monkey/ast"
	"github.com/kvalv/monkey/object"
)

func parseIfExpression(node *ast.IfExpression, env *object.Environment) object.Object {
	res := Eval(node.Cond, env)

	if isTruthy(res) {
		return Eval(node.Then, env)
	}
	if node.Else != nil {
		return Eval(node.Else, env)
	}
	return object.NULL
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case object.NULL:
		return false
	case object.FALSE:
		return false
	case object.TRUE:
		return true
	default:
		return true
	}
}
