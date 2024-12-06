package eval

import (
	"github.com/kvalv/monkey/ast"
	"github.com/kvalv/monkey/object"
)

func parseInfixExpression(node *ast.InfixExpression) object.Object {
	lhs := Eval(node.Lhs)
	rhs := Eval(node.Rhs)

	a, _ := lhs.(*object.Integer)
	b, _ := rhs.(*object.Integer)
	switch node.Op {
	case "<":
		return nativeBoolToBoolean(a.Value < b.Value)
	case ">":
		return nativeBoolToBoolean(a.Value > b.Value)
	case "+":
		return &object.Integer{Value: a.Value + b.Value}
	case "-":
		return &object.Integer{Value: a.Value - b.Value}
	case "*":
		return &object.Integer{Value: a.Value * b.Value}
	case "/":
		return &object.Integer{Value: a.Value / b.Value}
	case "==":
		// object.TRUE and object.FALSE are singletons so we check that they point
		// to the same object
		return nativeBoolToBoolean(lhs == rhs)
	case "!=":
		return nativeBoolToBoolean(lhs != rhs)
	default:
		return object.NULL
	}
}
func nativeBoolToBoolean(b bool) object.Object {
	if b {
		return object.TRUE
	} else {
		return object.FALSE
	}
}
