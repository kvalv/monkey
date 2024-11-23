package parser_test

import (
	"testing"

	"github.com/kvalv/monkey/ast"
	"github.com/kvalv/monkey/parser"
)

func TestParse(t *testing.T) {
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
