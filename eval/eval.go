package eval

import (
	"log"

	"github.com/kvalv/monkey/ast"
	"github.com/kvalv/monkey/object"
	"github.com/kvalv/monkey/tracer"
)

var trace = (&tracer.Tracer{}).Trace

func Eval(node ast.Node) object.Object {
	switch n := node.(type) {
	case *ast.Program:
		defer trace("evalStatements")(nil)
		return evalStatements(n.Statements)
	case *ast.BlockStatement:
		defer trace("evalBlockStatement")(nil)
		return evalBlockStatement(n.Statements)
	case *ast.ExpressionStatement:
		log.Printf("node is %s", node.TokenLiteral())
		defer trace("evalExpressionStatement")(nil)
		return Eval(n.Expr)
	case *ast.IfExpression:
		defer trace("evalIfExpression")(nil)
		return parseIfExpression(n)
	case *ast.ReturnExpression:
		defer trace("evalReturnExpression")(nil)
		return &object.Return{Object: Eval(n.Value)}
	case *ast.PrefixExpression:
		defer trace("evalPrefixExpression")(nil)
		return parsePrefixExpression(n)
	case *ast.InfixExpression:
		defer trace("evalInfixExpression")(nil)
		return parseInfixExpression(n)
	case *ast.Boolean:
		defer trace("evalBoolean")(nil)
		if n.Value {
			return object.TRUE
		}
		return object.FALSE
	case *ast.Number:
		defer trace("evalNumber")(nil)
		return &object.Integer{Value: int64(n.Value)}
	}
	return nil
}

func evalStatements(stmts []ast.Statement) object.Object {
	var res object.Object
	for _, s := range stmts {
		res = Eval(s)
		if res == nil {
			return nil
		}
		log.Printf("evaling %s", s)
		if res.Type() == object.RETURN_OBJ {
			return res.(*object.Return).Object
		}
	}
	return res
}
