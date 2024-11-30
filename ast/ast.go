package ast

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/kvalv/monkey/token"
)

type Node interface {
	TokenLiteral() string
	String() string
}

type Expression interface {
	Node
	expr()
}
type Statement interface {
	Node
	stmt()
}

type Program struct {
	token.Token
	Statements []Statement
}

func (n *Program) TokenLiteral() string { return n.Token.Literal }
func (n *Program) stmt()                {}
func (n *Program) String() string {
	buf := bytes.Buffer{}
	for _, stmt := range n.Statements {
		fmt.Fprintf(&buf, "%s", stmt)
	}
	return buf.String()
}

type LetStatement struct {
	token.Token
	Lhs *Identifier
	Rhs Expression
}

func (n *LetStatement) String() string       { return fmt.Sprintf("let %s = %s", n.Lhs, n.Rhs) }
func (n *LetStatement) TokenLiteral() string { return n.Token.Literal }
func (n *LetStatement) stmt()                {}

type ExpressionStatement struct {
	token.Token
	Expr Expression
}

func (n *ExpressionStatement) TokenLiteral() string { return n.Token.Literal }
func (n *ExpressionStatement) stmt()                {}
func (n *ExpressionStatement) String() string       { return fmt.Sprintf("%s", n.Expr) }

type Identifier struct {
	token.Token
	Value string
}

func (n *Identifier) TokenLiteral() string { return n.Token.Literal }
func (n *Identifier) expr()                {}
func (n *Identifier) String() string       { return n.Value }

type Number struct {
	token.Token
	Value int
}

func (n *Number) TokenLiteral() string { return n.Token.Literal }
func (n *Number) expr()                {}
func (n *Number) String() string       { return fmt.Sprintf("%d", n.Value) }

type Boolean struct {
	token.Token
	Value bool
}

func (n *Boolean) TokenLiteral() string { return n.Token.Literal }
func (n *Boolean) expr()                {}
func (n *Boolean) String() string       { return fmt.Sprintf("%t", n.Value) }

type PrefixExpression struct {
	token.Token
	Op  string
	Rhs Expression
}

func (n *PrefixExpression) TokenLiteral() string { return n.Token.Literal }
func (n *PrefixExpression) expr()                {}
func (n *PrefixExpression) String() string       { return fmt.Sprintf("(%s%s)", n.Op, n.Rhs) }

type InfixExpression struct {
	token.Token
	Op       string
	Lhs, Rhs Expression
}

func (n *InfixExpression) TokenLiteral() string { return n.Token.Literal }
func (n *InfixExpression) expr()                {}
func (n *InfixExpression) String() string       { return fmt.Sprintf("(%s %s %s)", n.Lhs, n.Op, n.Rhs) }

type BlockStatement struct {
	token.Token
	Statements []Statement
}

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

type IfExpression struct {
	token.Token
	Cond       Expression
	Then, Else *BlockStatement
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

type FunctionLiteral struct {
	token.Token
	Params []Identifier
	Body   *BlockStatement
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
