package eval

import (
	"github.com/kvalv/monkey/ast"
	"github.com/kvalv/monkey/object"
)

func evalIdentifier(id *ast.Identifier, env *object.Environment) object.Object {
	if value, ok := env.Get(id.Literal); ok {
		return value
	}
	if value, ok := builtin[id.Literal]; ok {
		return value
	}
	return object.Errorf("identifier '%s' not defined", id.Literal)

}
