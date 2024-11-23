package ast

import "github.com/kvalv/monkey/token"

type Program struct {
	Statements []Statement
}

type Node interface {
	TokenLiteral() string
}

type Expression interface {
	Node
	expr()
}
type Statement interface {
	Node
	stmt()
}

type LetStatement struct {
	token.Token
	Lhs *Identifier
	Rhs Expression
}

func (n *LetStatement) TokenLiteral() string { return n.Token.Literal }
func (n *LetStatement) stmt()                {}

type Identifier struct {
	token.Token
	Value string
}

func (n *Identifier) TokenLiteral() string { return n.Token.Literal }
func (n *Identifier) expr()                {}

type Number struct {
	token.Token
	Value int
}

func (n *Number) TokenLiteral() string { return n.Token.Literal }
func (n *Number) expr()                {}
