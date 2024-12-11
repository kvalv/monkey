package eval

import (
	"github.com/kvalv/monkey/ast"
	"github.com/kvalv/monkey/object"
)

func evalArray(arr *ast.Array, env *object.Environment) object.Object {
	out := &object.Array{}
	for _, elem := range arr.Elems {
		res := Eval(elem, env)
		if object.IsError(res) {
			return res
		}
		out.Elems = append(out.Elems, res)
	}
	return out
}
