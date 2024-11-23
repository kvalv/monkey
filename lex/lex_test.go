package lex_test

import (
	"testing"

	"github.com/kvalv/monkey/token"

	"github.com/kvalv/monkey/lex"
)

func TestNextToken(t *testing.T) {
	input := `=+-,! != foo fn return {} () 1 11 `
	l := lex.New(input)
	expected := []token.Token{
		{Ttype: token.EQ, Literal: "="},
		{Ttype: token.PLUS, Literal: "+"},
		{Ttype: token.MINUS, Literal: "-"},
		{Ttype: token.COMMA, Literal: ","},
		{Ttype: token.BANG, Literal: "!"},
		{Ttype: token.NEQ, Literal: "!="},
		{Ttype: token.IDENT, Literal: "foo"},
		{Ttype: token.FUNC, Literal: "fn"},
		{Ttype: token.RETURN, Literal: "return"},
		{Ttype: token.LBRACK, Literal: "{"},
		{Ttype: token.RBRACK, Literal: "}"},
		{Ttype: token.POPEN, Literal: "("},
		{Ttype: token.PCLOSE, Literal: ")"},
		{Ttype: token.INT, Literal: "1"},
		{Ttype: token.INT, Literal: "11"},
		{Ttype: token.EOF, Literal: ""},
	}
	for i, exp := range expected {
		if got := l.NextToken(); got.Ttype != exp.Ttype {
			t.Fatalf("%d: unexpected TokenType: expected %+v, got %+v", i, exp, got)
		} else if got.Literal != exp.Literal {
			t.Fatalf("%d: unexpected Literal: expected %+v, got %+v", i, exp, got)
		}
	}
}
