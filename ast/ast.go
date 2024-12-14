package ast

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/kvalv/monkey/token"
)

type (
	Node interface {
		TokenLiteral() string
		String() string
	}
	Expression interface {
		Node
		expr()
	}
	Statement interface {
		Node
		stmt()
	}
	Program struct {
		token.Token
		Statements []Statement
	}
	LetStatement struct {
		token.Token
		Lhs *Identifier
		Rhs Expression
	}
	AssignExpression struct {
		token.Token
		Lhs Expression
		Rhs Expression
	}
	ExpressionStatement struct {
		token.Token
		Expr Expression
	}
	Identifier struct {
		token.Token
		Value string
	}
	Number struct {
		token.Token
		Value int
	}
	Boolean struct {
		token.Token
		Value bool
	}
	String struct {
		token.Token
		Value string // we strip the quotes -> `"cat"` -> `cat`
	}
	PrefixExpression struct {
		token.Token
		Op  string
		Rhs Expression
	}
	InfixExpression struct {
		token.Token
		Op       string
		Lhs, Rhs Expression
	}
	BlockStatement struct {
		token.Token
		Statements []Statement
	}
	IfExpression struct {
		token.Token
		Cond       Expression
		Then, Else *BlockStatement
	}
	FunctionLiteral struct {
		token.Token
		Params []Identifier
		Body   *BlockStatement
	}
	CallExpression struct {
		token.Token
		Function Expression // identifier or FunctionLiteral
		Params   []Expression
	}
	ReturnExpression struct {
		token.Token
		Value Expression
	}
	Array struct {
		token.Token
		Elems []Expression
	}
	ArrayIndex struct {
		token.Token
		Array Expression // ident or array or Hash
		Index Expression // anything, but should evaluate to a number
	}
	// {"foo": "bar", true: false, 1: 3}
	HashLiteral struct {
		token.Token
		Pairs map[Expression]Expression
	}
)

func (n *Program) TokenLiteral() string { return n.Token.Literal }
func (n *Program) stmt()                {}
func (n *Program) String() string {
	buf := bytes.Buffer{}
	for _, stmt := range n.Statements {
		fmt.Fprintf(&buf, "%s", stmt)
	}
	return buf.String()
}

func (n *LetStatement) String() string {
	if n == nil {
		return "<LetStatement:nil>"
	}
	return fmt.Sprintf("let %s = %s", n.Lhs, n.Rhs)
}
func (n *LetStatement) TokenLiteral() string { return n.Token.Literal }
func (n *LetStatement) stmt()                {}

func (n *ExpressionStatement) TokenLiteral() string { return n.Token.Literal }
func (n *ExpressionStatement) stmt()                {}
func (n *ExpressionStatement) String() string       { return fmt.Sprintf("%s", n.Expr) }

func (n *Identifier) TokenLiteral() string { return n.Token.Literal }
func (n *Identifier) expr()                {}
func (n *Identifier) String() string {
	if n == nil {
		return "<Identifier:nil>"
	}
	return n.Value
}

func (n *Number) TokenLiteral() string { return n.Token.Literal }
func (n *Number) expr()                {}
func (n *Number) String() string       { return fmt.Sprintf("%d", n.Value) }

func (n *Boolean) TokenLiteral() string { return n.Token.Literal }
func (n *Boolean) expr()                {}
func (n *Boolean) String() string       { return fmt.Sprintf("%t", n.Value) }

func (n *String) TokenLiteral() string { return n.Token.Literal }
func (n *String) expr()                {}
func (n *String) String() string       { return fmt.Sprintf("%q", n.Value) }

func (n *PrefixExpression) TokenLiteral() string { return n.Token.Literal }
func (n *PrefixExpression) expr()                {}
func (n *PrefixExpression) String() string       { return fmt.Sprintf("(%s%s)", n.Op, n.Rhs) }

func (n *InfixExpression) TokenLiteral() string { return n.Token.Literal }
func (n *InfixExpression) expr()                {}
func (n *InfixExpression) String() string       { return fmt.Sprintf("(%s %s %s)", n.Lhs, n.Op, n.Rhs) }

func (n *BlockStatement) TokenLiteral() string { return n.Token.Literal }
func (n *BlockStatement) stmt()                {}
func (n *BlockStatement) String() string {
	if n == nil {
		return "<BlockStatement:nil>"
	}
	w := &bytes.Buffer{}
	fmt.Fprintf(w, "{")
	for _, stmt := range n.Statements {
		fmt.Fprintf(w, "%s", stmt)
	}
	fmt.Fprint(w, "}")
	return w.String()
}

func (n *IfExpression) TokenLiteral() string { return n.Token.Literal }
func (n *IfExpression) expr()                {}
func (n *IfExpression) String() string {
	w := &bytes.Buffer{}
	if n.Then == nil {
		panic("IfExpression: Then is nil")
	}
	fmt.Fprintf(w, "if %s %s", n.Cond.String(), n.Then.String())
	if n.Else != nil {
		fmt.Fprintf(w, " else %s", n.Else.String())
	}
	return w.String()
}

func (n *FunctionLiteral) TokenLiteral() string { return n.Token.Literal }
func (n *FunctionLiteral) expr()                {}
func (n *FunctionLiteral) String() string {
	var params []string
	for _, p := range n.Params {
		params = append(params, p.String())
	}
	return fmt.Sprintf("fn(%s) %s", strings.Join(params, ", "), n.Body.String())
}

func (n *CallExpression) TokenLiteral() string { return n.Token.Literal }
func (n *CallExpression) expr()                {}
func (n *CallExpression) String() string {
	if n == nil {
		return "<CallExpression:nil>"
	}
	var params []string
	for _, p := range n.Params {
		params = append(params, p.String())
	}
	return fmt.Sprintf("%s(%s)", n.Function, params)
}

func (n *ReturnExpression) TokenLiteral() string { return n.Token.Literal }
func (n *ReturnExpression) expr()                {}
func (n *ReturnExpression) String() string       { return fmt.Sprintf("return %s", n.Value.String()) }

func (a *Array) TokenLiteral() string { return a.Token.Literal }
func (a *Array) expr()                {}
func (a *Array) String() string {
	if a == nil {
		return "<Array:nil>"
	}
	var elems []string
	for _, elem := range a.Elems {
		elems = append(elems, elem.String())
	}
	return fmt.Sprintf("[%s]", strings.Join(elems, ", "))
}

func (a *ArrayIndex) TokenLiteral() string { return a.Token.Literal }
func (a *ArrayIndex) expr()                {}
func (a *ArrayIndex) String() string {
	if a == nil || a.Array == nil || a.Index == nil {
		return "<ArrayIndex:nil>"
	}
	return fmt.Sprintf("%s[%s]", a.Array.String(), a.Index.String())
}

func (h *HashLiteral) TokenLiteral() string { return h.Token.Literal }
func (h *HashLiteral) expr()                {}
func (h *HashLiteral) String() string {
	if h == nil {
		return "<HashLiteral:nil>"
	}
	var pairs []string
	for k, v := range h.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s: %s", k.String(), v.String()))
	}
	return fmt.Sprintf("{%s}", strings.Join(pairs, ", "))
}

func (ae *AssignExpression) TokenLiteral() string { return ae.Token.Literal }
func (ae *AssignExpression) expr()                {}
func (ae *AssignExpression) String() string {
	if ae == nil {
		return "<AssignExpression:nil>"
	}
	return fmt.Sprintf("%s=%s", ae.Lhs, ae.Rhs)
}
