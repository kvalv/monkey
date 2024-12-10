package eval

import (
	"github.com/kvalv/monkey/ast"
	"github.com/kvalv/monkey/object"
)

func evalString(node *ast.String, env *object.Environment) object.Object {
	return &object.String{Value: node.Value}
}
