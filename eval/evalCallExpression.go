package eval

import (
	"github.com/kvalv/monkey/ast"
	"github.com/kvalv/monkey/object"
)

func evalCallExpression(node *ast.CallExpression, env *object.Environment) object.Object {
	obj := Eval(node.Function, env)
	switch fn := obj.(type) {
	case *object.Error:
		return obj
	case *object.Function:
		return evalFunctionCallExpression(fn, node.Params, env)
	case *object.Builtin:
		return evalBuiltinCallExpression(fn, node.Params, env)
	default:
		return object.Errorf("evalCallExpression: unknown type %T", obj)
	}
}

func evalFunctionCallExpression(fn *object.Function, exprs []ast.Expression, env *object.Environment) object.Object {
	if exp, got := len(fn.Params), len(exprs); exp != got {
		return object.Errorf("Error invoking function: expected %d arguments but received %d", exp, got)
	}
	scoped := fn.Env.NewScope()
	params, ok := evalCallParams(exprs, env)
	if !ok {
		return params[0]
	}
	for i, p := range fn.Params {
		value := params[i]
		if object.IsError(value) {
			return value
		}
		scoped.Set(p.Literal, value)
	}
	return Eval(fn.Body, scoped)
}
func evalBuiltinCallExpression(fn *object.Builtin, exprs []ast.Expression, env *object.Environment) object.Object {
	params, ok := evalCallParams(exprs, env)
	if !ok {
		return params[0]
	}
	return fn.Fn(params...)
}

func evalCallParams(params []ast.Expression, env *object.Environment) ([]object.Object, bool) {
	var res []object.Object
	for _, p := range params {
		value := Eval(p, env)
		if object.IsError(value) {
			return []object.Object{value}, false
		}
		res = append(res, value)
	}
	return res, true
}
