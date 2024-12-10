package lex_test

import (
	"log"
	"os"
	"testing"

	"github.com/kvalv/monkey/token"

	"github.com/kvalv/monkey/lex"
)

func TestNextToken(t *testing.T) {
	input := `=+-,!*; != == foo fn return {} () 1 11 > < true false if else "hello" "hello world" "" x`
	l := lex.New(input)
	log.SetOutput(os.Stdout)
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
		{Type: token.IF, Literal: "if"},
		{Type: token.ELSE, Literal: "else"},
		{Type: token.STRING, Literal: `"hello"`},
		{Type: token.STRING, Literal: `"hello world"`},
		{Type: token.STRING, Literal: `""`},
		{Type: token.IDENT, Literal: `x`},
		{Type: token.EOF, Literal: ""},
	}
	for i, exp := range expected {
		if got := l.Next(); got.Type != exp.Type {
			t.Fatalf("%d: unexpected TokenType: expected %+v, got %+v", i, exp, got)
		} else if got.Literal != exp.Literal {
			t.Fatalf("%d: unexpected Literal: expected %+v, got %+v", i, exp, got)
		}
	}
}

func TestSpan(t *testing.T) {
	l := lex.New("the cat  sat fn @")
	expected := []token.Token{
		{Type: token.IDENT, Literal: "the", Span: token.Span{Start: 0, End: 3}},
		{Type: token.IDENT, Literal: "cat", Span: token.Span{Start: 4, End: 7}},
		{Type: token.IDENT, Literal: "sat", Span: token.Span{Start: 9, End: 12}},
		{Type: token.IDENT, Literal: "fn", Span: token.Span{Start: 13, End: 15}},
		{Type: token.IDENT, Literal: "@", Span: token.Span{Start: 16, End: 17}},
		{Type: token.EOF, Literal: "", Span: token.Span{Start: 17, End: 17}},
	}
	for _, exp := range expected {
		got := l.Next()
		if got.Span.Start != exp.Span.Start {
			t.Fatalf("span start mismatch: expected %d got %d - token = %+v", exp.Span.Start, got.Span.Start, got)
		}
		if got.Span.End != exp.Span.End {
			t.Fatalf("span end mismatch: expected %d got %d - token = %+v", exp.Span.End, got.Span.End, got)
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
		if got := l.Next(); got.Type != exp.Type {
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
		if got := l.Next(); got.Type != exp.Type {
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
		if got := l.Next(); got.Type != exp.Type {
			t.Fatalf("%d: unexpected TokenType: expected %+v, got %+v", i, exp, got)
		} else if got.Literal != exp.Literal {
			t.Fatalf("%d: unexpected Literal: expected %+v, got %+v", i, exp, got)
		}
	}
}
