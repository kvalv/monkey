package eval

import (
	"github.com/kvalv/monkey/ast"
	"github.com/kvalv/monkey/object"
	"github.com/kvalv/monkey/tracer"
)

var trace = (&tracer.Tracer{}).Trace

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch n := node.(type) {
	case *ast.LetStatement:
		defer trace("evalLetStatement")(nil)
		return evalLetStatement(n, env)
	case *ast.Program:
		defer trace("evalStatements")(nil)
		return evalStatements(n.Statements, env)
	case *ast.BlockStatement:
		defer trace("evalBlockStatement")(nil)
		return evalBlockStatement(n.Statements, env)
	case *ast.ExpressionStatement:
		defer trace("evalExpressionStatement")(nil)
		return Eval(n.Expr, env)
	case *ast.IfExpression:
		defer trace("evalIfExpression")(nil)
		return evalIfExpression(n, env)
	case *ast.ReturnExpression:
		defer trace("evalReturnExpression")(nil)
		return &object.Return{Object: Eval(n.Value, env)}
	case *ast.PrefixExpression:
		defer trace("evalPrefixExpression")(nil)
		return evalPrefixExpression(n, env)
	case *ast.InfixExpression:
		defer trace("evalInfixExpression")(nil)
		return evalInfixExpression(n, env)
	case *ast.Boolean:
		defer trace("evalBoolean")(nil)
		if n.Value {
			return object.TRUE
		}
		return object.FALSE
	case *ast.Number:
		defer trace("evalNumber")(nil)
		return &object.Integer{Value: int64(n.Value)}
	case *ast.Identifier:
		defer trace("evalIdentifier")(nil)
		return evalIdentifier(n, env)
	case *ast.FunctionLiteral:
		defer trace("evalFunctionLiteral")
		return evalFunctionLiteral(n, env)
	case *ast.CallExpression:
		defer trace("evalCallExpression")
		return evalCallExpression(n, env)
	}
	return object.Errorf("unable to evaluate node of type %T", node)
}

func evalStatements(stmts []ast.Statement, env *object.Environment) object.Object {
	var res object.Object
	for _, s := range stmts {
		res = Eval(s, env)
		if res == nil {
			return nil
		}
		if res.Type() == object.RETURN_OBJ {
			return res.(*object.Return).Object
		}
	}
	return res
}
