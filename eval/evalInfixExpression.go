package eval

import (
	"fmt"

	"github.com/kvalv/monkey/ast"
	"github.com/kvalv/monkey/object"
)

func evalStringInfixExpression(op string, left, right object.Object) object.Object {
	a := left.(*object.String).Value
	b := right.(*object.String).Value
	switch op {
	case "+":
		return &object.String{Value: fmt.Sprintf("%s%s", a, b)}
	case "==":
		return nativeBoolToBoolean(a == b)
	default:
		return object.Errorf("unknown operator: STRING %s STRING", op)
	}
}

func evalIntegerInfixExpression(op string, left, right object.Object) object.Object {
	a := left.(*object.Integer).Value
	b := right.(*object.Integer).Value
	switch op {
	case "+":
		return &object.Integer{Value: a + b}
	case "-":
		return &object.Integer{Value: a - b}
	case "*":
		return &object.Integer{Value: a * b}
	case "/":
		return &object.Integer{Value: a / b}
	case ">":
		return &object.Boolean{Value: a > b}
	case "<":
		return &object.Boolean{Value: a < b}
	default:
		return object.Errorf("unknown operator: %s", op)
	}
}

func evalInfixExpression(node *ast.InfixExpression, env *object.Environment) object.Object {
	lhs := Eval(node.Lhs, env)
	rhs := Eval(node.Rhs, env)

	if lhs.Type() != rhs.Type() {
		return object.Errorf("type mismatch: %s %s %s", lhs.Type(), node.Op, rhs.Type())
	}

	switch {
	case lhs.Type() == object.INTEGER_OBJ && rhs.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(node.Op, lhs, rhs)
	case lhs.Type() == object.STRING_OBJ && rhs.Type() == object.STRING_OBJ:
		return evalStringInfixExpression(node.Op, lhs, rhs)
	case node.Op == "==":
		return nativeBoolToBoolean(lhs == rhs)
	case node.Op == "!=":
		return nativeBoolToBoolean(lhs != rhs)
	default:
		return object.Errorf("unknown operator: %s %s %s", lhs.Type(), node.Op, rhs.Type())
	}
}
func nativeBoolToBoolean(b bool) object.Object {
	if b {
		return object.TRUE
	} else {
		return object.FALSE
	}
}
