package eval

import (
	"github.com/kvalv/monkey/ast"
	"github.com/kvalv/monkey/object"
)

func evalFunctionLiteral(node *ast.FunctionLiteral, env *object.Environment) object.Object {
	paramNames := make(map[string]struct{})
	for _, p := range node.Params {
		if _, ok := paramNames[p.Literal]; ok {
			return object.Errorf("repeated argument %q", p.Literal)
		}
		paramNames[p.Literal] = struct{}{}
	}

	fn := &object.Function{
		Env:    env,
		Params: node.Params,
		Body:   node.Body,
	}
	return fn
}
