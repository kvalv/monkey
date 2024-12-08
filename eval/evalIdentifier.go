package eval

import (
	"github.com/kvalv/monkey/ast"
	"github.com/kvalv/monkey/object"
)

func evalIdentifier(id *ast.Identifier, env *object.Environment) object.Object {
	value, ok := env.Get(id.Literal)
	if !ok {
		return object.Errorf("identifier '%s' not defined", id.Literal)
	}
	return value
}
