package parser_test

import (
	"testing"

	"github.com/kvalv/monkey/ast"
	"github.com/kvalv/monkey/parser"
)

func TestParsePrefix(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{
		{input: "!3;", want: "(!3);"},
		{input: "-foo;", want: "(-foo);"},
	}
	for i, tc := range cases {
		p := parser.New(tc.input)
		prog, errs := p.Parse()
		if len(errs) > 0 {
			t.Fatalf("got %d errors: %+v", len(errs), errs)
		}
		got := prog.Statements[0].String()
		if got != tc.want {
			t.Fatalf("[%d]: mismatch: got %q, want %q", i, got, tc.want)
		}
	}
}

func TestParseExpression(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{
		{input: "1 + 2 + 3", want: "((1 + 2) + 3);"},
		{input: "1 + 2 * 3", want: "(1 + (2 * 3));"},
	}
	for i, tc := range cases {
		p := parser.New(tc.input)
		prog, errs := p.Parse()
		if len(errs) > 0 {
			t.Fatalf("got %d errors: %+v", len(errs), errs)
		}
		if n := len(prog.Statements); n != 1 {
			t.Fatalf("expected 1 statement, got %d: prog.Statements=%+v", n, prog.Statements)
		}
		got := prog.Statements[0].String()
		if got != tc.want {
			t.Fatalf("[%d]: mismatch: got %q, want %q", i, got, tc.want)
		}
	}

}

func TestParse(t *testing.T) {
	t.Skip()
	input := `
    let x = y;
    let a = 3;
    let b = 123 + x;
    `
	p := parser.New(input)
	prog, errs := p.Parse()
	if len(errs) > 0 {
		t.Fatalf("got %d errors: %+v", len(errs), errs)
	}

	if prog == nil {
		t.Fatal("program is nil")
	}

	if n := len(prog.Statements); n != 3 {
		t.Fatalf("expected 2 statements, got %d", n)
	}
	expectLetStatement(t, prog.Statements[0], "x")
	expectLetStatement(t, prog.Statements[1], "a")
	expectLetStatement(t, prog.Statements[2], "b")
}

func expectLetStatement(t *testing.T, got ast.Statement, name string) {
	if got.TokenLiteral() != "let" {
		t.Fatalf("TokenLiteral mismatch: expected 'let' got %q", got.TokenLiteral())
	}

	stmt, ok := got.(*ast.LetStatement)
	if !ok {
		t.Fatal("not a let statement")
	}
	if stmt.Lhs.Value != name {
		t.Fatalf("value mismatch: expected %q, got %q", name, stmt.Lhs.Value)
	}
	// we won't check the rhs for now
}
