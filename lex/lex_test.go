package lex_test

import (
	"testing"

	"github.com/kvalv/monkey/token"

	"github.com/kvalv/monkey/lex"
)

func TestNextToken(t *testing.T) {
	input := `=+-,!*; != == foo fn return {} () 1 11 > < true false`
	l := lex.New(input)
	expected := []token.Token{
		{Type: token.ASSIGN, Literal: "="},
		{Type: token.PLUS, Literal: "+"},
		{Type: token.MINUS, Literal: "-"},
		{Type: token.COMMA, Literal: ","},
		{Type: token.BANG, Literal: "!"},
		{Type: token.MUL, Literal: "*"},
		{Type: token.SEMICOLON, Literal: ";"},
		{Type: token.NEQ, Literal: "!="},
		{Type: token.EQ, Literal: "=="},
		{Type: token.IDENT, Literal: "foo"},
		{Type: token.FUNC, Literal: "fn"},
		{Type: token.RETURN, Literal: "return"},
		{Type: token.LBRACK, Literal: "{"},
		{Type: token.RBRACK, Literal: "}"},
		{Type: token.POPEN, Literal: "("},
		{Type: token.PCLOSE, Literal: ")"},
		{Type: token.INT, Literal: "1"},
		{Type: token.INT, Literal: "11"},
		{Type: token.GT, Literal: ">"},
		{Type: token.Lt, Literal: "<"},
		{Type: token.TRUE, Literal: "true"},
		{Type: token.FALSE, Literal: "false"},
		{Type: token.EOF, Literal: ""},
	}
	for i, exp := range expected {
		if got := l.NextToken(); got.Type != exp.Type {
			t.Fatalf("%d: unexpected TokenType: expected %+v, got %+v", i, exp, got)
		} else if got.Literal != exp.Literal {
			t.Fatalf("%d: unexpected Literal: expected %+v, got %+v", i, exp, got)
		}
	}
}

func TestPrefix(t *testing.T) {
	l := lex.New("!3;")
	expected := []token.Token{
		{Type: token.BANG, Literal: "!"},
		{Type: token.INT, Literal: "3"},
		{Type: token.SEMICOLON, Literal: ";"},
		{Type: token.EOF, Literal: ""},
	}
	for i, exp := range expected {
		if got := l.NextToken(); got.Type != exp.Type {
			t.Fatalf("%d: unexpected TokenType: expected %+v, got %+v", i, exp, got)
		} else if got.Literal != exp.Literal {
			t.Fatalf("%d: unexpected Literal: expected %+v, got %+v", i, exp, got)
		}
	}
}

func TestJustAIdentAndSemicolon(t *testing.T) {
	l := lex.New("foo;")
	expected := []token.Token{
		{Type: token.IDENT, Literal: "foo"},
		{Type: token.SEMICOLON, Literal: ";"},
	}
	for i, exp := range expected {
		if got := l.NextToken(); got.Type != exp.Type {
			t.Fatalf("%d: unexpected TokenType: expected %+v, got %+v", i, exp, got)
		} else if got.Literal != exp.Literal {
			t.Fatalf("%d: unexpected Literal: expected %+v, got %+v", i, exp, got)
		}
	}
}

func TestMore(t *testing.T) {
	input := `let x = y;`
	l := lex.New(input)
	expected := []token.Token{
		{Type: token.LET, Literal: "let"},
		{Type: token.IDENT, Literal: "x"},
		{Type: token.ASSIGN, Literal: "="},
		{Type: token.IDENT, Literal: "y"},
		{Type: token.SEMICOLON, Literal: ";"},
		{Type: token.EOF, Literal: ""},
	}
	for i, exp := range expected {
		if got := l.NextToken(); got.Type != exp.Type {
			t.Fatalf("%d: unexpected TokenType: expected %+v, got %+v", i, exp, got)
		} else if got.Literal != exp.Literal {
			t.Fatalf("%d: unexpected Literal: expected %+v, got %+v", i, exp, got)
		}
	}
}
