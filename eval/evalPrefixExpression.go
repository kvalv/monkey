package eval

import (
	"github.com/kvalv/monkey/ast"
	"github.com/kvalv/monkey/object"
)

func evalPrefixExpression(node *ast.PrefixExpression, env *object.Environment) object.Object {
	rhs := Eval(node.Rhs, env)
	switch node.Op {
	case "-":
		return evalMinusPrefixOperator(rhs)
	case "!":
		return evalBangPrefixOperator(rhs)
	default:
		return object.NULL
	}
}

func evalBangPrefixOperator(obj object.Object) object.Object {
	switch v := obj.(type) {
	case *object.Boolean:
		if v.Value == true {
			return object.FALSE
		}
		return object.TRUE
	case *object.Integer:
		if v.Value > 0 {
			return object.FALSE
		}
		return object.TRUE
	default:
		return object.NULL
	}
}
func evalMinusPrefixOperator(obj object.Object) object.Object {
	switch v := obj.(type) {
	case *object.Integer:
		return &object.Integer{Value: -v.Value}
	default:
		return object.NULL
	}
}
