package eval

import (
	"github.com/kvalv/monkey/ast"
	"github.com/kvalv/monkey/object"
)

func evalBlockStatement(stmts []ast.Statement, env *object.Environment) object.Object {
	var res object.Object
	for _, s := range stmts {
		res = Eval(s, env)
		if res.Type() == object.RETURN_OBJ {
			return res
		}
	}
	return res
}
